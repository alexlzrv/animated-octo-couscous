package grpc

import (
	"context"
	pb "github.com/mayr0y/animated-octo-couscous.git/api/server"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	Address      string
	metricsStore storage.Store
	pb.UnimplementedMetricsServer
}

func (s *Server) Start(ctx context.Context, storage storage.Store) error {
	s.metricsStore = storage
	listen, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMetricsServer(grpcServer, s)
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	return grpcServer.Serve(listen)
}
