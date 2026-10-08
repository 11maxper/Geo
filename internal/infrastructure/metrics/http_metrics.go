package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status_code"},
	)

	httpRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)
)

// HTTPMetricsMiddleware создает middleware для сбора метрик HTTP запросов
// Middleware измеряет время выполнения и подсчитывает количество запросов
func HTTPMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()

		method := r.Method
		statusCode := strconv.Itoa(ww.Status())

		endpoint := chi.RouteContext(r.Context()).RoutePattern()
		if endpoint == "" {
			endpoint = "unknown"
		}

		httpRequestDuration.WithLabelValues(method, endpoint, statusCode).Observe(duration)
		httpRequestTotal.WithLabelValues(method, endpoint, statusCode).Inc()
	})
}
