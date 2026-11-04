# Xsolla School — Interactive Office Explorer (Backend)

Gamified onboarding backend service for Xsolla School (Baku Crescent Business Center, 31st Floor). Built with Go and PostgreSQL following Clean Architecture and SOLID principles.

## Tech Stack

- **Language:** Go (`net/http`)
- **Database:** PostgreSQL (`pgx/v5`)
- **Architecture:** Layered Clean Architecture (Domain, UseCase, Interface Adapters, Infrastructure)

## Project Structure

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
├── Dockerfile
└── go.mod
```

## REST API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/explorer/init` | Fetches map zones, hotspots, and user progress (masks locked zone assets) |
| `GET` | `/api/v1/explorer/zones/{zoneId}` | Fetches unmasked details and hotspots for an unlocked zone |
| `GET` | `/api/v1/explorer/quests/{zoneId}` | Returns shuffled quest code blocks without solution keys |
| `POST` | `/api/v1/explorer/quests/{zoneId}/verify` | Validates submitted code block sequence and unlocks the zone |
| `POST` | `/api/v1/explorer/easter-eggs/{hotspotId}/claim` | Claims a hidden Xsolla logo and updates completion metrics |

## Running Locally

```powershell
cd backend
go test ./...
go run ./cmd/api/main.go
```
