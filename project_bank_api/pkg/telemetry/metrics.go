package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "snap_api_http_requests_total",
			Help: "Total count of HTTP requests by method, path, and status code",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "snap_api_http_request_duration_seconds",
			Help:    "Latency of HTTP requests in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	BankTransactionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "snap_api_business_transactions_total",
			Help: "Total count of SNAP business transactions by type and status",
		},
		[]string{"type", "status"},
	)

	BankTransactionAmountTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "snap_api_business_transfer_amount_total_idr",
			Help: "Total cumulative balance transferred across all SNAP transactions (IDR)",
		},
	)

	BankAccountsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "snap_api_business_accounts_created_total",
			Help: "Total count of SNAP accounts created",
		},
	)

	BankTransactionAmountDistribution = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "snap_api_business_transfer_amount_distribution",
			Help:    "Distribution of individual SNAP transfer amounts",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 1000000, 10000000},
		},
	)
)
