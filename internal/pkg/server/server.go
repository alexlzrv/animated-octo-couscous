package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"net/http"
)

type MetricsServer struct {
	MetricsStore storage.Metric
	context      context.Context
}

func (s *MetricsServer) StartListener() {
	mux := chi.NewRouter()
	RegisterHandlers(mux, s.MetricsStore)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}
