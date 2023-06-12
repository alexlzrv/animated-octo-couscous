package server

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"
)

var tmpl = template.Must(template.New("index.html").Parse("html/index.gohtml"))

const (
	metricType     = "metricType"
	metricName     = "metricName"
	requestTimeout = 1 * time.Second
)

func RegisterHandlers(mux *chi.Mux, s storage.Store) {
	mux.Route("/", getAllMetricsHandler(s))
	mux.Route("/value/", getMetricHandler(s))
	mux.Route("/update/", updateHandler(s))
	mux.Route("/updates/", updatesBatchHandler(s))
	mux.Route("/ping", pingHandler(s))
}

func updateHandler(s storage.Store) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", updateMetricJSON(s))
		r.Post("/{metricType}/{metricName}/{metricValue}", updateMetricHandler(s))
	}
}

func getMetricHandler(s storage.Store) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", getMetricJSON(s))
		r.Get("/{metricType}/{metricName}", getMetric(s))
	}
}

func pingHandler(s storage.Store) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			if err := s.Ping(); err != nil {
				w.WriteHeader(http.StatusNotImplemented)
			}
			w.WriteHeader(http.StatusOK)
		})
	}
}

func getMetricJSON(s storage.Store) func(w http.ResponseWriter, r *http.Request) {
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
		requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
		defer requestCancel()

		logrus.Infof("Try get metric...%v %v", metric.ID, metric.MType)

		m, ok := s.GetMetric(requestContext, metric.ID, metric.MType)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			logrus.Errorf("Metric not found: %s", metric.ID)
			return
		}
		logrus.Infof("Get metric: %v %v", m.ID, m.MType)

		b, err := json.Marshal(m)
		if err != nil {
			logrus.Errorf("Cannot encode metric data: %q", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(b)
		if err != nil {
			logrus.Errorf("Cannot send request: %q", err)
		}
	}
}

func updatesBatchHandler(s storage.Store) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var metricBatch []*metrics.Metrics
			body, err := io.ReadAll(r.Body)
			if err != nil {
				logrus.Errorf("Error: %s", err)
			}

			err = json.Unmarshal(body, &metricBatch)
			if err != nil {
				logrus.Infof("Cannot decode provided data: %s", err)
				return
			}

			requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
			defer requestCancel()

			err = s.UpdateMetrics(requestContext, metricBatch)
			if err != nil {
				http.Error(w, "Failed to update metrics", http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
		})
	}
}

func updateMetricJSON(s storage.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
		defer requestCancel()

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
			if metric.Delta == nil {
				http.Error(w, "Delta is required field", http.StatusBadRequest)
			}
			err = s.UpdateCounterMetric(requestContext, metric.ID, *metric.Delta)
			if err != nil {
				http.Error(w, metric.MType, http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
		case metrics.GaugeMetricName:
			if metric.Value == nil {
				http.Error(w, "Value is required field", http.StatusBadRequest)
			}
			err = s.UpdateGaugeMetric(requestContext, metric.ID, *metric.Value)
			if err != nil {
				http.Error(w, metric.MType, http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, metric.MType, http.StatusNotImplemented)
		}
	}
}

func getAllMetricsHandler(s storage.Store) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
			defer requestCancel()

			w.Header().Set("Content-Type", "text/html")
			metricsData, err := s.GetMetrics(requestContext)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			err = tmpl.Execute(w, metricsData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
	}
}

func getMetric(s storage.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, metricType)
		metricName := chi.URLParam(r, metricName)

		var metricData string

		requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
		defer requestCancel()

		switch metricType {
		case metrics.CounterMetricName:
			metricDataCounter, ok := s.GetMetric(requestContext, metricName, metricType)
			if ok {
				metricData = strconv.FormatInt(int64(*metricDataCounter.Delta), 10)
			} else {
				http.Error(w, metricName, http.StatusNotFound)
				return
			}
		case metrics.GaugeMetricName:
			metricDataGauge, ok := s.GetMetric(requestContext, metricName, metricType)
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

func updateMetricHandler(s storage.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, metricType)
		metricName := chi.URLParam(r, metricName)
		metricValue := chi.URLParam(r, "metricValue")

		requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
		defer requestCancel()

		var err error
		switch metricType {
		case metrics.CounterMetricName:
			err = updateCounterMetric(requestContext, metricName, metricValue, s)
		case metrics.GaugeMetricName:
			err = updateGaugeMetric(requestContext, metricName, metricValue, s)
		default:
			http.Error(w, metricType, http.StatusNotImplemented)
		}
		if err != nil {
			http.Error(w, metricValue, http.StatusBadRequest)
		}
	}
}

func updateGaugeMetric(ctx context.Context, metricName string, valueMetric string, s storage.Store) error {
	val, err := strconv.ParseFloat(valueMetric, 64)
	if err == nil {
		return s.UpdateGaugeMetric(ctx, metricName, metrics.Gauge(val))
	}

	return err
}

func updateCounterMetric(ctx context.Context, metricName string, valueMetric string, s storage.Store) error {
	val, err := strconv.ParseInt(valueMetric, 10, 64)
	if err == nil {
		return s.UpdateCounterMetric(ctx, metricName, metrics.Counter(val))
	}

	return err
}
