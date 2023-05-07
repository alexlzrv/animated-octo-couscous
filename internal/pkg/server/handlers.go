package server

import (
	_ "embed"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"html/template"
	"net/http"
	"strconv"
)

//go:embed html/index.html
var templateFile string

func RegisterHandlers(mux *chi.Mux, metric storage.Metric) {
	mux.Route("/update/{metricType}/{metricName}/{metricValue}", updateMetricHandler(metric))
	mux.Route("/value/{metricType}/{metricName}", getMetricHandler(metric))
	mux.Route("/", getAllMetricsHandler(metric))
}

func getAllMetricsHandler(metric storage.Metric) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			metricsData := struct {
				Gauge   map[string]storage.Gauge
				Counter map[string]storage.Counter
			}{
				Gauge:   metric.GetGaugeMetrics(),
				Counter: metric.GetCounterMetrics(),
			}

			tmpl, err := template.New("index.html").Parse(templateFile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = tmpl.Execute(w, metricsData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		})
	}
}

func getMetricHandler(metric storage.Metric) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			metricType := chi.URLParam(r, "metricType")
			metricName := chi.URLParam(r, "metricName")

			var ok bool
			var metricData string
			switch metricType {
			case storage.CounterMetricName:
				var metricDataCounter storage.Counter
				metricDataCounter, ok = metric.GetCounterMetric(metricName)
				metricData = fmt.Sprintf("%d", metricDataCounter)
			case storage.GaugeMetricName:
				var metricDataGauge storage.Gauge
				metricDataGauge, ok = metric.GetGaugeMetric(metricName)
				metricData = fmt.Sprintf("%g", metricDataGauge)
			default:
				http.Error(w, metricType, http.StatusNotImplemented)
			}
			if ok {
				_, err := w.Write([]byte(metricData))
				if err != nil {
					http.Error(w, metricName, http.StatusInternalServerError)
				}
			} else {
				http.Error(w, metricName, http.StatusNotFound)
			}
		})
	}
}

func updateMetricHandler(metric storage.Metric) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			metricType := chi.URLParam(r, "metricType")
			metricName := chi.URLParam(r, "metricName")
			metricValue := chi.URLParam(r, "metricValue")

			var err error
			switch metricType {
			case storage.CounterMetricName:
				err = updateCounterMetric(metricName, metricValue, metric)
			case storage.GaugeMetricName:
				err = updateGaugeMetric(metricName, metricValue, metric)
			default:
				http.Error(w, metricType, http.StatusNotImplemented)
			}
			if err != nil {
				http.Error(w, metricValue, http.StatusBadRequest)
			}
		})
	}
}

func updateGaugeMetric(metricName string, valueMetric string, metricsStore storage.Metric) error {
	if val, err := strconv.ParseFloat(valueMetric, 64); err == nil {
		metricsStore.UpdateGaugeMetric(metricName, storage.Gauge(val))
	} else {
		return err
	}

	return nil
}

func updateCounterMetric(metricName string, valueMetric string, metricsStore storage.Metric) error {
	if val, err := strconv.ParseInt(valueMetric, 10, 64); err == nil {
		metricsStore.UpdateCounterMetric(metricName, storage.Counter(val))
	} else {
		return err
	}

	return nil
}
