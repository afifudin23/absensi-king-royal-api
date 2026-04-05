# Absensi King Royal API

## Overview Changelog

All notable changes to this project are documented in this file.

## [0.6.1] - 2026-04-05

### Added

- Leave request status update endpoint: `PATCH /api/v1/leave-requests/:leave_id/status` (admin only).

### Changed

- Payroll setting bulk update now accepts wrapper payload `{ "settings": [...] }` (legacy array payload still supported).

---

## [0.6.0] - 2026-03-21

### Added

- Payroll settings module:
  - `GET /api/v1/payroll-settings`
  - `POST /api/v1/payroll-settings`
  - `PATCH /api/v1/payroll-settings/:payroll_id`
  - `PUT /api/v1/payroll-settings/bulk`
  - `DELETE /api/v1/payroll-settings/:payroll_id`
- Payroll settings domain implementation:
  - Model: `internal/model/payroll_setting_model.go`
  - Repository: `internal/repository/payroll_setting_repository.go`
  - Service: `internal/service/payroll_setting_service.go`
  - Handler: `internal/delivery/http/handler/payroll_setting_handler.go`
  - Request DTO: `internal/delivery/http/request/payroll_setting_request.go`
  - Response DTO: `internal/delivery/http/response/payroll_setting_response.go`
  - Router registration: `internal/delivery/http/router/payroll_setting_route.go`
- Payroll settings migration:
  - `migrations/20260321125730_create_payroll_settings_table.up.sql`
  - `migrations/20260321125730_create_payroll_settings_table.down.sql`
- Attendance update endpoint:
  - `PATCH /api/v1/attendance/:attendance_id`

### Changed

- Attendance service now supports manual update payloads for `check_in_at` and `check_out_at`.
- Attendance repository now supports lookup by attendance ID for update flow.
- Attendance migration now includes `status` column with default `present`.
- Main router now registers payroll settings routes and removes the unused photo route registration.
- User profile model no longer marks `user_id` as a primary key in GORM tags.

### Fixed

- Attendance patch route now calls the update handler instead of falling through to logs.
- Attendance update flow now persists updated timestamps without returning a nil object.
- Leave request update now uses the shared not-found helper for consistent 404 handling.
- Payroll setting bulk update now returns updated rows correctly instead of appending the GORM result object.
- Payroll setting bulk update now reports the missing `config_key` in not-found errors.
- Attendance check-in invalid file response now returns `Invalid file_id` consistently.

### Migration Required

- Yes.
- Run:
  - `make migrate-up`

---

## [0.5.1] - 2026-03-17

### Added

- User payload/response now include allowances:
  - position_allowance
  - other_allowance

### Changed

- User profile fields moved from `users` into new 1:1 table `user_profiles` and are eager-loaded via GORM `Preload("Profile")`.
- Migration `20260301153134` now creates `user_profiles` and links `user_profiles.profile_picture_id` to `files(id)` (replacing the previous FK on `users`).

### Fixed

- Seeder now generates UUID `users.id` and ensures `user_profiles` row exists; prevents empty-id inserts.

### Migration Required

- Yes.
- Run:
  - `make migrate-up`

---

## [0.5.0] - 2026-03-16

### Added

- Files module:
  - `POST /api/v1/files`
  - `DELETE /api/v1/files/:file_id`
  - Model: `internal/model/file_model.go`
  - Repository: `internal/repository/file_repository.go`
  - Service: `internal/service/file_service.go`
  - Handler: `internal/delivery/http/handler/file_handler.go`
  - Router registration: `internal/delivery/http/router/file_route.go`
  - Migrations:
    - `migrations/20260301151022_create_files_table.up.sql`
    - `migrations/20260301151022_create_files_table.down.sql`
    - `migrations/20260301153134_create_user_profiles.up.sql`
    - `migrations/20260301153134_create_user_profiles.down.sql`
- Unit tests:
  - `internal/service/user_service_test.go`

### Changed

- File references now use ID-only payloads:
  - Attendance check-in/out request now only sends `file_id`; URL is fetched from `files`.
  - Leave request evidence now only sends `evidence_file_id`; URL is fetched from `files`.
  - User profile picture now only sends `profile_picture_id`; URL is fetched from `files`.
- File validations added across modules (exists, ownership, expected file type) to avoid FK 1452 and return proper 4xx errors.
- Attendance DB columns standardized to `check_in_*` / `check_out_*` naming.
- Evidence storage path now groups leave evidence under `files/evidence/<date>/...` while keeping `files.type` unchanged.
- Dependency injection standardized: repositories no longer instantiate DB internally; routers wire DB, repos, services, handlers.
- Context propagation added across handler → service → repository; GORM queries now use `db.WithContext(ctx)`.

### Fixed

- Attendance date-only storage no longer shifts day due to time zone conversions.
- Error message casing standardized (capitalize) for HTTP responses.

### Migration Required

- Yes.
- Run:
  - `make migrate-up`

---

## [0.4.0] - 2026-03-13

### Added

- Leave request module:
  - `POST /api/v1/leave-requests`
  - `GET /api/v1/leave-requests`
  - `GET /api/v1/leave-requests/me`
  - `GET /api/v1/leave-requests/:leave_id`
  - `PUT /api/v1/leave-requests/:leave_id`
  - `DELETE /api/v1/leave-requests/:leave_id`
