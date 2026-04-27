@echo off
title Taxas SEJUSP - Publicar Frontend
set "DIST_DIR=..\publish\frontend"

echo ==========================================
echo   Publicando Frontend (Vite/React)
echo ==========================================

echo [INFO] Limpando pasta de destino: %DIST_DIR%
if exist "%DIST_DIR%" rd /s /q "%DIST_DIR%"
mkdir "%DIST_DIR%"

echo [INFO] Executando build do Vite
call npm run build

if %ERRORLEVEL% NEQ 0 (
    echo [ERRO] Falha no build do Frontend!
    pause
    exit /b %ERRORLEVEL%
)

echo [INFO] Copiando arquivos gerados de dist para a pasta publish
xcopy /s /i /y dist "%DIST_DIR%"

if exist web.config (
    echo [INFO] Copiando web.config (IIS)
    copy web.config "%DIST_DIR%\web.config"
)

echo.
echo ============================================================
echo   FRONTEND PUBLICADO COM SUCESSO!
echo ============================================================
echo   Os arquivos ja estao prontos na pasta:
echo   "%DIST_DIR%"
echo.
echo   [DICA] Agora voce pode copiar esta pasta para o servidor.
echo ============================================================
echo.
pause
