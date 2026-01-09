# Project Bulky Backend

Backend API untuk platform recommerce B2B menggunakan Gin Framework (Golang) dengan arsitektur clean code dan best practices. Menyediakan API untuk manajemen produk, transaksi, user management, dan fitur-fitur pendukung lainnya.

## Tech Stack

- **Framework**: Gin Web Framework v1.9+
- **Database**: PostgreSQL 15+ dengan GORM ORM v2
- **Migration**: golang-migrate
- **Validation**: go-playground/validator
- **Authentication**: JWT (golang-jwt/jwt)
- **Hot Reload**: Air
- **Environment**: godotenv

## Struktur Project

```
project-bulky-be/
├── cmd/
│   ├── api/
│   │   └── main.go                    # Entry point aplikasi
│   └── seed/
│       └── main.go                    # Database seeder
├── docs/                              # Dokumentasi lengkap (API spec, schema, dll)
├── internal/
│   ├── config/
│   │   └── config.go                  # Konfigurasi aplikasi
│   ├── controllers/                   # HTTP handlers (25+ controllers)
│   │   ├── auth_controller.go
│   │   ├── auth_v2_controller.go
│   │   ├── admin_controller.go
│   │   ├── buyer_controller.go
│   │   ├── produk_controller.go
│   │   ├── kategori_produk_controller.go
│   │   ├── merek_produk_controller.go
│   │   ├── kondisi_produk_controller.go
│   │   ├── kondisi_paket_controller.go
│   │   ├── sumber_produk_controller.go
│   │   ├── tipe_produk_controller.go
│   │   ├── diskon_kategori_controller.go
│   │   ├── warehouse_controller.go
│   │   ├── banner_tipe_produk_controller.go
│   │   ├── banner_event_promo_controller.go
│   │   ├── hero_section_controller.go
│   │   ├── ulasan_controller.go
│   │   ├── alamat_buyer_controller.go
│   │   ├── force_update_controller.go
│   │   ├── mode_maintenance_controller.go
│   │   ├── app_status_controller.go
│   │   ├── public_controller.go
│   │   ├── master_controller.go
│   │   ├── helper.go
│   │   └── example_controller.go
│   ├── dto/                           # Data Transfer Objects
│   ├── middleware/
│   │   └── auth_middleware.go         # JWT Authentication & Authorization
│   ├── models/                        # Database models & entities (40+ files)
│   │   ├── admin.go
│   │   ├── admin_session.go
│   │   ├── buyer.go
│   │   ├── alamat_buyer.go
│   │   ├── kategori_produk.go
│   │   ├── merek_produk.go
│   │   ├── kondisi_produk.go
│   │   ├── kondisi_paket.go
│   │   ├── sumber_produk.go
│   │   ├── tipe_produk.go
│   │   ├── produk.go
│   │   ├── produk_gambar.go
│   │   ├── produk_dokumen.go
│   │   ├── warehouse.go
│   │   ├── diskon_kategori.go
│   │   ├── banner_tipe_produk.go
│   │   ├── banner_event_promo.go
│   │   ├── hero_section.go
│   │   ├── metode_pembayaran.go
│   │   ├── metode_pembayaran_group.go
│   │   ├── ppn.go
│   │   ├── pesanan.go
│   │   ├── pesanan_item.go
│   │   ├── pesanan_pembayaran.go
│   │   ├── pesanan_status_history.go
│   │   ├── informasi_pickup.go
│   │   ├── jadwal_gudang.go
│   │   ├── dokumen_kebijakan.go
│   │   ├── ulasan.go
│   │   ├── disclaimer.go
│   │   ├── force_update_app.go
│   │   ├── mode_maintenance.go
│   │   ├── role.go
│   │   ├── permission.go
│   │   ├── refresh_token.go
│   │   ├── activity_log.go
│   │   ├── request.go
│   │   ├── response.go
│   │   ├── sistem_kontrol_request.go
│   │   ├── sistem_kontrol_response.go
│   │   ├── ulasan_request.go
│   │   └── ulasan_response.go
│   ├── repositories/                  # Database operations layer (27+ files)
│   │   ├── admin_repository.go
│   │   ├── admin_session_repository.go
│   │   ├── auth_repository.go
│   │   ├── buyer_repository.go
│   │   ├── alamat_buyer_repository.go
│   │   ├── kategori_produk_repository.go
│   │   ├── merek_produk_repository.go
│   │   ├── kondisi_produk_repository.go
│   │   ├── kondisi_paket_repository.go
│   │   ├── sumber_produk_repository.go
│   │   ├── tipe_produk_repository.go
│   │   ├── produk_repository.go
│   │   ├── produk_gambar_repository.go
│   │   ├── produk_dokumen_repository.go
│   │   ├── warehouse_repository.go
│   │   ├── diskon_kategori_repository.go
│   │   ├── banner_tipe_produk_repository.go
│   │   ├── banner_event_promo_repository.go
│   │   ├── hero_section_repository.go
│   │   ├── pesanan_repository.go
│   │   ├── pesanan_item_repository.go
│   │   ├── ulasan_repository.go
│   │   ├── force_update_repository.go
│   │   ├── mode_maintenance_repository.go
│   │   ├── role_repository.go
│   │   ├── permission_repository.go
│   │   └── activity_log_repository.go
│   ├── routes/
│   │   ├── routes.go                  # Route definitions
│   │   └── auth_v2_routes.go          # Auth v2 routes
│   └── services/                      # Business logic layer (27+ files)
│       ├── admin_service.go
│       ├── auth_service.go
│       ├── auth_v2_service.go
│       ├── buyer_service.go
│       ├── alamat_buyer_service.go
│       ├── kategori_produk_service.go
│       ├── merek_produk_service.go
│       ├── kondisi_produk_service.go
│       ├── kondisi_paket_service.go
│       ├── sumber_produk_service.go
│       ├── tipe_produk_service.go
│       ├── produk_service.go
│       ├── produk_gambar_service.go
│       ├── produk_dokumen_service.go
│       ├── warehouse_service.go
│       ├── diskon_kategori_service.go
│       ├── banner_tipe_produk_service.go
│       ├── banner_event_promo_service.go
│       ├── hero_section_service.go
│       ├── ulasan_service.go
│       ├── force_update_service.go
│       ├── mode_maintenance_service.go
│       ├── master_service.go
│       ├── role_service.go
│       ├── permission_service.go
│       └── activity_log_service.go
├── logs/                              # Application logs
├── migrations/                        # Database migrations (50+ files)
├── pkg/
│   ├── database/
│   │   └── database.go                # Database connection
│   └── utils/                         # Helper functions
│       ├── response.go                # Response formatter
│       ├── jwt.go                     # JWT utilities
│       ├── password.go                # Password hashing
│       ├── slug.go                    # Slug generator
│       ├── validator.go               # Custom validators
│       ├── privacy.go                 # Privacy utilities
│       └── version.go                 # Version utilities
├── tmp/                               # Air hot-reload temp files
├── .air.toml                          # Air configuration
├── .env
├── .env.example                       # Environment template
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Arsitektur

Project ini menggunakan **Clean Architecture** dengan pemisahan layer:

```
┌─────────────────────────────────────────────────────────────┐
│                    External Request (HTTP)                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Controller (HTTP Handler)                  │
│  - Menerima HTTP request                                    │
│  - Validasi input                                           │
│  - Memanggil Service                                        │
│  - Return HTTP response                                     │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Service (Business Logic)                  │
│  - Business logic                                           │
│  - Orchestration                                            │
│  - Tidak tahu tentang HTTP                                  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Repository (Data Access)                   │
│  - CRUD operations                                          │
│  - Database queries                                         │
│  - Tidak tahu tentang business logic                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       Model (Entity)                         │
│  - Data structures                                          │
│  - Core business entities                                   │
│  - Request/Response DTO                                     │
└─────────────────────────────────────────────────────────────┘
```

**Layer Responsibilities:**

- **Controllers**: Menerima HTTP requests dan mengembalikan responses
- **Services**: Business logic dan validasi
- **Repositories**: Interaksi dengan database
- **Models**: Entity/schema database
- **DTO**: Data Transfer Objects untuk request/response
- **Middleware**: Authentication & Authorization
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
   
   DB_HOST=127.0.0.1
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=new_bulky_db
   
   JWT_ACCESS_EXPIRY=24h
   JWT_SECRET=your-secret-key-change-this
   ```

