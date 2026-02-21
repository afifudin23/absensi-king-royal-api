# Absensi King Royal API

Absensi King Royal backend API built with Go + Gin.

## Quick Setup

### 1) Prerequisites

- Go `1.25+`
- MySQL / MariaDB
- `make` (optional but recommended)

### 2) Clone and enter project

```bash
git clone <repo-url>
cd absensi-king-royal-api
```

### 3) Setup environment

```bash
cp .env.example .env
```

Edit `.env` for your local setup:

```env
APP_NAME=absensi-king-royal-api
DATABASE_URL=mysql://user:pass@tcp(localhost:3306)/absensi_king_royal_db
PORT=8080
```

### 4) Install dependencies and tools

```bash
make setup
```

If you do not use `make`:

```bash
go mod tidy
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/air-verse/air@latest
```

### 5) Run the app

```bash
make run
```

Hot reload:

```bash
make air
```

### 6) Database migration

Create a new migration (timestamp format):

```bash
make migrate-create name=create_users_table
```

Run migration:

```bash
make migrate-up
```

Rollback 1 step:

```bash
make migrate-down
```

### 7) Release version

```bash
go run ./scripts/release.go
```

The release script will:
- reject release if working tree is not clean,
- bump `internal/config/version.go`,
- auto run `git add .` + `git commit chore(release): vX.Y.Z`.

Before running release, always update `CHANGELOG.md` for feature/fix/breaking changes.

## Folder Overview

- `cmd/api`: app entry point (`main.go`).
- `configs`: environment/config files.
- `deployments`: deployment assets (Docker, etc).
- `docs`: project/API documentation.
- `internal`: internal business logic implementation.
- `migrations`: database migration files.
- `pkg`: reusable utilities across layers.
- `scripts`: development/ops helper scripts.
- `test`: test directories.
