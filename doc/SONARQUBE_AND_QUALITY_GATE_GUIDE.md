# 📊 Panduan Lengkap SonarQube & Quality Gate Setup
**Sistem**: Cafe ERP Management System (Go Chi Backend + Vue 3 Frontend)  
**Versi**: 1.0.0 Enterprise  
**Standar**: SonarQube Community Edition / SonarCloud + Strict Quality Gate  
**Tanggal**: 27 September 2026  

---

## 1. Arsitektur & Alur Kerja Analisis Kode

SonarQube bertindak sebagai inspektur kualitas kode otomatis (*Static Code Analysis & Security Linter*) yang memverifikasi kepatuhan arsitektur, celah keamanan, dan cakupan tes sebelum kode dapat di-merge ke branch utama.

```
       +-----------------------------------------------------------+
       |                  Developer / Git Commit                   |
       +-----------------------------------------------------------+
                                     |
                                     v
       +-----------------------------------------------------------+
       |   Backend Unit Tests (Go)   |   Frontend Build (Vue 3)    |
       |  `go test -coverprofile`    |      `npm run build`        |
       +-----------------------------------------------------------+
                                     |
                         coverage.out (Laporan Uji)
                                     |
                                     v
       +-----------------------------------------------------------+
       |       SonarScanner CLI (Monorepo Scanner Engine)          |
       |             Membaca `sonar-project.properties`            |
       +-----------------------------------------------------------+
                                     |
                     Analisis AST, Rules & Keamanan
                                     |
                                     v
       +-----------------------------------------------------------+
       |         SonarQube Server (Port :9000 / SonarCloud)        |
       |         Evaluasi Kriteria "Cafe ERP Strict Gate"          |
       +-----------------------------------------------------------+
                        /                         \
                       /                           \
         [ Status: PASSED ]                    [ Status: FAILED ]
                |                                      |
                v                                      v
       +------------------+                   +------------------+
       |  Build Berhasil  |                   |  Build Gagal &   |
       |    PR Di-merge   |                   |  Block PR / Push |
       +------------------+                   +------------------+
```

---

## 2. Server Lokal: Docker Compose (`docker-compose.sonar.yml`)

