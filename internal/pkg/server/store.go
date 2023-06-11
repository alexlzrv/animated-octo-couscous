package server

import (
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server/config"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
	"time"
)

func InitStore(c *config.ServerConfig) (storage.Store, error) {
	var (
		store storage.Store
		err   error
	)

	if c.DatabaseDSN != "" {
		store, err = storage.NewDBMetrics(c.DatabaseDSN)
		if err != nil {
			logrus.Errorf("Failed to create database store: %v", err)
		}
	} else if c.FileStoragePath != "" {
		store, err = storage.NewMetricsFile(c.FileStoragePath, time.Duration(c.StoreInterval)*time.Second)
		if err != nil {
			logrus.Errorf("Failed to create file storage: %v", err)
		}
	} else {
		store = storage.NewMetrics()
	}

	return store, nil
}
