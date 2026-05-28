# Feature Summary — Absensi King Royal

Last updated: 2026-05-13

---

## Legend
- ✅ Sudah ada & berfungsi
- ⚠️ Ada tapi belum sempurna / ada bug
- ❌ Belum ada

---

## 1. Auth

| Fitur | API | Flutter | Catatan |
|---|---|---|---|
| Login (email + password) | ✅ `POST /api/v1/auth/login` | ✅ Terhubung | |
| Logout | ✅ `POST /api/v1/auth/logout` | ✅ Terhubung | Token tidak di-blacklist di server |
| Restore session (auto-login) | ✅ `GET /api/v1/users/me` | ✅ Terhubung | |
| Register akun baru | ✅ `POST /api/v1/auth/register` | ❌ Tidak ada halaman register | |
| Ganti password | ❌ Belum ada endpoint | ⚠️ Ada halamannya tapi hardcoded sukses | Perlu `POST /auth/change-password` |
| Reset password (admin) | ❌ Belum ada endpoint | ⚠️ Ada tombol di admin tapi tidak terhubung | Perlu `POST /users/:id/reset-password` |

---

## 2. Profil Karyawan

| Fitur | API | Flutter | Catatan |
|---|---|---|---|
| Lihat profil sendiri | ✅ `GET /api/v1/users/me` | ✅ Terhubung | |
| Edit profil sendiri | ✅ `PUT /api/v1/users/me` | ❌ Belum terhubung | Ada form di `profile_page` tapi tidak ada |
| Upload foto profil | ✅ `POST /api/v1/files` | ❌ Belum terhubung | Flutter upload ke galeri lokal saja |
| Hapus foto profil | ✅ `DELETE /api/v1/files/:file_id` | ❌ Belum terhubung | Hapus hanya lokal |
| Lihat riwayat slip gaji (user) | ❌ Belum ada endpoint | ⚠️ Ada halaman tapi pakai data mock | Perlu `GET /payrolls/me` |

---

## 3. Absen Masuk / Pulang

> Flow yang benar: Upload foto dulu → dapat `file_id` → kirim check-in/out dengan `file_id`

| Fitur | API | Flutter | Catatan |
|---|---|---|---|
| Upload foto absen | ✅ `POST /api/v1/files` | ❌ Belum terhubung | Flutter kirim `Uint8List` tapi API butuh multipart |
| Check-in dengan `file_id` | ✅ `POST /api/v1/attendance/check-in` | ❌ Belum terhubung | |
| Check-out dengan `file_id` | ✅ `POST /api/v1/attendance/check-out` | ❌ Belum terhubung | |
| Cek status absen hari ini | ❌ Belum ada endpoint | ⚠️ Status disimpan di state lokal saja | Perlu `GET /attendance/today` |
| Validasi duplikasi check-in | ✅ Ada di service (cek per tanggal) | — | |

---

## 4. Dashboard — Info Bulan Ini

| Fitur | API | Flutter | Catatan |
|---|---|---|---|
| Total hadir/off/cuti/sakit/lembur bulan ini | ❌ Belum ada endpoint | ⚠️ Hardcoded angka (22/2/1/2/0/7) | Perlu `GET /attendance/stats?month=&year=` |
| Status absen hari ini (sudah masuk/pulang?) | ❌ Belum ada endpoint | ⚠️ State lokal, reset tiap buka app | Perlu `GET /attendance/today` |
| Daftar pengajuan izin bulan ini | ✅ `GET /api/v1/attendance-requests/me` | ❌ Belum terhubung | Flutter pakai data mock |

---

## 5. Riwayat Absensi

| Fitur | API | Flutter | Catatan |
|---|---|---|---|
| List riwayat absensi sendiri | ✅ `GET /api/v1/attendance/logs` | ❌ Belum terhubung | Flutter pakai data mock |
| Filter per bulan/tahun | ❌ Belum support query params | ⚠️ Filter client-side dari mock data | Tambah `?month=&year=` ke `/logs` |
| Foto check-in/out di riwayat | ✅ Ada `check_in_file_url` di response | ❌ Belum terhubung | Flutter pakai asset lokal |
| Ringkasan per status (hadir/off/cuti/dll) | ❌ Belum ada | ⚠️ Dihitung client-side dari mock | Bisa dihitung client-side kalau data sudah nyata |

> **Catatan penting**: Flutter pakai enum `hadir/off/cuti/sakit/alfa/extraOff`, API pakai `present/off/leave/sick/absent/extra_off` — perlu mapping saat integrasi.

---

## 6. Ajukan Izin / Cuti / Lembur

| Fitur | API | Flutter | Catatan |
|---|---|---|---|
| Buat pengajuan | ✅ `POST /api/v1/attendance-requests` | ❌ Belum terhubung | Flutter pakai callback lokal |
| Upload bukti (foto dokter dll) | ✅ `POST /api/v1/files` | ❌ Belum terhubung | Flutter pilih file lokal saja |
| Lihat riwayat pengajuan sendiri | ✅ `GET /api/v1/attendance-requests/me` | ❌ Belum terhubung | Flutter pakai mock di `leaveHistory` |
| Edit pengajuan (pending only) | ✅ `PUT /api/v1/attendance-requests/:id` | ❌ Tidak ada halaman edit | |
| Batalkan pengajuan | ✅ `DELETE /api/v1/attendance-requests/:id` | ❌ Tidak ada tombol | |

