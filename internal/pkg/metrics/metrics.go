package metrics

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
}
