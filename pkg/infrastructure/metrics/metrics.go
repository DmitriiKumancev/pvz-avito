package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type HTTPMetrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}

func NewHTTPMetrics() *HTTPMetrics {
	requestsTotal := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Общее количество HTTP запросов",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Длительность HTTP запроса в секундах",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	return &HTTPMetrics{
		RequestsTotal:   requestsTotal,
		RequestDuration: requestDuration,
	}
}

type Registry struct {
	HTTP *HTTPMetrics
}

func NewRegistry() *Registry {
	return &Registry{
		HTTP: NewHTTPMetrics(),
	}
}
