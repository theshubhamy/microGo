# Quick Commerce Backend

## Project Architecture (Go + gRPC + GraphQL)

This document outlines a production-grade microservices architecture for a **Quick Commerce Platform** like Blinkit, built using **Go (Golang)** with **gRPC for internal communication** and **GraphQL as the API gateway**. It includes:

- Microservices breakdown
- Required tools/libraries per service
- Folder structure with package descriptions
- Scalability & DevOps practices

---

## 📂 Microservices & Their Responsibilities

### 1. `account`

**Responsibilities:**

- User signup/login (OTP/email)
- JWT-based authentication
- Wallet and address management

**Tools & Packages:**

- `bcrypt` for password hashing
- `jwt-go` for JWT tokens
- `gRPC` for internal calls
- `PostgreSQL` via `pgx`
- `Redis` for session cache

---

### 2. `catalog`

**Responsibilities:**

- Products, categories, and inventory management
- Price and stock updates

**Tools & Packages:**

- `gorm` or `pgx` for DB layer
- `gRPC`
- `Kafka` for stock sync with inventory service
- `ElasticSearch` (optional) for search indexing

---

### 3. `order`

**Responsibilities:**

- Cart and order placement
- Dynamic coupon and pricing engine
- Order state management

**Tools & Packages:**

- `gRPC`
- `Redis` (for cart state)
- `PostgreSQL`
- `Kafka` for event dispatch (order_created, order_paid)

---

### 4. `payment`

**Responsibilities:**

- Integrate Razorpay/Stripe
- Payment initiation, webhook validation
- Transaction logging

**Tools & Packages:**

- `Go HTTP client`
- `Kafka` (for payment success/failure events)
- `PostgreSQL`
- `gRPC`

---

### 5. `notification`

**Responsibilities:**

- Handle email, SMS, and push notifications
- Send OTPs, order updates, promotions

**Tools & Packages:**

- `Kafka` (consume events)
- `SendGrid`, `Twilio`, or Firebase Cloud Messaging
- `cron` or scheduler for retries

---

### 6. `search`

**Responsibilities:**

- Full-text search for products, categories
- Filtering and sorting support

**Tools & Packages:**

- `ElasticSearch` or `Meilisearch`
- `gRPC`
- `Redis` for caching

---

## 🗒️ Folder Structure with Tools and Utilities

```

.
├── graphql/ # GraphQL Gateway
│ ├── main.go
│ ├── schema.graphql # GraphQL Schema
│ ├── gqlgen.yml
│ └── resolvers/ # Query/Mutation Resolvers
│
├── services/ # Microservices
│ ├── account/
│ ├── catalog/
│ ├── order/
│ ├── payment/
│ ├── notification/
│ └── search/
│
├── proto/ # Protobuf files
├── kafka/ # Kafka topic setup & producers/consumers
├── monitoring/ # Prometheus + Grafana setup
├── docker-compose.yml
├── Makefile
├── README.md

```

## Service Architecture

```
main.go
  └── Loads config, connects DB, starts gRPC server

server.go (transport layer)
  └── Implements gRPC interface
      └── Calls service.go (business logic)
              └── Calls repository.go (DB logic)

client.go
  └── Allows other services to call this via gRPC

pb/
  └── Auto-generated from account.proto

Dockerfiles
  └── Build containers for app and DB

db.sql
  └── Bootstraps schema (users, sessions, etc.)

```

### `graphql/`

**Uses:**

- `gqlgen` for GraphQL server
- Converts schema to Go types and resolvers
- Acts as public API Gateway

### `services/*`

Each service is a Go module with its own:

- `main.go`: Entry point
- `server.go`: gRPC server
- `handlers/`, `models/`, `clients/`
- Dockerfile for containerization

### `proto/`

**Uses:**

- Shared `.proto` files
- Compiled with `protoc-gen-go` and `protoc-gen-go-grpc`

### `docker-compose.yml`

**Uses:**

- Service orchestration
- Runs DBs (Postgres, Redis), Kafka, Prometheus

### `Makefile`

**Uses:**

- Build, lint, run proto commands
- Common local dev tasks

---

## ⚙️ DevOps + Production Setup

### 📬 Kafka Message Queue for Asynchronous Workflows

- **Kafka** for async events

- Topics:
  - `order.created`
  - `order.paid`
  - `inventory.updated`
  - `wallet.topup`
  - `notification.send`
- Consumers in `order`, `catalog`, `payment`, and a new `notification` service
- Ensures reliable background processing, decouples services

### 3. 📊 Monitoring with Prometheus + Grafana

- Prometheus scrapes metrics from services (`/metrics` endpoint using `promhttp`)
- Grafana dashboards:
  - Request latency & throughput per service
  - DB query performance
  - Kafka lag metrics

### Observability

- **Prometheus** + **Grafana** for metrics
- **Sentry** for error tracking
- **Jaeger** for distributed tracing

### Security

- HTTPS with TLS (via NGINX ingress)
- JWT for auth
- API rate limiting via Envoy or Istio

### Scalability

- **Kubernetes** deployment with:

  - **HPA (Horizontal Pod Autoscaler)**
  - Auto-scaling based on CPU/RAM
  - Liveness/readiness probes

### Alerts

- Grafana alerts for latency, DB errors, traffic spikes
- PagerDuty/Slack integration

---

## 🌟 Optional Features

- **Machine Learning service** for product recommendations
- **Admin portal** (React) to manage products, orders, coupons
- **A/B testing framework** for experiments

---

## 📌 Future Enhancements

- Replace Kafka with NATS JetStream or Redis Streams (lighter options)
- Use gRPC Gateway for REST fallback
- Add CI/CD pipeline using GitHub Actions or GitLab CI
