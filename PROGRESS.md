# Progress Development

## 2026-05-31

### Todo

- [X] Tambah field `OvertimeMinutes *int` di `internal/model/attendance_model.go` (setelah field `Note`)
- [X] Fix overtime case di `applyApprovedRequestToAttendance` — isi `overtime_minutes` ke attendance saat diapprove
- [X] Panggil `applyApprovedRequestToAttendance` di `UpdateStatus` saat status jadi `approved`
- [X] Tambah `OvertimeMinutes` di `AttendanceUpdateRequest` supaya admin bisa edit via PATCH
- [X] Jalankan `make migrate-up` untuk apply migration `add_overtime_minutes_to_attendances`
