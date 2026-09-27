@echo off
setlocal enabledelayedexpansion

echo ==============================================================================
echo [SonarQube Analysis Runner] - Cafe ERP Backend ^& Frontend
echo ==============================================================================

set "PROJECT_ROOT=%~dp0.."
cd /d "%PROJECT_ROOT%"

echo.
echo [Step 1/4] Running Go Unit Tests ^& Generating Coverage Profile...
cd /d "%PROJECT_ROOT%\backend"
set "TARGET_PKGS=cafe-erp-system/backend/internal/delivery/http/middleware,cafe-erp-system/backend/internal/usecase/auth,cafe-erp-system/backend/internal/usecase/pos,cafe-erp-system/backend/internal/usecase/inventory,cafe-erp-system/backend/pkg/...,cafe-erp-system/backend/internal/config"
go test -v "-coverprofile=coverage.out" "-coverpkg=%TARGET_PKGS%" ./tests/unit
if %errorlevel% neq 0 (
    echo [ERROR] Backend unit testing failed! Please resolve test failures.
    exit /b %errorlevel%
)
go tool cover "-func=coverage.out" | findstr /i "total:"
echo [OK] Coverage report generated at backend\coverage.out

echo.
echo [Step 2/4] Running Frontend Unit Tests (Vitest) ^& Generating Coverage...
cd /d "%PROJECT_ROOT%\frontend"
call npm run test:coverage
if %errorlevel% neq 0 (
    echo [ERROR] Frontend unit testing failed!
    exit /b %errorlevel%
)
echo [OK] Frontend coverage report generated at frontend\coverage\lcov.info

echo.
echo [Step 3/4] Checking SonarQube Server Availability...
curl.exe -s http://localhost:9000/api/system/status | findstr /i "UP" >nul
if %errorlevel% neq 0 (
    echo [WARNING] SonarQube Server at http://localhost:9000 is not UP yet.
    echo Please make sure SonarQube is started using:
    echo   scripts\start_sonarqube_server.bat
    echo.
)

echo.
echo [Step 4/4] Executing SonarScanner CLI for Frontend...
cd /d "%PROJECT_ROOT%\frontend"
if exist "D:\sonarqube\scanner\bin\sonar-scanner.bat" (
    call "D:\sonarqube\scanner\bin\sonar-scanner.bat" -Dsonar.host.url=http://localhost:9000 -Dsonar.token=squ_cafe_erp_system_sonar_token_2026
) else (
    call sonar-scanner -Dsonar.host.url=http://localhost:9000 -Dsonar.token=squ_cafe_erp_system_sonar_token_2026
)

echo.
echo ==============================================================================
echo Analysis Completed!
echo Review Frontend metrics at: http://localhost:9000/dashboard?id=cafe-erp-frontend
echo Review Backend metrics at:  http://localhost:9000/dashboard?id=cafe-erp-backend
echo ==============================================================================
