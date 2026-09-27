# ==============================================================================
# SonarQube Analysis Runner (PowerShell) - Cafe ERP Monorepo
# ==============================================================================

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Resolve-Path "$ScriptDir\.."

Write-Host "==============================================================================" -ForegroundColor Cyan
Write-Host " [SonarQube Analysis Runner] - Cafe ERP Monorepo" -ForegroundColor Cyan
Write-Host "==============================================================================" -ForegroundColor Cyan

# Step 1: Run Go Unit Tests
Write-Host "`n[Step 1/3] Running Go Unit Tests & Generating Coverage Profile..." -ForegroundColor Yellow
Push-Location "$ProjectRoot\backend"
try {
    $pkgs = "cafe-erp-system/backend/internal/delivery/http/middleware,cafe-erp-system/backend/internal/usecase/auth,cafe-erp-system/backend/internal/usecase/pos,cafe-erp-system/backend/internal/usecase/inventory,cafe-erp-system/backend/pkg/...,cafe-erp-system/backend/internal/config"
    go test -v "-coverprofile=coverage.out" "-coverpkg=$pkgs" ./tests/unit
    $cov = go tool cover "-func=coverage.out" | Select-String "total:"
    Write-Host "[OK] Coverage report generated: $cov" -ForegroundColor Green
} finally {
    Pop-Location
}

# Step 2: Check SonarQube Server Availability
Write-Host "`n[Step 2/3] Checking SonarQube Server Availability..." -ForegroundColor Yellow
try {
    $res = Invoke-RestMethod -Uri "http://localhost:9000/api/system/status" -TimeoutSec 3 -ErrorAction SilentlyContinue
    if ($res.status -eq "UP") {
        Write-Host "[OK] SonarQube Server is UP and Healthy on port 9000 (v$($res.version))" -ForegroundColor Green
    } else {
        Write-Host "[WARNING] SonarQube status is $($res.status). Waiting..." -ForegroundColor Yellow
    }
} catch {
    Write-Host "[WARNING] Could not connect to http://localhost:9000. Proceeding anyway..." -ForegroundColor Yellow
}

# Step 3: Execute SonarScanner
Write-Host "`n[Step 3/3] Executing SonarScanner CLI..." -ForegroundColor Yellow
Push-Location "$ProjectRoot"
try {
    $scannerPath = "D:\sonarqube\scanner\bin\sonar-scanner.bat"
    $token = "squ_cafe_erp_system_sonar_token_2026"
    if (Test-Path $scannerPath) {
        & $scannerPath "-Dsonar.host.url=http://localhost:9000" "-Dsonar.token=$token"
    } else {
        sonar-scanner "-Dsonar.host.url=http://localhost:9000" "-Dsonar.token=$token"
    }
    Write-Host "`n==============================================================================" -ForegroundColor Cyan
    Write-Host "Analysis Completed Successfully!" -ForegroundColor Green
    Write-Host "Review Dashboard & Quality Gate at: http://localhost:9000/dashboard?id=cafe-erp-system" -ForegroundColor Green
    Write-Host "==============================================================================" -ForegroundColor Cyan
} finally {
    Pop-Location
}
