package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

func NewDBStore(db *sql.DB) *DBStore {
	return &DBStore{connection: db}
}

func NewDBMetrics(databaseDSN string) (*DBStore, error) {
	db, err := sql.Open(driverName, databaseDSN)
	if err != nil {
		return nil, err
	}

	dbCon := &DBStore{
		connection: db,
	}

	err = dbCon.createDB()
	if err != nil {
		logrus.Errorf("Failed to create db %v", err)
		return nil, err
	}

	return dbCon, nil
}

func (db *DBStore) createDB() error {
	_, err := db.connection.Exec(`CREATE TABLE IF NOT EXISTS gauge(
    									metric_id VARCHAR (50) PRIMARY KEY,
    									metric_value DOUBLE PRECISION);`)
	if err != nil {
		logrus.Errorf("Error with create gauge db: %v", err)
		return err
	}

	_, err = db.connection.Exec(`CREATE TABLE IF NOT EXISTS counter(
    									metric_id VARCHAR (50) PRIMARY KEY,
    									metric_delta BIGINT);`)
	if err != nil {
		logrus.Errorf("Error with create counter db: %v", err)
		return err
	}

	return nil
}

func (db *DBStore) UpdateMetrics(ctx context.Context, metricsBatch []*metrics.Metrics) error {
	tx, err := db.connection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryGauge := `INSERT INTO gauge (metric_id, metric_value) VALUES %s
						ON CONFLICT (metric_id) DO UPDATE SET metric_value = $2
						RETURNING metric_id, metric_value`

	queryCounter := `INSERT INTO counter (metric_id, metric_delta) VALUES %s
						ON CONFLICT (metric_id) DO UPDATE SET metric_delta = EXCLUDED.metric_delta + counter.metric_delta
						RETURNING metric_id, metric_delta`

	metricMap := make(map[string]*metrics.Metrics, len(metricsBatch))
	var metricArgsCounter []string
	var metricArgsGauge []string
	var argsCounter []interface{}
	var argsGauge []interface{}

	for _, metric := range metricsBatch {
		if value, ok := metricMap[metric.ID]; ok && metric.MType == metrics.CounterMetricName {
			counter := *metric.Delta + *value.Delta
			metrics := metrics.Metrics{
				ID:    metric.ID,
				MType: metric.MType,
				Delta: &counter,
				Value: metric.Value,
			}
			*metricMap[metric.ID] = metrics
			continue
		}
		metricMap[metric.ID] = metric
	}

	counterI := 0
	gaugeI := 0
	for _, v := range metricMap {
		switch {
		case v.MType == metrics.CounterMetricName:
			metricArgsCounter = append(metricArgsCounter, fmt.Sprintf("($%d, $%d)", counterI*4+1, counterI*4+2))
			argsCounter = append(argsCounter, v.ID)
			argsCounter = append(argsCounter, v.Delta)
			counterI++
		case v.MType == metrics.GaugeMetricName:
			metricArgsGauge = append(metricArgsGauge, fmt.Sprintf("($%d, $%d)", gaugeI*4+1, gaugeI*4+2))
			argsGauge = append(argsGauge, v.ID)
			argsGauge = append(argsGauge, v.Value)
			gaugeI++
		}
	}

	for _, metric := range metricMap {
		switch {
		case metric.MType == metrics.GaugeMetricName:
			queryGauge = fmt.Sprintf(queryGauge, strings.Join(metricArgsGauge, ","))
			_, err = tx.ExecContext(ctx, queryGauge, argsGauge...)
			if err != nil {
				return err
			}
		case metric.MType == metrics.CounterMetricName:
			queryCounter = fmt.Sprintf(queryCounter, strings.Join(metricArgsCounter, ","))
			_, err = tx.ExecContext(ctx, queryCounter, argsCounter...)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (db *DBStore) UpdateCounterMetric(ctx context.Context, name string, value metrics.Counter) error {
	tx, err := db.connection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertCounter := `INSERT INTO counter (metric_id, metric_delta) VALUES ($1, $2)
						ON CONFLICT (metric_id) DO UPDATE SET metric_delta = EXCLUDED.metric_delta + counter.metric_delta
						RETURNING metric_id, metric_delta`

	if _, err = tx.ExecContext(ctx, insertCounter, name, value); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DBStore) ResetCounterMetric(ctx context.Context, name string) error {
	tx, err := db.connection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var zero metrics.Counter
	insertCounter := `INSERT INTO counter (metric_id, metric_delta) VALUES ($1, $2)
				ON CONFLICT (metric_id) DO UPDATE SET metric_delta = $2
						RETURNING metric_id, metric_delta`

	if _, err = tx.ExecContext(ctx, insertCounter, name, zero); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DBStore) UpdateGaugeMetric(ctx context.Context, name string, value metrics.Gauge) error {
	tx, err := db.connection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertGauge := `INSERT INTO gauge (metric_id, metric_value) VALUES ($1, $2)
				ON CONFLICT (metric_id) DO UPDATE SET metric_value = $2
				RETURNING metric_id, metric_value`

	if _, err = tx.ExecContext(ctx, insertGauge, name, value); err != nil {
		return err
	}

	return tx.Commit()
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
		if !errors.Is(err, nil) && !errors.Is(err, sql.ErrNoRows) {
			logrus.Errorf("Error with get counter: %v", err)
			return nil, false
		}
		metric.Delta = &counter

	case metrics.GaugeMetricName:
		var gauge metrics.Gauge
		row := db.connection.QueryRowContext(ctx,
			`SELECT metric_value FROM gauge WHERE metric_id = $1`, name)

		err := row.Scan(&gauge)
		if !errors.Is(err, nil) && !errors.Is(err, sql.ErrNoRows) {
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
		if !errors.Is(err, nil) && !errors.Is(err, sql.ErrNoRows) {
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
		if !errors.Is(err, nil) && !errors.Is(err, sql.ErrNoRows) {
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
