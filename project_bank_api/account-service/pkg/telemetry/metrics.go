package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "account_svc_http_requests_total",
			Help: "Total HTTP requests by method, path, and status",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "account_svc_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	BankTransactionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "account_svc_business_transactions_total",
			Help: "Total business transactions by type and status",
		},
		[]string{"type", "status"},
	)

	BankAccountsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "account_svc_accounts_created_total",
			Help: "Total accounts created",
		},
	)
)
