package main

import (
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, handlers.UpdateMetricHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
