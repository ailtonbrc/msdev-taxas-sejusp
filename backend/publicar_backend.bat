@echo off
title Taxas SEJUSP - Publicar Backend
set DIST_DIR=..\publish\backend

echo ==========================================
echo   Publicando Servidor Backend
echo ==========================================
echo [INFO] Limpando pasta de destino: %DIST_DIR%
if exist %DIST_DIR% rd /s /q %DIST_DIR%
mkdir %DIST_DIR%

echo [INFO] Compilando executavel (Windows)...
go build -o %DIST_DIR%\taxas-sejusp.exe main.go

if %ERRORLEVEL% NEQ 0 (
    echo [ERRO] Falha na compilacao do Go!
    pause
    exit /b %ERRORLEVEL%
)

echo [INFO] Copiando arquivos de configuracao...
if exist .env (
    copy .env %DIST_DIR%\.env.example
    echo [AVISO] O .env foi copiado como .env.example para seguranca.
)

echo.
echo ==========================================
echo   [SUCESSO] Backend publicado em:
echo   %DIST_DIR%
echo ==========================================
echo.
pause
