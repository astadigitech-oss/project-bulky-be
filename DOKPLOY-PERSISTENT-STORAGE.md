# Persistent Storage Setup untuk Dokploy

## Overview

File upload (icon kategori, gambar produk, dll) akan disimpan ke **persistent storage** agar tidak hilang saat redeploy.

## Checklist Deployment ke Dokploy

### ✅ 1. Environment Variables

Di Dokploy UI > Application > **Environment**, set environment variables.

**Cara cepat:** Copy dari `.env.staging.example` dan sesuaikan nilai-nilainya:

```env
UPLOAD_PATH=/app/uploads
BASE_URL=https://api-staging.bulky.id
DB_HOST=<your-db-host>
DB_PASSWORD=<your-db-password>
JWT_SECRET=<generate-with-openssl-rand-base64-32>
# ... dst
```

**Catatan:** 
- Sesuaikan `BASE_URL` dengan domain staging/production Anda
- Gunakan secret key yang berbeda untuk staging dan production
- Lihat `.env.staging.example` untuk template lengkap

### ✅ 2. Volume Mount

Di Dokploy UI > Application > **Volumes & Mounts**:

1. Klik **Add Volume**
2. Isi:
   - **Volume Name:** `bulky-uploads`
   - **Mount Path:** `/app/uploads`
3. **Save**

### ✅ 3. Deploy

Setelah setup volume, klik **Redeploy**. File upload akan tersimpan persistent di volume `bulky-uploads`.

### ✅ 4. Testing

Upload file via endpoint (contoh kategori produk):
```
PUT /api/v1/panel/kategori-produk/{id}/upload
```

File akan tersimpan di volume dan accessible via:
```
https://api-staging.bulky.id/uploads/product-categories/{filename}
```

---

## Setup di Dokploy (Detail)

### 1. Environment Variables

Di Dokploy web UI, tambahkan environment variables:

```env
UPLOAD_PATH=/app/uploads
BASE_URL=https://api.bulky.id/uploads
```

**Penjelasan:**
- `UPLOAD_PATH`: Path di dalam container tempat file disimpan
- `BASE_URL`: URL publik untuk akses file (served via nginx/static route)

### 2. Volume Mapping

Di Dokploy web UI, buat dan mount volume:

**Langkah-langkah:**
1. Masuk ke aplikasi > **Volumes & Mounts**
2. Klik **Add Volume**
3. Isi field:
   ```
   Volume Name: bulky-uploads
   Mount Path: /app/uploads
   ```
4. **Save** dan **Redeploy**

**Penjelasan:**
- Volume Name: nama volume di Dokploy (managed volume)
- Mount Path: path di dalam container tempat volume di-mount
- Dokploy otomatis handle persistent storage di host server
- Data persist meskipun container di-redeploy

### 3. Static File Serving

**Option 1: Go Application (Recommended untuk Dokploy)**

Go application sudah handle serving static files via routes. Pastikan `BASE_URL` di environment variables sesuai dengan domain staging/production:

```env
BASE_URL=https://api-staging.bulky.id
```

**Option 2: Nginx Reverse Proxy (Optional)**

Jika ingin nginx serve static files untuk performance, konfigurasi nginx di Dokploy:

```nginx
# Proxy untuk API
location /api/ {
    proxy_pass http://container:8080;
    proxy_set_header Host $host;
}

# Direct serve untuk uploads (butuh akses ke named volume)
location /uploads/ {
    # Ini memerlukan nginx container yang share volume yang sama
    alias /app/uploads/;
    access_log off;
    expires 30d;
    add_header Cache-Control "public, immutable";
}
```

Untuk Dokploy, **Option 1** (Go serve) lebih simple dan tidak butuh konfigurasi nginx tambahan.

## Development Setup

Untuk development lokal:

```env
# .env
UPLOAD_PATH=./uploads
BASE_URL=http://localhost:8080/uploads
```

File akan disimpan di folder `./uploads` di project root.

## File Structure

```
/app/uploads/                          (atau ./uploads di development)
├── product-categories/
│   ├── uuid1.png                     (icon kategori)
│   ├── uuid2.jpg
│   └── kondisi/
│       ├── uuid3.png                 (gambar kondisi tambahan)
│       └── uuid4.jpg
├── products/
│   ├── uuid5.jpg                     (gambar produk)
│   └── uuid6.png
└── documents/
    ├── uuid7.pdf                     (dokumen produk)
    └── uuid8.pdf
```

## Testing Upload

### Via Apidog

**Request:**
```
Method: PUT
URL: http://localhost:8080/api/v1/panel/kategori-produk/{id}/upload
Headers:
  Authorization: Bearer {token}
  Content-Type: multipart/form-data

Body (form-data):
  nama: "Elektronik Updated" (Text)
  icon: [pilih file image] (File)
  is_active: "true" (Text)
```

**Response:**
```json
{
  "status": "success",
  "message": "Kategori produk berhasil diupdate",
  "data": {
    "id": "uuid-here",
    "nama": "Elektronik Updated",
    "slug": "elektronik-updated",
    "icon_url": "product-categories/uuid.png",  // Relative path
    "is_active": true
  }
}
```

**Full URL:** `https://api.bulky.id/uploads/product-categories/uuid.png`

## Backup Strategy

File uploads tersimpan di:
```
Host: /var/lib/dokploy/volumes/bulky-uploads
```

Untuk backup:
```bash
# Backup
tar -czf bulky-uploads-$(date +%Y%m%d).tar.gz /var/lib/dokploy/volumes/bulky-uploads

# Restore
tar -xzf bulky-uploads-20260107.tar.gz -C /
```

## Migration dari Local Storage

Jika sebelumnya file disimpan di container (akan hilang), migrate ke persistent storage:

```bash
# 1. Backup file dari container
docker cp {container-id}:/app/uploads ./uploads-backup

# 2. Copy ke persistent volume
cp -r ./uploads-backup/* /var/lib/dokploy/volumes/bulky-uploads/

# 3. Set permissions
chown -R dokploy:dokploy /var/lib/dokploy/volumes/bulky-uploads
chmod -R 755 /var/lib/dokploy/volumes/bulky-uploads
```

## Troubleshooting

### File tidak bisa di-upload

**Cek permission folder:**
```bash
ls -la /var/lib/dokploy/volumes/bulky-uploads
```

**Fix permission:**
```bash
sudo chown -R dokploy:dokploy /var/lib/dokploy/volumes/bulky-uploads
sudo chmod -R 755 /var/lib/dokploy/volumes/bulky-uploads
```

### File hilang setelah redeploy

Cek apakah volume sudah di-mount dengan benar:
```bash
docker inspect {container-id} | grep Mounts -A 10
```

Harus ada:
```json
{
  "Type": "bind",
  "Source": "/var/lib/dokploy/volumes/bulky-uploads",
  "Destination": "/app/uploads"
}
```

### URL file tidak bisa diakses

Pastikan:
1. `BASE_URL` di env sudah correct
2. Nginx config serve static files (optional)
3. Atau Go application serve static files (default)

## Notes

- File disimpan dengan UUID filename untuk menghindari collision
- Supported formats: jpg, png, webp, svg
- Old file otomatis dihapus saat update (jika ada)
- Relative path disimpan di database, full URL di-generate saat response
