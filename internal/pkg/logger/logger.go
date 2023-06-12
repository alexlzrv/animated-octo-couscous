package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t := time.Now()
		logrus.Infof("requset with method %s: %s, duration %s", r.RequestURI, r.Method, time.Since(t))

		logrus.Infof("response status : %d, size : %d", ww.Status(), ww.BytesWritten())
		next.ServeHTTP(ww, r)
	})
}

func HTTPRequestLogger() func(next http.Handler) http.Handler {
	logger := httplog.NewLogger("metrics")
	return httplog.RequestLogger(logger)
}
