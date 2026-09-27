@echo off
setlocal

echo ==============================================================================
echo [SonarQube Server Launcher] - Cafe ERP Monorepo
echo ==============================================================================

set "JAVA_HOME=D:\sonarqube\jdk-21"
set "SONAR_JAVA_PATH=D:\sonarqube\jdk-21\bin\java.exe"
set "PATH=D:\sonarqube\jdk-21\bin;%PATH%"

if not exist "%JAVA_HOME%\bin\java.exe" (
    echo [ERROR] JDK 21 not found at %JAVA_HOME%
    pause
    exit /b 1
)

echo Starting SonarQube Server with Java 21...
echo Host URL: http://localhost:9000
echo.
cd /d "D:\sonarqube-26.9.0.129388\bin\windows-x86-64"
call "D:\sonarqube-26.9.0.129388\bin\windows-x86-64\StartSonar.bat"
