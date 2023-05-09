package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server/config"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"log"
	"net/http"
)

type MetricsServer struct {
	MetricsStore storage.Metric
}

func (s *MetricsServer) StartListener(c *config.ServerConfig) {
	mux := chi.NewRouter()
	RegisterHandlers(mux, s.MetricsStore)
	err := http.ListenAndServe(c.ServerAddress, mux)

	if err != nil {
		log.Fatal(err)
	}
}
