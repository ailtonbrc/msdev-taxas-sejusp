package main

import (
	"database/sql"
	"taxas-sejusp/backend/config"
	"taxas-sejusp/backend/handlers"
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
		log.Fatal("Erro fatal ao pingar o SQL Server: ", err.Error())
	}
	log.Println("Conexão com SQL Server estabelecida com sucesso.")

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
			authGroup.GET("/user", handlers.VerificarUsuario)
		}

		taxasGroup := api.Group("/taxas")
		// taxasGroup.Use(handlers.AuthRequired()) // Desabilitado temporariamente para produção rápida
		{
			// Rota de Dados (JSON)
			taxasGroup.GET("/dados", handlers.ListarDadosTaxas(db))
			// Rota de Exportação (Excel)
			taxasGroup.GET("/exportar", handlers.ExportarDadosExcel(db))
		}
	}

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)

	// ✅ CORREÇÃO: LOG EXPLICITO DA PORTA
	log.Printf("Servidor escutando em http://localhost%s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Erro server: ", err.Error())
	}
}

