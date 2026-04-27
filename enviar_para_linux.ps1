# Script de Deploy Otimizado - Taxas SEJUSP -> Linux (AlmaLinux)
# Compacta tudo para evitar multiplas solicitacoes de senha.

$ServerIP = "172.20.20.26"
$ServerUser = "ajbranco@FAZENDA.MS"
$RemotePath = "/var/www/taxasejusp"

echo "===================================================="
echo "    Iniciando Build e Deploy para $ServerIP"
echo "===================================================="

# 1. Build do Backend (Linux)
echo "[1/5] Compilando Backend para Linux..."
cd backend
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o bin/taxas_backend .
if ($LASTEXITCODE -ne 0) { echo "ERRO: Falha no build do Go"; exit }
cd ..

# 2. Build do Frontend
echo "[2/5] Compilando Frontend (React/Vite)..."
cd frontend
npm run build
if ($LASTEXITCODE -ne 0) { echo "ERRO: Falha no build do Vite"; exit }
cd ..

# 3. Preparar Pacote
echo "[3/5] Empacotando arquivos para envio unico..."
if (Test-Path "deploy_package") { Remove-Item -Recurse -Force "deploy_package" }
New-Item -ItemType Directory -Path "deploy_package/backend" -Force
New-Item -ItemType Directory -Path "deploy_package/frontend" -Force

Copy-Item "backend/bin/taxas_backend" "deploy_package/backend/"
Copy-Item "backend/.env" "deploy_package/backend/"
Copy-Item "backend/taxas-sejusp.service" "deploy_package/backend/"
Copy-Item "backend/deploy_linux.sh" "deploy_package/backend/"
Copy-Item -Recurse "frontend/dist/*" "deploy_package/frontend/"
Copy-Item "taxasejusp.conf" "deploy_package/"

# Criar arquivo comprimido (tar.gz) - nativo no Windows 10+ e Linux
tar -czf deploy.tar.gz -C deploy_package .

# 3. Organizar pacotes e enviar via SCP
echo "[3/5] Enviando pacote compactado para /tmp..."
# Usar --% para evitar que o PowerShell tente interpretar o @ do dominio nos comandos externos
scp deploy.tar.gz "${ServerUser}@${ServerIP}:/tmp/"
if ($LASTEXITCODE -ne 0) { echo "ERRO: Falha no envio via SCP"; exit }

# 5. Extrair e Ativar (Uma solicitacao de senha para o Sudo aqui)
echo "[5/5] Extraindo e Ativando no servidor remoto..."
# Comando remoto: Criar pasta, extrair e rodar deploy script
$RemoteCmd = "sudo mkdir -p $RemotePath && sudo tar -xzf /tmp/deploy.tar.gz -C $RemotePath && sudo bash ${RemotePath}/backend/deploy_linux.sh"
ssh -t "${ServerUser}@${ServerIP}" $RemoteCmd

# Limpeza local
Remove-Item "deploy.tar.gz"
Remove-Item -Recurse "deploy_package"

echo "===================================================="
echo "    DEPLOY FINALIZADO COM SUCESSO!"
echo "    URL: http://taxasejusp.fazenda.ms.gov.br/"
echo "===================================================="