4. **Buat database**
   ```sql
   CREATE DATABASE new_bulky_db;
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

### Base URL
```
Development: http://localhost:8080/api/v1
Production:  https://api.yourdomain.com/api/v1
```

### Authentication

Gunakan JWT Bearer Token untuk endpoint yang memerlukan autentikasi:

```
Authorization: Bearer <access-token>
```

### Endpoint Groups

#### 1. Authentication & Authorization
- `POST /api/v1/auth/login` - Login (Admin/Buyer)
- `POST /api/v1/auth/logout` - Logout
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/profile` - Get user profile
- `PUT /api/v1/auth/profile` - Update user profile

#### 2. Master Data - Produk
- `GET /api/v1/kategori-produk` - List kategori produk
- `POST /api/v1/kategori-produk` - Create kategori
- `GET /api/v1/kategori-produk/:id` - Get detail kategori
- `PUT /api/v1/kategori-produk/:id` - Update kategori
- `DELETE /api/v1/kategori-produk/:id` - Delete kategori

- `GET /api/v1/merek-produk` - List merek produk
- `GET /api/v1/kondisi-produk` - List kondisi produk
- `GET /api/v1/kondisi-paket` - List kondisi paket
- `GET /api/v1/sumber-produk` - List sumber produk
- `GET /api/v1/tipe-produk` - List tipe produk (Read-only: Paletbox, Container, Truckload)
- `GET /api/v1/tipe-produk/:id` - Get detail tipe produk by ID
- `GET /api/v1/tipe-produk/slug/:slug` - Get detail tipe produk by slug

