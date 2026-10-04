@echo off
title Cafe ERP Cluster ^& Load Balancer Launcher
echo =================================================================
echo   CAFE ERP - CLUSTER ^& SMART LOAD BALANCER LAUNCHER (WINDOWS)
echo =================================================================
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0start_load_balancer.ps1"
pause
