# Absensi King Royal API

## Overview Changelog

All notable changes to this project are documented in this file.

## [0.1.1] - 2026-02-24

### Added

- Authentication module:
  - `POST /api/v1/auth/register`
  - `POST /api/v1/auth/login`
  - `POST /api/v1/auth/logout`
- User domain implementation:
  - User model (`internal/model/user_model.go`)
  - User repository (`internal/repository/user_repositoy.go`)
  - Auth service (`internal/service/auth_service.go`)
  - Auth request/response DTOs
  - Password hashing helper (Argon2id)
  - Access token generator (HMAC SHA-256 JWT)
- Database bootstrap flow:
  - App context initialization for env + DB connection (`internal/config/app_context.go`)
  - MySQL URL normalization for `DATABASE_URL` format `mysql://...` (`internal/config/database.go`)
- Users migration:
  - `migrations/20260222055214_create_users_table.up.sql`
  - `migrations/20260222055214_create_users_table.down.sql`
- Logging utility and request logging middleware scaffolding:
  - `pkg/logger/logger.go`
  - `internal/middleware/logging.go`

### Changed

- Server startup now initializes app context and DB connection before running HTTP server (`cmd/api/main.go`).
- Env config now supports `JWT_SECRET` and includes default fallback for local development.
- Router now registers auth routes under `/api/v1/auth`.
- Common app error codes expanded (including auth-specific codes), with improved centralized error handling and validation error formatting.
- `.env.example` updated to include JWT secret and DB-related vars.
- `Makefile` migration commands improved to normalize non-`tcp(...)` MySQL DSN forms.

### Removed

- Removed scaffold placeholder README files across project directories to keep repository clean.

### Migration Required

- Yes.
- Run:
  - `make migrate-up`

## [0.1.0] - 2026-02-21

### Added

- Initial project structure for Go + Gin (`cmd`, `internal`, `pkg`, `migrations`, `scripts`, `test`).
- Versioned API routing base `/api/v1`.
- Basic endpoints: `GET /api/v1/` and `GET /api/v1/health`.
- Common API response schema (`success`, `data`, `error`).
- Common app error model and error handler.
- Environment config loader (`APP_NAME`, `DATABASE_URL`, `PORT`).
- Makefile commands for setup, run, build, and migrations.
- Release helper script (`scripts/release.go`).
- Changelog writing guide (`CHANGELOG_GUIDE.md`).

### Changed

- Standardized HTTP response format across handlers.

### Fixed

- Cleaned release flow to reject dirty working tree before version bump.

### Breaking Changes

- None.