File [`docker-compose.sonar.yml`](file:///d:/Project/cafe-erp-system/docker-compose.sonar.yml) telah disiapkan di root proyek untuk menjalankan server SonarQube lokal lengkap dengan database PostgreSQL independen:

### Perintah Menjalankan Server:
```bash
# Menjalankan SonarQube + Database di background
docker compose -f docker-compose.sonar.yml up -d

# Memeriksa status log container
docker compose -f docker-compose.sonar.yml logs -f sonarqube

# Menghentikan server
docker compose -f docker-compose.sonar.yml down
```

### Akses Awal Web Console:
- **URL**: `http://localhost:9000`
- **Username Default**: `admin`
- **Password Default**: `admin` *(Sistem akan mewajibkan penggantian kata sandi pada login pertama)*.

---

## 3. Konfigurasi Strict Quality Gate (Kriteria Ketat)

Sesuai spesifikasi proyek, berikut adalah parameter **Quality Gate Ketat** yang dikonfigurasikan di SonarQube:

| Metrik Kualitas | Syarat / Batas Ambang | Kategori | Penjelasan & Tindakan |
|---|:---:|:---:|---|
| **Code Coverage** | **>= 75.0%** | Keandalan | Cakupan unit test minimum 75% pada kode baru dan keseluruhan kode. |
| **Vulnerabilities** | **= 0 (Zero Tolerance)** | Keamanan | 0 celah keamanan tingkat *Blocker*, *Critical*, atau *Major*. |
| **Hardcoded Secrets** | **= 0** | Keamanan | 0 kunci rahasia (API key, JWT secret, database password) yang tertanam di kode. |
| **Security Hotspots** | **100% Reviewed** | Keamanan | Seluruh bagian kode sensitif (kriptografi, auth) wajib sudah direviu. |
| **Maintainability Rating** | **A (Technical Debt < 5%)** | Pemeliharaan | Rasio hutang teknis maksimal 5%. |
| **Reliability Rating** | **A (0 Bugs Blocker/Critical)** | Keandalan | Bebas dari bug kritis yang dapat merusak alur aplikasi. |
| **Duplicated Lines** | **< 3.0%** | Kebersihan | Batas toleransi duplikasi kode maksimal 3%. |

### Langkah Pembuatan Quality Gate di SonarQube UI:
1. Buka `http://localhost:9000` &rarr; Masuk sebagai Administrator.
2. Navigasi ke menu **Quality Gates** &rarr; Klik tombol **Create**.
3. Beri nama: `Cafe ERP Strict Gate`.
4. Tambahkan kondisi sesuai tabel di atas:
   - Klik **Add Condition** &rarr; Pilih `Coverage` &rarr; Operator: `is less than` &rarr; Value: `75%`.
   - Klik **Add Condition** &rarr; Pilih `Vulnerabilities` &rarr; Operator: `is greater than` &rarr; Value: `0`.
   - Klik **Add Condition** &rarr; Pilih `Security Hotspots Reviewed` &rarr; Operator: `is less than` &rarr; Value: `100%`.
   - Klik **Add Condition** &rarr; Pilih `Duplicated Lines (%)` &rarr; Operator: `is greater than` &rarr; Value: `3%`.
5. Klik **Set as Default** agar otomatis diterapkan ke proyek `cafe-erp-system`.

---

## 4. Konfigurasi Proyek: `sonar-project.properties`

File konfigurasi monorepo [`sonar-project.properties`](file:///d:/Project/cafe-erp-system/sonar-project.properties) di root direktori telah dikonfigurasikan untuk mengenali arsitektur **Go + Vue 3**:

```properties
# Identifikasi Proyek
sonar.projectKey=cafe-erp-system
sonar.projectName=Cafe ERP Management System
sonar.projectVersion=1.0.0
sonar.sourceEncoding=UTF-8

# Sumber Kode Monorepo (Go Backend & Vue 3 Frontend)
sonar.sources=backend,frontend/src
sonar.tests=backend/tests,backend/pkg
sonar.test.inclusions=**/*_test.go,**/*.spec.ts,**/*.test.ts

# Pengecualian File (Dependensi, Aset, Build, Dokumen)
sonar.exclusions=\
  **/node_modules/**,\
  **/dist/**,\
  **/vendor/**,\
  **/doc/**,\
  **/*.docx,\
  **/*.png,\
  **/*.jpg,\
  **/*.log,\
  backend/migrations/**

# Laporan Coverage Go
sonar.go.coverage.reportPaths=backend/coverage.out

# Konfigurasi TypeScript / Frontend
sonar.typescript.tsconfigPath=frontend/tsconfig.json
sonar.javascript.lcov.reportPaths=frontend/coverage/lcov.info

# Menunggu evaluasi Quality Gate (otomatis exit non-zero jika gagal)
sonar.qualitygate.wait=true
```

---

## 5. Integrasi CI/CD Pipeline (Build Break Enforcement)

### 5.1. GitHub Actions (`.github/workflows/sonar-quality-gate.yml`)
Workflow telah dibuat di [`.github/workflows/sonar-quality-gate.yml`](file:///d:/Project/cafe-erp-system/.github/workflows/sonar-quality-gate.yml). 

**Cara Kerja Pipeline:**
1. Berjalan otomatis setiap ada `push` atau `pull_request` ke branch `dev` dan `main`.
2. Mengeksekusi unit test Go dan membuat `backend/coverage.out`.
3. Memeriksa apakah cakupan unit test memenuhi ambang minimum **75%**.
4. Menjalankan analisis kode SonarQube Scanner.
5. Memverifikasi status Quality Gate via `SonarSource/sonarqube-quality-gate-action`:
   - Jika Quality Gate **FAILED**, step ini akan **menggagalkan build GitHub Actions (exit code 1)** dan **memblokir Pull Request** dari penggabungan ke branch `main`.

**Setup Secret di GitHub Repository:**
- Buka **Settings &rarr; Secrets and variables &rarr; Actions &rarr; New repository secret**.
- Tambahkan:
  - `SONAR_TOKEN`: Token otentikasi proyek dari SonarQube (*Security &rarr; Users &rarr; Tokens*).
  - `SONAR_HOST_URL`: URL server SonarQube Anda (misal `http://sonarqube.yourdomain.com:9000` atau `https://sonarcloud.io`).

---

### 5.2. GitLab CI/CD (`.gitlab-ci.yml`)
Bagi tim yang menggunakan GitLab, pipeline telah disiapkan di [`.gitlab-ci.yml`](file:///d:/Project/cafe-erp-system/.gitlab-ci.yml):
- Tahap `test`: Menjalankan unit tests Go dan menyimpan artifact `coverage.out`.
- Tahap `sonarqube`: Menjalankan container `sonarsource/sonar-scanner-cli` dengan flag `-Dsonar.qualitygate.wait=true` sehingga pipeline GitLab langsung ditandai **FAILED** jika ada pelanggaran standar keamanan/coverage.

---

## 6. Panduan Menjalankan Analisis Lokal

### Menggunakan Helper Script Otomatis:

#### Di Windows (Command Prompt):
```cmd
scripts\run_sonar_analysis.bat
```

#### Di Windows (PowerShell):
```powershell
.\scripts\run_sonar_analysis.ps1
```

Script ini akan secara otomatis:
1. Menjalankan unit test Go dan memastikan seluruh 28 test lulus.
2. Menghasilkan file laporan coverage di `backend/coverage.out`.
3. Memeriksa SonarScanner CLI dan menjalankan pemindaian menyeluruh.
4. Menampilkan link langsung ke dashboard hasil analisis: `http://localhost:9000/dashboard?id=cafe-erp-system`.

---

## 7. Alternatif Tanpa Docker: SonarQube Standalone (Lokal Windows)

Karena komputer Anda sudah memiliki **Java 17 LTS**, Anda dapat menjalankan SonarQube lokal tanpa Docker:
1. Unduh **SonarQube Community Edition (ZIP)** dari [downloads.sonarsource.com](https://www.sonarsource.com/products/sonarqube/downloads/).
2. Ekstrak file zip ke `C:\sonarqube`.
3. Buka terminal atau jalankan: `C:\sonarqube\bin\windows-x86-64\StartSonar.bat`.
4. Tunggu pesan `SonarQube is operational`, lalu buka peramban di `http://localhost:9000`.
5. Unduh **SonarScanner CLI Windows (ZIP)**, ekstrak ke `C:\sonar-scanner`, dan masukkan folder `C:\sonar-scanner\bin` ke dalam System Environment Variables `PATH`.
