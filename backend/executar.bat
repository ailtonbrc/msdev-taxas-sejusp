@echo off
title Taxas SEJUSP - Backend (Go)

echo [INFO] Limpando o cache de build do Go...
go clean -cache

echo.
echo [INFO] Verificando diretorio bin...
if not exist "bin" mkdir bin

echo.
echo [INFO] Compilando o backend para bin\taxas-sejusp.exe...
go build -o bin\taxas-sejusp.exe .

if %ERRORLEVEL% NEQ 0 (
    echo [ERRO] Falha na compilacao!
    pause
    exit /b %ERRORLEVEL%
)

echo.
echo [INFO] Iniciando o servidor backend...
echo [INFO] Pressione CTRL+C para parar o servidor.
echo.

:: Executa o binário compilado
bin\taxas-sejusp.exe

echo.
echo [INFO] Servidor finalizado.
pause