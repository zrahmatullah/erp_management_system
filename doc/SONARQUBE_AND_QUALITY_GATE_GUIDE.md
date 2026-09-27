# Panduan Komprehensif Setup SonarQube & Strict Quality Gate

Dokumen ini menjelaskan arsitektur, instalasi, konfigurasi monorepo, dan penerapan **Strict Quality Gate** untuk sistem **Cafe ERP Management System** (Go Backend + Vue 3 Frontend).

---

## 1. Arsitektur & Spesifikasi Lingkungan

| Komponen | Spesifikasi / Lokasi | Keterangan |
|---|---|---|
| **SonarQube Server** | Community Edition 26.9.0.129388 (`D:\sonarqube`) | Berjalan di port `9000` |
| **Java Runtime** | OpenJDK 21 LTS (`D:\sonarqube\jdk-21`) | Wajib Java 21 untuk SonarQube 26+ |
| **Search Engine** | Embedded Elasticsearch 9.4.3 | Port HTTP: `9002` (resolusi konflik port 9001 dari service eksternal seperti Herd) |
| **SonarScanner CLI** | SonarScanner 8.1.0.6389 (`D:\sonarqube\scanner`) | Digunakan untuk eksekusi inspeksi lokal |
| **Monorepo Scope** | Go (`backend/`) + Vue 3 / TypeScript (`frontend/src/`) | Analisis terpadu dalam 1 project dashboard |
| **Test Coverage** | Go Coverprofile (`backend/coverage.out`) | 28 Unit Test passing, terhubung otomatis ke Sonar |

---

## 2. Definisi Strict Quality Gate

Berdasarkan spesifikasi Prompt 8, Quality Gate diset dengan kriteria ketat berikut:

```
[Strict Quality Gate Rules]
---------------------------------------------------------------------------------
1. Overall Code Coverage              >= 75.0%
2. Critical & Blocker Vulnerabilities = 0
3. Hardcoded Secrets Detected         = 0
4. Security Hotspots Reviewed         = 100.0%
5. Technical Debt / Maintainability   = Grade 'A' (Debt ratio <= 5%)
6. Code Duplication                   <= 3.0%
---------------------------------------------------------------------------------
Hasil Evaluasi: Jika salah satu syarat tidak terpenuhi -> BUILD GAGAL (Exit Code 1)
```

---

## 3. Konfigurasi Monorepo (`sonar-project.properties`)

File konfigurasi root `sonar-project.properties` telah dikonfigurasi khusus untuk Go Backend & Vue 3 Frontend:

```properties
# Project Identification
sonar.projectKey=cafe-erp-system
sonar.projectName=Cafe ERP Management System
sonar.projectVersion=1.0.0

# Server Connection
sonar.host.url=http://localhost:9000
sonar.token=squ_cafe_erp_system_sonar_token_2026
sonar.sourceEncoding=UTF-8

# Source Code Paths (Go backend + Vue 3 frontend)
sonar.sources=backend,frontend/src
sonar.tests=backend/tests

# Exclusions
sonar.exclusions=**/node_modules/**,**/dist/**,**/temp/**,**/vendor/**,**/*.min.js,**/*.svg,backend/coverage.out,backend/tests/**

# Go Coverage Integration
sonar.go.coverage.reportPaths=backend/coverage.out

# Frontend Environment
sonar.javascript.environments=browser,node
sonar.typescript.tsconfigPaths=frontend/tsconfig.json
```

---

## 4. Cara Menjalankan SonarQube Server

### Opsi A: Menjalankan Server Native (Sudah Terpasang di Sistem)
Server SonarQube dan Java 21 sudah terpasang dan siap digunakan:
```cmd
scripts\start_sonarqube_server.bat
```
Atau akses langsung melalui browser saat service aktif:
- **URL Dashboard**: [http://localhost:9000](http://localhost:9000)
- **Kredensial Default**:
  - Username: `admin`
  - Password: `admin` *(akan diminta mengganti password saat login pertama kali)*

### Opsi B: Menggunakan Docker Compose
Tersedia juga konfigurasi containerisasi dengan PostgreSQL 15:
```bash
docker compose -f docker-compose.sonar.yml up -d
```

---

## 5. Cara Menjalankan Analisis Kode Monorepo

Untuk menjalankan pengujian unit Go, kalkulasi test coverage, dan mengirimkan hasil analisis ke SonarQube:

### Menggunakan Batch Script (Windows CMD)
```cmd
scripts\run_sonar_analysis.bat
```

### Menggunakan PowerShell
```powershell
.\scripts\run_sonar_analysis.ps1
```

Tahapan yang dijalankan otomatis:
1. Menjalankan seluruh test suite Go (`go test -v -coverprofile=coverage.out -coverpkg=./... ./tests/unit`).
2. Melakukan health check koneksi SonarQube server (`http://localhost:9000/api/system/status`).
3. Menjalankan `sonar-scanner` CLI yang memetakan kode Go dan Vue 3 ke SonarQube.
4. Menghasilkan link dashboard visual hasil scan.

---

## 6. Integrasi CI/CD & Build Breaker

### A. GitHub Actions (`.github/workflows/sonar-quality-gate.yml`)
Workflow otomatis berjalan saat push/PR ke branch `main` atau `dev`:
- Menjalankan test Go & coverage.
- Melakukan SonarQube scan.
- **Memblokir build/merge** jika Quality Gate gagal menggunakan `sonarsource/sonarqube-quality-gate-action`.

### B. GitLab CI (`.gitlab-ci.yml`)
Pipeline 3 stage (`test` -> `analyze` -> `quality-gate`):
- Stage `quality-gate` memeriksa endpoint REST API SonarQube (`/api/qualitygates/project_status?projectKey=cafe-erp-system`).
- Apabila status respon bukan `OK`, pipeline otomatis mengembalikan `exit 1` sehingga deployment dibatalkan.

---

## 7. Catatan Teknis & Resolusi Isu

1. **Persyaratan Java 21**:
   - SonarQube versi 26.x ke atas membutuhkan runtime Java 21 LTS. Mesin default sebelumnya menggunakan Java 17, sehingga disiapkan Java 21 di `D:\sonarqube\jdk-21` tanpa mengganggu instalasi Java global user.
2. **Konflik Port Elasticsearch**:
   - Secara default Elasticsearch internal SonarQube menggunakan port `9001`. Port ini sering dipakai aplikasi lokal lain (seperti Laravel Herd).
   - Port search diubah ke `9002` melalui parameter `sonar.search.port=9002` pada `D:\sonarqube\conf\sonar.properties` untuk menjamin stabilitas.
