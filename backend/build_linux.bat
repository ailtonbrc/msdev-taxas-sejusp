@echo off
echo [INFO] Configurando variaveis de ambiente para Linux (amd64)...
set GOOS=linux
set GOARCH=amd64

echo [INFO] Compilando backend...
if not exist "bin" mkdir bin
go build -o bin/taxas_backend .

if %ERRORLEVEL% EQU 0 (
    echo [SUCESSO] Binario criado em: bin\taxas_backend
    echo [INFO] Este arquivo esta pronto para ser enviado para o servidor Linux.
) else (
    echo [ERRO] Falha na compilacao.
)
pause
