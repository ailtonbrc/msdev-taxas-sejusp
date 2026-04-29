@echo off
title Taxas SEJUSP - Backend (DEV)
echo ==========================================
echo   Iniciando Servidor Backend (DEV)
echo ==========================================

:: Garante que estamos compilando para Windows (evita herdar GOOS=linux do deploy)
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0

:: Cria a pasta bin se nao existir
if not exist bin mkdir bin

echo [INFO] Compilando para Windows...
go build -o bin\backend_dev.exe main.go
if %errorlevel% neq 0 (
    echo [ERRO] Compilacao falhou!
    pause
    exit /b 1
)
echo [INFO] Executando servidor...
echo [INFO] Pressione CTRL+C para parar.
echo.

bin\backend_dev.exe

echo.
echo [AVISO] Servidor encerrado.
pause