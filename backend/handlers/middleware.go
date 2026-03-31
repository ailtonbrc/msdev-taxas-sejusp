package handlers

import (
	"taxas-sejusp/backend/config"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// AuthRequired é o nosso middleware de autenticação.
// Ele decide se usa a lógica "mock" ou "real" com base no .env.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.AppConfig

		if cfg.ModoAutenticacao == "mock" {
			// --- MOCK ---
			_ = godotenv.Load()
			cpfCnpjMockado := os.Getenv("MOCK_CPF_CNPJ")
			if cpfCnpjMockado == "" {
				cpfCnpjMockado = cfg.MockCpfCnpj
			}
			if cpfCnpjMockado == "" {
				cpfCnpjMockado = "36828238168"
				log.Println("AVISO: MOCK_CPF_CNPJ não definido no .env, usando fallback.")
			}

			// Injeta no contexto
			c.Set(ChaveUsuarioCpfCnpj, cpfCnpjMockado)
			c.Next()
			return
		}

		// --- MODO REAL ---
		// Verifica se existe o cookie com o CPF/CNPJ (definido no Callback)
		// Em produção, deveríamos validar o token JWT (cookie auth_token) a cada request
		// ou verificar uma sessão em banco/memória.
		// Para este MVP, verificamos a presença do cookie do usuário.
		cpfCookie, err := c.Cookie("user_cpf")
		if err != nil || cpfCookie == "" {
			// Não autenticado
			// Se for chamada AJAX (JSON), retorna 401.
			// Se fosse navegação normal, poderia redirecionar.
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Usuário não autenticado"})
			return
		}

		// Opcional: Validar se o token também existe
		_, errToken := c.Cookie("auth_token")
		if errToken != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Sessão inválida"})
			return
		}

		// Usuário autenticado, injeta no contexto
		c.Set(ChaveUsuarioCpfCnpj, cpfCookie)
		c.Next()
	}
}

