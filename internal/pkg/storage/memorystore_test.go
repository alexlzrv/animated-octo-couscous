package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/stretchr/testify/assert"
)

const (
	testMetrics      = `{"Alloc":{"id":"Alloc","type":"gauge","value":1336312}}`
	testMetrics2     = `{"Alloc":{"id":"Alloc","type":"gauge","value":1336313}}`
	testMetricValue  = 145544
	testMetricValue2 = 457855
)

func TestInMemoryStore_UpdateCounterMetric(t *testing.T) {
	metricsCache := make(map[string]*metrics.Metrics)
	testMetricName := "Alloc"
	metricValue := metrics.Counter(testMetricValue)

	type fields struct {
		metricsCache map[string]*metrics.Metrics
	}
	type args struct {
		in0        context.Context
		metricName string
		metricData metrics.Counter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "TestUpdateCounterMetric",
			fields: fields{
				metricsCache: metricsCache,
			},
			args: args{
				in0:        context.Background(),
				metricName: testMetricName,
				metricData: metricValue,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &storage.MemoryStore{
				Metrics: tt.fields.metricsCache,
			}
			tt.wantErr(t, m.UpdateCounterMetric(tt.args.in0, tt.args.metricName, tt.args.metricData), fmt.Sprintf("UpdateCounterMetric(%v, %v, %v)", tt.args.in0, tt.args.metricName, tt.args.metricData))
		})
	}
}

func TestInMemoryStore_ResetCounterMetric(t *testing.T) {
	metricsCache := make(map[string]*metrics.Metrics)
	testMetricName := "Alloc"

	type fields struct {
		metricsCache map[string]*metrics.Metrics
	}
	type args struct {
		in0        context.Context
		metricName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "TestResetCounterMetric",
			fields: fields{
				metricsCache: metricsCache,
			},
			args: args{
				in0:        context.Background(),
				metricName: testMetricName,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &storage.MemoryStore{
				Metrics: tt.fields.metricsCache,
			}
			tt.wantErr(t, m.ResetCounterMetric(tt.args.in0, tt.args.metricName), fmt.Sprintf("ResetCounterMetric(%v, %v)", tt.args.in0, tt.args.metricName))
		})
	}
}

func TestInMemoryStore_UpdateGaugeMetric(t *testing.T) {
	metricsCache := make(map[string]*metrics.Metrics)
	testMetricName := "Alloc"
	testMetricValue := metrics.Gauge(testMetricValue)

	type fields struct {
		metricsCache map[string]*metrics.Metrics
	}
	type args struct {
		in0        context.Context
		metricName string
		metricData metrics.Gauge
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "TestUpdateGaugeMetric",
			fields: fields{
				metricsCache: metricsCache,
			},
			args: args{
				in0:        context.Background(),
				metricName: testMetricName,
				metricData: testMetricValue,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &storage.MemoryStore{
				Metrics: tt.fields.metricsCache,
			}
			tt.wantErr(t, m.UpdateGaugeMetric(tt.args.in0, tt.args.metricName, tt.args.metricData), fmt.Sprintf("UpdateGaugeMetric(%v, %v, %v)", tt.args.in0, tt.args.metricName, tt.args.metricData))
		})
	}
}

