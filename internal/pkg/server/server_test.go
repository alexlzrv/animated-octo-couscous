package server

import (
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server/config"
	"testing"
)

func TestStartListener(t *testing.T) {
	type args struct {
		c *config.ServerConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartListener(tt.args.c)
		})
	}
}
