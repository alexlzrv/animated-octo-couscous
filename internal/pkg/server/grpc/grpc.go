package grpc

import (
	"errors"
	"fmt"
	pb "github.com/mayr0y/animated-octo-couscous.git/api/server"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"io"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer
}

func (s *Server) UpdateMetrics(stream pb.Metrics_UpdateMetricsServer) error {
	var metric metrics.Metrics
	metricsSlice := make([]*metrics.Metrics, 0)
	for {
		message, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		switch message.Metric.MType {
		case metrics.GaugeMetricName:
			gaugeValue := metrics.Gauge(message.Metric.Value)
			metric = metrics.Metrics{
				ID:    message.Metric.ID,
				MType: message.Metric.MType,
				Value: &gaugeValue,
			}
		case metrics.CounterMetricName:
			counterValue := metrics.Counter(message.Metric.Delta)
			metric = metrics.Metrics{
				ID:    message.Metric.ID,
				MType: message.Metric.MType,
				Delta: &counterValue,
			}
		default:
			err := fmt.Errorf("unknown metric type: %s", message.Metric.MType)
			return stream.SendAndClose(&pb.UpdateMetricResponse{Error: err.Error()})
		}

		metricsSlice = append(metricsSlice, &metric)
	}

	if err := s.metricsStore.UpdateMetrics(stream.Context(), metricsSlice); err != nil {
		return stream.SendAndClose(&pb.UpdateMetricResponse{Error: err.Error()})
	}

	return stream.SendAndClose(&pb.UpdateMetricResponse{Error: "Metrics are updated"})
}