> **Catatan field**: Flutter `lembur` → API `overtime`, Flutter `cuti` → API `leave`, Flutter `extraOff` → API `extra_off`, Flutter `sakit` → API `sick`.

> **Bug kritis di API**: Saat pengajuan di-*approve*, `applyApprovedRequestToAttendance()` **tidak pernah dipanggil**. Data absensi tidak ter-update otomatis. Perlu diperbaiki di `attendance_request_service.go`.

---

## 7. Admin — Manajemen Karyawan

| Fitur | API | Flutter | Catatan |
|---|---|---|---|
| List semua karyawan | ✅ `GET /api/v1/users` | ❌ Belum terhubung | Flutter pakai 3 data mock |
| Tambah karyawan baru | ✅ `POST /api/v1/users` | ❌ Belum terhubung | |
| Edit data karyawan | ✅ `PUT /api/v1/users/:user_id` | ❌ Belum terhubung | |
| Nonaktifkan karyawan (soft delete) | ✅ `DELETE /api/v1/users/:user_id` | ❌ Belum terhubung | |
| Reset password karyawan | ❌ Belum ada endpoint | ⚠️ Ada tombol tapi tidak terhubung | Perlu `POST /users/:id/reset-password` |
| Filter/search karyawan | ❌ Belum support query params | ⚠️ Filter client-side dari mock | Tambah `?search=&role=&status=` |

---

## 8. Admin — Rekap Absensi

| Fitur | API | Flutter | Catatan |
|---|---|---|---|
| Rekap absensi semua karyawan per bulan | ❌ Belum ada endpoint | ⚠️ Pakai mock data (3 karyawan, 2 bulan) | Perlu `GET /attendance/recap?month=&year=` |
| Detail absensi harian per karyawan | ❌ Belum ada endpoint terpisah | ⚠️ Pakai mock data harian | Bagian dari rekap |
| Edit absensi manual (admin) | ✅ `PATCH /api/v1/attendance/:attendance_id` | ❌ Belum terhubung | |
| Filter rekap per bulan/karyawan | ❌ Belum ada endpoint | ⚠️ Filter client-side | |
| Export Excel/PDF | ❌ Belum ada endpoint | ⚠️ Tombol ada, action mock | |

---

## 9. Admin — Approval Pengajuan Izin

| Fitur | API | Flutter | Catatan |
|---|---|---|---|
| List semua pengajuan | ✅ `GET /api/v1/attendance-requests` | ❌ Belum terhubung | Flutter pakai 5 mock request |
| Approve pengajuan | ✅ `PATCH /api/v1/attendance-requests/:id/status` | ❌ Belum terhubung | |
| Reject pengajuan | ✅ `PATCH /api/v1/attendance-requests/:id/status` | ❌ Belum terhubung | |
| Filter per status/tanggal | ❌ Belum support query params | ⚠️ Filter client-side | |
| Update absensi otomatis saat approved | ⚠️ Logic ada tapi **tidak dipanggil** | — | **Bug kritis** di API |

---

## 10. Admin — Payroll (Slip Gaji)

| Fitur | API | Flutter | Catatan |
|---|---|---|---|
| List semua slip gaji | ✅ `GET /api/v1/payrolls` | ❌ Belum terhubung | Flutter pakai `_SalarySlip` mock lokal |
| Generate slip 1 karyawan | ✅ `POST /api/v1/payrolls/generate/:employee_id` | ❌ Belum terhubung | |
| Generate slip semua karyawan | ✅ `POST /api/v1/payrolls/generate-all` | ❌ Belum terhubung | |
| Edit komponen slip gaji | ✅ `PUT /api/v1/payrolls/:payroll_id` | ❌ Belum terhubung | |
| Kirim slip via email | ✅ `POST /api/v1/payrolls/:payroll_id/send` | ❌ Belum terhubung | |
| Filter slip per karyawan | ❌ Belum support query params | ⚠️ Filter client-side | Tambah `?employee_id=` |
| Konfigurasi komponen gaji | ✅ `GET/POST/PUT /api/v1/payroll-settings` | ❌ Tidak ada halaman di Flutter | |

> **Catatan field payroll**: Flutter pakai `int` (rupiah bulat), API pakai `DECIMAL(15,2)` float. Perlu konversi.

> **Catatan overtime**: Flutter hardcode 25.000/jam, API ambil dari `payroll_settings.config_key = 'hourly_overtime_rate'` × 2. Kalkulasi aktual belum pakai data jam lembur nyata dari absensi.

---

## 11. Admin — Log Aktivitas

