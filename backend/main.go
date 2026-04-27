package main

import (
	"database/sql"
	"taxas-sejusp/backend/config"
	"taxas-sejusp/backend/handlers"
	"taxas-sejusp/backend/servicos"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

var db *sql.DB

func main() {
	fmt.Println("---------------------------------------------------------")
	fmt.Println("   VERSÃO 1.0 - TAXAS SEJUSP - CONSULTA DIRETA           ")
	fmt.Println("   BUILD: 2025-12-04")
	fmt.Println("---------------------------------------------------------")

	config.CarregarConfig()
	cfg := config.AppConfig

	var err error

	db, err = sql.Open("sqlserver", cfg.DbConnectionString)
	if err != nil {
		log.Fatal("Erro conexao: ", err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Printf("AVISO: SQL Server inacessível no momento (%v). O servidor continuará em execução e tentará reconectar em background.", err)
	}
	// --- Inicialização do Cache de Dados (Performance) ---
	log.Println("Carregando banco de dados para a memória (cache)...")
	if err := servicos.GlobalCache.Carregar(db); err != nil {
		log.Printf("Aviso: Falha ao carregar cache inicial: %v", err)
	}

	// ✅ Inicializa o banco de dados de usuários e permissões
	if err := servicos.IniciarServicoUsuario(); err != nil {
		log.Fatalf("Falha crítica ao iniciar serviço de usuários: %v", err)
	}
	
	// ✅ Ativa o agendador de atualização (06h, 12h, 16h)
	servicos.IniciarSincronizacaoAutomatica(db)

	// --- Desliga o modo Debug do Gin ---
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Configuração de CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api := router.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.GET("/login", handlers.IniciarLogin)
			authGroup.GET("/callback", handlers.CallbackLogin)
			authGroup.GET("/user", handlers.AuthRequired(), handlers.VerificarUsuario) // ✅ Middleware adicionado
			authGroup.GET("/logout", handlers.SairLogin)
			authGroup.POST("/mock-login", handlers.MockLoginSimulado)
		}

		// Gestão de Usuários
		adminGroup := api.Group("/usuarios")
		adminGroup.Use(handlers.AuthRequired())
		{
			adminGroup.GET("/", handlers.ListarUsuarios)
			adminGroup.POST("/", handlers.SalvarUsuario)
			adminGroup.DELETE("/:id", handlers.ExcluirUsuario)
		}

		taxasGroup := api.Group("/taxas")
		taxasGroup.Use(handlers.AuthRequired()) // Reativado para proteger os dados
		{
			taxasGroup.GET("/dados", handlers.ListarDadosTaxas(db))
			taxasGroup.GET("/exportar", handlers.ExportarDadosExcel(db))
			taxasGroup.GET("/exportar-csv", handlers.ExportarDadosCSV(db))
			taxasGroup.GET("/opcoes-filtros", handlers.ObterOpcoesFiltros)
			taxasGroup.POST("/refresh", handlers.RecarregarCache(db))
		}
	}

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)

	// ✅ CORREÇÃO: LOG EXPLICITO DA PORTA
	log.Printf("Servidor escutando em http://localhost%s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Erro server: ", err.Error())
	}
}

