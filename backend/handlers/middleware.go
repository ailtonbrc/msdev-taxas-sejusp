package handlers

import (
	"taxas-sejusp/backend/config"
	"taxas-sejusp/backend/modelos"
	"taxas-sejusp/backend/servicos"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// AuthRequired é o nosso middleware de autenticação.
// Ele decide se usa a lógica "mock" ou "real" com base no .env.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.AppConfig
		var cpfUsuario string

		// 1. Tenta pegar o CPF do cookie de sessão (Funciona para Real e Mock após o login)
		cpfCookie, err := c.Cookie("user_cpf")
		
		if err == nil && cpfCookie != "" {
			cpfUsuario = cpfCookie
		} else if cfg.ModoAutenticacao == "mock" {
			// 2. Se for MOCK e não houver cookie, usa o fallback do .env ou config
			cpfUsuario = os.Getenv("MOCK_CPF_CNPJ")
			if cpfUsuario == "" {
				cpfUsuario = cfg.MockCpfCnpj
			}
			if cpfUsuario == "" {
				cpfUsuario = "admin" 
			}
		} else {
			// 3. Se for REAL e não houver cookie, barramos o acesso
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Usuário não autenticado"})
			return
		}

		// Validação de token (Opcional no Mock, obrigatória no Real)
		if cfg.ModoAutenticacao != "mock" {
			_, errToken := c.Cookie("auth_token")
			if errToken != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Sessão inválida"})
				return
			}
		}

		// ✅ Busca o usuário no banco de dados local para pegar as permissões
		user, err := servicos.GlobalUsuarios.BuscarPorCPF(cpfUsuario)
		if err != nil {
			log.Printf("[ERRO] Falha ao buscar usuário %s: %v", cpfUsuario, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"erro": "Erro interno ao validar permissões"})
			return
		}

		if user == nil {
			// ✅ Se o usuário não existe no DB, tentamos pegar o nome real que veio do e-Fazenda (via cookie)
			nomeExibicao, errNome := c.Cookie("user_name")
			if errNome != nil || nomeExibicao == "" {
				nomeExibicao = "Usuário não identificado"
			}

			log.Printf("[AVISO] Usuário autenticado mas não cadastrado: %s (%s). Acesso restrito.", nomeExibicao, cpfUsuario)
			guestUser := modelos.Usuario{
				Nome:         nomeExibicao,
				CPF:          cpfUsuario,
				Role:         "guest",
				Instituicoes: []string{}, // Nenhuma permissão
			}
			c.Set(ChaveUsuarioCpfCnpj, guestUser.CPF)
			c.Set(ChaveUsuarioInfo, guestUser)
			c.Next()
			return
		}

		// Injeta os dados no contexto para uso nos handlers
		c.Set(ChaveUsuarioCpfCnpj, user.CPF)
		c.Set(ChaveUsuarioInfo, *user)
		c.Next()
	}
}

