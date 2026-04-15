# SNAP Bank API

Production-ready Banking API built with Go, fully compliant with the **SNAP (Standard National Open API Payment)** protocol.

## Tech Stack

| Category | Technology |
|---|---|
| **Language** | Go 1.24 |
| **Database** | PostgreSQL 15 |
| **Cache / Rate Limit** | Redis 7 |
| **Message Broker** | Apache Kafka (KRaft mode) |
| **Tracing** | OpenTelemetry → Tempo |
| **Metrics** | Prometheus |
| **Logging** | Zap → Loki (via Promtail) |
| **Telemetry Pipeline** | OpenTelemetry Collector |
| **Dashboards** | Grafana |
| **Containerization** | Docker Compose |

## Architecture

```
project_bank_api/
├── main.go                         # Application entrypoint
├── config/                         # Configuration & infrastructure setup
│   ├── config.go                   #   Database & Redis connections
│   ├── kafka.go                    #   Kafka writer
│   ├── logger.go                   #   Zap logger (dev/prod)
│   ├── migration.go                #   Auto database migrations
│   └── topic.go                    #   Kafka topic creation
├── internal/                       # Private application code
│   ├── adapter/                    #   External service adapters
│   │   └── event_publisher.go      #     Kafka event publisher
│   ├── constant/                   #   SNAP codes & Kafka topics
│   │   ├── service_code.go
│   │   └── kafka_topics.go
│   ├── domain/                     #   Domain models
│   │   ├── account.go
│   │   └── transaction.go
│   ├── dto/                        #   Request/Response DTOs
│   │   ├── account_dto.go
│   │   ├── event_dto.go
│   │   ├── snap_response.go
│   │   └── transaction_dto.go
│   ├── handler/                    #   HTTP handlers
│   │   ├── account_handler.go
│   │   └── snap_middleware.go      #     SNAP header & signature validation
│   ├── middleware/                  #   Global middleware chain
│   │   ├── core.go                 #     Middleware orchestration
│   │   ├── cache.go                #     Redis response caching
│   │   ├── idempotency.go          #     Redis idempotency (SetNX)
│   │   ├── metrics.go              #     Prometheus HTTP metrics
│   │   ├── rate_limit.go           #     Redis rate limiter
│   │   ├── request_logger.go       #     Structured request logging
│   │   └── tracing.go              #     OpenTelemetry span creation
│   ├── repository/                 #   Database access layer
│   │   ├── account_repository.go
│   │   └── transaction_repository.go
│   └── service/                    #   Business logic layer
│       └── account_service.go
├── pkg/                            # Shared utilities
│   ├── telemetry/                  #   Tracer & Prometheus metrics
│   │   ├── tracer.go
│   │   └── metrics.go
│   ├── response.go                 #   Response helpers
│   └── snap_util.go                #   HMAC signature & ID generators
├── migrations/                     #   SQL migration files
│   └── 001_init.sql
├── docker-compose.yml              # Full infrastructure stack
├── Dockerfile                      # Multi-stage build
├── otel-collector.yaml             # OpenTelemetry Collector config
├── prometheus.yaml                 # Prometheus scrape config
├── grafana-datasources.yaml        # Grafana datasource provisioning
├── tempo.yaml                      # Tempo config
└── promtail.yaml                   # Promtail log shipping config
```

### Request Flow

```
Client Request
    │
    ▼
┌─────────────────────────────────────────────────┐
│  Middleware Chain                                │
│  OTel → Prometheus → Logger → CORS → JSON       │
│  → RateLimit → Idempotency → Timeout            │
│  → SNAPMiddleware (header + signature validation)│
└────────────────────┬────────────────────────────┘
                     ▼
              ┌─────────────┐
              │   Handler   │  ← Parse request, map SNAP response codes
              └──────┬──────┘
                     ▼
              ┌─────────────┐
              │   Service   │  ← Business logic, validation, Kafka events
              └──────┬──────┘
                     ▼
              ┌─────────────┐
              │ Repository  │  ← Database queries (PostgreSQL)
              └─────────────┘
```

