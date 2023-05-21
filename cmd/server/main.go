package main

import (
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server/config"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
)

func main() {
	cfg := config.NewServerConfig()
	config.Init(cfg)

	metricServer := server.MetricsServer{
		MetricsStore: storage.NewMetrics(),
	}
	metricServer.StartListener(cfg)
}
