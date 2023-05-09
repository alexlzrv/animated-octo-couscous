package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type test struct {
	name   string
	method string
	metric string
	want   want
}

type want struct {
	code int
	data string
}

var tests = []test{
	{
		name:   "OK gauge update",
		metric: "/update/gauge/test/100.000000",
		method: http.MethodPost,
		want: want{
			code: http.StatusOK,
		},
	},
	{
		name:   "OK counter update",
		metric: "/update/counter/test/100",
		method: http.MethodPost,
		want: want{
			code: http.StatusOK,
		},
	},
	{
		name:   "Test gauge post 1",
		metric: "/update/gauge/testSetGet134/96969.519",
		method: http.MethodPost,
		want: want{
			code: http.StatusOK,
		},
	},
	{
		name:   "Test gauge post 2",
		metric: "/update/gauge/testSetGet135/156519.255",
		method: http.MethodPost,
		want: want{
			code: http.StatusOK,
		},
	},
	{
		name:   "Test gauge get 1",
		metric: "/value/gauge/testSetGet134",
		method: http.MethodGet,
		want: want{
			code: http.StatusOK,
			data: "96969.519",
		},
	},
	{
		name:   "Test gauge get 2",
		metric: "/value/gauge/testSetGet135",
		method: http.MethodGet,
		want: want{
			code: http.StatusOK,
			data: "156519.255",
		},
	},
	{
		name:   "BAD gauge update",
		metric: "/update/gauge/test/none",
		method: http.MethodPost,
		want: want{
			code: http.StatusBadRequest,
		},
	},
	{
		name:   "BAD counter update",
		metric: "/update/counter/test/none",
		method: http.MethodPost,
		want: want{
			code: http.StatusBadRequest,
		},
	},
	{
		name:   "NotFound gauge update",
		metric: "/update/gauge/",
		method: http.MethodPost,
		want: want{
			code: http.StatusNotFound,
		},
	},
	{
		name:   "NotFound counter update",
		metric: "/update/counter/",
		method: http.MethodPost,
		want: want{
			code: http.StatusNotFound,
		},
	},
	{
		name:   "NotImplemented update",
		metric: "/update/unknown/test/1001",
		method: http.MethodPost,
		want: want{
			code: http.StatusNotImplemented,
		},
	},
	{
		name:   "Get gauge metric",
		metric: "/value/gauge/test",
		method: http.MethodGet,
		want: want{
			code: http.StatusOK,
			data: "100",
		},
	},
	{
		name:   "Get counter metric",
		metric: "/value/counter/test",
		method: http.MethodGet,
		want: want{
			code: http.StatusOK,
			data: "100",
		},
	},
}

func TestRouter(t *testing.T) {
	mux := chi.NewRouter()
	RegisterHandlers(mux, storage.NewMetrics())
	ts := httptest.NewServer(mux)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRequest(t, ts, tt)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, testData test) {
	req, err := http.NewRequest(testData.method, ts.URL+testData.metric, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, testData.want.code, resp.StatusCode)
	require.NoError(t, err)

	err = resp.Body.Close()
	require.NoError(t, err)
}
