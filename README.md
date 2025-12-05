# Project Bulky Backend

<!-- Backend API menggunakan Gin Framework (Golang) dengan arsitektur clean code dan best practices. -->

## Tech Stack

- **Framework**: Gin Web Framework
- **Database**: PostgreSQL dengan GORM ORM
- **Hot Reload**: Air
- **Environment**: godotenv

## Struktur Project

```
project-bulky-be/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point aplikasi
├── internal/
│   ├── config/
│   │   └── config.go                  # Konfigurasi aplikasi
│   ├── controllers/
│   │   └── example_controller.go      # HTTP handlers
│   ├── middleware/
│   │   └── auth_middleware.go         # Middleware (auth, logging, dll)
│   ├── models/
│   │   └── user.go                    # Database models
│   ├── repositories/
│   │   └── user_repository.go         # Database operations layer
│   ├── routes/
│   │   └── routes.go                  # Route definitions
│   └── services/
│       └── user_service.go            # Business logic layer
├── pkg/
│   ├── database/
│   │   └── database.go                # Database connection
│   └── utils/
│       └── response.go                # Helper functions
├── migrations/                         # Database migrations
├── logs/                              # Application logs
├── .air.toml                          # Air configuration
├── .env.example                       # Environment template
├── .gitignore
├── go.mod
└── README.md
```

## Arsitektur

Project ini menggunakan **Clean Architecture** dengan pemisahan layer:

- **Controllers**: Menerima HTTP requests dan mengembalikan responses
- **Services**: Business logic dan validasi
- **Repositories**: Interaksi dengan database
- **Models**: Entity/schema database
- **Middleware**: Auth, logging, CORS, dll
- **Utils**: Helper functions yang reusable

## Prerequisites

- Go 1.21 atau lebih tinggi
- PostgreSQL
- Air (untuk hot reload development)

## Installation

1. **Clone repository**
   ```bash
   git clone <repository-url>
   cd project-bulky-be
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Setup environment**
   ```bash
   copy .env.example .env
   ```
   
   Edit file `.env` sesuai konfigurasi database Anda:
   ```env
   APP_ENV=development
   APP_PORT=8080
   
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=bulky_db
   
   JWT_SECRET=your-secret-key-change-this
   ```

4. **Buat database**
   ```sql
   CREATE DATABASE bulky_db;
   ```

5. **Install Air untuk hot reload** (opsional, untuk development)
   ```bash
   go install github.com/air-verse/air@latest
   ```

## Running the Application

### Development Mode (dengan hot reload)

```bash
air
```

Server akan berjalan di `http://localhost:8080` dan otomatis reload saat ada perubahan code.

### Production Mode

```bash
go run cmd/api/main.go
```

Atau build terlebih dahulu:

```bash
go build -o bin/app.exe cmd/api/main.go
./bin/app.exe
```

## API Endpoints

### Health Check
```
GET /health
```

### API v1

#### Public Routes
```
GET /api/v1/example
```

#### Protected Routes (Memerlukan Bearer Token - Sementara belum ada logicnya)
```
GET /api/v1/protected
```

**Header:**
```
Authorization: Bearer <your-token>
```

## Development

### Menambahkan Endpoint Baru

1. **Buat Model** di `internal/models/`
2. **Buat Repository** di `internal/repositories/`
3. **Buat Service** di `internal/services/`
4. **Buat Controller** di `internal/controllers/`
5. **Daftarkan Route** di `internal/routes/routes.go`

### Contoh Response Format

**Success Response:**
```json
{
  "success": true,
  "message": "Success message",
  "data": {...}
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Error message",
  "error": "Error details"
}
```

## Database Migration

Untuk menjalankan migration (belum diimplementasikan):
```bash
# TODO: Setup migration tool
```

## Testing

```bash
go test ./...
```

Dengan coverage:
```bash
go test -cover ./...
```

## Project Configuration

### Air Configuration (`.air.toml`)

Hot reload dikonfigurasi untuk:
- Monitor perubahan file `.go`, `.html`, `.tpl`, `.tmpl`
- Auto rebuild saat ada perubahan
- Exclude folder: `tmp`, `vendor`, `logs`, `migrations`

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| APP_ENV | Application environment | development |
| APP_PORT | Server port | 8080 |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | - |
| DB_NAME | Database name | bulky_db |
| JWT_SECRET | JWT secret key | - |

