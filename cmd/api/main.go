// Package main is the entry point of the Absensi King Royal API.
//
//	@title			Absensi King Royal API
//	@version		1.0
//	@description	REST API untuk manajemen absensi, karyawan, dan penggajian King Royal.
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	Afifudin
//	@contact.email	afifudin23@example.com
//
//	@license.name	MIT
//
//	@host		localhost:8080
//	@BasePath	/api/v1
//
//	@tag.name		General
//	@tag.description	Root dan health check
//	@tag.name		Auth
//	@tag.description	Registrasi, login, dan manajemen password
//	@tag.name		Users
//	@tag.description	Manajemen data pengguna dan profil
//	@tag.name		Attendance
//	@tag.description	Check-in, check-out, dan rekap absensi
//	@tag.name		Attendance Requests
//	@tag.description	Pengajuan cuti, sakit, lembur, dan libur tambahan
//	@tag.name		Files
//	@tag.description	Upload dan manajemen file
//	@tag.name		Payroll Settings
//	@tag.description	Konfigurasi komponen penggajian
//	@tag.name		Payroll
//	@tag.description	Slip gaji dan penggajian karyawan
//	@tag.name		Activity Logs
//	@tag.description	Log aktivitas sistem (Admin only)
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Masukkan token JWT dengan format: **Bearer &lt;token&gt;** — contoh: `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
package main

import (
	"log"

	_ "github.com/afifudin23/absensi-king-royal-api/docs"
	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/router"
	"github.com/afifudin23/absensi-king-royal-api/pkg/logger"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("failed to initialize app context: %v", err)
	}
	defer func() {
		if err := config.CloseDB(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	env := config.GetEnv()
	if env == nil {
		log.Fatalf("failed to read app env")
	}
	logger.Configure(env.AppName, env.Environment)

	r := router.New()
	r.Static("/files", "./files")

	log.Printf("starting %s on %s", env.AppName, env.Port)
	if err := r.Run(env.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
