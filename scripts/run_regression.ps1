# ==============================================================================
# Cafe ERP System - Full Regression Test Suite Orchestrator
# Compatible with Windows PowerShell 5.1 & PowerShell Core 7+
# ==============================================================================

$ErrorActionPreference = "Stop"
$StartTime = Get-Date

Write-Host "======================================================================" -ForegroundColor Cyan
Write-Host "🚀 CAFE ERP SYSTEM - REGRESSION TESTING SUITE" -ForegroundColor Cyan
Write-Host "======================================================================" -ForegroundColor Cyan
Write-Host "Timestamp: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor Gray
Write-Host ""

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Resolve-Path "$ScriptDir\.."
$BackendDir = "$ProjectRoot\backend"
$FrontendDir = "$ProjectRoot\frontend"

$TestResults = @()

function Run-Step {
    param (
        [string]$Name,
        [scriptblock]$Command
    )

    Write-Host "----------------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host "▶ Running $Name..." -ForegroundColor Yellow
    $StepStart = Get-Date

    try {
        & $Command
        $Duration = (Get-Date) - $StepStart
        Write-Host "✅ $Name PASSED ($([math]::Round($Duration.TotalSeconds, 2))s)" -ForegroundColor Green
        $global:TestResults += [PSCustomObject]@{
            Stage    = $Name
            Status   = "PASS"
            Duration = "$([math]::Round($Duration.TotalSeconds, 2))s"
        }
    }
    catch {
        $Duration = (Get-Date) - $StepStart
        Write-Host "❌ $Name FAILED ($([math]::Round($Duration.TotalSeconds, 2))s): $_" -ForegroundColor Red
        $global:TestResults += [PSCustomObject]@{
            Stage    = $Name
            Status   = "FAIL"
            Duration = "$([math]::Round($Duration.TotalSeconds, 2))s"
        }
        Show-Summary
        exit 1
    }
}

function Show-Summary {
    $TotalDuration = (Get-Date) - $StartTime
    Write-Host ""
    Write-Host "======================================================================" -ForegroundColor Cyan
    Write-Host "📊 REGRESSION TEST SUMMARY" -ForegroundColor Cyan
    Write-Host "======================================================================" -ForegroundColor Cyan
    $TestResults | Format-Table -AutoSize
    
    $Failed = $TestResults | Where-Object { $_.Status -ne "PASS" }
    if ($Failed) {
        Write-Host "🚨 REGRESSION SUITE FAILED in $([math]::Round($TotalDuration.TotalSeconds, 2))s" -ForegroundColor Red
        exit 1
    } else {
        Write-Host "🎉 ALL REGRESSION SUITES PASSED in $([math]::Round($TotalDuration.TotalSeconds, 2))s" -ForegroundColor Green
        exit 0
    }
}

# -----------------------------------------------------------------------------
# 1. SMOKE TESTING (Sanity Gate)
# -----------------------------------------------------------------------------
Run-Step "1. Smoke & Sanity Testing" {
    Push-Location $BackendDir
    try {
        go test -count=1 -v ./tests/smoke/...
        if ($LASTEXITCODE -ne 0) { throw "Smoke tests failed with exit code $LASTEXITCODE" }
    } finally {
        Pop-Location
    }
}

# -----------------------------------------------------------------------------
# 2. UNIT TESTING
# -----------------------------------------------------------------------------
Run-Step "2. Unit Testing Suite" {
    Push-Location $BackendDir
    try {
        go test -count=1 -v ./pkg/crypto/...
        if ($LASTEXITCODE -ne 0) { throw "Unit tests failed with exit code $LASTEXITCODE" }
    } finally {
        Pop-Location
    }
}

# -----------------------------------------------------------------------------
# 3. INTEGRATION TESTING (POS, Stock, Ledger, PO)
# -----------------------------------------------------------------------------
Run-Step "3. Integration Testing Suite" {
    Push-Location $BackendDir
    try {
        go test -count=1 -v ./tests/integration/...
        if ($LASTEXITCODE -ne 0) { throw "Integration tests failed with exit code $LASTEXITCODE" }
    } finally {
        Pop-Location
    }
}

# -----------------------------------------------------------------------------
# 4. END-TO-END (E2E) TESTING (Dine-In, Takeaway, KDS)
# -----------------------------------------------------------------------------
Run-Step "4. End-to-End (E2E) Testing Suite" {
    Push-Location $BackendDir
    try {
        go test -count=1 -v ./tests/e2e/...
        if ($LASTEXITCODE -ne 0) { throw "E2E tests failed with exit code $LASTEXITCODE" }
    } finally {
        Pop-Location
    }
}

# -----------------------------------------------------------------------------
# 5. FRONTEND TYPECHECK & BUILD LINT
# -----------------------------------------------------------------------------
Run-Step "5. Frontend Quality Gate (Typecheck & Build)" {
    Push-Location $FrontendDir
    try {
        npx vue-tsc --noEmit
        if ($LASTEXITCODE -ne 0) { throw "Frontend TypeScript check failed with exit code $LASTEXITCODE" }
        npm run build
        if ($LASTEXITCODE -ne 0) { throw "Frontend production build failed with exit code $LASTEXITCODE" }
    } finally {
        Pop-Location
    }
}

# -----------------------------------------------------------------------------
# FINAL REPORT
# -----------------------------------------------------------------------------
Show-Summary

