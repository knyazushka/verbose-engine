// pkg/metrics/metrics.go
package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	// HTTP метрики
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	HTTPRequestSize     *prometheus.SummaryVec
	HTTPResponseSize    *prometheus.SummaryVec

	// Бизнес метрики
	FileUploadsTotal   prometheus.Counter
	FileDownloadsTotal prometheus.Counter
	FileDeletesTotal   prometheus.Counter
	FileUploadSize     prometheus.Histogram
	FileStorageUsage   prometheus.Gauge
}

func New() *Metrics {
	return &Metrics{
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),

		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),

		FileUploadsTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "file_uploads_total",
				Help: "Total number of file uploads",
			},
		),

		FileDownloadsTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "file_downloads_total",
				Help: "Total number of file downloads",
			},
		),

		FileUploadSize: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "file_upload_size_bytes",
				Help:    "Size of uploaded files in bytes",
				Buckets: prometheus.ExponentialBuckets(1024, 2, 10), // 1KB to 1MB
			},
		),

		FileStorageUsage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "file_storage_usage_bytes",
				Help: "Total storage usage in bytes",
			},
		),
	}
}

func (m *Metrics) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		m.HTTPRequestsTotal.WithLabelValues(
			r.Method,
			r.URL.Path,
			http.StatusText(rw.statusCode),
		).Inc()

		m.HTTPRequestDuration.WithLabelValues(
			r.Method,
			r.URL.Path,
		).Observe(duration)
	})
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.Handler()
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
