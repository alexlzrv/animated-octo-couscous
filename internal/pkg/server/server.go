package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/logger"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server/config"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func StartListener(c *config.ServerConfig) {
	logrus.Info("Server is running...")
	metricStore := storage.NewMetrics()
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(logger.LoggingMiddleware)
	log.Println(metricStore)
	RegisterHandlers(mux, metricStore)

	err := http.ListenAndServe(c.ServerAddress, mux)

	if err != nil {
		logrus.Fatal(err)
	}
}
