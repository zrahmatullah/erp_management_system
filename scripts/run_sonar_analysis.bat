@echo off
echo ===============================================================================
echo            Cafe ERP System - Local SonarQube Scanner Runner
echo ===============================================================================
echo.

cd /d "%~dp0\.."

echo [1/3] Menjalankan Unit Tests Go & Menghasilkan Coverage Report...
cd backend
go test -v ./tests/unit/...
if %errorlevel% neq 0 (
    echo [ERROR] Unit test gagal! Analisis SonarQube dibatalkan.
    exit /b %errorlevel%
)
go test -coverprofile=coverage.out ./pkg/poscalc/... ./pkg/invcalc/... ./pkg/p2pstate/... ./pkg/crypto/...
if %errorlevel% neq 0 (
    echo [ERROR] Gagal membuat coverage report!
    exit /b %errorlevel%
)
cd ..

echo [OK] Coverage report berhasil dibuat di backend/coverage.out
echo.

echo [2/3] Memeriksa keberadaan SonarScanner CLI...
where sonar-scanner >nul 2>nul
if %errorlevel% neq 0 (
    echo [PERHATIAN] 'sonar-scanner' belum terdaftar di system PATH Windows.
    echo.
    echo Pilihan cara menjalankan:
    echo 1. Jika SonarQube berjalan via Docker:
    echo    docker run --rm --network host -v "%cd%:/usr/src" sonarsource/sonar-scanner-cli
    echo.
    echo 2. Atau unduh SonarScanner CLI for Windows dari:
    echo    https://docs.sonarsource.com/sonarqube/latest/analyzing-source-code/scanners/sonarscanner/
    echo    dan tambahkan folder bin ke PATH.
    echo.
    pause
    exit /b 1
)

echo [3/3] Menjalankan SonarScanner Analisis Kode...
sonar-scanner -Dsonar.projectKey=cafe-erp-system -Dsonar.sources=backend,frontend/src -Dsonar.go.coverage.reportPaths=backend/coverage.out -Dsonar.qualitygate.wait=true

echo.
echo ===============================================================================
echo Analisis selesai! Lihat hasil lengkap di: http://localhost:9000/dashboard?id=cafe-erp-system
echo ===============================================================================
pause
