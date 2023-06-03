package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"net/http"
)

type StorageReport interface {
	GetMetrics() map[string]*metrics.Metrics
}

func SendMetrics(s StorageReport, serverAddress string) error {
	url := fmt.Sprintf("http://%s/update/", serverAddress)

	for _, v := range s.GetMetrics() {
		err := createPostRequest(url, v)
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
