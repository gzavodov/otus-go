package queuemonitoring

import (
	"log"
	"net/http"

	"github.com/gzavodov/otus-go/calendar/pkg/queue"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func NewMiddleware(serviceName string, logger *zap.Logger) *Middleware {
	m := &Middleware{serviceName: serviceName, logger: logger}

	m.registry = prometheus.NewRegistry()

	m.sendingCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: serviceName,
			Subsystem: "queue",
			Name:      "total_sent",
			Help:      "The total quantity of sent messages.",
		},
		[]string{},
	)

	if err := m.registry.Register(m.sendingCounter); err != nil {
		log.Fatal(err)
	}

	return m
}

func WrapChannel(c queue.NotificationChannel, mdlw *Middleware) queue.NotificationChannel {
	if mdlw == nil {
		return c
	}

	return mdlw.PrepareChannelWrapper(c)
}

type Middleware struct {
	serviceName    string
	registry       *prometheus.Registry
	sendingCounter *prometheus.CounterVec
	logger         *zap.Logger
}

func (m *Middleware) PrepareMetricExportHandler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		m.registry, promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}),
	)
}

func (m *Middleware) PrepareChannelWrapper(c queue.NotificationChannel) queue.NotificationChannel {
	return NewNotificationChannelSink(c, m.ObserveSending)
}

func (m *Middleware) ObserveSending(info *SendingInfo) {
	m.sendingCounter.WithLabelValues().Add(float64(info.Count))
}
