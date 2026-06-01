# Go API — Best Practice Project Structure
> Dipakai di perusahaan besar seperti Grab, Gojek, Tokopedia, dll.
> Pola: Clean Architecture + Domain Driven Design (DDD)

---

## Struktur Folder Lengkap

```
myproject/
│
├── cmd/                            ← Entry point aplikasi
│   ├── api/
│   │   └── main.go                 ← Start HTTP server
│   └── worker/
│       └── main.go                 ← Start background worker (opsional)
│
├── internal/                       ← Kode inti, tidak bisa diimport project lain
│   │
│   ├── config/                     ← Konfigurasi aplikasi
│   │   └── config.go               ← Baca env, setup DB, validasi config
│   │
│   ├── domain/                     ← Entity & interface bisnis (jantung project)
│   │   ├── user.go                 ← Struct User + method bisnis
│   │   ├── order.go
│   │   └── errors.go               ← Error bisnis: ErrNotFound, ErrUnauthorized
│   │
│   ├── repository/                 ← Interface akses data
│   │   ├── repository.go           ← Kumpulan semua interface repository
│   │   └── mysql/                  ← Implementasi MySQL
│   │       ├── user_repository.go
│   │       └── order_repository.go
│   │
│   ├── service/                    ← Logic bisnis
│   │   ├── service.go              ← Kumpulan semua interface service
│   │   ├── user_service.go
│   │   └── order_service.go
│   │
│   ├── handler/                    ← Delivery layer (HTTP/gRPC/dll)
│   │   └── http/
│   │       ├── router.go           ← Daftarkan semua route
│   │       ├── user_handler.go     ← Handler per domain
│   │       ├── order_handler.go
│   │       └── middleware/
│   │           ├── auth.go         ← JWT validation
│   │           ├── logging.go      ← Request logging
│   │           └── recovery.go     ← Handle panic
│   │
│   └── dto/                        ← Data Transfer Object (request & response shape)
│       ├── request/
│       │   ├── user_request.go     ← Struct input + validasi
│       │   └── order_request.go
│       └── response/
│           ├── base.go             ← Format standar: {success, data, error}
│           ├── user_response.go
│           └── order_response.go
│
├── pkg/                            ← Utility shared, boleh dipakai project lain
│   ├── errors/
│   │   └── errors.go               ← HTTP error helper: BadRequest, NotFound, dll
│   ├── logger/
│   │   └── logger.go               ← Structured logger (zerolog/zap)
│   ├── utils/
│   │   ├── jwt.go
│   │   ├── hash.go
│   │   ├── pagination.go
│   │   └── string.go
│   └── validator/
│       └── validator.go            ← Custom validation rules
│
├── migrations/                     ← SQL migration, urut berdasarkan timestamp
│   ├── 20260101000000_create_users_table.up.sql
│   ├── 20260101000000_create_users_table.down.sql
│   └── ...
│
├── config/                         ← File konfigurasi (bukan kode)
│   ├── config.yaml                 ← Config default
│   ├── config.development.yaml
│   └── config.production.yaml
│
├── scripts/                        ← Script bantu, tidak masuk binary
│   ├── seeder/
│   │   └── main.go
│   └── migration/
│       └── run.go
│
├── docs/                           ← Dokumentasi
│   ├── API.md
│   └── swagger.yaml
│
├── test/                           ← Integration & E2E test
│   ├── integration/
│   │   └── user_test.go
│   └── fixture/
│       └── user.json               ← Data dummy untuk test
│
├── .env.example                    ← Template env (wajib commit)
├── .env                            ← Env lokal (jangan commit)
├── .gitignore
├── .air.toml                       ← Hot reload config
├── go.mod
├── go.sum
├── Makefile                        ← Shortcut command
└── README.md
```

---

## Penjelasan Tiap Layer

### 1. `domain/` — Paling Penting
Isi entity utama dan interface. Tidak boleh import package lain selain standard library.
```go
// internal/domain/user.go
type User struct {
    ID        string
    Name      string
    Email     string
    CreatedAt time.Time
}

// internal/domain/errors.go
var (
    ErrUserNotFound   = errors.New("user not found")
    ErrEmailDuplicate = errors.New("email already exists")
    ErrUnauthorized   = errors.New("unauthorized")
)
```

