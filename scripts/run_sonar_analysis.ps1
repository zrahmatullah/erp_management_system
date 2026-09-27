# ===============================================================================
# Cafe ERP System - Local SonarQube Scanner Runner (PowerShell)
# ===============================================================================

Write-Host "===============================================================================" -ForegroundColor Cyan
Write-Host "           Cafe ERP System - Local SonarQube Scanner Runner" -ForegroundColor Cyan
Write-Host "===============================================================================" -ForegroundColor Cyan
Write-Host ""

$rootDir = Split-Path -Parent $PSScriptRoot
Set-Location $rootDir

Write-Host "[1/3] Menjalankan Unit Tests Go & Menghasilkan Coverage Report..." -ForegroundColor Yellow
Set-Location "$rootDir\backend"
go test -v ./tests/unit/...
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Unit test gagal! Analisis SonarQube dibatalkan." -ForegroundColor Red
    Set-Location $rootDir
    exit $LASTEXITCODE
}
go test -coverprofile=coverage.out ./pkg/poscalc/... ./pkg/invcalc/... ./pkg/p2pstate/... ./pkg/crypto/...
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Gagal membuat coverage report!" -ForegroundColor Red
    Set-Location $rootDir
    exit $LASTEXITCODE
}
Set-Location $rootDir

Write-Host "[OK] Coverage report berhasil dibuat di backend\coverage.out" -ForegroundColor Green
Write-Host ""

Write-Host "[2/3] Memeriksa keberadaan SonarScanner CLI..." -ForegroundColor Yellow
$scannerExists = Get-Command "sonar-scanner" -ErrorAction SilentlyContinue

if (-not $scannerExists) {
    Write-Host "[INFO] 'sonar-scanner' belum terdaftar di system PATH Windows." -ForegroundColor Yellow
    Write-Host "Jika Anda menggunakan Docker Desktop, Anda bisa menjalankan scanner container:" -ForegroundColor White
    Write-Host "docker run --rm --network host -v `"${rootDir}:/usr/src`" sonarsource/sonar-scanner-cli" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Atau unduh SonarScanner CLI for Windows dari:" -ForegroundColor White
    Write-Host "https://docs.sonarsource.com/sonarqube/latest/analyzing-source-code/scanners/sonarscanner/" -ForegroundColor Cyan
    exit 1
}

Write-Host "[3/3] Menjalankan SonarScanner Analisis Kode..." -ForegroundColor Yellow
sonar-scanner `
    -Dsonar.projectKey=cafe-erp-system `
    -Dsonar.sources=backend,frontend/src `
    -Dsonar.go.coverage.reportPaths=backend/coverage.out `
    -Dsonar.qualitygate.wait=true

Write-Host ""
Write-Host "===============================================================================" -ForegroundColor Green
Write-Host "Analisis selesai! Lihat hasil lengkap di: http://localhost:9000/dashboard?id=cafe-erp-system" -ForegroundColor Green
Write-Host "===============================================================================" -ForegroundColor Green
