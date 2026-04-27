#!/bin/bash

# Script de Configuracao Automatica - Taxas SEJUSP
# Para ser executado no servidor AlmaLinux 9.6

echo "===================================================="
echo "    Configurando Taxas SEJUSP no Servidor"
echo "===================================================="

# 1. Definir permissoes de execucao para o binario
echo "[1/4] Ajustando permissoes do binario..."
chmod +x /var/www/taxasejusp/backend/taxas_backend

# 2. Configurar o Systemd (Servico)
echo "[2/4] Configurando o servico Systemd..."
cp /var/www/taxasejusp/backend/taxas-sejusp.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable taxas-sejusp
systemctl restart taxas-sejusp

# 3. Configurar o Nginx
echo "[3/4] Configurando o Nginx..."
# Copia o arquivo para o path do Nginx (AlmaLinux costuma usar /etc/nginx/conf.d/)
cp /var/www/taxasejusp/taxasejusp.conf /etc/nginx/conf.d/

# Valida a sintaxe do Nginx
nginx -t
if [ $? -eq 0 ]; then
    echo "Sintaxe do Nginx OK, recarregando..."
    systemctl reload nginx
else
    echo "ERRO na sintaxe do Nginx. Verifique os logs."
fi

# 4. Status final
echo "[4/4] Verificando status do servico..."
systemctl is-active --quiet taxas-sejusp && echo "Backend: ONLINE" || echo "Backend: FALHA"
systemctl is-active --quiet nginx && echo "Nginx: ONLINE" || echo "Nginx: FALHA"

echo "===================================================="
echo "    Deploy Concluido com Sucesso!"
echo "    URL: http://taxasejusp.fazenda.ms.gov.br/"
echo "===================================================="
