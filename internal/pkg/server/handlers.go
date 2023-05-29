package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

var tmpl = template.Must(template.New("index.html").Parse("html/index.gohtml"))

const (
	metricType = "metricType"
	metricName = "metricName"
)

func RegisterHandlers(mux *chi.Mux, metricStorage storage.Storage) {
	mux.Route("/", getAllMetricsHandler(metricStorage))
	mux.Route("/value/", getMetricHandler(metricStorage))
	mux.Route("/update/", updateHandler(metricStorage))
}

func updateHandler(metricStorage storage.Storage) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", updateMetricJSON(metricStorage))
		r.Post("/{metricType}/{metricName}/{metricValue}", updateMetricHandler(metricStorage))
	}
}

func getMetricHandler(metricStorage storage.Storage) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", getMetricJSON(metricStorage))
		r.Get("/{metricType}/{metricName}", getMetric(metricStorage))
	}
}

func getMetricJSON(metricStorage storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric *metrics.Metrics
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logrus.Errorf("Error: %s", err)
		}

		err = json.Unmarshal(body, &metric)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			logrus.Errorf("Cannot decode provided data: %s", err)
			return
		}

		m, ok := metricStorage.GetMetric(metric.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			logrus.Errorf("Metric not found: %s", metric.ID)
			return
		}

		b, err := json.Marshal(m)
		if err != nil {
			logrus.Errorf("Cannot encode metric data: %q", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(b)
		if err != nil {
			logrus.Errorf("Cannot send request: %q", err)
		}
	}
}

func updateMetricJSON(metricStorage storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric metrics.Metrics
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logrus.Errorf("Error: %s", err)
		}
		err = json.Unmarshal(body, &metric)
		if err != nil {
			logrus.Infof("Cannot decode provided data: %s, %s", metric.ID, err)
			return
		}

		switch metric.MType {
		case metrics.CounterMetricName:
			err = metricStorage.UpdateCounterMetric(metric.ID, *metric.Delta)
		case metrics.GaugeMetricName:
			err = metricStorage.UpdateGaugeMetric(metric.ID, *metric.Value)
		default:
			http.Error(w, metric.MType, http.StatusNotImplemented)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func getAllMetricsHandler(metricStorage storage.Storage) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			metricsData := metricStorage.GetMetrics()

			err := tmpl.Execute(w, metricsData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
	}
}

func getMetric(metricStorage storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, metricType)
		metricName := chi.URLParam(r, metricName)

		var metricData string

		switch metricType {
		case metrics.CounterMetricName:
			metricDataCounter, ok := metricStorage.GetMetric(metricName)
			if ok {
				metricData = fmt.Sprintf("%d", *metricDataCounter.Delta)
			} else {
				http.Error(w, metricName, http.StatusNotFound)
				return
			}
		case metrics.GaugeMetricName:
			metricDataGauge, ok := metricStorage.GetMetric(metricName)
			if ok {
				metricData = fmt.Sprintf("%g", *metricDataGauge.Value)
			} else {
				http.Error(w, metricName, http.StatusNotFound)
				return
			}
		default:
			http.Error(w, metricType, http.StatusNotImplemented)
			return
		}
		_, err := w.Write([]byte(metricData))
		if err != nil {
			http.Error(w, metricName, http.StatusInternalServerError)
		}
	}
}

func updateMetricHandler(metricStorage storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, metricType)
		metricName := chi.URLParam(r, metricName)
		metricValue := chi.URLParam(r, "metricValue")
		var err error
		switch metricType {
		case metrics.CounterMetricName:
			err = updateCounterMetric(metricName, metricValue, metricStorage)
		case metrics.GaugeMetricName:
			err = updateGaugeMetric(metricName, metricValue, metricStorage)
		default:
			http.Error(w, metricType, http.StatusNotImplemented)
		}
		if err != nil {
			http.Error(w, metricValue, http.StatusBadRequest)
		}
	}
}

func updateGaugeMetric(metricName string, valueMetric string, metricsStore storage.Storage) error {
	val, err := strconv.ParseFloat(valueMetric, 64)
	if err == nil {
		return metricsStore.UpdateGaugeMetric(metricName, metrics.Gauge(val))
	}

	return err
}

func updateCounterMetric(metricName string, valueMetric string, metricsStore storage.Storage) error {
	val, err := strconv.ParseInt(valueMetric, 10, 64)
	if err == nil {
		return metricsStore.UpdateCounterMetric(metricName, metrics.Counter(val))
	}

	return err
}
