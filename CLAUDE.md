# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**ParkLah!** is a Singapore parking assistant app. This repo (`parklah-go`) is the Go backend — it runs an ingestion pipeline that aggregates carpark data from three sources and stores it in PostgreSQL. The eventual product includes an Android client (Kotlin/Jetpack Compose) and a REST API layer; the current codebase is Phase 1: data ingestion only.

**Three data sources:**
- **URA** — token-based auth, daily static details + every 3–5 min availability
- **HDB** (data.gov.sg) — no auth, static carpark info + real-time availability
- **LTA DataMall** — API key auth, real-time availability (not yet implemented)

## Commands

```bash
# Run the ingestion pipeline
go run ./cmd/ingestion/main.go

# Start PostgreSQL (port 5433)
docker-compose up -d

# Connect to the database directly
psql -h localhost -U parklah -d parklah_db -p 5433

# Run database migrations (requires Goose installed)
bash scripts/migrate_up.sh

# Check migration status
source .env && goose -dir ./db/migrations postgres "$DB_URL" status

# Install/sync dependencies
go mod tidy
```

## Environment Setup

Requires a `.env` file with:
- `POSTGRES_*` vars (user, password, db, host, port, sslmode, container_name)
- `URA_ACCESS_KEY` — URA API access key
- `DB_URL` — full PostgreSQL DSN (used by Goose migrations)
- `ENV` — `development` or `production`

Config is loaded via `godotenv` (the only external dependency). `internal/config/utils.go` provides `getEnvOrFail` (panics if missing) and `getEnvOrDefault`.

## Architecture

### Intended Folder Structure

The project follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout). Planned structure as the API layer is built out:

```
cmd/
  ingestion/   # current: data ingestion worker
  api/         # future: REST API server
internal/
  config/      # env loading, Config struct, DSN builder
  ura/         # URA API client (token auth + generic callURAAPI[T] wrapper)
  hdb/         # HDB/data.gov.sg API client
  handlers/    # future: HTTP handlers
  services/    # future: background jobs
  repositories/ # future: DB access layer (currently internal/repository/)
  models/      # future: shared domain models
  db/          # future: DB connection setup
db/
  migrations/  # Goose SQL migration files
```

### Data Flow (Current)

```
cmd/ingestion/main.go
  └─ fetch.go (parallel sync.WaitGroup)
       ├─ internal/ura/ → URA API (token-based auth, 24h token refresh)
       └─ internal/hdb/ → data.gov.sg API (no auth required)
            ↓
       PostgreSQL (via db/migrations/ schema)
```

### Key Packages

**`internal/ura/`** — URA API client with token-based auth. `utils.go` provides a generic `callURAAPI[T]()` wrapper that handles token refresh via `ensureValidToken()`. The token is cached in `URAClient` and refreshed when expired.

**`internal/hdb/`** — Simple HTTP client for data.gov.sg. Fetches static carpark info and real-time availability.

**`internal/config/`** — Loads env vars into a `Config` struct; `Config.DSN()` builds the PostgreSQL connection string.

**`cmd/ingestion/fetch.go`** — Runs URA detail + season detail fetches in parallel. First error wins; no silent partial failures.

### Database Schema

Four tables managed by Goose migrations in `db/migrations/`:

| Table | Purpose |
|-------|---------|
| `carparks` | Master records; UNIQUE on `(carpark_code, data_source)` |
| `carpark_rates` | Rates by vehicle type and time window |
| `carpark_availability` | Time-series snapshots; indexed on `(carpark_code, snapshot_time DESC)` |
| `carpark_metadata` | HDB-specific info (decks, gantry height, basement flag) |

All tables use a `data_source` column (`ura`, `hdb`, or `lta`) so the same carpark code can exist across sources. Use `INSERT ... ON CONFLICT DO UPDATE` (UPSERT) for all writes.

**Coordinate system:** URA and HDB use SVY21 coordinates (`location_x`, `location_y`). These must be converted to WGS84 (lat/lon) before storing, so the mobile client can use them directly. LTA DataMall already provides WGS84.

### API Response Patterns

URA responses use a generic wrapper:
```go
type URAResponse[T any] struct {
    Status  string
    Message string
    Result  T
}
```

HDB models use nested structs matching the data.gov.sg response shape. LTA DataMall returns an OData envelope (`odata.metadata`, `value` array).

### Ingestion Cadence

- **Daily:** URA carpark details, season parking details, HDB static carpark info
- **Every 3–5 min:** URA availability, HDB availability, LTA availability

### Logging Convention

Log lines are prefixed with the package/operation context in brackets: `[hdb]`, `[ura]`, `[get-token]`, `[config]`.

## Planned Stack (Beyond This Repo)

| Layer | Technology |
|-------|-----------|
| Mobile | Kotlin + Jetpack Compose + Room + Retrofit |
| Cache | Redis (ElastiCache) |
| Auth | Amazon Cognito / JWT |
| Infra | AWS ECS + RDS + S3 + CloudWatch |
| CI/CD | GitHub Actions + Terraform |