### 2. `repository/` — Akses Data
Interface dulu di `repository.go`, implementasi di subfolder DB-nya.
```go
// internal/repository/repository.go
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    Create(ctx context.Context, user *domain.User) error
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id string) error
}
```

### 3. `service/` — Logic Bisnis
Interface dulu, implementasi di file yang sama.
```go
// internal/service/service.go
type UserService interface {
    GetByID(ctx context.Context, id string) (*domain.User, error)
    Create(ctx context.Context, payload dto.CreateUserRequest) (*domain.User, error)
    Update(ctx context.Context, id string, payload dto.UpdateUserRequest) (*domain.User, error)
    Delete(ctx context.Context, id string) error
}
```

### 4. `handler/http/` — HTTP Layer
Tidak ada logic bisnis di sini, cuma terima request → panggil service → kirim response.
```go
// internal/handler/http/user_handler.go
type UserHandler struct {
    userService service.UserService
}

func (h *UserHandler) GetByID(c *gin.Context) {
    user, err := h.userService.GetByID(c.Request.Context(), c.Param("id"))
    if err != nil {
        // handle error
        return
    }
    c.JSON(200, response.Success(user))
}
```

### 5. `dto/` — Input & Output Shape
Pisah dari domain supaya perubahan API tidak merusak logic bisnis.
```go
// internal/dto/request/user_request.go
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required,min=3"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// internal/dto/response/base.go
type Response struct {
    Success bool   `json:"success"`
    Data    any    `json:"data,omitempty"`
    Error   *Error `json:"error,omitempty"`
}
```

---

## Arah Dependency (Tidak Boleh Dilanggar)

```
cmd/api/main.go
      │
      ▼
  handler          ← tahu tentang HTTP
      │
      ▼
  service          ← tahu tentang bisnis
      │
      ▼
  repository       ← tahu tentang database
      │
      ▼
  database

  domain           ← tidak tahu siapa pun, semua tahu domain
```

**Aturan:**
- `handler` boleh pakai `service` dan `dto`
- `service` boleh pakai `repository` dan `domain`
- `repository` boleh pakai `domain`
- `domain` tidak boleh import layer lain

---

## Makefile (Wajib Ada)

```makefile
.PHONY: run build test migrate-up migrate-down seed lint

run:
	air

build:
	go build -o bin/api ./cmd/api

test:
	go test ./... -v -cover

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

seed:
	go run scripts/seeder/main.go

lint:
	golangci-lint run ./...

tidy:
	go mod tidy
```

---

## Format Response Standar

```json
// Sukses
{
  "success": true,
  "data": { "id": "...", "name": "..." }
}

// List dengan pagination
{
  "success": true,
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "limit": 20
  }
}

// Error
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Email tidak valid",
    "details": {
      "email": "must be a valid email"
    }
  }
}
```

---

## go.mod Minimal

```go
module github.com/username/myproject

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1       // HTTP framework
    gorm.io/gorm v1.25.0                  // ORM
    gorm.io/driver/mysql v1.5.0           // MySQL driver
    github.com/golang-jwt/jwt/v5 v5.0.0  // JWT
    github.com/joho/godotenv v1.5.1       // Load .env
    github.com/google/uuid v1.3.0         // UUID generator
    go.uber.org/zap v1.24.0               // Structured logger
)
```

---

## .env.example

```env
# App
APP_NAME="My Project"
APP_PORT=8080
APP_ENV=development        # development | staging | production

# Database
DATABASE_URL=mysql://user:pass@localhost:3306/mydb

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24

# SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=noreply@example.com
SMTP_PASSWORD=your-app-password
```

---

## Tips Scaling

| Ukuran Tim | Yang Perlu Ditambah |
|---|---|
| 1-2 dev | Struktur di atas sudah cukup |
| 3-5 dev | Tambah `test/integration/`, CI/CD pipeline |
| 5-10 dev | Pisah `internal/` per modul/fitur |
| 10+ dev | Pertimbangkan microservice, pisah repo |