#### 3. Produk
- `GET /api/v1/produk` - List produk dengan filter & pagination
- `POST /api/v1/produk` - Create produk
- `GET /api/v1/produk/:id` - Get detail produk
- `PUT /api/v1/produk/:id` - Update produk
- `DELETE /api/v1/produk/:id` - Delete produk

#### 4. Admin Management
- `GET /api/v1/admin` - List admin
- `POST /api/v1/admin` - Create admin
- `GET /api/v1/admin/:id` - Get admin detail
- `PUT /api/v1/admin/:id` - Update admin
- `DELETE /api/v1/admin/:id` - Delete admin

#### 5. Buyer Management
- `GET /api/v1/buyer` - List buyer
- `POST /api/v1/buyer` - Create buyer
- `GET /api/v1/buyer/:id` - Get buyer detail
- `PUT /api/v1/buyer/:id` - Update buyer
- `DELETE /api/v1/buyer/:id` - Delete buyer
- `GET /api/v1/buyer/:id/alamat` - Get buyer addresses
- `POST /api/v1/buyer/:id/alamat` - Create buyer address

#### 6. Marketing
- `GET /api/v1/hero-section` - List hero section
- `GET /api/v1/banner-tipe-produk` - List banner tipe produk
- `GET /api/v1/banner-event-promo` - List banner event/promo
- `GET /api/v1/diskon-kategori` - List diskon per kategori

#### 7. Warehouse & Operational
- `GET /api/v1/warehouse` - List warehouse
- `POST /api/v1/warehouse` - Create warehouse

#### 8. Ulasan (Reviews)
- `GET /api/v1/ulasan` - List ulasan
- `POST /api/v1/ulasan` - Create ulasan
- `GET /api/v1/ulasan/:id` - Get ulasan detail

#### 9. System Control
- `GET /api/v1/app-status` - Check app status
- `GET /api/v1/force-update` - Check force update status
- `GET /api/v1/mode-maintenance` - Check maintenance mode

#### 10. Public Endpoints
- `GET /api/v1/public/master-data` - Get all master data

### Health Check
```
GET /health
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
  "message": "Operasi berhasil",
  "data": {...}
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Deskripsi error",
  "errors": [
    {
      "field": "nama",
      "message": "Nama wajib diisi"
    }
  ]
}
```

## Database Migration

Project ini menggunakan **golang-migrate** untuk database migration.

### Install golang-migrate

**Windows (Chocolatey):**
```bash
choco install golang-migrate
```

