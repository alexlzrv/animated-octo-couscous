package storage

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

const (
	driverName = "pgx"
)

type DBStore struct {
	connection *sql.DB
}

func NewDBMetrics(databaseDSN string) (*DBStore, error) {
	var dbCon DBStore

	db, err := sql.Open(driverName, databaseDSN)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer db.Close()

	dbCon = DBStore{
		connection: db,
	}
	err = dbCon.createDB()
	if err != nil {
		logrus.Errorf("Failed to create db: %v", err)
		return nil, err
	}

	return &dbCon, nil
}

func (db *DBStore) createDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		CREATE TABLE IF NOT EXISTS metrics(
			id TEXT NOT NULL,
			mtype TEXT NOT NULL,
			delta BIGINT,
			value DOUBLE PRECISION,
			PRIMARY KEY(id, mtype)
		);`

	_, err := db.connection.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *DBStore) UpdateCounterMetric(ctx context.Context, name string, value metrics.Counter) error {
	var counter metrics.Counter

	row := db.connection.QueryRowContext(ctx,
		`SELECT metric_delta FROM counter WHERE metric_id = $1`, name)

	err := row.Scan(&counter)
	if !errors.Is(err, nil) && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	counter += value
	_, err = db.connection.ExecContext(ctx,
		`INSERT INTO counter (metric_id, metric_delta) VALUES ($1, $2)
				ON CONFLICT (metric_id) DO UPDATE SET metric_delta = $2`,
		name, counter)

	return err
}

func (db *DBStore) ResetCounterMetric(ctx context.Context, name string) error {
	var zero metrics.Counter
	_, err := db.connection.ExecContext(ctx,
		`INSERT INTO counter (metric_id, metric_delta) VALUES ($1, $2) 
			ON CONFLICT (metric_id) DO UPDATE SET metric_delta = $2`,
		name, zero)

	return err
}

func (db *DBStore) UpdateGaugeMetric(ctx context.Context, name string, value metrics.Gauge) error {
	_, err := db.connection.ExecContext(ctx,
		`INSERT INTO gauge (metric_id, metric_value) VALUES ($1, $2)
				ON CONFLICT (metric_id) DO UPDATE SET metric_value = $2`,
		name, value)

	return err
}

func (db *DBStore) GetMetric(ctx context.Context, name string, metricType string) (*metrics.Metrics, bool) {
	metric := metrics.Metrics{
		ID:    name,
		MType: metricType,
	}

	switch metricType {
	case metrics.CounterMetricName:
		var counter metrics.Counter
		row := db.connection.QueryRowContext(ctx,
			`SELECT metric_delta FROM counter WHERE metric_id = $1`, name)

		err := row.Scan(&counter)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, false
		case !errors.Is(err, nil):
			return nil, false
		}
		metric.Delta = &counter
	case metrics.GaugeMetricName:
		var gauge metrics.Gauge
		row := db.connection.QueryRowContext(ctx,
			`SELECT metric_value FROM gauge WHERE metric_id = $1`, name)

		err := row.Scan(&gauge)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, false
		case !errors.Is(err, nil):
			return nil, false
		}
		metric.Value = &gauge
	default:
		return nil, false
	}

	return &metric, true
}

func (db *DBStore) GetMetrics(ctx context.Context) (map[string]*metrics.Metrics, error) {
	metricsMap := make(map[string]*metrics.Metrics)

	counters, err := db.connection.QueryContext(ctx,
		`SELECT metric_id,metric_delta FROM counter`)

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Couldn't close rows: %q", err)
		}
	}(counters)

	for counters.Next() {
		var counter metrics.Counter
		metric := metrics.Metrics{
			MType: metrics.CounterMetricName,
			Delta: &counter,
		}
		err = counters.Scan(&metric.ID, metric.Delta)
		if err != nil {
			return nil, err
		}

		metricsMap[metric.ID] = &metric
	}

	err = counters.Err()
	if err != nil {
		return nil, err
	}

	gauges, err := db.connection.QueryContext(ctx,
		`SELECT metric_id,metric_value FROM gauge`)

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Couldn't close rows: %q", err)
		}
	}(gauges)

	for gauges.Next() {
		var gauge metrics.Gauge
		metric := metrics.Metrics{
			MType: metrics.GaugeMetricName,
			Value: &gauge,
		}

		err = gauges.Scan(&metric.ID, metric.Value)
		if err != nil {
			return nil, err
		}

		metricsMap[metric.ID] = &metric
	}

	err = gauges.Err()
	if err != nil {
		return nil, err
	}

	return metricsMap, nil
}

func (db *DBStore) Ping(ctx context.Context) error {
	return db.connection.PingContext(ctx)
}

func (db *DBStore) LoadMetrics(_ string) error {
	return nil
}

func (db *DBStore) SaveMetrics(_ string) error {
	return nil
}
