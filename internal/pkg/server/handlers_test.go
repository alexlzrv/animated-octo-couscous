package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type want struct {
	code int
	data string
}

type test struct {
	name   string
	method string
	metric string
	want   want
}

type testJSON struct {
	name   string
	method string
	url    string
	metric *metrics.Metrics
	want   want
}

var gaugeValue metrics.Gauge = 96969.519

var testsJSON = []testJSON{
	{
		name:   "Post JSON metric",
		method: http.MethodPost,
		url:    "/update/",
		metric: &metrics.Metrics{
			ID:    "Alloc",
			MType: metrics.GaugeMetricName,
			Value: &gaugeValue,
		},
		want: want{
			code: http.StatusOK,
			data: "",
		},
	},
	{
		name:   "Get JSON metric",
		method: http.MethodPost,
		url:    "/value/",
		metric: &metrics.Metrics{
			ID:    "Alloc",
			MType: metrics.GaugeMetricName,
		},
		want: want{
			code: http.StatusOK,
			data: "{\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":96969.519}\n",
		},
	},
}

var tests = []test{
	{
		name:   "OK gauge update",
		metric: "/update/gauge/test1/100.000000",
		method: http.MethodPost,
		want: want{
			code: http.StatusOK,
		},
	},
	{
		name:   "OK counter update",
		metric: "/update/counter/test2/100",
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
		metric: "/value/gauge/test1",
		method: http.MethodGet,
		want: want{
			code: http.StatusOK,
			data: "100",
		},
	},
	{
		name:   "Get counter metric",
		metric: "/value/counter/test2",
		method: http.MethodGet,
		want: want{
			code: http.StatusOK,
			data: "100",
		},
	},
}

func TestRouter(t *testing.T) {
	mux := chi.NewRouter()
	server.RegisterHandlers(mux, storage.NewMetrics())
	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRequest(t, ts, tt)
		})
	}

	for _, tt := range testsJSON {
		t.Run(tt.name, func(t *testing.T) {
			testJSONRequest(t, ts, tt)
		})
	}
}

func BenchmarkRouter(b *testing.B) {
	mux := chi.NewRouter()
	server.RegisterHandlers(mux, storage.NewMetrics())
	ts := httptest.NewServer(mux)
	defer ts.Close()

	benchJSONRequest(b, ts, testsJSON[0])

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchJSONRequest(b, ts, testsJSON[1])
	}
}

func testJSONRequest(t *testing.T, ts *httptest.Server, testData testJSON) {
	body, _ := testData.metric.EncodeMetric()
	req, err := http.NewRequest(testData.method, ts.URL+testData.url, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, testData.want.code, resp.StatusCode)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	if string(respBody) != "" {
		assert.JSONEq(t, testData.want.data, string(respBody))
	} else {
		assert.Equal(t, testData.want.data, string(respBody))
	}

	require.NoError(t, err)

	err = resp.Body.Close()
	if err != nil {
		return
	}
}

func testRequest(t *testing.T, ts *httptest.Server, testData test) {
	req, err := http.NewRequest(testData.method, ts.URL+testData.metric, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, testData.want.code, resp.StatusCode)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	if testData.method == http.MethodGet {
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, testData.want.data, string(respBody))
		require.NoError(t, err)
	}
}

func benchJSONRequest(b *testing.B, ts *httptest.Server, testData testJSON) {
	body, _ := testData.metric.EncodeMetric()
	req, err := http.NewRequest(testData.method, ts.URL+testData.url, body)
	require.NoError(b, err)

	resp, err := http.DefaultClient.Do(req)
	assert.Equal(b, testData.want.code, resp.StatusCode)
	require.NoError(b, err)

	respBody, err := io.ReadAll(resp.Body)
	if string(respBody) != "" {
		assert.JSONEq(b, testData.want.data, string(respBody))
	} else {
		assert.Equal(b, testData.want.data, string(respBody))
	}

	require.NoError(b, err)

	err = resp.Body.Close()
	if err != nil {
		return
	}
}
