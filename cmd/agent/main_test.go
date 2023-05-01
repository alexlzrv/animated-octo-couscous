package main

import (
	"reflect"
	"testing"
)

func Test_createPostRequest(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createPostRequest(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("createPostRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createURL(t *testing.T) {
	type args struct {
		address     string
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createURL(tt.args.address, tt.args.metricType, tt.args.metricName, tt.args.metricValue); got != tt.want {
				t.Errorf("createURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMetrics(t *testing.T) {
	tests := []struct {
		name string
		want map[string]gauge
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sendMetrics(t *testing.T) {
	type args struct {
		metric map[string]gauge
		count  int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := sendMetrics(tt.args.metric, tt.args.count); (err != nil) != tt.wantErr {
				t.Errorf("sendMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
