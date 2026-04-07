package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HttpRequestsTotal tracks the total number of HTTP requests
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bank_api_http_requests_total",
			Help: "Total count of HTTP requests by method, path, and status code",
		},
		[]string{"method", "path", "status"},
	)

	// HttpRequestDuration records processing latency for every HTTP request
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bank_api_http_request_duration_seconds",
			Help:    "Latency of HTTP requests in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	// BankTransactionsTotal tracks the number of business transactions
	BankTransactionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bank_api_business_transactions_total",
			Help: "Total count of business transactions by type and status",
		},
		[]string{"type", "status"},
	)

	// BankTransactionAmountTotal tracks the cumulative sum of money transferred
	BankTransactionAmountTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "bank_api_business_transfer_amount_total_usd",
			Help: "Total cumulative balance transferred across all transactions",
		},
	)

	// BankAccountsTotal tracks the total number of bank accounts created
	BankAccountsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "bank_api_business_accounts_created_total",
			Help: "Total count of business bank accounts created",
		},
	)

	// BankTransactionAmountDistribution records the distribution of money transferred per transaction
	BankTransactionAmountDistribution = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "bank_api_business_transfer_amount_distribution",
			Help:    "Distribution of individual transfer amounts",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000},
		},
	)
)
