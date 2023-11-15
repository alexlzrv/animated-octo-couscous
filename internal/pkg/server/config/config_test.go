package config

import (
	"reflect"
	"testing"
)

func TestNewServerConfig(t *testing.T) {
	type fields struct {
		ServerAddress   string
		StoreInterval   int
		FileStoragePath string
		Restore         bool
		DatabaseDSN     string
		SignKey         string
	}

	tests := []struct {
		name   string
		fields fields
		want   *ServerConfig
	}{
		{
			name: "check server config",
			fields: fields{
				ServerAddress:   "localhost:8080",
				StoreInterval:   300,
				FileStoragePath: "/tmp/metrics-db.json",
				Restore:         true,
				DatabaseDSN:     "",
				SignKey:         "",
			},
			want: &ServerConfig{
				ServerAddress:   "localhost:8080",
				StoreInterval:   300,
				FileStoragePath: "/tmp/metrics-db.json",
				Restore:         true,
				DatabaseDSN:     "",
				SignKey:         "",
			},
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := NewServerConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
