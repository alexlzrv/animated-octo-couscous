package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func SendMetrics(m storage.Storage, serverAddress string) error {
	url := fmt.Sprintf("http://%s/update/", serverAddress)

	for _, v := range m.GetMetrics() {
		err := createPostRequest(url, v)
		if err != nil {
			return fmt.Errorf("error create post request %s", err)
		}
	}
	return nil
}

func createPostRequest(url string, metric *metrics.Metrics) error {
	body, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("error encoding metric %s", err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err = gz.Write(body); err != nil {
		return fmt.Errorf("error %s", err)
	}

	gz.Close()

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return fmt.Errorf("error send request %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error client %s", err)
	}

	io.Copy(io.Discard, resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Infof("%v", metric)
		return errors.New(resp.Status)
	}

	return nil
}
