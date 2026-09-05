# Panduan Deployment - Cafe ERP Management System

Dokumen ini adalah panduan resmi deployment produksi untuk sistem **Cafe ERP Management System** yang mencakup Backend (Go Chi/Gin REST API), Frontend (Vue 3 + Vite + Tailwind CSS), dan Basis Data (PostgreSQL 16+).

---

## 1. Arsitektur Produksi & Prasyarat

### Komponen Sistem:
- **Backend Service**: Go binary (port `8080`), stateless, performa tinggi, koneksi pooling via `pgx/v5`.
- **Frontend SPA**: Static assets hasil build Vite yang dilayani melalui Web Server (Nginx / Caddy) dengan gzip & cache headers.
- **Database**: PostgreSQL 15/16+ dengan schema 35 tabel relasional terindeks.
- **Reverse Proxy**: Nginx menangani SSL termination (HTTPS Let's Encrypt), forwarding request API `/api/v1/*` ke backend Go, dan serving SPA routes `/` ke `index.html`.

### Prasyarat Server Minimum (VPS / Cloud):
- **CPU**: 2 Core (vCPU)
- **RAM**: 2 GB (Rekomendasi 4 GB untuk PostgreSQL + Redis)
- **Storage**: 25 GB SSD
- **OS**: Ubuntu 22.04 LTS / Debian 12 / Rocky Linux 9
- **Domain**: Domain / Subdomain aktif (contoh: `cafe-erp.yourdomain.com`)

---

## 2. Struktur Konfigurasi Environment (`.env`)

> [!CAUTION]
> **PENTING: JANGAN PERNAH MENYIMPAN ATAU MELAKUKAN COMMIT BERKAS `.env` KE DALAM GIT REPOSITORY!**
> Seluruh berkas `.env`, `.env.*`, dan `.env_database` telah dikecualikan secara ketat melalui `.gitignore`.

Gunakan berkas template [`.env.example`](file:///d:/Project/cafe-erp-system/.env.example) sebagai acuan:

```env
# Database Configuration
DB_CONNECTION=pgsql
DB_HOST=127.0.0.1
DB_PORT=5432
DB_DATABASE=cafe_erp
DB_USERNAME=cafe_user
DB_PASSWORD=GANTI_DENGAN_PASSWORD_DATABASE_SANGAT_KUAT
DB_SSLMODE=disable

# Backend API Configuration
APP_NAME="Cafe ERP System"
APP_ENV=production
APP_PORT=8080
APP_URL=https://cafe-erp.yourdomain.com
JWT_SECRET=GANTI_DENGAN_RANDOM_SECRET_MINIMAL_32_KARAKTER_HEX
JWT_EXPIRE_HOURS=24

# Redis (Opsional)
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=

# Frontend
VITE_API_BASE_URL=/api/v1
VITE_APP_TITLE="Cafe ERP Management System"
```

---

## 3. Opsi Deployment 1: Docker & Docker Compose (Direkomendasikan)

### A. Persiapan Berkas Docker

#### 1. `backend/Dockerfile`:
```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server cmd/server/main.go

# Runtime stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Jakarta
WORKDIR /app
COPY --from=builder /app/server .
COPY migrations/ ./migrations/
EXPOSE 8080
CMD ["./server"]
```

#### 2. `frontend/Dockerfile`:
```dockerfile
# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Production Nginx stage
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

#### 3. `docker-compose.yml` (Root):
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    container_name: cafe_erp_db
    restart: always
    environment:
      POSTGRES_DB: cafe_erp
      POSTGRES_USER: cafe_user
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./backend/migrations:/docker-entrypoint-initdb.d
    ports:
      - "127.0.0.1:5432:5432"

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: cafe_erp_backend
    restart: always
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_DATABASE: cafe_erp
      DB_USERNAME: cafe_user
      DB_PASSWORD: ${DB_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}
      APP_PORT: 8080
    depends_on:
      - postgres
    ports:
      - "127.0.0.1:8080:8080"

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: cafe_erp_frontend
    restart: always
    ports:
      - "80:80"
    depends_on:
      - backend

volumes:
  pgdata:
```

### B. Menjalankan Docker:
```bash
# 1. Salin template env
cp .env.example .env

# 2. Build dan jalankan seluruh container
docker compose up -d --build

# 3. Cek log status
docker compose logs -f
```

---

## 4. Opsi Deployment 2: Linux VPS / VM (Systemd & Nginx)

### Langkah 1: Setup Database PostgreSQL
```bash
sudo apt update && sudo apt install -y postgresql postgresql-contrib

# Masuk ke prompt postgres
sudo -u postgres psql

# Buat user dan database
CREATE USER cafe_user WITH ENCRYPTED PASSWORD 'password_sangat_aman_123';
CREATE DATABASE cafe_erp OWNER cafe_user;
GRANT ALL PRIVILEGES ON DATABASE cafe_erp TO cafe_user;
\q

# Jalankan skrip migrasi skema & seed
PGPASSWORD='password_sangat_aman_123' psql -U cafe_user -d cafe_erp -h 127.0.0.1 -f backend/migrations/000001_full_cafe_erp_schema.sql
PGPASSWORD='password_sangat_aman_123' psql -U cafe_user -d cafe_erp -h 127.0.0.1 -f backend/migrations/000002_complete_seed_data.sql
PGPASSWORD='password_sangat_aman_123' psql -U cafe_user -d cafe_erp -h 127.0.0.1 -f backend/migrations/000003_operational_seed_data.sql
```

---

### Langkah 2: Build & Jalankan Backend Service (Systemd)
```bash
# 1. Kompilasi binary Go untuk arsitektur Linux x86_64
cd /var/www/cafe-erp-system/backend
go build -ldflags="-w -s" -o /var/www/cafe-erp-system/backend/bin/server cmd/server/main.go
chmod +x /var/www/cafe-erp-system/backend/bin/server

# 2. Buat berkas unit Systemd
sudo nano /etc/systemd/system/cafe-erp-backend.service
```

Isi berkas `cafe-erp-backend.service`:
```ini
[Unit]
Description=Cafe ERP Backend API Daemon
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/var/www/cafe-erp-system/backend
ExecStart=/var/www/cafe-erp-system/backend/bin/server
Restart=always
RestartSec=5s
EnvironmentFile=/var/www/cafe-erp-system/backend/.env

# Security sandbox
LimitNOFILE=65535
StandardOutput=append:/var/log/cafe-erp-backend.log
StandardError=append:/var/log/cafe-erp-backend.error.log

[Install]
WantedBy=multi-user.target
```

```bash
# 3. Aktifkan dan jalankan daemon backend
sudo systemctl daemon-reload
sudo systemctl enable --now cafe-erp-backend
sudo systemctl status cafe-erp-backend
```

---

### Langkah 3: Build Frontend SPA
```bash
cd /var/www/cafe-erp-system/frontend

# Install dependencies & build
npm ci
npm run build

# Pastikan folder dist terbuat di: /var/www/cafe-erp-system/frontend/dist
```

---

### Langkah 4: Konfigurasi Nginx & SSL (Certbot)
```bash
sudo apt install -y nginx certbot python3-certbot-nginx
sudo nano /etc/nginx/sites-available/cafe-erp.conf
```

Isi berkas `cafe-erp.conf`:
```nginx
server {
    listen 80;
    server_name cafe-erp.yourdomain.com;

    # Root frontend build
    root /var/www/cafe-erp-system/frontend/dist;
    index index.html;

    # Gzip Compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript image/svg+xml;

    # Frontend SPA router handling
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Reverse Proxy ke Backend Go
    location /api/v1/ {
        proxy_pass http://127.0.0.1:8080/api/v1/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check endpoint
    location /health {
        proxy_pass http://127.0.0.1:8080/health;
    }

    # Static assets caching
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)$ {
        expires 1y;
        add_header Cache-Control "public, no-transform";
    }
}
```

```bash
# Aktifkan site Nginx
sudo ln -s /etc/nginx/sites-available/cafe-erp.conf /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

# Pasang SSL Gratis Let's Encrypt
sudo certbot --nginx -d cafe-erp.yourdomain.com
```

---

## 5. Kredensial Akun Bawaan (Default Login)

| Role | Username | Email | Password | Akses Utama |
| :--- | :--- | :--- | :--- | :--- |
| **Super Admin** | `admin` | `admin@cafe-erp.com` | `Admin@123` | Seluruh Modul + Master Data CRUD Hub |
| **Store Manager** | `manager` | `manager@cafe-erp.com` | `Admin@123` | Laporan, Inventory, POS, Presensi |
| **Cashier (Kasir)** | `kasir` | `kasir@cafe-erp.com` | `Admin@123` | POS Meja, Billing Tagihan, Pesanan |
| **Barista** | `barista` | `barista@cafe-erp.com` | `Admin@123` | Kitchen Display System (KDS) |
| **Warehouse** | `gudang` | `gudang@cafe-erp.com` | `Admin@123` | Stok Masuk, Purchase Order, Opname |

---

## 6. Prosedur Backup & Pemeliharaan Berkala

### Backup Otomatis Basis Data PostgreSQL (Cronjob):
```bash
# Buat script backup di /usr/local/bin/backup-cafe-erp.sh
#!/bin/bash
BACKUP_DIR="/var/backups/cafe_erp"
mkdir -p $BACKUP_DIR
DATE=$(date +\%Y\%m\%d_\%H\%M\%S)
pg_dump -U cafe_user -h 127.0.0.1 cafe_erp | gzip > "$BACKUP_DIR/db_$DATE.sql.gz"
# Hapus backup lebih tua dari 14 hari
find $BACKUP_DIR -type f -name "*.sql.gz" -mtime +14 -delete
```
Jalankan setiap hari jam 02:00 pagi:
```bash
0 2 * * * /usr/local/bin/backup-cafe-erp.sh
```
