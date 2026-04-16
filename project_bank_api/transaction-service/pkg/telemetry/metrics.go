package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "txn_svc_http_requests_total",
			Help: "Total HTTP requests by method, path, and status",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "txn_svc_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	BankTransactionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "txn_svc_business_transactions_total",
			Help: "Total business transactions by type and status",
		},
		[]string{"type", "status"},
	)

	BankTransactionAmountTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "txn_svc_transfer_amount_total_idr",
			Help: "Total cumulative transfer amount (IDR)",
		},
	)

	BankTransactionAmountDistribution = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "txn_svc_transfer_amount_distribution",
			Help:    "Distribution of transfer amounts",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 1000000, 10000000},
		},
	)
)
