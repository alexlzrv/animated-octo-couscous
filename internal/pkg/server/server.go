package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"net/http"
)

type Config struct {
	ServerAddress string `env:"ADDRESS"`
}

func NewServerConfig() *Config {
	return &Config{}
}

type MetricsServer struct {
	MetricsStore storage.Metric
}

func (s *MetricsServer) StartListener(c *Config) {
	mux := chi.NewRouter()
	RegisterHandlers(mux, s.MetricsStore)
	err := http.ListenAndServe(c.ServerAddress, mux)
	if err != nil {
		return
	}
}
