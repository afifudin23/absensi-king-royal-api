# Feature Summary — Absensi King Royal API

Last updated: 2026-06-01 (rev 4)

---

## Keterangan

- ✅ Sudah ada & berfungsi
- ⚠️ Ada tapi belum sempurna / ada bug
- ❌ Belum ada

---

## 1. Auth

| Fitur                        | Endpoint                                    | Status | Catatan                                                    |
| ---------------------------- | ------------------------------------------- | ------ | ---------------------------------------------------------- |
| Login (email + password)     | `POST /api/v1/auth/login`                 | ✅     |                                                            |
| Logout                       | `POST /api/v1/auth/logout`                | ✅     | Token tidak di-blacklist (by design, butuh Redis)          |
| Register akun baru           | `POST /api/v1/auth/register`              | ✅     | Tidak dipakai langsung — karyawan dibuat admin             |
| Ganti password               | `PUT /api/v1/users/me/password`           | ✅     | Min 3 karakter (dev), kembalikan ke 8 sebelum production   |
| Lupa password — kirim OTP    | `POST /api/v1/auth/forgot-password`       | ✅     | OTP 6 digit, expire 10 menit                               |
| Reset password via OTP       | `POST /api/v1/auth/reset-password`        | ✅     | Min 3 karakter (dev)                                       |

---

## 2. User / Karyawan

| Fitur                          | Endpoint                                     | Status | Catatan                                               |
| ------------------------------ | -------------------------------------------- | ------ | ----------------------------------------------------- |
| Lihat profil sendiri           | `GET /api/v1/users/me`                     | ✅     |                                                       |
| Edit profil sendiri            | `PUT /api/v1/users/me`                     | ✅     | Nama & jabatan hanya admin yang bisa ubah             |
| Ganti foto profil              | `PUT /api/v1/users/me` (field `profile_picture_file_id`) | ✅ | Via `FileApi.upload` + update profile |
| Hapus foto profil              | `PUT /api/v1/users/me` (field null)         | ✅     |                                                       |
| List semua karyawan (admin)    | `GET /api/v1/users`                        | ✅     | `?search=&role=` tersedia                             |
| Tambah karyawan (admin)        | `POST /api/v1/users`                       | ✅     |                                                       |
| Edit karyawan (admin)          | `PUT /api/v1/users/:user_id`               | ✅     |                                                       |
| Hapus karyawan (admin)         | `DELETE /api/v1/users/:user_id`            | ✅     |                                                       |
| Toggle aktif/nonaktif          | —                                            | ❌     | Field `is_active` belum ada di model User — perlu migration |
| Reset password karyawan (admin)| `POST /api/v1/users/:user_id/reset-password` | ✅   | Admin only, tidak bisa reset password sendiri         |
| Status kepegawaian             | Field `employment_status` di `user_profiles` | ✅   | permanent / contract / internship / freelance          |

---

## 3. Absensi

| Fitur                        | Endpoint                                       | Status | Catatan                                      |
| ---------------------------- | ---------------------------------------------- | ------ | -------------------------------------------- |
| Absen masuk (check-in)       | `POST /api/v1/attendances/check-in`           | ✅     | Butuh `file_id` dari upload foto             |
| Absen pulang (check-out)     | `POST /api/v1/attendances/check-out`          | ✅     | Butuh `file_id` dari upload foto             |
| List log absensi             | `GET /api/v1/attendances/logs`                | ✅     | `?start_date=&end_date=&user_id=`            |
| Rekap absensi per bulan      | `GET /api/v1/attendances/recap`               | ✅     | `?month=&year=` — admin only                 |
| Edit absensi manual (admin)  | `PATCH /api/v1/attendances/:attendance_id`    | ✅     | Status, jam, catatan, lembur — admin only    |
| Validasi duplikasi check-in  | —                                              | ✅     | Validasi di service layer                    |

---

## 4. Pengajuan Izin / Cuti / Lembur

