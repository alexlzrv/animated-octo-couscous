package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
	"io"
	"log"
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

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error send request %s", err)
	}
	log.Println(req)
	req.Header.Set("Content-Type", "application/json")

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
