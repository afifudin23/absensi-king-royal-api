# CHANGELOG Guide

A short guide to keep `CHANGELOG.md` clean and consistent.

## Format

Use one section per release version:

```md
## [1.0.1] - 2026-02-22

### Added
- Added endpoint `GET /api/v1/users`.

### Changed
- Updated login response format.

### Fixed
- Fixed request validation issue.

### Breaking Changes
- `error` field format changed.
```

## Rules

- One bullet = one change.
- Write user/API impact, not low-level implementation detail.
- Put backward-incompatible changes in `Breaking Changes`.

## Required Before Release

Before release:

1. Update `CHANGELOG.md`.
2. Commit changelog together with code changes.

## Quick Template

```md
## [x.y.z] - YYYY-MM-DD

### Added
- 

### Changed
- 

### Fixed
- 

### Breaking Changes
- 
```