**Linux/Mac:**
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
mv migrate /usr/local/bin/migrate
```

### Menjalankan Migration

**Migration Up (Apply):**
```bash
migrate -path migrations -database "postgresql://postgres:root@localhost:5432/new_bulky_db?sslmode=disable" up
```

**Migration Down (Rollback):**
```bash
migrate -path migrations -database "postgresql://postgres:root@localhost:5432/new_bulky_db?sslmode=disable" down
```

**Migration ke Versi Tertentu:**
```bash
migrate -path migrations -database "postgresql://postgres:root@localhost:5432/new_bulky_db?sslmode=disable" goto <version>
```

**Cek Versi Migration Saat Ini:**
```bash
migrate -path migrations -database "postgresql://postgres:root@localhost:5432/new_bulky_db?sslmode=disable" version
```

### Membuat Migration Baru

```bash
migrate create -ext sql -dir migrations -seq <migration_name>
```

Contoh:
```bash
migrate create -ext sql -dir migrations -seq create_new_table
```

Ini akan membuat 2 file:
- `000XXX_create_new_table.up.sql` - Migration untuk apply
- `000XXX_create_new_table.down.sql` - Migration untuk rollback

### Migration Files

Project ini memiliki 50+ migration files yang mencakup:
- Master data (kategori, merek, kondisi, sumber, tipe produk)
- Admin & Buyer management
- Produk & warehouse
- Marketing (banner, hero section, diskon)
- Transaksi & pembayaran
- Ulasan & review
- Authentication & authorization (role, permission)
- Activity logging
- System control (force update, maintenance mode)

## Testing

```bash
go test ./...
```

Dengan coverage:
```bash
go test -cover ./...
```

## Features

Project ini mendukung fitur-fitur berikut:

### ✅ Implemented

- **Authentication & Authorization**
  - JWT-based authentication (access token & refresh token)
  - Role-based access control (RBAC)
  - Permission management
  - Activity logging & audit trail
  
- **Master Data Management**
  - Kategori Produk
  - Merek Produk
  - Kondisi Produk
  - Kondisi Paket
  - Sumber Produk
  - Tipe Produk
  - Warehouse
  
- **Product Management**
  - CRUD Produk
  - Upload gambar & dokumen produk
  - Filter & search produk
  - Pagination
  
- **User Management**
  - Admin management
  - Buyer management
  - Profile management
  - Address management (for buyers)
  
- **Marketing Features**
  - Hero Section
  - Banner Tipe Produk
  - Banner Event/Promo
  - Diskon per Kategori
  
- **Review & Rating**
  - Ulasan produk
  - Rating system
  
- **System Control**
  - Force Update App
  - Maintenance Mode
  - App Status Check
  
- **Operational**
  - Warehouse management
  - Jadwal pickup
  - Informasi pickup

<!-- ### 🚧 Planned/In Progress

- Transaction & Order Management
- Payment Gateway Integration
- Notification System
- Email Service
- File Storage (MinIO/S3)
- Redis Caching
- API Rate Limiting
- Swagger Documentation

## Common Response Format

### Success Response -->

<!-- ```json
{
    "success": true,
    "message": "Operasi berhasil",
    "data": { ... }
}
```

### Success Response (List dengan Pagination)

```json
{
    "success": true,
    "message": "Data berhasil diambil",
    "data": [ ... ],
    "meta": {
        "halaman": 1,
        "per_halaman": 10,
        "total_data": 100,
        "total_halaman": 10
    }
}
```

### Error Response

```json
{
    "success": false,
    "message": "Deskripsi error",
    "errors": [
        {
            "field": "nama",
            "message": "Nama wajib diisi"
        }
    ]
}
``` -->

## HTTP Status Codes

| Code | Deskripsi |
|------|-----------|
| 200 | OK - Request berhasil |
| 201 | Created - Data berhasil dibuat |
| 400 | Bad Request - Validasi gagal |
| 401 | Unauthorized - Token tidak valid |
| 403 | Forbidden - Tidak punya akses |
| 404 | Not Found - Data tidak ditemukan |
| 409 | Conflict - Data sudah ada (duplicate) |
| 500 | Internal Server Error |

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
| DB_NAME | Database name | new_bulky_db |
| JWT_ACCESS_EXPIRY | JWT access token expiry | 24h |
| JWT_SECRET | JWT secret key | - |

