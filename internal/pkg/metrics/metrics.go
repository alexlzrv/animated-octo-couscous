package metrics

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
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
	Delta *Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // хеш метрики
}

func (m *Metrics) SetHash(key string) {
	if key == "" {
		return
	}

	m.Hash = m.getHash(key)
}

func (m *Metrics) ValidateHash(key string) bool {
	if key == "" {
		return true
	}
	return m.Hash == m.getHash(key)
}

func (m *Metrics) getHash(key string) string {
	var metric string
	switch m.MType {
	case CounterMetricName:
		metric = fmt.Sprintf("%s:%s:%d", m.ID, CounterMetricName, *m.Delta)
	case GaugeMetricName:
		metric = fmt.Sprintf("%s:%s:%f", m.ID, CounterMetricName, *m.Value)
	default:
		logrus.Errorf("unknow metric type: %s", m.MType)
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(metric))

	return hex.EncodeToString(h.Sum(nil))
}
