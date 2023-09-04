package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/stretchr/testify/assert"
)

func TestDBStore_UpdateCounterMetric(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewDBStore(db)

	type args struct {
		ctx   context.Context
		name  string
		value metrics.Counter
	}

	type mockBehavior func(args args, value metrics.Counter)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    metrics.Counter
		wantErr bool
	}{
		{
			name: "ok",
			input: args{
				ctx:   ctx,
				name:  "test",
				value: 1,
			},
			want: 1,
			mock: func(args args, value metrics.Counter) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO counter").
					WithArgs(args.name, args.value).
					WillReturnRows(mock.NewRows([]string{"value"}).AddRow(value))

				mock.ExpectCommit()
			},
		},
		{
			name: "failed test",
			input: args{
				ctx:   ctx,
				name:  "test",
				value: 1,
			},
			want: 1,
			mock: func(args args, value metrics.Counter) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO counter").
					WithArgs(args.name, args.value).
					WillReturnRows(mock.NewRows([]string{"value"}).
						AddRow(value).
						RowError(0, errors.New("insert error")))

				mock.ExpectCommit()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.UpdateCounterMetric(tt.input.ctx, tt.input.name, tt.input.value)
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBStore_UpdateGaugeMetric(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewDBStore(db)

	type args struct {
		ctx   context.Context
		name  string
		value metrics.Gauge
	}

	type mockBehavior func(args args, value metrics.Gauge)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    metrics.Gauge
		wantErr bool
	}{
		{
			name: "ok",
			input: args{
				ctx:   ctx,
				name:  "test",
				value: 1.0,
			},
			want: 1.0,
			mock: func(args args, value metrics.Gauge) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO gauge").
					WithArgs(args.name, args.value).
					WillReturnRows(mock.NewRows([]string{"value"}).AddRow(value))

				mock.ExpectCommit()
			},
		},
		{
			name: "failed test",
			input: args{
				ctx:   ctx,
				name:  "test",
				value: 1.0,
			},
			want: 1,
			mock: func(args args, value metrics.Gauge) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO gauge").
					WithArgs(args.name, args.value).
					WillReturnRows(mock.NewRows([]string{"value"}).
						AddRow(value).
						RowError(0, errors.New("insert error")))

				mock.ExpectCommit()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.UpdateGaugeMetric(tt.input.ctx, tt.input.name, tt.input.value)
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBStore_Ping(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewDBStore(db)

	tests := []struct {
		name    string
		args    *sqlmock.ExpectedPing
		wantErr bool
	}{
		{
			name: "ok",
			args: mock.ExpectPing(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.Ping()
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBStore_Close(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewDBStore(db)

	tests := []struct {
		name    string
		args    *sqlmock.ExpectedClose
		wantErr bool
	}{
		{
			name: "ok",
			args: mock.ExpectClose(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.Close()
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
