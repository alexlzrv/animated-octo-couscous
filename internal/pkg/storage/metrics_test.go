package storage

import (
	"reflect"
	"testing"
)

func TestMetrics_GetCounterMetric(t *testing.T) {
	type fields struct {
		GaugeMetrics   map[string]Gauge
		CounterMetrics map[string]Counter
	}
	type args struct {
		metricName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Counter
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				GaugeMetrics:   tt.fields.GaugeMetrics,
				CounterMetrics: tt.fields.CounterMetrics,
			}
			got, got1 := m.GetCounterMetric(tt.args.metricName)
			if got != tt.want {
				t.Errorf("GetCounterMetric() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetCounterMetric() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMetrics_GetCounterMetrics(t *testing.T) {
	type fields struct {
		GaugeMetrics   map[string]Gauge
		CounterMetrics map[string]Counter
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]Counter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				GaugeMetrics:   tt.fields.GaugeMetrics,
				CounterMetrics: tt.fields.CounterMetrics,
			}
			if got := m.GetCounterMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCounterMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_GetGaugeMetric(t *testing.T) {
	type fields struct {
		GaugeMetrics   map[string]Gauge
		CounterMetrics map[string]Counter
	}
	type args struct {
		metricName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Gauge
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				GaugeMetrics:   tt.fields.GaugeMetrics,
				CounterMetrics: tt.fields.CounterMetrics,
			}
			got, got1 := m.GetGaugeMetric(tt.args.metricName)
			if got != tt.want {
				t.Errorf("GetGaugeMetric() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetGaugeMetric() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMetrics_GetGaugeMetrics(t *testing.T) {
	type fields struct {
		GaugeMetrics   map[string]Gauge
		CounterMetrics map[string]Counter
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]Gauge
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				GaugeMetrics:   tt.fields.GaugeMetrics,
				CounterMetrics: tt.fields.CounterMetrics,
			}
			if got := m.GetGaugeMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGaugeMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_UpdateCounterMetric(t *testing.T) {
	type fields struct {
		GaugeMetrics   map[string]Gauge
		CounterMetrics map[string]Counter
	}
	type args struct {
		metricName  string
		metricValue Counter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				GaugeMetrics:   tt.fields.GaugeMetrics,
				CounterMetrics: tt.fields.CounterMetrics,
			}
			m.UpdateCounterMetric(tt.args.metricName, tt.args.metricValue)
		})
	}
}

func TestMetrics_UpdateGaugeMetric(t *testing.T) {
	type fields struct {
		GaugeMetrics   map[string]Gauge
		CounterMetrics map[string]Counter
	}
	type args struct {
		metricName  string
		metricValue Gauge
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				GaugeMetrics:   tt.fields.GaugeMetrics,
				CounterMetrics: tt.fields.CounterMetrics,
			}
			m.UpdateGaugeMetric(tt.args.metricName, tt.args.metricValue)
		})
	}
}

func TestNewMetrics(t *testing.T) {
	tests := []struct {
		name string
		want *Metrics
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
