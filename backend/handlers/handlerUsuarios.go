package handlers

import (
	"net/http"
	"strconv"
	"taxas-sejusp/backend/modelos"
	"taxas-sejusp/backend/servicos"

	"github.com/gin-gonic/gin"
)

func ListarUsuarios(c *gin.Context) {
	users, err := servicos.GlobalUsuarios.ListarTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao listar usuários"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func SalvarUsuario(c *gin.Context) {
	var u modelos.Usuario
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	if err := servicos.GlobalUsuarios.Salvar(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao salvar usuário"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensagem": "Usuário salvo com sucesso"})
}

func ExcluirUsuario(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := servicos.GlobalUsuarios.Excluir(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao excluir usuário"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensagem": "Usuário excluído com sucesso"})
}

// MockLoginSimulado permite entrar no sistema apenas informando o CPF
func MockLoginSimulado(c *gin.Context) {
	var req struct {
		CPF string `json:"cpf" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "CPF é obrigatório"})
		return
	}

	// Verificar se o usuário existe no nosso cadastro
	user, err := servicos.GlobalUsuarios.BuscarPorCPF(req.CPF)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao validar usuário"})
		return
	}
	if user == nil {
		// ✅ Se não estiver cadastrado, logamos como convidado (guest)
		guestUser := &modelos.Usuario{
			Nome:         "Usuário convidado",
			CPF:          req.CPF,
			Role:         "guest",
			Instituicoes: []string{},
		}
		
		c.SetCookie("auth_token", "mock-session-" + req.CPF, 3600*8, "/", "", false, true)
		c.SetCookie("user_cpf", req.CPF, 3600*8, "/", "", false, false)
		
		c.JSON(http.StatusOK, gin.H{"mensagem": "Login simulado como convidado", "usuario": guestUser})
		return
	}

	// Criar os cookies de sessão idênticos ao e-Fazenda
	c.SetCookie("auth_token", "mock-session-" + req.CPF, 3600*8, "/", "", false, true)
	c.SetCookie("user_cpf", req.CPF, 3600*8, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"mensagem": "Login simulado com sucesso", "usuario": user})
}
