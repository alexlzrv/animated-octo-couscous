package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		name   string
		metric string
		want   want
	}{
		{
			name:   "OK gauge update",
			metric: "/update/gauge/test/100.0",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "OK counter update",
			metric: "/update/counter/test/100",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "BAD gauge update",
			metric: "/update/gauge/test/none",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "BAD counter update",
			metric: "/update/counter/test/none",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "NotFound gauge update",
			metric: "/update/gauge/",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "NotFound counter update",
			metric: "/update/counter/",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "NotImplemented update",
			metric: "/update/unknown/test/1001",
			want: want{
				code: http.StatusNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.metric, nil)
			request.Header.Add("Content-Type", "text/plain")

			w := httptest.NewRecorder()
			UpdateMetricHandler(w, request)

			res := w.Result()
			assert.Equal(t, res.StatusCode, tt.want.code)

			defer res.Body.Close()

			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}
