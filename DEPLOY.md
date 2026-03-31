# Guia de Deploy - DIMP Municípios

Este guia descreve os passos para compilar e implantar a aplicação no servidor de produção.

**Servidor de Destino:** `172.20.20.25`
**Porta Backend:** `8000`
**Porta Frontend:** `80` (ou `8080` via servidor web)

## 1. Preparação do Frontend

O frontend deve ser compilado para arquivos estáticos (HTML/CSS/JS).

1.  Navegue até a pasta do frontend:
    ```powershell
    cd d:\desenvolvimento\dimp-municipios\frontend
    ```
2.  Instale as dependências (se necessário):
    ```powershell
    npm install
    ```
3.  Execute o build de produção:
    ```powershell
    npm run build
    ```
    *Isso criará uma pasta `dist` com os arquivos compilados.*
    *O arquivo `.env.production` já configura a URL da API para `http://172.20.20.25:8000/api`.*

## 2. Preparação do Backend

O backend deve ser compilado para um executável Windows.

1.  Navegue até a pasta do backend:
    ```powershell
    cd d:\desenvolvimento\dimp-municipios\backend
    ```
2.  Compile o executável:
    ```powershell
    go build -o dimp_server.exe main.go
    ```
3.  Verifique o arquivo `.env` de produção (no servidor):
    *   Certifique-se de que `SERVER_PORT=8000`.
    *   Certifique-se de que `MODO_AUTENTICACAO` está correto (não "mock" se for produção real).
    *   Configure as credenciais do banco de dados.

## 3. Estrutura no Servidor

Crie uma pasta no servidor (ex: `C:\Sistemas\DIMP`) e organize os arquivos da seguinte forma:

```
C:\Sistemas\DIMP\
├── backend\
│   ├── dimp_server.exe
│   ├── .env
│   └── config\ (se houver arquivos de config externos)
└── frontend\
    └── (conteúdo da pasta 'dist' gerada no passo 1)
```

## 4. Executando a Aplicação

### Backend
Execute o `dimp_server.exe`. Ele abrirá na porta 8000.
Recomenda-se usar um gerenciador de serviços (como NSSM) para manter o executável rodando como serviço do Windows.

### Frontend
Os arquivos da pasta `frontend` (o conteúdo de `dist`) devem ser servidos por um servidor web (IIS, Nginx, Apache) ou por um servidor estático simples.

Se for usar o **IIS**:
1.  Crie um novo Site.
2.  Aponte o caminho físico para a pasta `C:\Sistemas\DIMP\frontend`.
3.  Configure a porta para `80` (ou `8080`).
4.  **Importante:** Configure o "URL Rewrite" ou página de erro 404 para redirecionar para `index.html`, pois é uma SPA (Single Page Application).

Se quiser testar rapidamente sem IIS, pode usar o `serve`:
```powershell
npm install -g serve
serve -s C:\Sistemas\DIMP\frontend -l 8080
```

## 5. Verificação

1.  Acesse `http://172.20.20.25:8080` (ou a porta configurada para o frontend).
2.  O frontend deve carregar.
3.  Ao tentar fazer login ou carregar dados, ele deve chamar `http://172.20.20.25:8000/api/...`.
