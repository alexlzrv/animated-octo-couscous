package storage

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/sirupsen/logrus"
)

const (
	driverName = "pgx"
)

type DBStore struct {
	connection *sql.DB
}

func NewDBMetrics(databaseDSN string) (*DBStore, error) {
	logrus.Info("Try open connect db...")

	db, err := sql.Open(driverName, databaseDSN)
	if err != nil {
		return nil, err
	}
	logrus.Info("Open connect db successfully.")

	dbCon := &DBStore{
		connection: db,
	}

	logrus.Info("Create database...")

	err = dbCon.createDB()
	if err != nil {
		logrus.Errorf("Failed to create db %v", err)
		return nil, err
	}
	logrus.Info("Create db successfully")

	return dbCon, nil
}

func (db *DBStore) createDB() error {
	logrus.Info("Create db gauge...")
	_, err := db.connection.Exec(`CREATE TABLE IF NOT EXISTS gauge(
    									metric_id VARCHAR (50) PRIMARY KEY,
    									metric_value DOUBLE PRECISION);`)
	if err != nil {
		return err
	}
	logrus.Info("Create db gauge successfully")

	logrus.Info("Create db counter...")
	_, err = db.connection.Exec(`CREATE TABLE IF NOT EXISTS counter(
    									metric_id VARCHAR (50) PRIMARY KEY,
    									metric_delta BIGINT);`)
	if err != nil {
		return err
	}
	logrus.Info("Create db counter successfully")

	return nil
}

func (db *DBStore) UpdateMetrics(ctx context.Context, metricsBatch []*metrics.Metrics) error {
	tx, err := db.connection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stateInsertGauge, err := tx.Prepare(`INSERT INTO gauge (metric_id, metric_value) VALUES ($1, $2)
						ON CONFLICT (metric_id) DO UPDATE SET metric_value = $2`)
	if err != nil {
		logrus.Errorf("Error with insert gauge: %v", err)
		return err
	}
	defer func(stateInsertGauge *sql.Stmt) {
		err = stateInsertGauge.Close()
		if err != nil {
			logrus.Errorf("Failed to close insert statement: %v", err)
		}
	}(stateInsertGauge)

	stateSelectCounter, err := tx.Prepare(`SELECT metric_delta FROM counter WHERE metric_id = $1`)
	if err != nil {
		logrus.Errorf("Error with select counter: %v", err)
		return err
	}
	defer func(stateSelectCounter *sql.Stmt) {
		err = stateSelectCounter.Close()
		if err != nil {
			logrus.Errorf("Failed to close select statement: %v", err)
		}
	}(stateSelectCounter)

	stateInsertCounter, err := tx.Prepare(`INSERT INTO counter (metric_id, metric_delta) VALUES ($1, $2)
						ON CONFLICT (metric_id) DO UPDATE SET metric_delta = $2`)
	if err != nil {
		logrus.Errorf("Error with insert counter: %v", err)
		return err
	}
	defer func(stateInsertCounter *sql.Stmt) {
		err = stateInsertCounter.Close()
		if err != nil {
			logrus.Errorf("Failed to close insert statement: %v", err)
		}
	}(stateInsertCounter)

	for _, metric := range metricsBatch {
		switch {
		case metric.MType == metrics.GaugeMetricName:
			if _, err = stateInsertGauge.Exec(metric.ID, *(metric.Value)); err != nil {
				if err = tx.Rollback(); err != nil {
					logrus.Errorf("enable rollback transaction: %v", err)
					return err
				}
				logrus.Errorf("Error with update gauge: %v", err)
				return err
			}
		case metric.MType == metrics.CounterMetricName:
			var counter metrics.Counter
			query := stateSelectCounter.QueryRow(metric.ID)
			err = query.Scan(&counter)
			if !errors.Is(err, nil) && !errors.Is(err, sql.ErrNoRows) {
				if err = tx.Rollback(); err != nil {
					logrus.Errorf("enable rollback transaction: %v", err)
					return err
				}
				logrus.Errorf("Error with scan select counter: %v", err)
				return err
			}

			counter += *(metric.Delta)

			if _, err = stateInsertCounter.Exec(metric.ID, counter); err != nil {
				if err = tx.Rollback(); err != nil {
					logrus.Errorf("enable rollback transaction: %v", err)
					return err
				}
				logrus.Errorf("Error with update counter: %v", err)
				return err
			}
		}
	}

	return tx.Commit()
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
		if err != nil {
			logrus.Errorf("Error with get counter: %v", err)
			return nil, false
		}
		metric.Delta = &counter

	case metrics.GaugeMetricName:
		var gauge metrics.Gauge
		row := db.connection.QueryRowContext(ctx,
			`SELECT metric_value FROM gauge WHERE metric_id = $1`, name)

		err := row.Scan(&gauge)
		if err != nil {
			logrus.Errorf("Error with get gauge: %v", err)
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
		err = rows.Close()
		if err != nil {
			logrus.Errorf("Couldn't close rows: %v", err)
		}
	}(counters)

	for counters.Next() {
		var counter metrics.Counter
		metric := metrics.Metrics{
			MType: metrics.GaugeMetricName,
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
		err = rows.Close()
		if err != nil {
			logrus.Errorf("Couldn't close rows: %v", err)
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

func (db *DBStore) Ping() error {
	return db.connection.Ping()
}

func (db *DBStore) Close() error {
	logrus.Info("Close database connection")
	return db.connection.Close()
}

func (db *DBStore) LoadMetrics(_ string) error {
	return nil
}

func (db *DBStore) SaveMetrics(_ string) error {
	return nil
}
