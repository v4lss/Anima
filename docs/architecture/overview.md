# Architecture Overview

Animas follows **Clean Architecture** in the Go backend.

```
Domain      - pure business logic, no external dependencies
Application - use-cases, orchestrates domain + ports
Infrastructure - adapters: MongoDB, Redis, HTTP
Workers     - background goroutines: scheduler, http_checker, tcp_checker
```

Dependencies always point inward: Infrastructure → Application → Domain.