| Fitur                        | Endpoint                                            | Status | Catatan                                      |
| ---------------------------- | --------------------------------------------------- | ------ | -------------------------------------------- |
| Buat pengajuan               | `POST /api/v1/attendance-requests`                | ✅     | Tipe: sick, leave, extra_off, overtime        |
| Lihat pengajuan sendiri      | `GET /api/v1/attendance-requests/me`              | ✅     |                                              |
| List semua pengajuan (admin) | `GET /api/v1/attendance-requests`                 | ✅     | `?status=&type=&start_date=&end_date=`       |
| Approve / Reject (admin)     | `PATCH /api/v1/attendance-requests/:id/status`    | ✅     | Admin tidak bisa approve pengajuan sendiri   |
| Update absensi saat approved | — (internal)                                        | ✅     | `applyApprovedRequestToAttendance` otomatis  |
| Edit pengajuan               | `PUT /api/v1/attendance-requests/:id`             | ✅     | Belum ada UI Flutter                         |
| Hapus pengajuan (bulk)       | `DELETE /api/v1/attendance-requests`              | ✅     | `{ "ids": [...] }` — belum ada UI Flutter    |

---

## 5. Payroll / Slip Gaji

| Fitur                        | Endpoint                                          | Status | Catatan                                      |
| ---------------------------- | ------------------------------------------------- | ------ | -------------------------------------------- |
| List semua slip (admin)      | `GET /api/v1/payrolls`                          | ✅     |                                              |
| Lihat slip sendiri (user)    | `GET /api/v1/payrolls/me`                       | ✅     | Hanya slip dengan status `sent`             |
| Generate slip 1 karyawan     | `POST /api/v1/payrolls/generate/:employee_id`   | ✅     |                                              |
| Generate slip semua          | `POST /api/v1/payrolls/generate-all`            | ✅     |                                              |
| Edit komponen slip           | `PUT /api/v1/payrolls/:payroll_id`              | ✅     |                                              |
| Kirim slip via email         | `POST /api/v1/payrolls/:payroll_id/send`        | ✅     |                                              |
| Konfigurasi komponen gaji    | `GET/POST/PUT /api/v1/payroll-settings`         | ✅     | Belum ada UI Flutter                         |

---

## 6. File / Foto

| Fitur                    | Endpoint                          | Status | Catatan                            |
| ------------------------ | --------------------------------- | ------ | ---------------------------------- |
| Upload gambar            | `POST /api/v1/files`            | ✅     | Multipart/form-data, maks 5 MB     |
| Hapus file               | `DELETE /api/v1/files/:file_id` | ✅     | Hanya owner yang bisa hapus        |
| Validasi tipe file       | —                                 | ✅     | Sniff content-type (hanya gambar)  |

---

## 7. Log Aktivitas

| Fitur                          | Endpoint                         | Status | Catatan                                              |
| ------------------------------ | -------------------------------- | ------ | ---------------------------------------------------- |
| List log aktivitas             | `GET /api/v1/activity-logs`    | ✅     | Admin only. `?page=&limit=&user_id=&method=&search=` |
| Auto-record setiap request     | Middleware `StructuredLoggingMiddleware` | ✅ | Hanya POST/PUT/PATCH/DELETE — GET tidak disimpan |
| Deskripsi human-readable (ID)  | `describeActivity()` di middleware | ✅  | Bahasa Indonesia — "Absen masuk", "Generate slip gaji", dll. |
| Simpan user name               | — (lookup saat save)             | ✅     | `full_name` dari tabel `users`                       |

---

## Bug yang Masih Ada

| # | Bug                                          | File                  | Status |
| - | -------------------------------------------- | --------------------- | ------ |
| 1 | Token tidak di-blacklist saat logout         | `auth_service.go`    | ⏸️ By design — butuh Redis blacklist |
| 2 | Field `is_active` tidak ada di User model    | `user_model.go`       | ❌ Perlu migration + update handler |
| 3 | Min password 3 karakter (harusnya 8 di prod) | `auth_request.go`, `user_request.go` | ⚠️ Development only |

---

## Yang Perlu Dikerjakan

| # | Item                                          | Prioritas |
| - | --------------------------------------------- | --------- |
| 1 | Migration tambah `is_active` ke tabel `users` | 🟡 Sedang |
| 2 | Kembalikan min password ke 8 sebelum production | 🔴 Penting |
| 3 | Blacklist token saat logout (Redis)           | 🟢 Rendah |
