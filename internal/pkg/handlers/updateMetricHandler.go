package handlers

import (
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"net/http"
	"strconv"
	"strings"
)

func UpdateMetricHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	res.Header().Set("content-type", "text/plain")
	statusCode := parseURLandUpdateMetrics(req)
	res.WriteHeader(statusCode)
}

func parseURLandUpdateMetrics(req *http.Request) int {
	var body string
	body += req.URL.Path
	bodySplit := strings.Split(body, "/")

	if len(bodySplit) < 5 {
		return http.StatusNotFound
	}

	typeMetric := bodySplit[2]
	nameMetric := bodySplit[3]
	valueMetric := bodySplit[4]

	return updateMetrics(typeMetric, nameMetric, valueMetric)
}

func updateMetrics(typeMetric string, nameMetric string, valueMetric string) int {

	memStorage := storage.NewMetrics()

	switch typeMetric {
	default:
		return http.StatusNotImplemented
	case storage.CounterMetricName:
		if val, err := strconv.ParseInt(valueMetric, 8, 64); err == nil {
			memStorage.CounterMetrics[nameMetric] += storage.Counter(val)
		} else {
			return http.StatusBadRequest
		}
	case storage.GaugeMetricName:
		if val, err := strconv.ParseFloat(valueMetric, 64); err == nil {
			memStorage.GaugeMetrics[nameMetric] = storage.Gauge(val)
		} else {
			return http.StatusBadRequest
		}
	}

	return http.StatusOK
}
