package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"net/http"
	"strconv"
)

func RegisterHandlers(mux *chi.Mux, store storage.Metric) {
	mux.Route("/update/{metricType}/{metricName}/{metricValue}", UpdateMetricHandler(store))
}

func UpdateMetricHandler(metric storage.Metric) func(r chi.Router) {
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
	if val, err := strconv.ParseInt(valueMetric, 8, 64); err == nil {
		metricsStore.UpdateCounterMetric(metricName, storage.Counter(val))
	} else {
		return err
	}

	return nil
}
