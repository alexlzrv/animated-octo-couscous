package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	GaugeMetricName   = "gauge"
	CounterMetricName = "counter"
)

type Gauge float64
type Counter int64

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Value *Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Delta *Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
}

func (m *Metrics) EncodeMetric() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	jsonEncoder := json.NewEncoder(&buf)

	if err := jsonEncoder.Encode(m); err != nil {
		return nil, fmt.Errorf("error encode %w", err)
	}

	return &buf, nil
}

func (m *Metrics) String() string {
	switch m.MType {
	case GaugeMetricName:
		return fmt.Sprintf("%g", *(m.Value))
	case CounterMetricName:
		return fmt.Sprintf("%d", *(m.Delta))
	default:
		return ""
	}
}
