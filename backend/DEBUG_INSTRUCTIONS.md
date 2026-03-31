# Instruções de Debug

O comando `go run .` não cria um arquivo `.exe` na pasta atual. Ele cria um executável temporário em uma pasta do sistema e o executa de lá. Por isso você não encontra o arquivo.

Para garantir que você está rodando a versão nova (com os logs):

1.  **Feche TODOS os terminais** do backend.
2.  Abra o **Gerenciador de Tarefas** (Ctrl+Shift+Esc).
3.  Procure por processos chamados `main.exe` ou `go.exe` e **finalize-os**.
4.  Abra um novo terminal na pasta `backend`.
5.  Execute `.\executar.bat`.

Se os logs ainda não aparecerem, tente rodar manualmente o comando de build para gerar um executável fixo:

```powershell
go build -o backend.exe .
.\backend.exe
```

Isso criará um arquivo `backend.exe` na pasta, que você pode rodar diretamente.
