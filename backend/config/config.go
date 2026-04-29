package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Constantes
const (
	ViewTaxas = "bdfaz.dbo.view_taxas_sejusp" // View completa com o nome do banco
)

// Config é uma struct que armazena todas as configurações da aplicação
type Config struct {
	ServerPort           string
	ModoAutenticacao     string
	MockCpfCnpj          string // Novo campo para o Mock
	DbConnectionString   string // Será montada no CarregarConfig
	EfazendaClientID     string
	EfazendaClientSecret string
	EfazendaRedirectURI  string
	OidcProviderURL      string // Novo campo para a URL base do provedor OIDC
	PortalApiURL         string // Novo campo para a URL da API do Portal e-Fazenda
}

// AppConfig é uma instância global (singleton) de nossas configurações.
var AppConfig *Config

// CarregarConfig é responsável por ler o arquivo .env
func CarregarConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado. Usando variáveis de ambiente do sistema.")
	}

	// 1. Carregar valores brutos
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	// ✅ NOVO: Leitura do Timeout como String
	dbTimeout := os.Getenv("DB_TIMEOUT_SECONDS")
	if dbTimeout == "" {
		dbTimeout = "60" // Valor padrão seguro
	}

	// 2. Montar a string de conexão completa (SQL Server)
	// Adiciona "connection timeout=DB_TIMEOUT_SECONDS" e keepAlive
	dbConnString := fmt.Sprintf(
		"server=%s;port=%s;database=%s;user id=%s;password=%s;connection timeout=%s;keepAlive=30",
		dbHost, dbPort, dbName, dbUser, dbPass, dbTimeout,
	)

	// 3. Popular a struct AppConfig
	AppConfig = &Config{
		ServerPort:           os.Getenv("SERVER_PORT"),
		ModoAutenticacao:     os.Getenv("MODO_AUTENTICACAO"),
		MockCpfCnpj:          os.Getenv("MOCK_CPF_CNPJ"), // Lê do .env
		DbConnectionString:   dbConnString,
		EfazendaClientID:     os.Getenv("EFAZENDA_CLIENT_ID"),
		EfazendaClientSecret: os.Getenv("EFAZENDA_CLIENT_SECRET"),
		EfazendaRedirectURI:  os.Getenv("EFAZENDA_REDIRECT_URI"),
		OidcProviderURL:      os.Getenv("OIDC_PROVIDER_URL"),
		PortalApiURL:         os.Getenv("PORTAL_API_URL"),
	}

	// Validação
	if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" {
		log.Fatal("FATAL: Variáveis de banco não definidas.")
	}
	if AppConfig.EfazendaClientID == "" {
		log.Println("AVISO: EFAZENDA_CLIENT_ID não definida no .env. O login via Portal e-Fazenda ficará inativo.")
	}

	if AppConfig.ModoAutenticacao == "mock" {
		log.Println("-----------------------------------------------------")
		log.Println("AVISO: Servidor rodando em MODO DE AUTENTICAÇÃO MOCK.")
		log.Printf("MOCK CPF/CNPJ: %s", AppConfig.MockCpfCnpj) // Loga o CPF que será usado
		log.Println("-----------------------------------------------------")
	}
}

