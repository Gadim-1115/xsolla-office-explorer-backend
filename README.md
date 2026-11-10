# Xsolla School — Interactive Office Explorer (Backend)

Gamified onboarding backend API for Xsolla School (Baku Crescent Business Center, 31st Floor). Built with Go and PostgreSQL following Clean Architecture and SOLID principles.

## Prerequisites

- **Go:** 1.24 or newer
- **PostgreSQL:** 15+ (Local installation or Docker Desktop)

## Quick Start for Teammates

### Option A: Run Database with Docker (Automatic Migrations)

From the project root directory:

```powershell
docker compose up -d
```

This starts PostgreSQL on `localhost:5432` (`user: postgres`, `password: 12345`, `db: xsolla_explorer`) and automatically executes `001_init_schema.up.sql` and `002_seed_crescent_office.up.sql`.

### Option B: Run with Local PostgreSQL

1. Create a database named `xsolla_explorer` in your local PostgreSQL instance.
2. Execute `backend/migrations/001_init_schema.up.sql` followed by `backend/migrations/002_seed_crescent_office.up.sql`.
3. If your local `postgres` user password differs from `12345`, set `DATABASE_URL` before running the server:

```powershell
$env:DATABASE_URL="postgres://postgres:YOUR_PASSWORD@localhost:5432/xsolla_explorer?sslmode=disable"
```

## Running the API Server

```powershell
cd backend
go mod download
go run ./cmd/api/main.go
```

The server listens on `http://localhost:8080`.

## Running Unit Tests

```powershell
cd backend
go test -v ./...
```

## Environment Variables

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `HTTP_PORT` | `8080` | Port for the HTTP API server |
| `DATABASE_URL` | `postgres://postgres:12345@localhost:5432/xsolla_explorer?sslmode=disable` | PostgreSQL connection string |
| `ALLOWED_ORIGIN` | `http://localhost:5173` | CORS allowed origin for the React frontend |
| `LMS_JWT_SECRET` | *(empty)* | The LMS backend's `JWT_ACCESS_SECRET`. Set it to let students use their LMS sign-in (see below). Empty: everyone is a guest |

## Who is the student? (guest or LMS account)

Every request is tied to one student's progress:

1. **LMS account** — `Authorization: Bearer <LMS access token>`. The token is checked with `LMS_JWT_SECRET` (signature and expiry). Progress is stored under a fixed UUID made from the LMS user id (a version-5 UUID), so it follows the student to any browser and needs no database change. A bad or expired token gets **401**: refresh it with the LMS (`POST /auth/refresh`) or sign in again.
2. **Guest** — `X-User-ID: <uuid>`, a random ID the frontend keeps per browser (as before).
3. Neither header — the shared default guest.

Without `LMS_JWT_SECRET` the `Authorization` header is ignored and everyone is a guest, so nothing changes for setups that do not use the LMS.

## REST API Endpoints

All endpoints return a standardized JSON envelope:

```json
{
  "success": true,
  "data": {},
  "error": null
}
```

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/explorer/init` | Returns all 20 office zones, map hotspots, and student progress (masks locked zone media) |
| `GET` | `/api/v1/explorer/zones/{zoneId}` | Returns unmasked details and hotspots for an unlocked zone |
| `GET` | `/api/v1/explorer/quests/{zoneId}` | Returns Fisher-Yates shuffled Go code blocks without solution keys |
| `POST` | `/api/v1/explorer/quests/{zoneId}/verify` | Validates `{ "submitted_order": ["b1", "b2", ...] }` and unlocks the zone |
| `POST` | `/api/v1/explorer/easter-eggs/{hotspotId}/claim` | Claims a hidden Xsolla logo and returns updated completion metrics |

## Clean Architecture Directory Structure

```text
backend/
├── cmd/api/main.go
├── internal/
│   ├── config/
│   ├── domain/
│   │   ├── entity/
│   │   ├── errors.go
│   │   └── repository.go
│   ├── usecase/
│   ├── adapter/
│   │   ├── handler/
│   │   │   ├── dto/
│   │   │   ├── middleware/
│   │   │   ├── explorer_handler.go
│   │   │   ├── quest_handler.go
│   │   │   └── router.go
│   │   └── storage/postgres/
│   └── pkg/shuffler/
├── migrations/
│   ├── 001_init_schema.up.sql
│   ├── 001_init_schema.down.sql
│   └── 002_seed_crescent_office.up.sql
├── Dockerfile
└── go.mod
```