- Leave request domain implementation:
  - Model: `internal/model/leave_request_model.go`
  - Repository: `internal/repository/leave_request_repository.go`
  - Service: `internal/service/leave_request_service.go`
  - Handler: `internal/delivery/http/handler/leave_request_handler.go`
  - Request DTO: `internal/delivery/http/request/leave_request.go`
  - Response DTO: `internal/delivery/http/response/leave_response.go`
  - Router registration: `internal/delivery/http/router/leave_request_route.go`
- Leave request migration:
  - `migrations/20260312034410_create_leave_requests_table.up.sql`
  - `migrations/20260312034410_create_leave_requests_table.down.sql`
- Seeder bootstrap:
  - `internal/database/seeder/user_seed.go`
  - `scripts/seeder/main.go`

### Changed

- App database bootstrap moved into dedicated package `internal/database/database.go`, and app context now uses the centralized DB initializer.
- Air config now watches `.env`, so env changes trigger reload during development.
- Common success response now includes reusable action payload helper:
  - `common.ToSuccessResponse(...)`
- User and attendance handlers now reuse shared `utils.GetCurrentUserID(...)` to remove duplicated auth context parsing.
- User response flow simplified to use shared action success response payloads.
- User service/repository cleanup:
  - parameter naming normalized,
  - redundant delete variable removed,
  - update flow kept consistent with existing partial update rules.
- Attendance migration naming standardized from `add_attendance_table` to `create_attendances_table`.

### Fixed

- Validation error field names now use API-facing `snake_case` keys such as `start_date` and `end_date`.
- Validation error messages no longer capitalize field names and now show explicit enum options for `oneof` failures.
- Leave request list response now returns `[]` instead of `null` when empty.
- Leave request create/update handlers now validate date-only input (`YYYY-MM-DD`) consistently and stop execution after returning bad-request responses.
- Leave request update flow now supports safe partial updates without nil pointer dereference and without resetting status unintentionally.
- Leave request `GetByID`, update, and delete flows now use explicit UUID filtering and proper not-found handling.
- User seed passwords are now hashed consistently for both default seeded accounts.

### Migration Required

- Yes.
- Run:
  - `make migrate-up`

---

## [0.3.0] - 2026-03-03

### Added

- Attendance module:
  - `POST /api/v1/attendance/check-in`
  - `POST /api/v1/attendance/check-out`
  - `GET /api/v1/attendance/logs`
- Attendance domain implementation:
  - Model: `internal/model/attendance_model.go`
  - Repository: `internal/repository/attendance_repository.go`
  - Service: `internal/service/attendance_service.go`
  - Handler: `internal/delivery/http/handler/attendance_handler.go`
  - Response DTO: `internal/delivery/http/response/attendance_response.go`
  - Router registration: `internal/delivery/http/router/attendance_route.go`
- Attendance migration:
  - `migrations/20260302162641_add_attendance_table.up.sql`
  - `migrations/20260302162641_add_attendance_table.down.sql`
  - Includes unique key per user per date (`uq_attendances_user_id_date`).
- Centralized structured logger package with simple API:
  - `logger.Configure(...)`
  - `logger.Info(...)`
  - `logger.Warn(...)`
  - `logger.Error(...)`

### Changed

- Request logging middleware now uses centralized structured logger and emits consistent JSON log fields:
  - `timestamp`, `level`, `service`, `environment`, `request_id`, `user_id`, `message`, `logger_name`, `http`
- Router middleware stack now uses `StructuredLoggingMiddleware()` (replacing `gin.Logger()`).
- Auth middleware now propagates `user_id` into logger context for downstream logs.
- App startup now configures logger using env config (`cmd/api/main.go`).
- Environment config now includes `ENVIRONMENT` (`internal/config/env.go`).
- User repository/service refactored to pointer-based signatures for create/read/update consistency with GORM best practice.
- User response helper names standardized:
  - `ToUserResponse`
  - `ToUserListResponse`
  - `ToUserSuccessResponse`
- `.env.example` updated and standardized:
  - Adds `ENVIRONMENT`
  - Uses `ACCESS_KEY`
  - Uses quoted values for consistency.

### Migration Required

- Yes.
- Run:
  - `make migrate-up`

---

## [0.2.0] - 2026-03-02

### Added

- User management module:
  - `POST /api/v1/users`
  - `GET /api/v1/users`
  - `GET /api/v1/users/:user_id`
  - `PUT /api/v1/users/:user_id`
  - `DELETE /api/v1/users/:user_id`
  - `GET /api/v1/users/me`
  - `PUT /api/v1/users/me`
- User service and request/response DTOs:
  - `internal/service/user_service.go`
  - `internal/delivery/http/request/user_request.go`
  - `internal/delivery/http/response/user_response.go`
- Global middleware package:
  - `internal/middleware/auth.go`
  - `internal/middleware/error.go`
- Shared service error helpers:
  - `internal/service/errors.go`

### Changed

- Authentication flow now uses bearer token verification with centralized env access key (`pkg/utils/jwt.go`, `internal/service/auth_service.go`).
- Auth middleware now validates token format, verifies JWT, and checks active (non-deleted) user existence before continuing request (`internal/middleware/auth.go`).
- Router structure updated to register user routes and new root route file naming (`internal/delivery/http/router/router.go`, `internal/delivery/http/router/user_route.go`, `internal/delivery/http/router/root_route.go`).
- Auth handlers and user handlers improved for validation handling and payload normalization order.
- Repository and service layers improved for duplicate email error mapping and password hashing on create/update.

### Fixed

- Prevented nil pointer panic on deleted-account login handling.
- Prevented inconsistent error response behavior for middleware-thrown errors.
- Fixed update endpoint binding issues that previously caused `422 Invalid request body` for valid JSON payloads.

---

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

---

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
