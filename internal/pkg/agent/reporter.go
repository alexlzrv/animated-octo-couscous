package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

func SendMetrics(ctx context.Context, s storage.Store, serverAddress string) error {
	url := fmt.Sprintf("http://%s/update/", serverAddress)

	metricMap, err := s.GetMetrics(ctx)
	if err != nil {
		logrus.Errorf("Error with get metrics: %v", err)
	}

	for _, v := range metricMap {
		err = createPostRequest(url, v)
		if err != nil {
			return fmt.Errorf("error create post request %v", err)
		}
	}
	return nil
}

func createPostRequest(url string, metric *metrics.Metrics) error {
	body, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("error encoding metric %v", err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err = gz.Write(body); err != nil {
		return fmt.Errorf("error %s", err)
	}

	gz.Close()

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return fmt.Errorf("error send request %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error client %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("status code not 200")
	}

	return nil
}
