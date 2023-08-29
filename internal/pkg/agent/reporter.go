package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent/config"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
)

func RunSendMetric(ctx context.Context, reportTicker *time.Ticker, c *config.AgentConfig, s storage.Store) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-reportTicker.C:
			ok, err := SendMetrics(ctx, s, c.ServerAddress, c.SignKeyByte)
			if err != nil {
				logrus.Errorf("Error send metrics %v", err)
			}
			if ok {
				if err = s.ResetCounterMetric(ctx, "PollCount"); err != nil {
					logrus.Errorf("Error reset metrics %v", err)
				}
			}
		}
	}
}

func SendMetrics(ctx context.Context, s storage.Store, serverAddress string, signKey []byte) (bool, error) {
	metricsMap, err := s.GetMetrics(ctx)
	if err != nil {
		logrus.Errorf("Some error ocured during metrics get: %q", err)
		return false, err
	}

	metricsBatch := make([]*metrics.Metrics, 0)
	for _, v := range metricsMap {
		metricsBatch = append(metricsBatch, v)
	}

	url := fmt.Sprintf("http://%s/updates/", serverAddress)

	if err = SendBatchJSON(url, metricsBatch, signKey); err != nil {
		return false, fmt.Errorf("error create post request %w", err)
	}
	return true, nil
}

func SendBatchJSON(url string, metricsBatch []*metrics.Metrics, signKey []byte) error {
	body, err := json.Marshal(metricsBatch)
	if err != nil {
		return fmt.Errorf("error encoding metric %w", err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err = gz.Write(body); err != nil {
		return fmt.Errorf("error %w", err)
	}

	err = gz.Close()
	if err != nil {
		return fmt.Errorf("error close %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return fmt.Errorf("error send request %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	if signKey != nil {
		h := hmac.New(sha256.New, signKey)
		h.Write(body)
		serverHash := hex.EncodeToString(h.Sum(nil))
		req.Header.Set("HashSHA256", serverHash)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error client %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %v and not 200", resp.StatusCode)
	}

	return nil
}
