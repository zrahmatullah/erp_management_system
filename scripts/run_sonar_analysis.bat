@echo off
setlocal enabledelayedexpansion

echo ==============================================================================
echo [SonarQube Analysis Runner] - Cafe ERP Monorepo
echo ==============================================================================

set "PROJECT_ROOT=%~dp0.."
cd /d "%PROJECT_ROOT%"

echo.
echo [Step 1/3] Running Go Unit Tests ^& Generating Coverage Profile...
cd /d "%PROJECT_ROOT%\backend"
go test -v "-coverprofile=coverage.out" "-coverpkg=./..." ./tests/unit
if %errorlevel% neq 0 (
    echo [ERROR] Unit testing failed! Please resolve test failures before running SonarQube analysis.
    exit /b %errorlevel%
)
echo [OK] Coverage report generated at backend\coverage.out

echo.
echo [Step 2/3] Checking SonarQube Server Availability...
curl.exe -s http://localhost:9000/api/system/status | findstr /i "UP" >nul
if %errorlevel% neq 0 (
    echo [WARNING] SonarQube Server at http://localhost:9000 is not UP yet.
    echo Please make sure SonarQube is started using:
    echo   scripts\start_sonarqube_server.bat
    echo.
)

echo.
echo [Step 3/3] Executing SonarScanner CLI...
cd /d "%PROJECT_ROOT%"
if exist "D:\sonarqube\scanner\bin\sonar-scanner.bat" (
    call "D:\sonarqube\scanner\bin\sonar-scanner.bat" -Dsonar.host.url=http://localhost:9000 -Dsonar.token=squ_cafe_erp_system_sonar_token_2026
) else (
    call sonar-scanner -Dsonar.host.url=http://localhost:9000 -Dsonar.token=squ_cafe_erp_system_sonar_token_2026
)

echo.
echo ==============================================================================
echo Analysis Completed!
echo Review your project metrics ^& Quality Gate at:
echo http://localhost:9000/dashboard?id=cafe-erp-system
echo ==============================================================================
