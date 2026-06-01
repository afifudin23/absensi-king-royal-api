# Fitur yang Di-generate AI

Daftar fitur/endpoint yang dibuat oleh AI selama sesi development.

---

## API

1. **Query params `?start_date=&end_date=` di `GET /attendance/logs`**
   Supaya Flutter bisa filter log absensi per rentang tanggal — dipakai untuk riwayat, dashboard stats bulan ini, dan cek status absen hari ini.

2. **`GET /attendance/recap?month=&year=`**
   Endpoint baru untuk rekap absensi semua karyawan per bulan (aggregate per user). Dipakai di admin dashboard section rekap.

3. **`GET /payrolls/me`**
   Endpoint untuk karyawan lihat slip gaji sendiri — hanya slip dengan status `sent`.

4. **`POST /users/:user_id/reset-password`**
   Admin reset password karyawan lain. Tidak bisa reset password sendiri (ada validasi).

5. **Filter params `?status=&type=&start_date=&end_date=` di `GET /attendance-requests`**
   Admin bisa filter pengajuan izin berdasarkan status, tipe, dan rentang tanggal.

6. **Filter params `?search=&role=` di `GET /users`**
   Admin bisa search karyawan by nama/email dan filter by role.

7. **Validasi ukuran file maks 5 MB di `POST /files`**
   Cek `fileHeader.Size` sebelum proses upload — return 400 jika melebihi batas.

8. **`POST /auth/forgot-password` + `POST /auth/reset-password`**
   Flow lupa password via OTP email. OTP 6 digit, expire 10 menit, max 3x request — freeze 10 menit jika melebihi. Table `user_otps`.

