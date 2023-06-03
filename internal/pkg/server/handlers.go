package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

type StorageHandlers interface {
	UpdateCounterMetric(name string, value metrics.Counter) error
	ResetCounterMetric(name string) error
	UpdateGaugeMetric(name string, value metrics.Gauge) error
	GetMetric(name string) (*metrics.Metrics, bool)
	GetMetrics() map[string]*metrics.Metrics
}

var tmpl = template.Must(template.New("index.html").Parse("html/index.gohtml"))

const (
	metricType = "metricType"
	metricName = "metricName"
)

func RegisterHandlers(mux *chi.Mux, s StorageHandlers) {
	mux.Route("/", getAllMetricsHandler(s))
	mux.Route("/value/", getMetricHandler(s))
	mux.Route("/update/", updateHandler(s))
}

func updateHandler(s StorageHandlers) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", updateMetricJSON(s))
		r.Post("/{metricType}/{metricName}/{metricValue}", updateMetricHandler(s))
	}
}

func getMetricHandler(s StorageHandlers) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", getMetricJSON(s))
		r.Get("/{metricType}/{metricName}", getMetric(s))
	}
}

func getMetricJSON(s StorageHandlers) func(w http.ResponseWriter, r *http.Request) {
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

		m, ok := s.GetMetric(metric.ID)
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

func updateMetricJSON(s StorageHandlers) func(w http.ResponseWriter, r *http.Request) {
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
			err = s.UpdateCounterMetric(metric.ID, *metric.Delta)
		case metrics.GaugeMetricName:
			err = s.UpdateGaugeMetric(metric.ID, *metric.Value)
		default:
			http.Error(w, metric.MType, http.StatusNotImplemented)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func getAllMetricsHandler(s StorageHandlers) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			metricsData := s.GetMetrics()

			err := tmpl.Execute(w, metricsData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
	}
}

func getMetric(s StorageHandlers) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, metricType)
		metricName := chi.URLParam(r, metricName)

		var metricData string

		switch metricType {
		case metrics.CounterMetricName:
			metricDataCounter, ok := s.GetMetric(metricName)
			if ok {
				metricData = strconv.FormatInt(int64(*metricDataCounter.Delta), 10)
			} else {
				http.Error(w, metricName, http.StatusNotFound)
				return
			}
		case metrics.GaugeMetricName:
			metricDataGauge, ok := s.GetMetric(metricName)
			if ok {
				metricData = strconv.FormatFloat(float64(*metricDataGauge.Value), 'f', -1, 64)
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
			logrus.Errorf("error %v", err)
			return
		}
	}
}

func updateMetricHandler(s StorageHandlers) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, metricType)
		metricName := chi.URLParam(r, metricName)
		metricValue := chi.URLParam(r, "metricValue")
		var err error
		switch metricType {
		case metrics.CounterMetricName:
			err = updateCounterMetric(metricName, metricValue, s)
		case metrics.GaugeMetricName:
			err = updateGaugeMetric(metricName, metricValue, s)
		default:
			http.Error(w, metricType, http.StatusNotImplemented)
		}
		if err != nil {
			http.Error(w, metricValue, http.StatusBadRequest)
		}
	}
}

func updateGaugeMetric(metricName string, valueMetric string, s StorageHandlers) error {
	val, err := strconv.ParseFloat(valueMetric, 64)
	if err == nil {
		return s.UpdateGaugeMetric(metricName, metrics.Gauge(val))
	}

	return err
}

func updateCounterMetric(metricName string, valueMetric string, s StorageHandlers) error {
	val, err := strconv.ParseInt(valueMetric, 10, 64)
	if err == nil {
		return s.UpdateCounterMetric(metricName, metrics.Counter(val))
	}

	return err
}