func TestInMemoryStore_UpdateMetrics(t *testing.T) {
	metricsCache := make(map[string]*metrics.Metrics)
	testMetricName := "Alloc"
	testMetricValue := metrics.Gauge(testMetricValue)

	type fields struct {
		metricsCache map[string]*metrics.Metrics
	}
	type args struct {
		in0          context.Context
		metricsBatch []*metrics.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "TestUpdateMetrics",
			fields: fields{
				metricsCache: metricsCache,
			},
			args: args{
				in0: context.Background(),
				metricsBatch: []*metrics.Metrics{
					{
						ID:    testMetricName,
						Value: &testMetricValue,
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &storage.MemoryStore{
				Metrics: tt.fields.metricsCache,
			}
			tt.wantErr(t, m.UpdateMetrics(tt.args.in0, tt.args.metricsBatch), fmt.Sprintf("UpdateMetrics(%v, %v)", tt.args.in0, tt.args.metricsBatch))
		})
	}
}

func TestInMemoryStore_GetMetric(t *testing.T) {
	metricsCache := make(map[string]*metrics.Metrics)
	testMetricName := "Alloc"
	testMetricValue := metrics.Gauge(testMetricValue)

	type fields struct {
		metricsCache map[string]*metrics.Metrics
	}
	type args struct {
		in0        context.Context
		metricName string
		in2        string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *metrics.Metrics
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "TestGetMetric",
			fields: fields{
				metricsCache: metricsCache,
			},
			args: args{
				in0:        context.Background(),
				metricName: testMetricName,
			},
			want: &metrics.Metrics{
				ID:    testMetricName,
				Value: &testMetricValue,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &storage.MemoryStore{
				Metrics: tt.fields.metricsCache,
			}
			got, ok := m.GetMetric(tt.args.in0, tt.args.metricName, tt.args.in2)
			if !ok {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMetric(%v, %v, %v)", tt.args.in0, tt.args.metricName, tt.args.in2)
		})
	}
}

func TestInMemoryStore_GetMetrics(t *testing.T) {
	metricsCache := make(map[string]*metrics.Metrics)

	type fields struct {
		metricsCache map[string]*metrics.Metrics
	}
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*metrics.Metrics
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "TestGetMetrics",
			fields: fields{
				metricsCache: metricsCache,
			},
			args: args{
				in0: context.Background(),
			},
			want: metricsCache,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &storage.MemoryStore{
				Metrics: tt.fields.metricsCache,
			}
			got, err := m.GetMetrics(tt.args.in0)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMetrics(%v)", tt.args.in0)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMetrics(%v)", tt.args.in0)
		})
	}
}

func TestFileStore_LoadMetrics(t *testing.T) {
	f := "tests"
	testMetricName := "Alloc"
	metricValue := metrics.Gauge(testMetricValue)

	metricsCache := make(map[string]*metrics.Metrics)
	type fields struct {
		file         string
		metricsCache map[string]*metrics.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		want   metrics.Gauge
	}{
		{
			name: testMetricName,
			fields: fields{
				file:         f,
				metricsCache: metricsCache,
			},
			want: metricValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &storage.MemoryStore{
				FileStoragePath: tt.fields.file,
				Metrics:         tt.fields.metricsCache,
			}

			if err := fs.LoadMetrics(f); err != nil {
				t.Errorf("LoadMetrics() failed (error = %v)", err)
			}
		})
	}
}

func TestFileStore_SaveMetrics(t *testing.T) {
	f := "tests"
	testMetricName := "Alloc"
	metricValue := metrics.Gauge(testMetricValue)

	metricsCache := make(map[string]*metrics.Metrics)
	metric := metrics.Metrics{
		ID:    testMetricName,
		MType: metrics.GaugeMetricName,
		Value: &metricValue,
	}
	metricsCache[testMetricName] = &metric

	metricValue2 := metrics.Gauge(testMetricValue2)
	metricsCache2 := make(map[string]*metrics.Metrics)
	metric2 := metrics.Metrics{
		ID:    testMetricName,
		MType: metrics.GaugeMetricName,
		Value: &metricValue2,
	}
	metricsCache2[testMetricName] = &metric2

	type fields struct {
		file         string
		metricsCache map[string]*metrics.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Empty test",
			fields: fields{
				file:         f,
				metricsCache: metricsCache2,
			},
			want: testMetrics2,
		},
		{
			name: testMetricName,
			fields: fields{
				file:         f,
				metricsCache: metricsCache,
			},
			want: testMetrics,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &storage.MemoryStore{
				FileStoragePath: tt.fields.file,
				Metrics:         tt.fields.metricsCache,
			}
			err := fs.SaveMetrics(f)
			if err != nil {
				t.Errorf("SaveMetrics() failed (error = %v)", err)
			}
		})
	}
}
