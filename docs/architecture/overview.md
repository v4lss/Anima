# Architecture Overview

Animas follows **Clean Architecture** in the Go backend.

```
Domain      - pure business logic, no external dependencies
Application - use-cases, orchestrates domain + ports
Infrastructure - adapters: MongoDB, Redis, HTTP
Workers     - background goroutines: scheduler, http_checker, tcp_checker
```

Dependencies always point inward: Infrastructure, Application, Domain.

## Components

### Backend (Go)

**Domain Layer** (`internal/domain/`)
- `user/` - User entity, password hashing service
- `monitor/` - Monitor and Check entities, validation service
- `alert/` - Alert and AlertConfig entities

**Application Layer** (`internal/application/`)
- `auth/` - Register, Login use-cases
- `monitor/` - Create, Delete, List monitors use-cases

**Infrastructure Layer** (`internal/infrastructure/`)
- `mongodb/` - Repository implementations (User, Monitor, Check, Alert)
- `redis/` - Queue and Cache implementations
- `http/` - Chi router, handlers, middleware (JWT, RateLimit)

**Workers** (`internal/workers/`)
- `scheduler.go` - Enqueues monitors for checking based on interval
- `http_checker.go` - Performs HTTP/HTTPS probes, sends alerts on status change
- `tcp_checker.go` - Performs TCP dial checks, sends alerts on status change

**Packages** (`pkg/`)
- `jwt/` - JWT signing and verification
- `logger/` - Structured logging
- `response/` - JSON response helpers
- `notifier/` - Discord webhook notifications

### Frontend (React + TypeScript)

**Pages** (`web/src/pages/`)
- `Login.tsx` - User authentication
- `Register.tsx` - User registration
- `Dashboard.tsx` - Monitor list and creation
- `MonitorDetail.tsx` - Monitor stats, history, alerts

**Services** (`web/src/services/`)
- `api.ts` - API wrapper with proxy to :8080

**Hooks** (`web/src/hooks/`)
- `useAuth.ts` - Authentication state management
- `useAlerts.ts` - Alert configuration management

## Data Flow

1. **Monitor Creation**: Frontend, API, MongoDB, Scheduler enqueues jobs
2. **Checking**: Scheduler, Redis Queue, HTTP/TCP Checker Workers, MongoDB
3. **Alerts**: Workers detect status change, Discord webhook, Alert record
4. **Monitoring**: Frontend polls API for monitor status and history

## Infrastructure

- **MongoDB**: Stores users, monitors, checks, alerts, alert configs
- **Redis**: Job queue for workers, rate limiting counters
- **HTTP**: Chi router with JWT auth and rate limiting middleware
