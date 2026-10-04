<#
.SYNOPSIS
    Menghentikan seluruh proses Cluster Backend dan Load Balancer Cafe ERP.
#>

Write-Host "Menghentikan seluruh service Cafe ERP (Port 8080, 8081, 8082)..." -ForegroundColor Yellow

$Ports = @(8080, 8081, 8082)
foreach ($Port in $Ports) {
    $conns = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue
    if ($conns) {
        foreach ($conn in $conns) {
            Write-Host "  -> Menghentikan PID $($conn.OwningProcess) di port $Port..." -ForegroundColor Green
            Stop-Process -Id $conn.OwningProcess -Force -ErrorAction SilentlyContinue
        }
    }
}

# Bersihkan file cluster_pids.json jika ada
$pidFile = Join-Path $PSScriptRoot "..\backend\bin\cluster_pids.json"
if (Test-Path $pidFile) {
    Remove-Item $pidFile -Force -ErrorAction SilentlyContinue
}

Write-Host "Seluruh worker cluster telah dihentikan." -ForegroundColor Green
