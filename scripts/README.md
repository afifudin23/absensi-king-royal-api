# scripts

Development, CI/CD, and operations helper scripts.

## release.go

Updates app version in `internal/config/version.go`.

Examples:

```bash
go run ./scripts/release.go
go run ./scripts/release.go patch
go run ./scripts/release.go minor
go run ./scripts/release.go major
go run ./scripts/release.go set 1.2.0
go run ./scripts/release.go current
```
