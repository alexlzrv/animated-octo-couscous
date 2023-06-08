package server

import (
	"database/sql"
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
	metricStore := initStore(c)
	mux := chi.NewRouter()
	mux.Use(
		middleware.Logger,
		logger.LoggingMiddleware,
		compress.CompressMiddleware,
	)

	db, err := sql.Open("pgx", c.DatabaseDSN)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	defer db.Close()

	RegisterHandlers(mux, metricStore, db)

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
		logrus.Fatal(err)
	}
}

func initStore(c *config.ServerConfig) *storage.MemoryStore {
	var metricStore *storage.MemoryStore
	if c.FileStoragePath != "" {
		metricStore = storage.NewMetricsFile(c.FileStoragePath, time.Duration(c.StoreInterval)*time.Second)
		defer metricStore.Close()
	} else {
		metricStore = storage.NewMetrics()
	}
	return metricStore
}
