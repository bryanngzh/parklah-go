# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**ParkLah!** is a Singapore parking assistant app. This repo (`parklah-go`) is the Go backend — it runs a data ingestion pipeline and serves a REST API for the Android client.

**Two phases complete:**
- **Phase 1** — Ingestion pipeline: aggregates carpark data from URA and HDB into PostgreSQL
- **Phase 2** — REST API: serves carpark discovery, availability, and rates to the mobile client

**Data sources:**
- **URA** — token-based auth, daily static details + every 3–5 min availability
- **HDB** (data.gov.sg) — API key for availability, static carpark info
- **LTA DataMall** — not yet implemented

## Commands

```bash
# Run the REST API server (port 8080)
go run ./cmd/api/main.go

# Run the ingestion pipeline (one-shot)
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
- `DATA_GOV_API_KEY` — data.gov.sg API key (HDB availability)
- `DB_URL` — full PostgreSQL DSN (used by Goose migrations)
- `API_PORT` — API server port (default: `8080`)
- `ENV` — `development` or `production`

`internal/config/utils.go` provides `getEnvOrFail` (panics if missing) and `getEnvOrDefault`.

## Architecture

```
cmd/
  api/           # REST API server (chi router, graceful shutdown)
  ingestion/     # One-shot ingestion pipeline worker
internal/
  config/        # Env loading, Config struct, DSN builder
  db/            # pgxpool connection setup (Connect)
  ura/           # URA API client + transform (SVY21→WGS84, normalize)
  hdb/           # HDB API client + transform + rate derivation
  models/        # Shared domain structs (Carpark, Availability, etc.)
  repository/    # DB read/write functions (pgx.Batch upserts, CopyFrom, queries)
  services/      # Business logic layer (joins data from multiple repo calls)
  handlers/      # HTTP handlers (param parsing, validation, JSON response)
  util/          # SVY21→WGS84 coordinate conversion
db/
  migrations/    # Goose SQL migration files
```

### Data Flow

**Ingestion:**
```
cmd/ingestion/main.go
  ├─ URA: fetch Details + SeasonDetails → transform → UpsertCarparks + UpsertRates
  ├─ HDB: fetch CarparkInfo → transform → UpsertCarparks + UpsertFeatures + DeriveRates
  ├─ URA: fetch Availability → transform → InsertAvailabilityBatch (CopyFrom)
  └─ HDB: fetch Availability → transform → InsertAvailabilityBatch (CopyFrom)
```

**API:**
```
HTTP request → chi router → handlers/ → services/ → repository/ → PostgreSQL
```

### REST API Endpoints

All responses use `{"data": ..., "meta": ...}` envelope. Base path: `/v1`.

| Method | Path | Key params |
|--------|------|------------|
| GET | `/v1/carparks/nearby` | `lat`, `lon`, `radius` (default 600m, max 2000m), `vehicle_type` (C/M/H), `limit` (default 20, max 50) |
| GET | `/v1/carparks/{code}` | `?source=ura\|hdb` |
| GET | `/v1/carparks/{code}/availability` | `?source=ura\|hdb` |
| GET | `/v1/carparks/{code}/rates` | `?source=ura\|hdb` |

`/nearby` returns the N nearest carparks ordered by distance, joined with the latest availability snapshot for the requested `vehicle_type`. All carparks are returned regardless of snapshot age — `snapshot_time` is included so clients can show staleness.

### Database Schema

Five tables managed by Goose migrations in `db/migrations/`:

| Table | Purpose | Write pattern |
|-------|---------|---------------|
| `carparks` | Master records; UNIQUE on `(carpark_code, data_source)` | UPSERT |
| `carpark_short_term_rates` | Hourly rates by vehicle type + day type + time window | UPSERT |
| `carpark_season_rates` | Monthly season rates by vehicle type + ticket type | UPSERT |
| `carpark_availability` | Time-series lot snapshots (append-only) | INSERT via CopyFrom |
| `carpark_features` | HDB-specific metadata (decks, gantry, basement flags) | UPSERT |

All tables use `data_source` (`ura` or `hdb`) so the same carpark code can coexist across sources.

**Coordinate system:** URA and HDB provide SVY21 coordinates. `internal/util/coordinates.go` converts to WGS84 (lat/lon) at ingest time using the SLA Transverse Mercator formula. Everything stored is WGS84.

**Nearby query:** Uses Haversine formula with a bounding-box pre-filter (no PostGIS required). `LEAST(1.0, ...)` guards against `acos` domain errors from floating-point rounding.

### Key Packages

**`internal/ura/`** — API client (`client.go`) with 24h token auto-refresh, plus `transform.go` that converts raw API responses to domain models (normalises parking system, vehicle type, parses `"$0.60"` rate strings and `"07.00 AM"` time strings).

**`internal/hdb/`** — API client, `transform.go`, and `rates.go`. HDB has no rates API — `rates.go` hardcodes the standardised rate table derived from the HDB website, with lookup maps for the 16 Central Area carparks and 12 peak-hour carparks.

**`internal/repository/`** — Split into `carparks.go`, `rates.go`, `availability.go`, `features.go`. Each file contains both write functions (upserts used by ingestion) and read functions (queries used by the API). Uses `pgx.Batch`+`SendBatch` for upserts and `pgx.CopyFrom` for availability bulk inserts.

**`internal/services/carparks.go`** — Joins data from multiple repository calls. `GetNearby` fetches nearby carparks then batch-fetches their latest availability and merges by `(carpark_code, data_source)`. `GetCarparkDetail` fetches carpark + features + availability in sequence.

### Ingestion Cadence

- **Daily:** URA carpark details, season parking details, HDB static carpark info + derived rates
- **Every 3–5 min:** URA availability, HDB availability

The ingestion binary is currently one-shot. Scheduling (ticker loop or cron) is not yet implemented.

### Logging Convention

Log lines are prefixed with the package/operation context in brackets: `[hdb]`, `[ura]`, `[api]`, `[main]`, `[config]`.

## Planned Stack (Beyond This Repo)

| Layer | Technology |
|-------|-----------|
| Mobile | Kotlin + Jetpack Compose + Room + Retrofit |
| Cache | Redis (ElastiCache) |
| Auth | Amazon Cognito / JWT |
| Infra | AWS ECS + RDS + S3 + CloudWatch |
| CI/CD | GitHub Actions + Terraform |
