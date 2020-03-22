package monitoring

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func NewMiddleware(serviceName string, logger *zap.Logger) *Middleware {
	m := &Middleware{serviceName: serviceName, logger: logger}

	m.registry = prometheus.NewRegistry()

	m.requestDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: serviceName,
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "The latency of the HTTP requests.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"service", "url", "method", "code", "code_group"},
	)

	if err := m.registry.Register(m.requestDurationHistogram); err != nil {
		log.Fatal(err)
	}

	m.requestTotalCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: serviceName,
			Subsystem: "http",
			Name:      "request_total",
			Help:      "The total quantity of HTTP requests.",
		},
		[]string{"service", "url", "method", "code", "code_group"},
	)

	if err := m.registry.Register(m.requestTotalCounter); err != nil {
		log.Fatal(err)
	}

	return m
}

func WrapHandler(hdlr http.Handler, mdlw *Middleware) http.Handler {
	if mdlw == nil {
		return hdlr
	}

	return mdlw.PrepareHandlerWrapper(hdlr)
}

type Middleware struct {
	serviceName              string
	registry                 *prometheus.Registry
	requestDurationHistogram *prometheus.HistogramVec
	requestTotalCounter      *prometheus.CounterVec
	logger                   *zap.Logger
}

func (m *Middleware) PrepareMetricExportHandler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		m.registry, promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}),
	)
}

func (m *Middleware) PrepareHandlerWrapper(hdlr http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info(m.serviceName, zap.String("URL Path", r.URL.Path))

		writerSink := &ResponseWriterSink{
			StatusCode:     http.StatusOK,
			ResponseWriter: w,
		}

		start := time.Now()
		requestInfo := &RequestInfo{ServiceName: m.serviceName, URL: r.URL.Path, Method: r.Method, StatusCode: writerSink.StatusCode}
		defer func() {
			m.ObserveRequest(r.Context(), requestInfo, time.Since(start))
		}()

		hdlr.ServeHTTP(writerSink, r)
	})
}

func (m *Middleware) ObserveRequest(ctx context.Context, info *RequestInfo, duration time.Duration) {
	latency := duration.Seconds()

	m.requestDurationHistogram.WithLabelValues(
		info.ServiceName,
		info.URL,
		info.Method,
		info.GetStatusCodeLabel(),
		info.GetStatusCodeGroupLabel(),
	).Observe(latency)

	m.requestTotalCounter.WithLabelValues(
		info.ServiceName,
		info.URL,
		info.Method,
		info.GetStatusCodeLabel(),
		info.GetStatusCodeGroupLabel(),
	).Inc()
}
