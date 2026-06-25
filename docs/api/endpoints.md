# API Endpoints

Base URL: `http://localhost:8080`

## Authentication

### POST /api/auth/register
Register a new user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (201):**
```json
{
  "data": {
    "id": "user123",
    "email": "user@example.com"
  }
}
```

### POST /api/auth/login
Authenticate and receive JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

## Monitors

All monitor endpoints require JWT authentication via `Authorization: Bearer <token>` header.

### POST /api/monitors
Create a new monitor.

**Request Body:**
```json
{
  "name": "My API",
  "target": "https://api.example.com",
  "type": "HTTPS",
  "interval": 60
}
```

**Response (201):**
```json
{
  "data": {
    "id": "monitor123",
    "name": "My API",
    "target": "https://api.example.com",
    "type": "HTTPS",
    "interval": 60,
    "enabled": true,
    "lastStatus": "UP"
  }
}
```

### GET /api/monitors
List all monitors for the authenticated user.

**Response (200):**
```json
{
  "data": [
    {
      "id": "monitor123",
      "name": "My API",
      "target": "https://api.example.com",
      "type": "HTTPS",
      "interval": 60,
      "enabled": true,
      "lastStatus": "UP"
    }
  ]
}
```

### GET /api/monitors/:id
Get a specific monitor by ID.

**Response (200):**
```json
{
  "data": {
    "id": "monitor123",
    "name": "My API",
    "target": "https://api.example.com",
    "type": "HTTPS",
    "interval": 60,
    "enabled": true,
    "lastStatus": "UP"
  }
}
```

### DELETE /api/monitors/:id
Delete a monitor.

**Response (204):** No content

### GET /api/monitors/:id/history
Get check history for a monitor.

**Response (200):**
```json
{
  "data": [
    {
      "id": "check123",
      "monitorid": "monitor123",
      "status": "UP",
      "responsetime": 125,
      "checkedat": "2026-06-24T20:00:00Z"
    }
  ]
}
```

## Alerts

### POST /api/monitors/:monitorId/alerts
Create an alert configuration for a monitor.

**Request Body:**
```json
{
  "type": "DISCORD",
  "webhook": "https://discord.com/api/webhooks/..."
}
```

**Response (201):**
```json
{
  "data": {
    "id": "alert123",
    "monitorid": "monitor123",
    "userid": "user123",
    "type": "DISCORD",
    "webhook": "https://discord.com/api/webhooks/...",
    "enabled": true,
    "createdat": "2026-06-24T20:00:00Z",
    "updatedat": "2026-06-24T20:00:00Z"
  }
}
```

### GET /api/monitors/:monitorId/alerts
List alert configurations for a monitor.

**Response (200):**
```json
{
  "data": [
    {
      "id": "alert123",
      "monitorid": "monitor123",
      "userid": "user123",
      "type": "DISCORD",
      "webhook": "https://discord.com/api/webhooks/...",
      "enabled": true,
      "createdat": "2026-06-24T20:00:00Z",
      "updatedat": "2026-06-24T20:00:00Z"
    }
  ]
}
```

### DELETE /api/monitors/:monitorId/alerts/:id
Delete an alert configuration.

**Response (204):** No content

## Health

### GET /health
Health check endpoint (no authentication required).

**Response (200):**
```json
{
  "status": "ok"
}
```

## Error Responses

All endpoints return errors in the following format:

```json
{
  "error": "error message"
}
```

Common status codes:
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing or invalid token)
- `404` - Not Found
- `429` - Too Many Requests (rate limit exceeded)
- `500` - Internal Server Error