| Fitur | API | Flutter | Catatan |
|---|---|---|---|
| List log aktivitas | ❌ Belum ada endpoint | ⚠️ Ada UI lengkap tapi pakai mock | Perlu `GET /activity-logs` |
| Log otomatis saat ada aksi (approve, edit, dll) | ❌ Tidak ada sistem logging | ⚠️ Dibuat manual di Flutter state | |

---

## 12. File / Foto

| Fitur | API | Flutter | Catatan |
|---|---|---|---|
| Upload gambar (JPEG/PNG) | ✅ `POST /api/v1/files` | ❌ Belum terhubung | Harus multipart/form-data |
| Hapus file | ✅ `DELETE /api/v1/files/:file_id` | ❌ Belum terhubung | |
| Validasi hanya gambar | ✅ Ada (sniff content-type) | — | |
| Batas ukuran file | ❌ Tidak ada validasi size | — | |

---

## Ringkasan: Endpoint yang Perlu Dibuat di API

| # | Endpoint | Dibutuhkan Oleh | Prioritas |
|---|---|---|---|
| 1 | `GET /api/v1/attendance/today` | Dashboard — status absen hari ini | 🔴 Tinggi |
| 2 | `GET /api/v1/attendance/stats?month=&year=` | Dashboard — info bulan ini | 🔴 Tinggi |
| 3 | `GET /api/v1/attendance/logs?month=&year=` | Riwayat absensi (filter) | 🔴 Tinggi |
| 4 | `GET /api/v1/attendance/recap?month=&year=` | Admin rekap absensi bulanan | 🔴 Tinggi |
| 5 | `POST /api/v1/auth/change-password` | Halaman ganti password | 🟡 Sedang |
| 6 | `GET /api/v1/payrolls/me` | User lihat riwayat slip gaji | 🟡 Sedang |
| 7 | `POST /api/v1/users/:user_id/reset-password` | Admin reset password karyawan | 🟡 Sedang |
| 8 | `GET /api/v1/attendance-requests` + filter params | Admin filter pengajuan | 🟡 Sedang |
| 9 | `GET /api/v1/payrolls` + filter `?employee_id=` | Admin filter slip per karyawan | 🟡 Sedang |
| 10 | `GET /api/v1/activity-logs` | Admin log aktivitas | 🟢 Rendah |
| 11 | `GET /api/v1/users` + filter `?search=&role=` | Admin search karyawan | 🟢 Rendah |

---

## Bug Kritis di API yang Perlu Diperbaiki Sebelum Integrasi

| # | Bug | File | Dampak |
|---|---|---|---|
| 1 | `applyApprovedRequestToAttendance()` didefinisikan tapi **tidak pernah dipanggil** di `UpdateStatus()` | `attendance_request_service.go` | Approve izin tidak memperbarui data absensi |
| 2 | `PayrollSetting.IsActive` diabaikan saat generate payroll | `payroll_service.go` | Setting yang dinonaktifkan tetap dipakai |
| 3 | Token tidak di-blacklist saat logout | `auth_service.go` | Token lama masih valid sampai expired |
| 4 | Tidak ada validasi `CheckInAt < CheckOutAt` | `attendance_service.go` | Check-out bisa sebelum check-in |

---

## Mapping Enum Flutter ↔ API

### Tipe Absensi
| Flutter | API |
|---|---|
| `hadir` | `present` |
| `off` | `off` |
| `cuti` | `leave` |
| `sakit` | `sick` |
| `alfa` | `absent` |
| `extraOff` | `extra_off` |

### Tipe Pengajuan Izin
| Flutter | API |
|---|---|
| `sakit` | `sick` |
| `cuti` | `leave` |
| `extraOff` | `extra_off` |
| `lembur` | `overtime` |

### Status Pengajuan
| Flutter | API |
|---|---|
| `pending` | `pending` |
| `approved` | `approved` |
| `rejected` | `rejected` |

### Role
| Flutter | API |
|---|---|
| `'admin'` | `admin` |
| `'staff'` / `'user'` | `user` |

---

## Status Integrasi Keseluruhan

| Modul | API | Flutter | Integrasi |
|---|---|---|---|
| Auth (login/logout/session) | ✅ | ✅ | ✅ Selesai |
| Profil (lihat) | ✅ | ✅ | ✅ Selesai |
| Profil (edit + foto) | ✅ | ✅ | ❌ Belum |
| Absen masuk/pulang | ✅ | ✅ | ❌ Belum |
| Dashboard stats bulan ini | ❌ | ✅ | ❌ Tunggu API |
| Riwayat absensi | ✅ | ✅ | ❌ Belum |
| Ajukan izin | ✅ | ✅ | ❌ Belum |
| Admin - karyawan | ✅ | ✅ | ❌ Belum |
| Admin - rekap absensi | ❌ | ✅ | ❌ Tunggu API |
| Admin - approval | ⚠️ Bug | ✅ | ❌ Tunggu bug fix |
| Admin - payroll | ✅ | ✅ | ❌ Belum |
| Admin - log aktivitas | ❌ | ✅ | ❌ Tunggu API |
| Ganti password | ❌ | ✅ | ❌ Tunggu API |
| File upload | ✅ | ❌ | ❌ Belum |
