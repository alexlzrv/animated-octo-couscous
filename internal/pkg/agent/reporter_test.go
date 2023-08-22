package agent_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
)

const (
	updatePathLength = 3
	gaugeBitSize     = 64
	counterBase      = 10
	counterBitSize   = 64
)

func TestSendReport(t *testing.T) {
	secretKey := []byte("secret")
	mtr := storage.NewMetrics()

	agent.UpdateMetrics(context.Background(), mtr)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		tokens := strings.FieldsFunc(req.URL.Path, func(c rune) bool {
			return c == '/'
		})
		if len(tokens) != updatePathLength {
			t.Errorf("Metric value not provided or url malformed: %s", req.URL.Path)
			http.Error(
				rw,
				fmt.Sprintf("Metric value not provided or url malformed: %s", req.URL.Path),
				http.StatusNotFound,
			)

			return
		}
		metricType := tokens[0]
		metricData := tokens[len(tokens)-1]

		switch {
		case metricType == metrics.GaugeMetricName:
			if _, err := strconv.ParseFloat(metricData, gaugeBitSize); err != nil {
				t.Errorf("Cannot save provided data: %s", metricData)
				http.Error(rw, fmt.Sprintf("Cannot save provided data: %s", metricData), http.StatusBadRequest)
			}

		case metricType == metrics.CounterMetricName:
			if _, err := strconv.ParseInt(metricData, counterBase, counterBitSize); err != nil {
				t.Errorf("Cannot save provided data: %s", metricData)
				http.Error(rw, fmt.Sprintf("Cannot save provided data: %s", metricData), http.StatusBadRequest)
			}
		default:
			t.Errorf("Metric type not implemented: %s", metricType)
			http.Error(
				rw,
				fmt.Sprintf("Metric type not implemented: %s", metricType),
				http.StatusNotImplemented,
			)
		}
	}))
	defer server.Close()

	agent.SendMetrics(context.Background(), mtr, server.URL, secretKey)
}
