# SNAP Bank API (Microservices Architecture)

Production-ready Banking API built with Go, fully compliant with the **SNAP (Standard National Open API Payment)** protocol. The application has been refactored from a monolith into a robust event-driven microservices architecture.

## Tech Stack

| Category | Technology |
|---|---|
| **Language** | Go 1.24 (Alpine) |
| **Databases** | PostgreSQL 15 (3 separate instances) |
| **Cache / Rate Limit** | Redis 7 |
| **Message Broker** | Apache Kafka (KRaft mode) |
| **Communication** | REST (HTTP), gRPC, Async (Kafka Event-Driven) |
| **Tracing** | OpenTelemetry → Tempo |
| **Metrics** | Prometheus |
| **Logging** | Zap → Loki (via Promtail) |
| **Telemetry Pipeline** | OpenTelemetry Collector |
| **Dashboards** | Grafana |
| **Containerization** | Docker Compose |

## Microservices Architecture

The system is broken down into three main services, fulfilling specific business domains, each with its own independent PostgreSQL database to maintain complete data isolation:

1. **Account Service (`8081`)**: Manages account creation and balance inquiries. Exposes a internal gRPC Server (`50051`) for synchronous balance validations.
2. **Transaction Service (`8082`)**: Handles intra-bank transfers and transaction history. Communicates via REST for clients and synchronously requests balance updates from the Account Service via gRPC. 
3. **Notification Service (`8083`)**: Consumes Kafka events seamlessly and sends fully-automated asynchronous HTTP callbacks to partner webhooks.

### Project Structure
```text
project_bank_api/
├── account-service/                # Account management & gRPC Server   
├── transaction-service/            # Transfer management & gRPC Client
├── notification-service/           # Webhook Async Callback dispatcher
├── pb/                             # Shared Protobufs generated code
├── proto/                          # .proto definition files
├── docker-compose.yml              # Full infrastructure stack
├── otel-collector.yaml             # OpenTelemetry Collector config
├── prometheus.yaml                 # Prometheus scrape config
├── grafana-datasources.yaml        # Grafana datasource provisioning
├── tempo.yaml                      # Tempo config
└── promtail.yaml                   # Promtail log shipping config
```

### Request Flow (Transfer Example)

```text
Client Request
    │
    ▼
┌───────────────────────┐
│ Transaction Service   │  ← Validates HTTP Request, checks SNAP headers
└───┬───────────────────┘
    │ (gRPC)
    ▼
┌───────────────────────┐
│ Account Service       │  ← Validates balance constraints, deducts/credits balances 
└───┬───────────────────┘
    │ (Kafka Event: snap.transfer.completed)
    ▼
┌───────────────────────┐
│ Notification Service  │  ← Consumes event
└───┬───────────────────┘
    │ (HTTP POST)
    ▼
┌───────────────────────┐
│ Partner Webhook URL   │  ← Third-party server receives transaction status (e.g Webhook.site)
└───────────────────────┘
```

## Getting Started

### Prerequisites

- Docker & Docker Compose

### Run the Stack

To build the Go binaries and start exactly 14 containers (Services, DBs, Kafka, Redis, Obs-stack):
```bash
# Start the full stack in detached mode
docker-compose up -d --build

# Check the status of all containers
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

### Services & Ports

| Component | URL / Port | Credentials / Note |
|---|---|---|
| **Account Service** | http://localhost:8081 | Internal gRPC on 50051 |
| **Transaction Service**| http://localhost:8082 | — |
| **Notification Service**| http://localhost:8083 | — |
| **Grafana** | http://localhost:3001 | admin / admin |
| **Kafka UI** | http://localhost:8090 | — |
| **Prometheus** | http://localhost:9091 | — |

## API Endpoints

All endpoints heavily enforce **SNAP required headers**.

### Account Service (`8081`)
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/snap/v1/account-creation` | Create a new account |
| `GET` | `/snap/v1/accounts` | List all accounts (Cached via Redis) |
| `GET` | `/snap/v1/accounts/{accountNo}` | Balance inquiry |

### Transaction Service (`8082`)
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/snap/v1/transfer-intrabank` | Transfer between accounts |
| `POST` | `/snap/v1/transaction-history` | Fetch transaction history |

### Required SNAP Headers

For every single request listed above, you must pass:

| Header | Example Value | Description |
|---|---|---|
| `X-TIMESTAMP` | `2024-03-15T10:00:00+07:00` | ISO 8601 timestamp |
| `X-SIGNATURE` | `any-string` | HMAC-SHA256 signature (Mocked for dev) |
| `X-PARTNER-ID` | `partner-xyz` | Partner identifier |
| `X-EXTERNAL-ID` | `req-12345` | Unique request check (used for Idempotency)|
| `Content-Type` | `application/json` | standard json |

## Testing Interactions (Postman & Webhook.site)

Because of the full lifecycle capability, you can test everything synchronously and asynchronously.

1. **Setup Webhook Callback Receiver**:
   - Go to [Webhook.site](https://webhook.site/) and copy "Your unique URL"
   - Modify the `PARTNER_CALLBACK_URL` under `notification-service` inside `docker-compose.yml`.
   - Run `docker-compose up -d notification-service` to apply the changes.

2. **Run the Create Account Request (Postman)**:
   - Hit `POST http://localhost:8081/snap/v1/account-creation` and include the mandatory headers above.
   - You will instantly see an asynchronous payload hit your Webhook.site!

3. **Run the Transfer Request (Postman)**:
   - Hit `POST http://localhost:8082/snap/v1/transfer-intrabank` and provide two valid Account Numbers.
   - The transaction service will interface with Account via gRPC. Once finished, you will immediately see a payload hit Webhook.site indicating `snap.transfer.completed`.

## Kafka Events

Topics are configured to auto-create and persist robustly.

| Topic | Publisher | Consumer | Trigger |
|---|---|---|---|
| `snap.account.created` | Account Service | Notification Service | Account creation |
| `snap.transfer.completed` | Transaction Service | Notification Service | Successful intra-bank transfer |

Browse messages dynamically via **Kafka UI**: http://localhost:8090

## Observability & Distributed Tracing

As the platform is now a distributed microservices network, **OpenTelemetry** propagates context continuously. Every HTTP request sent generates a Distributed Trace linking the `Transaction Service` directly to the `Account Service`'s gRPC invocation and the subsequent `Notification Service` execution!

To view trace graphs and see bottlenecks visually:
1. Open **Grafana**: http://localhost:3001
2. Go to **Explore** -> **Tempo**
3. Browse generated span traces for seamless cross-service insights.
