# Go API Project — Best Practice Folder Structure

Struktur ini menggunakan **Clean Architecture** yang umum dipakai di perusahaan besar.
Cocok untuk project jangka panjang, mudah dikembangkan dari kecil ke besar.

---

## Struktur Folder

```
myproject/
│
├── cmd/
│   └── api/
│       └── main.go                  ← Entry point, nyalain server
│
├── internal/                        ← Kode inti, tidak bisa diimport dari luar
│   │
│   ├── config/
│   │   └── env.go                   ← Baca .env, setup koneksi DB
│   │
│   ├── model/                       ← Struct GORM (representasi tabel DB)
│   │   ├── user_model.go
│   │   └── ...
│   │
│   ├── repository/                  ← Semua akses database (query/insert/update/delete)
│   │   ├── user_repository.go
│   │   └── ...
│   │
│   ├── service/                     ← Logic bisnis
│   │   ├── user_service.go
│   │   └── ...
│   │
│   ├── delivery/
│   │   └── http/
│   │       ├── handler/             ← Terima request, panggil service, kirim response
│   │       │   ├── user_handler.go
│   │       │   └── ...
│   │       ├── request/             ← Struct + validasi input dari client
│   │       │   ├── user_request.go
│   │       │   └── ...
│   │       ├── response/            ← Format output ke client
│   │       │   ├── user_response.go
│   │       │   └── common/
│   │       │       └── response.go  ← Format standar: {success, data, error}
│   │       └── router/              ← Daftarkan semua endpoint
│   │           ├── router.go
│   │           ├── user_route.go
│   │           └── ...
│   │
│   ├── middleware/                  ← Jalan sebelum/sesudah handler
│   │   ├── auth.go                  ← Validasi JWT
│   │   ├── logging.go               ← Log setiap request
│   │   └── error.go                 ← Handle panic/error global
│   │
│   └── database/
│       └── seeder/                  ← Data awal untuk development
│           └── user_seed.go
│
├── pkg/                             ← Utility yang bisa dipakai ulang (boleh diimport dari luar)
│   ├── logger/
│   │   └── logger.go                ← Structured logging
│   └── utils/
│       ├── jwt.go                   ← Generate/verify token
│       ├── hash.go                  ← Hash & verify password
│       ├── mailer.go                ← Kirim email via SMTP
│       └── string.go                ← Helper string
│
├── migrations/                      ← File SQL migration berurutan
│   ├── 20260101000000_create_users.up.sql
│   ├── 20260101000000_create_users.down.sql
│   └── ...
│
├── templates/                       ← Template HTML (email, PDF, dll)
│   ├── otp_email.html
│   └── welcome_email.html
│
├── assets/                          ← File statis (logo, gambar, font)
│   └── images/
│       └── logo.png
│
├── scripts/                         ← Script bantu (tidak masuk binary utama)
│   └── seeder/
│       └── main.go
│
├── files/                           ← Storage file upload (gitignore-kan)
│   ├── profile_picture/
│   ├── check_in/
│   └── ...
│
├── .env                             ← Konfigurasi lokal (jangan di-commit)
├── .env.example                     ← Template .env (wajib di-commit)
├── .air.toml                        ← Config hot reload (air)
├── .gitignore
├── go.mod
├── go.sum
└── Makefile                         ← Shortcut command: migrate, run, build, test
```

---

## Alur Dependency (Wajib Dijaga)

```
Request masuk
     ↓
  Router         ← daftarkan endpoint
     ↓
  Handler        ← terima request, validasi input, format response
     ↓
  Service        ← logic bisnis, kalkulasi, aturan
     ↓
 Repository      ← akses database
     ↓
  Database
```

**Aturan ketat:**
- Handler tidak boleh akses Repository langsung (harus lewat Service)
- Service tidak boleh tahu soal HTTP (tidak boleh pakai `*gin.Context`)
- Repository tidak boleh ada logic bisnis
- Kalau aturan ini dijaga, ganti database/framework tidak merusak yang lain

---

## Naming Convention

| Jenis File | Format | Contoh |
|---|---|---|
| Model | `nama_model.go` | `user_model.go` |
| Repository | `nama_repository.go` | `user_repository.go` |
| Service | `nama_service.go` | `user_service.go` |
| Handler | `nama_handler.go` | `user_handler.go` |
| Request | `nama_request.go` | `user_request.go` |
| Response | `nama_response.go` | `user_response.go` |
| Router | `nama_route.go` | `user_route.go` |
| Migration | `YYYYMMDDHHMMSS_aksi.up.sql` | `20260101000000_create_users.up.sql` |

---

## Pola Interface (Wajib di Repository & Service)

Setiap repository dan service harus punya interface-nya:

```go
// interface dulu
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*model.User, error)
    Create(ctx context.Context, user *model.User) error
}

// baru implementasinya
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}
```

Manfaatnya:
- Mudah di-mock saat testing
- Bisa ganti implementasi tanpa ubah yang lain
- Kode lebih mudah dibaca

---

## Format Response Standar

Semua endpoint pakai format yang sama:

```json
// Sukses
{
  "success": true,
  "data": { ... }
}

// Error
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Email tidak valid",
    "details": { ... }
  }
}
```

---

## Makefile (Sangat Direkomendasikan)

```makefile
run:
    air

build:
    go build -o bin/api cmd/api/main.go

migrate-up:
    migrate -path migrations -database "..." up

migrate-down:
    migrate -path migrations -database "..." down 1

seed:
    go run scripts/seeder/main.go

test:
    go test ./...

lint:
    golangci-lint run
```

---

## Tips Jangka Panjang

**Saat project kecil (1-2 developer):**
- Boleh satu file per domain (`user_service.go` handle semua logic user)
- Inject dependency manual di router

**Saat project mulai besar (3+ developer / fitur banyak):**
- Pertimbangkan pisah per fitur (feature-based) di dalam `internal/`
- Pertimbangkan `wire` atau `fx` untuk dependency injection otomatis
- Tambah folder `test/` untuk integration test

**Yang tidak boleh berubah:**
- Arah dependency (Handler → Service → Repository)
- Semua entity/model di `internal/model/`
- Interface di setiap layer
