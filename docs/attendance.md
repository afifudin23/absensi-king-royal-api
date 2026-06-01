# Attendance — Penjelasan Lengkap

Dokumentasi ini menjelaskan semua endpoint, logic, dan alur data untuk modul **attendance** dan **attendance-requests**.

---

## Daftar Isi

1. [Model & Tabel](#model--tabel)
2. [Attendance Endpoints](#attendance-endpoints)
   - [POST /check-in](#post-check-in)
   - [POST /check-out](#post-check-out)
   - [GET /logs](#get-logs)
   - [PATCH /:attendance_id (Admin)](#patch-attendance_id-admin)
3. [Attendance Request Endpoints](#attendance-request-endpoints)
   - [POST / (Buat Pengajuan)](#post--buat-pengajuan)
   - [GET / (Semua Pengajuan)](#get--semua-pengajuan)
   - [GET /me (Pengajuan Sendiri)](#get-me-pengajuan-sendiri)
   - [GET /:id (Detail)](#get-id-detail)
   - [PUT /:id (Edit Pengajuan)](#put-id-edit-pengajuan)
   - [PATCH /:id/status (Approve/Reject — Admin)](#patch-idstatus-approvereject--admin)
   - [DELETE / (Bulk Delete)](#delete--bulk-delete)
4. [applyApprovedRequestToAttendance](#applyapprovedrequesttoattendance)

---

## Model & Tabel

### `attendances`

```go
type Attendance struct {
    ID             string           // PK uuid
    UserID         string           // FK ke users
    Status         AttendanceStatus // present | off | sick | extra_off | absent | leave
    Date           time.Time        // tanggal (date only)
    CheckInAt      *time.Time       // jam masuk, nullable
    CheckInFileID  *string          // FK ke files (foto masuk), UNIQUE
    CheckOutAt     *time.Time       // jam pulang, nullable
    CheckOutFileID *string          // FK ke files (foto pulang), UNIQUE
    Note           *string          // catatan, nullable
    Source         AttendanceSource // self_service | admin_edit | approved_request | system
    OvertimeHours  *int             // jam lembur, diisi saat request overtime diapprove
    EvidenceFileID *string          // FK ke files (bukti admin edit / bukti izin), UNIQUE
    UpdatedBy      *string          // admin yang terakhir ubah, nullable
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

**Source** menjelaskan dari mana data berasal:
| Nilai | Artinya |
|---|---|
| `self_service` | Karyawan absen sendiri (check-in/out) |
| `admin_edit` | Admin edit langsung |
| `approved_request` | Otomatis dari pengajuan yang diapprove |
| `system` | Dibuat oleh sistem |

**Unique constraints:**
- `check_in_file_id` — satu foto masuk hanya bisa dipakai di satu attendance
- `check_out_file_id` — satu foto pulang hanya bisa dipakai di satu attendance
- `evidence_file_id` — satu file bukti hanya bisa dipakai di satu attendance

---

### `attendance_requests`

```go
type AttendanceRequest struct {
    ID                     string                  // PK uuid
    UserID                 string                  // FK ke users (pengaju)
    Type                   AttendanceRequestType   // sick | leave | extra_off | overtime
    Status                 AttendanceRequestStatus // pending | approved | rejected
    StartDate              time.Time               // tanggal mulai
    EndDate                time.Time               // tanggal selesai
    RequestedOvertimeHours *int                    // jam lembur yang diminta (overtime)
    Reason                 string                  // alasan pengajuan
    EvidenceFileID         *string                 // bukti (misal surat dokter), nullable
    ReviewedBy             *string                 // admin yang review, nullable
    ReviewedAt             *time.Time              // waktu review, nullable
    ReviewNote             *string                 // catatan dari reviewer, nullable
    CreatedAt              time.Time
    UpdatedAt              time.Time
}
```

**Tipe pengajuan:**
| Tipe | Artinya |
|---|---|
| `sick` | Izin sakit |
| `leave` | Cuti |
| `extra_off` | Libur tambahan |
| `overtime` | Lembur |

> `attendance_requests` adalah antrian pengajuan (sementara). Setelah diapprove → efeknya disalin ke `attendances`. Kalau ditolak → tidak ada yang berubah di `attendances`.

---

## Attendance Endpoints

Base URL: `/api/v1/attendance`
Semua endpoint butuh JWT token (AuthMiddleware).

---

### POST /check-in

**Akses:** Semua user

**Request body:**
```json
{ "file_id": "uuid-foto-masuk" }
```

**Logic:**
```go
// 1. Validasi file — harus ada, tipe harus check_in, milik user ini
file, err := s.fileRepo.GetByID(ctx, payload.FileID)
if file.Type != model.FileTypeCheckIn { ... }
if file.UploadedBy != userID { ... }

// 2. Cek apakah sudah ada attendance hari ini
attendance, err := s.attendanceRepo.GetByUserAndDate(ctx, userID, today)

// 3. Kalau belum ada → buat baru
// 4. Kalau sudah ada → update (overwrite jam masuk)
attendance.CheckInAt = &now
attendance.CheckInFileID = &payload.FileID
attendance.Source = "self_service"
```

---

### POST /check-out

**Akses:** Semua user

**Logic:**
```go
// 1. Validasi file — harus tipe check_out, milik user ini
// 2. Cari attendance hari ini — wajib ada
if attendance.CheckInAt == nil { return error("harus check-in dulu") }
if attendance.CheckOutAt != nil { return error("sudah check-out") }

// 3. Update jam pulang
attendance.CheckOutAt = &now
attendance.CheckOutFileID = &payload.FileID
```

---

### GET /logs

**Akses:** Semua user

Ambil semua attendance milik user, urut dari terbaru. Tidak ada filter status — hadir, sakit, cuti semua tampil dalam satu list. Include URL foto masuk, pulang, dan evidence.

---

### PATCH /:attendance_id (Admin)

**Akses:** Admin only

**Request body (semua opsional):**
```json
{
  "status": "sick",
  "check_in_at": "08:00",
  "check_out_at": "17:00",
  "note": "Dikoreksi oleh admin",
  "overtime_hours": 2,
  "evidence_file_id": "uuid-file-bukti"
}
```

**Logic:**
```go
// Update field yang dikirim saja (partial update)
// Validasi evidence_file_id harus ada di tabel files
existing.Source = "admin_edit"
existing.UpdatedBy = &updaterID
```

> Admin bisa ubah status dari `present` ke `sick`/`leave`/dll dan sebaliknya tanpa perlu pengajuan.

---

## Attendance Request Endpoints

Base URL: `/api/v1/attendance-requests`

---

### POST / (Buat Pengajuan)

**Akses:** Semua user

**Request body:**
```json
{
  "type": "sick",
  "start_date": "2026-03-12",
  "end_date": "2026-03-12",
  "reason": "Demam tinggi",
  "evidence_file_id": "uuid-surat-dokter",
  "requested_overtime_hours": 2
}
```

Status awal selalu `pending`.

---

### GET / (Semua Pengajuan)

**Akses:** Semua user ⚠️ (sebaiknya admin only)

---

### GET /me (Pengajuan Sendiri)

**Akses:** Semua user

---

### GET /:id (Detail)

**Akses:** Semua user

---

### PUT /:id (Edit Pengajuan)

**Akses:** Semua user (hanya milik sendiri)

Bisa edit: `start_date`, `end_date`, `reason`, `type`, `evidence_file_id`, `requested_overtime_hours`.

---

### PATCH /:id/status (Approve/Reject — Admin)

**Akses:** Admin only

**Request body:**
```json
{ "status": "approved" }
```

**Validasi:**
```go
// Admin tidak boleh approve pengajuan miliknya sendiri
if existing.UserID == reviewerID {
    return error("You cannot approve your own request")
}
```

**Logic setelah update status:**
```go
if payload.Status == "approved" {
    s.applyApprovedRequestToAttendance(ctx, existing, reviewerID)
}
```

---

### DELETE / (Bulk Delete)

**Akses:** Semua user

**Request body:**
```json
{ "ids": ["id1", "id2", "id3"] }
```

Response:
```json
{
  "success": true,
  "data": { "deleted_count": 3 }
}
```

---

## applyApprovedRequestToAttendance

Dipanggil otomatis saat admin approve request. Menerapkan efek ke tabel `attendances`.

### Case: `sick` / `leave` / `extra_off`

```go
// Loop setiap hari dari start_date sampai end_date
for d := req.StartDate; !d.After(req.EndDate); d = d.AddDate(0, 0, 1) {
    upsertAttendanceForDay(ctx, req.UserID, day, func(a *model.Attendance) {
        a.Status = AttendanceStatus(req.Type) // "sick" / "leave" / "extra_off"
        a.Note = &req.Reason
        a.CheckInAt = nil      // tidak ada jam masuk saat izin
        a.CheckOutAt = nil
        a.CheckInFileID = nil
        a.CheckOutFileID = nil
        a.EvidenceFileID = req.EvidenceFileID // salin bukti dari pengajuan
        a.Source = "approved_request"
        a.UpdatedBy = &reviewerID
    })
}
```

### Case: `overtime`

```go
// Loop setiap hari dari start_date sampai end_date
for d := req.StartDate; !d.After(req.EndDate); d = d.AddDate(0, 0, 1) {
    upsertAttendanceForDay(ctx, req.UserID, day, func(a *model.Attendance) {
        a.OvertimeHours = req.RequestedOvertimeHours // jam lembur per hari
        a.Source = "approved_request"
        a.UpdatedBy = &reviewerID
    })
}
```

> `upsertAttendanceForDay` = kalau record sudah ada → update, kalau belum ada → buat baru.

---

## Alur Lengkap: Karyawan Izin Sakit

```
1. Upload surat dokter → POST /files → dapat file_id

2. Ajukan izin →
   POST /attendance-requests
   { type: "sick", start_date, end_date, reason, evidence_file_id }
   → status: "pending"

3. Admin approve →
   PATCH /attendance-requests/:id/status { status: "approved" }

4. applyApprovedRequestToAttendance →
   Loop tiap hari: upsert attendance dengan status "sick",
   note = reason, evidence_file_id = surat dokter

5. Karyawan lihat riwayat →
   GET /attendance/logs → tampil hari izin dengan status "sick" + URL bukti
```

---

## Alur Lengkap: Karyawan Lembur

```
1. Ajukan lembur →
   POST /attendance-requests
   { type: "overtime", start_date, end_date, reason, requested_overtime_hours: 2 }
   → status: "pending"

2. Admin approve →
   PATCH /attendance-requests/:id/status { status: "approved" }

3. applyApprovedRequestToAttendance →
   Loop tiap hari: upsert attendance dengan overtime_hours = 2

4. Data lembur tersimpan di attendances.overtime_hours
```

---

## Catatan Desain

- **Koreksi jam masuk/pulang** → admin langsung edit via `PATCH /attendance/:id`, tidak perlu pengajuan
- **Pindah status (hadir → izin)** → admin langsung edit via `PATCH /attendance/:id`
- **Hari tanpa record attendance** → dianggap tidak hadir, dipotong di payroll (tidak ada cron job)
- **Unique file** → satu file hanya bisa dipakai di satu attachment (per kolom, per tabel)
