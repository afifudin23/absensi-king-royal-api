# Absensi King Royal API

## Overview Changelog

All notable changes to this project are documented in this file.

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
