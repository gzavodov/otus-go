package monitoring

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func NewMiddleware(serviceName string, logger *zap.Logger) *Middleware {

	m := &Middleware{ServiceName: serviceName, logger: logger}

	m.registry = prometheus.NewRegistry()
	m.requestDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		//Namespace: serviceName,
		//Subsystem: "http",
		Name:    "zg_request_duration_seconds",
		Help:    "The latency of the HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"service", "url", "method", "code"})

	m.logger.Info(serviceName, zap.String("Registring metric", "zg_request_duration_seconds"))
	err := m.registry.Register(m.requestDurationHistogram)
	if err != nil {
		panic(err)
	}
	return m
}

type RequestInfo struct {
	ServiceName string
	URL         string
	Method      string
	StatusCode  string
}

type ResponseWriterSink struct {
	http.ResponseWriter
	StatusCode int
}

func (w *ResponseWriterSink) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type Middleware struct {
	ServiceName              string
	registry                 *prometheus.Registry
	requestDurationHistogram *prometheus.HistogramVec
	logger                   *zap.Logger
}

func (m *Middleware) GetMetricHandler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		m.registry, promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}),
	)
}

func (m *Middleware) GetCollectorHandler(hdlr http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info(m.ServiceName, zap.String("Request URL", r.URL.Path))

		writerSink := &ResponseWriterSink{
			StatusCode:     http.StatusOK,
			ResponseWriter: w,
		}

		start := time.Now()
		defer func() {
			m.ObserveRequestDuration(
				r.Context(),
				&RequestInfo{ServiceName: m.ServiceName, URL: r.URL.Path, Method: r.Method, StatusCode: strconv.Itoa(writerSink.StatusCode)},
				time.Since(start),
			)
		}()

		hdlr.ServeHTTP(writerSink, r)
	})
}

func (m *Middleware) ObserveRequestDuration(ctx context.Context, info *RequestInfo, duration time.Duration) {
	m.requestDurationHistogram.WithLabelValues(info.ServiceName, info.URL, info.Method, info.StatusCode).Observe(duration.Seconds())
}
