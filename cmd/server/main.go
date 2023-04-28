package main

import (
	"net/http"
	"strconv"
	"strings"
)

type gauge float64
type counter int64

type MemStorage struct {
	gaugeData   map[string]gauge
	counterData map[string]counter
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, responseHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func responseHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	res.Header().Set("content-type", "text/plain")
	statusCode := parseURL(req)
	res.WriteHeader(statusCode)
}

func parseURL(req *http.Request) int {
	var body string
	body += req.URL.Path
	bodySplit := strings.Split(body, "/")

	if len(bodySplit) < 5 {
		return http.StatusNotFound
	}

	typeMetric := bodySplit[2]
	nameMetric := bodySplit[3]
	valueMetric := bodySplit[4]

	return checkMetric(typeMetric, nameMetric, valueMetric)
}

func checkMetric(typeMetric string, nameMetric string, valueMetric string) int {

	switch typeMetric {

	default:
		return http.StatusNotImplemented
	case "counter":
		if _, err := strconv.ParseInt(valueMetric, 8, 64); err != nil {
			return http.StatusBadRequest
		}

	case "gauge":
		if _, err := strconv.ParseFloat(valueMetric, 64); err != nil {
			return http.StatusBadRequest
		}
	}

	return http.StatusOK
}
