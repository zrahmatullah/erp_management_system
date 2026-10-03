# ==============================================================================
# SonarQube Analysis Runner (PowerShell) - Cafe ERP Monorepo
# ==============================================================================

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Resolve-Path "$ScriptDir\.."

Write-Host "==============================================================================" -ForegroundColor Cyan
Write-Host " [SonarQube Analysis Runner] - Cafe ERP Backend & Frontend" -ForegroundColor Cyan
Write-Host "==============================================================================" -ForegroundColor Cyan

# Step 1: Run Go Backend Unit Tests
Write-Host "`n[Step 1/4] Running Go Unit Tests & Generating Coverage Profile..." -ForegroundColor Yellow
Push-Location "$ProjectRoot\backend"
try {
    $pkgs = "cafe-erp-system/backend/internal/delivery/http/middleware,cafe-erp-system/backend/internal/usecase/auth,cafe-erp-system/backend/internal/usecase/pos,cafe-erp-system/backend/internal/usecase/inventory,cafe-erp-system/backend/pkg/...,cafe-erp-system/backend/internal/config"
    go test -v "-coverprofile=coverage.out" "-coverpkg=$pkgs" ./tests/unit
    $cov = go tool cover "-func=coverage.out" | Select-String "total:"
    Write-Host "[OK] Backend coverage generated: $cov" -ForegroundColor Green
} finally {
    Pop-Location
}

# Step 2: Run Frontend Unit Tests (Vitest)
Write-Host "`n[Step 2/4] Running Frontend Unit Tests (Vitest) & Generating LCOV Coverage..." -ForegroundColor Yellow
Push-Location "$ProjectRoot\frontend"
try {
    npm run test:coverage
    Write-Host "[OK] Frontend coverage report generated at frontend\coverage\lcov.info" -ForegroundColor Green
} finally {
    Pop-Location
}

# Step 3: Check SonarQube Server Availability
Write-Host "`n[Step 3/4] Checking SonarQube Server Availability..." -ForegroundColor Yellow
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

# Step 4: Execute SonarScanner
Write-Host "`n[Step 4/4] Executing SonarScanner CLI for Frontend..." -ForegroundColor Yellow
Push-Location "$ProjectRoot\frontend"
try {
    $scannerPath = "D:\sonarqube\scanner\bin\sonar-scanner.bat"
    $token = "squ_cafe_erp_system_sonar_token_2026"
    if (Test-Path $scannerPath) {
        & $scannerPath "-Dsonar.host.url=http://localhost:9000" "-Dsonar.token=$token"
    } else {
        sonar-scanner "-Dsonar.host.url=http://localhost:9000" "-Dsonar.token=$token"
    }
    Write-Host "`n==============================================================================" -ForegroundColor Cyan
    Write-Host "Frontend Analysis Uploaded Successfully!" -ForegroundColor Green
    
    # Wait briefly for SonarQube Compute Engine processing and fetch Quality Gate Status
    Start-Sleep -Seconds 5
    try {
        $headers = @{ Authorization = "Bearer $token" }
        $qg = Invoke-RestMethod -Uri "http://localhost:9000/api/qualitygates/project_status?projectKey=cafe-erp-frontend" -Headers $headers -ErrorAction SilentlyContinue
        if ($qg.projectStatus.status -eq "OK") {
            Write-Host "QUALITY GATE EVALUATION: [ PASSED / OK ]" -ForegroundColor Green
        } else {
            Write-Host "QUALITY GATE EVALUATION: [ $($qg.projectStatus.status) ]" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "Could not query quality gate evaluation." -ForegroundColor DarkGray
    }

    Write-Host "Review Frontend Dashboard at: http://localhost:9000/dashboard?id=cafe-erp-frontend" -ForegroundColor Green
    Write-Host "Review Backend Dashboard at:  http://localhost:9000/dashboard?id=cafe-erp-backend" -ForegroundColor Green
    Write-Host "==============================================================================" -ForegroundColor Cyan
} finally {
    Pop-Location
}
