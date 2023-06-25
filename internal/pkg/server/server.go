package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/compress"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/logger"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server/config"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func StartListener(c *config.ServerConfig) {
	logrus.Info("Init store...")
	logrus.Infof("ServerAddress: %v", c.ServerAddress)
	logrus.Infof("StoreInterval: %v", c.StoreInterval)
	logrus.Infof("Restore: %v", c.Restore)
	logrus.Infof("FileStoragePath: %v", c.FileStoragePath)

	var (
		metricStore storage.Store
		err         error
	)

	if c.DatabaseDSN != "" {
		metricStore, err = storage.NewDBMetrics(c.DatabaseDSN)
	} else if c.FileStoragePath != "" {
		metricStore, err = storage.NewMetricsFile(c.FileStoragePath, time.Duration(c.StoreInterval)*time.Second)
	} else {
		metricStore = storage.NewMetrics()
	}

	if err != nil {
		logrus.Errorf("Error init store: %v", err)
		return
	}

	defer metricStore.Close()

	logrus.Info("Init store successfully")

	mux := chi.NewRouter()
	mux.Use(
		middleware.Logger,
		logger.LoggingMiddleware,
		logger.HTTPRequestLogger(),
		compress.CompressMiddleware,
	)

	RegisterHandlers(mux, metricStore, c.SignKey)

	if c.Restore {
		if err = metricStore.LoadMetrics(c.FileStoragePath); err != nil {
			logrus.Errorf("Error update metric from file %v", err)
		}
	}
	if c.StoreInterval > 0 {
		storeInterval := time.NewTicker(time.Duration(c.StoreInterval) * time.Second)
		defer storeInterval.Stop()
		go func() {
			for range storeInterval.C {
				err = metricStore.SaveMetrics(c.FileStoragePath)
				if err != nil {
					logrus.Errorf("Error save metric from file %v", err)
				}
			}
		}()
	}

	logrus.Info("Server is running...")
	err = http.ListenAndServe(c.ServerAddress, mux)

	if err != nil {
		logrus.Fatalf("Error with server running: %v", err)
		return
	}
}