### Observability Pipeline

```
                 ┌──────────────┐
                 │  Go App      │
                 │  (snap-bank) │
                 └──────┬───────┘
                        │ OTLP
                        ▼
              ┌──────────────────┐
              │  OpenTelemetry   │
              │  Collector       │
              └──┬───────┬───┬──┘
                 │       │   │
        traces   │  metrics  │  logs
                 ▼       ▼   ▼
              ┌─────┐ ┌────┐ ┌─────┐
              │Tempo│ │Prom│ │Loki │
              └──┬──┘ └─┬──┘ └──┬──┘
                 │      │      │
              ┌──▼──────▼──────▼──┐
              │     Grafana       │
              │ trace↔log↔metric  │
              └───────────────────┘
```

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.24+ (for local development)

### Run

```bash
# Start the full stack
docker compose up --build -d

# Check logs
docker compose logs -f app
```

### Services

| Service | URL | Credentials |
|---|---|---|
| **API** | http://localhost:8081 | — |
| **Grafana** | http://localhost:3001 | admin / admin |
| **Kafka UI** | http://localhost:8090 | — |
| **Prometheus** | http://localhost:9091 | — |

### Health Check

```bash
curl http://localhost:8081/health
# {"status":"ok"}
```

## API Endpoints

All endpoints under `/snap/v1/` require SNAP headers.

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/snap/v1/account-creation` | Create a new account |
| `POST` | `/snap/v1/transfer-intrabank` | Transfer between accounts |
| `GET` | `/snap/v1/accounts` | List all accounts |
| `GET` | `/snap/v1/accounts/{accountNo}` | Balance inquiry |
| `POST` | `/snap/v1/transaction-history` | Transaction history |

### Required SNAP Headers

| Header | Description |
|---|---|
| `Authorization` | Bearer token |
| `X-TIMESTAMP` | ISO 8601 timestamp |
| `X-SIGNATURE` | HMAC-SHA256(`timestamp\|partnerID`, secret_key) |
| `X-PARTNER-ID` | Partner identifier |
| `X-EXTERNAL-ID` | Unique request ID |
| `CHANNEL-ID` | Channel identifier |

### Postman Collection

Import `SNAP_Bank_API.postman_collection.json` — all headers and signature generation are pre-configured.

## Kafka Events

Events are published after successful operations:

| Topic | Trigger | Key |
|---|---|---|
| `snap.account.created` | Account creation | accountNo |
| `snap.transfer.completed` | Transfer success | referenceNo |

Browse messages at **Kafka UI**: http://localhost:8090

## Deep Tracing

Every request produces a full trace visible in **Grafana Tempo**:

```
HTTP POST /snap/v1/transfer-intrabank
  └── middleware.SNAPHeaderValidation
        └── handler.Transfer
              └── service.Account.Transfer
                    ├── repository.Transaction.GetByPartnerReferenceNo
                    ├── repository.Account.GetByAccountNo (source)
                    ├── repository.Account.GetByAccountNo (beneficiary)
                    ├── repository.Transaction.Create
                    ├── repository.Account.UpdateBalance (debit)
                    ├── repository.Account.UpdateBalance (credit)
                    └── kafka.Publish (snap.transfer.completed)
```

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `SNAP_DB_HOST` | localhost | PostgreSQL host |
| `SNAP_DB_PORT` | 5433 | PostgreSQL port |
| `SNAP_DB_USER` | postgres | Database user |
| `SNAP_DB_PASSWORD` | postgres | Database password |
| `SNAP_DB_NAME` | snap_bank | Database name |
| `SNAP_SERVER_PORT` | 8081 | API server port |
| `SNAP_SECRET_KEY` | my-snap-secret-key-for-signature | HMAC signing key |
| `REDIS_HOST` | localhost | Redis host |
| `REDIS_PORT` | 6381 | Redis port |
| `KAFKA_BROKER` | localhost:9092 | Kafka broker address |
| `OTEL_COLLECTOR_HOST` | localhost | OpenTelemetry Collector host |
