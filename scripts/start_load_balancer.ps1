<#
.SYNOPSIS
    Memulai Cluster Backend Cafe ERP dengan Layer-7 Load Balancer pada Windows.
.DESCRIPTION
    Menjalankan:
    1. Worker Core (Port 8081) - Transaksi cepat POS, KDS, Auth, Master
    2. Worker Heavy (Port 8082) - Laporan intensif, Payroll Run, Opname
    3. Load Balancer (Port 8080) - Smart routing L7 + Health Checks
#>

$ErrorActionPreference = "Stop"

Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host "   CAFE ERP - CLUSTER & SMART LOAD BALANCER LAUNCHER (WINDOWS)   " -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

$BackendDir = Join-Path $PSScriptRoot "..\backend"
Set-Location $BackendDir

# Pastikan binary terkompilasi
Write-Host "[1/3] Memeriksa kompilasi binary server & load balancer..." -ForegroundColor Yellow
if (-not (Test-Path "bin\server.exe")) {
    Write-Host "  -> Mengompilasi bin\server.exe..." -ForegroundColor Gray
    go build -o bin\server.exe cmd\server\main.go
}

if (-not (Test-Path "bin\loadbalancer.exe")) {
    Write-Host "  -> Mengompilasi bin\loadbalancer.exe..." -ForegroundColor Gray
    go build -o bin\loadbalancer.exe cmd\loadbalancer\main.go
}

# Hentikan proses lama jika ada yang masih berjalan di port terkait
Write-Host "[2/3] Memastikan port 8080, 8081, 8082 bebas..." -ForegroundColor Yellow
$Ports = @(8080, 8081, 8082)
foreach ($Port in $Ports) {
    $conns = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue
    if ($conns) {
        foreach ($conn in $conns) {
            Write-Host "  -> Menutup proses lama PID $($conn.OwningProcess) di port $Port..." -ForegroundColor DarkGray
            Stop-Process -Id $conn.OwningProcess -Force -ErrorAction SilentlyContinue
        }
    }
}

Write-Host "[3/3] Menjalankan Node Cluster..." -ForegroundColor Green

# 1. Jalankan Core Worker (Port 8081)
$coreEnv = @{
    PORT = "8081"
}
$coreProcess = Start-Process -FilePath "bin\server.exe" -Environment $coreEnv -PassThru -WindowStyle Hidden
Write-Host "  ✅ Core Worker 1 berjalan di port 8081 (PID: $($coreProcess.Id))" -ForegroundColor Green

# 2. Jalankan Heavy Worker (Port 8082)
$heavyEnv = @{
    PORT = "8082"
    HEAVY_MAX_CONCURRENCY = "8"
}
$heavyProcess = Start-Process -FilePath "bin\server.exe" -Environment $heavyEnv -PassThru -WindowStyle Hidden
Write-Host "  ✅ Heavy Worker berjalan di port 8082 (PID: $($heavyProcess.Id))" -ForegroundColor Green

# Tunggu backend siap
Start-Sleep -Seconds 2

# 3. Jalankan Smart Load Balancer (Port 8080)
$lbEnv = @{
    LB_PORT = "8080"
    CORE_BACKENDS = "http://127.0.0.1:8081"
    HEAVY_BACKENDS = "http://127.0.0.1:8082"
}

Write-Host "`n🚀 Menjalankan Load Balancer di http://localhost:8080..." -ForegroundColor Cyan
Write-Host "🎯 POS/KDS/Auth/Master  --> http://127.0.0.1:8081" -ForegroundColor Gray
Write-Host "⚡ Reports/Payroll/Opname --> http://127.0.0.1:8082" -ForegroundColor Gray
Write-Host "📊 Pantau Metrik Status   --> http://localhost:8080/lb-status" -ForegroundColor Cyan
Write-Host "Tekan CTRL+C untuk menghentikan seluruh cluster.`n" -ForegroundColor Yellow

# Simpan PID untuk proses shutdown
$clusterState = @{
    CorePid = $coreProcess.Id
    HeavyPid = $heavyProcess.Id
}
$clusterState | ConvertTo-Json | Set-Content (Join-Path $BackendDir "bin\cluster_pids.json")

try {
    # Jalankan Load Balancer di foreground agar log terlihat
    $env:LB_PORT = "8080"
    $env:CORE_BACKENDS = "http://127.0.0.1:8081"
    $env:HEAVY_BACKENDS = "http://127.0.0.1:8082"
    & "bin\loadbalancer.exe"
}
finally {
    Write-Host "`n🛑 Menghentikan seluruh worker cluster..." -ForegroundColor Yellow
    Stop-Process -Id $coreProcess.Id -Force -ErrorAction SilentlyContinue
    Stop-Process -Id $heavyProcess.Id -Force -ErrorAction SilentlyContinue
    Remove-Item (Join-Path $BackendDir "bin\cluster_pids.json") -ErrorAction SilentlyContinue
    Write-Host "✅ Cluster berhasil dimatikan dengan aman." -ForegroundColor Green
}
