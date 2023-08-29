package config

import (
	"reflect"
	"testing"
)

func TestNewAgentConfig(t *testing.T) {
	type fields struct {
		ServerAddress  string
		SignKey        string
		ReportInterval int
		PollInterval   int
		RateLimit      int
	}

	tests := []struct {
		name   string
		fields fields
		want   *AgentConfig
	}{
		{
			name: "check agent config",
			fields: fields{
				ServerAddress:  "localhost:8080",
				SignKey:        "",
				ReportInterval: 10,
				PollInterval:   2,
				RateLimit:      3,
			},
			want: &AgentConfig{
				ServerAddress:  "localhost:8080",
				SignKey:        "",
				ReportInterval: 10,
				PollInterval:   2,
				RateLimit:      3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAgentConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAgentConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
