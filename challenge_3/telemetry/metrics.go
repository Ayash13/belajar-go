package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HttpRequestsTotal mencatat total jumlah request yang masuk ke API
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bank_api_http_requests_total",
			Help: "Total jumlah HTTP request yang masuk berdasarkan method, path, dan status",
		},
		[]string{"method", "path", "status"},
	)

	// HttpRequestDuration merekam durasi/latensi pemrosesan setiap request HTTP (dalam sekon)
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bank_api_http_request_duration_seconds",
			Help:    "Distribusi durasi pemrosesan HTTP request (dalam sekon)",
			Buckets: prometheus.DefBuckets, // Menggunakan batasan bucket standard milik Prometheus (contoh: 0.005s, 0.01s, 0.1s dst)
		},
		[]string{"method", "path", "status"},
	)
)
