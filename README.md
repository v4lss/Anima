<p align="center">
  <h1 align="center">Animas</h1>
  <p align="center">Open-source uptime monitoring for websites, APIs and TCP ports.</p>
</p>

---

## Features

- HTTP / HTTPS endpoint monitoring
- TCP port monitoring
- Real-time response time metrics
- Uptime history & check logs
- JWT-based authentication
- Discord & Email alerts *(v2)*
- Public status pages *(v2)*
- Multi-region checks *(v3)*

## Stack

| Layer     | Technology      |
|-----------|-----------------|
| Backend   | Go              |
| Frontend  | TypeScript      |
| Database  | MongoDB         |
| Cache/Queue | Redis         |

## Getting Started

```bash
# Clone
git clone https://github.com/v4lss/Anima.git
cd Anima

# Copy env
cp .env.example .env

# Run API (requires Go 1.22+)
cd api && go run ./cmd/animas

# Run Web
cd web && npm install && npm run dev
```

## Project Structure

```
animas/
├── api/        # Go backend (Clean Architecture)
├── web/        # TypeScript frontend
├── docs/       # Documentation
└── scripts/    # Dev & deployment scripts
```

## Roadmap

See [ROADMAP.md](./ROADMAP.md)

## License

[MIT](./LICENSE)
