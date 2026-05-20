# Feature Summary — Absensi King Royal API

Last updated: 2026-05-12

---

## ✅ Fitur yang Sudah Selesai

### 1. Auth
| Endpoint | Keterangan |
|---|---|
| `POST /api/v1/auth/register` | Daftar akun baru |
| `POST /api/v1/auth/login` | Login, return JWT token |
| `POST /api/v1/auth/logout` | Logout (invalidate token) |

- Password hashing dengan Argon2id
- JWT menggunakan HMAC SHA-256
- Auth middleware: validasi token + cek user aktif
- Role-based: `admin` dan `user`

---

### 2. User Management
| Endpoint | Keterangan | Role |
|---|---|---|
| `GET /api/v1/users` | List semua user | Admin |
| `POST /api/v1/users` | Buat user baru | Admin |
| `GET /api/v1/users/:user_id` | Detail user | Admin |
| `PUT /api/v1/users/:user_id` | Update user | Admin |
| `DELETE /api/v1/users/:user_id` | Hapus user (soft delete) | Admin |
| `GET /api/v1/users/me` | Lihat profil sendiri | User |
| `PUT /api/v1/users/me` | Update profil sendiri | User |

- Soft delete (akun yang dihapus tidak bisa login)
- User profile terpisah (relasi `user_profiles`)

---

### 3. Attendance (Absensi Harian)
| Endpoint | Keterangan | Role |
|---|---|---|
| `POST /api/v1/attendance/check-in` | Check-in dengan foto | User |
| `POST /api/v1/attendance/check-out` | Check-out dengan foto | User |
| `GET /api/v1/attendance/logs` | Riwayat absensi sendiri | User |
| `PATCH /api/v1/attendance/:attendance_id` | Edit data absensi manual | Admin |

- Status attendance: `present`, `off`, `sick`, `extra_off`, `absent`, `leave`
- Source tracking: `self_service`, `admin_edit`, `approved_request`, `system`
- Check-in/out bisa upload foto (file)
- Cek duplikasi: tidak bisa check-in 2x dalam 1 hari

---

### 4. Attendance Request (Pengajuan)
| Endpoint | Keterangan | Role |
|---|---|---|
| `GET /api/v1/attendance-requests` | List semua request | Admin |
| `POST /api/v1/attendance-requests` | Buat pengajuan baru | User |
| `GET /api/v1/attendance-requests/me` | Pengajuan milik sendiri | User |
| `GET /api/v1/attendance-requests/:id` | Detail pengajuan | User/Admin |
| `PUT /api/v1/attendance-requests/:id` | Edit pengajuan (pending only) | User |
| `PATCH /api/v1/attendance-requests/:id/status` | Approve / Reject | Admin |
| `DELETE /api/v1/attendance-requests/:id` | Hapus pengajuan | User |

- Tipe pengajuan: `sick`, `leave`, `extra_off`, `overtime`, `correction`
- Status: `pending` → `approved` / `rejected` / `cancelled`
- Saat di-approve, data absensi otomatis diupdate (upsert)
- Bisa lampir bukti foto (evidence file)
- Reviewer + review note tersimpan

---

### 5. File Upload
| Endpoint | Keterangan |
|---|---|
| `POST /api/v1/files` | Upload file (foto/bukti) |
| `DELETE /api/v1/files/:file_id` | Hapus file |

- Validasi: hanya gambar yang diterima
- Folder penyimpanan berdasarkan tipe file

---

### 6. Payroll Setting (Konfigurasi Gaji)
| Endpoint | Keterangan | Role |
|---|---|---|
| `GET /api/v1/payroll-settings` | List semua config | Admin |
| `POST /api/v1/payroll-settings` | Tambah config baru | Admin |
| `PATCH /api/v1/payroll-settings/:id` | Update satu config | Admin |
| `PUT /api/v1/payroll-settings/bulk` | Update banyak config sekaligus | Admin |
| `DELETE /api/v1/payroll-settings` | Hapus config | Admin |

- Config key unik (normalized)
- Contoh config: gaji pokok default, tarif lembur, pajak, tunjangan

---

### 7. Payroll (Penggajian)
| Endpoint | Keterangan | Role |
|---|---|---|
| `GET /api/v1/payrolls` | List semua payroll | Admin |
| `GET /api/v1/payrolls/:payroll_id` | Detail payroll | Admin |
| `POST /api/v1/payrolls/generate/:employee_id` | Generate payroll 1 karyawan | Admin |
| `POST /api/v1/payrolls/generate-all` | Generate payroll semua karyawan | Admin |
| `PUT /api/v1/payrolls/:payroll_id` | Update data payroll | Admin |
| `POST /api/v1/payrolls/:payroll_id/send` | Kirim slip gaji via email | Admin |

- Kalkulasi otomatis: gaji pokok + tunjangan + lembur - potongan - pajak
- Generate bersifat idempotent per bulan (update jika sudah ada)
- Generate PDF slip gaji
- Kirim slip via email (SMTP)
- Status: `unsent`, `sent`, `failed`

---

## ❌ Fitur yang Belum Ada

### High Priority
| Fitur | Keterangan |
|---|---|
| **Rekap Absensi (Report)** | Export/tampil rekap absensi per periode (bulanan/mingguan), belum ada endpoint laporan |
| **Kuota Cuti / Leave Balance** | Tidak ada tracking sisa jatah cuti per karyawan per tahun |
| **Dashboard / Statistik** | Tidak ada endpoint untuk ringkasan data (total hadir, absen, terlambat, dll) |
| **Filter & Pagination** | List endpoint belum support filter (by date, by user, by status) dan pagination |

### Medium Priority
| Fitur | Keterangan |
|---|---|
| **Jadwal Kerja / Shift** | Tidak ada manajemen shift, jam masuk/pulang standar, hari kerja |
| **Kalender Hari Libur** | Tidak ada data hari libur nasional / libur perusahaan |
| **Keterlambatan (Late)** | Tidak ada kalkulasi dan penanda telat masuk kerja |
| **Notifikasi** | Tidak ada notifikasi (email/push) saat request diapprove/reject |
| **Departemen / Divisi** | User tidak punya relasi ke department/divisi |

### Low Priority
| Fitur | Keterangan |
|---|---|
| **Refresh Token** | Auth hanya pakai access token, tidak ada refresh token |
| **Riwayat Payroll per Karyawan** | Tidak ada endpoint `GET /payrolls?employee_id=...` |
| **Audit Log** | Tidak ada log history perubahan data penting |
| **Export PDF / Excel** | Tidak ada export laporan ke file |

---

## Ringkasan Modul

| Modul | Status |
|---|---|
| Auth (Login/Register) | ✅ Selesai |
| User Management | ✅ Selesai |
| File Upload | ✅ Selesai |
| Absensi Harian (Check-in/out) | ✅ Selesai |
| Pengajuan Absensi (Request) | ✅ Selesai |
| Payroll Setting | ✅ Selesai |
| Penggajian (Payroll + Email) | ✅ Selesai |
| Kuota Cuti | ❌ Belum ada |
| Jadwal Kerja / Shift | ❌ Belum ada |
| Laporan / Rekap | ❌ Belum ada |
| Dashboard Statistik | ❌ Belum ada |
| Notifikasi | ❌ Belum ada |
| Departemen | ❌ Belum ada |
