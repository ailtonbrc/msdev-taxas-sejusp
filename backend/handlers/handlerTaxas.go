package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"taxas-sejusp/backend/modelos"
	"taxas-sejusp/backend/servicos"

	"github.com/gin-gonic/gin"
)

func ListarDadosTaxas(dbSQLServer *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params modelos.ParametrosBusca
		if err := c.ShouldBindQuery(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"erro": "Parâmetros inválidos"})
			return
		}
		params.Normalizar()

		// ✅ Recupera o usuário injetado pelo middleware para aplicar RBAC
		userObj, _ := c.Get(ChaveUsuarioInfo)
		user := userObj.(modelos.Usuario)
		isAdmin := user.Role == "admin"

		// Busca otimizada usando o cache em memória
		dados, total := servicos.GlobalCache.Listar(params, isAdmin, user.Instituicoes)

		// Obter status atual do cache
		carregado, _, erroSync := servicos.GlobalCache.ObterStatus()

		c.JSON(http.StatusOK, modelos.RespostaDadosTaxas{
			TotalRegistros:  total,
			PaginaAtual:     params.Pagina,
			LimitePorPagina: params.Limite,
			Dados:           dados,
			Sincronizado:    carregado,
			ErroSync:        erroSync,
		})
	}
}

// ObterOpcoesFiltros retorna as listas de valores únicos filtradas pelo contexto atual e permissões
func ObterOpcoesFiltros(c *gin.Context) {
	var params modelos.ParametrosBusca
	if err := c.ShouldBindQuery(&params); err != nil {
		params = modelos.ParametrosBusca{}
	}
	params.Normalizar()

	// ✅ Recupera o usuário injetado pelo middleware para aplicar RBAC nas opções
	userObj, _ := c.Get(ChaveUsuarioInfo)
	user := userObj.(modelos.Usuario)
	isAdmin := user.Role == "admin"

	opcoes := servicos.GlobalCache.ObterOpcoes(params, isAdmin, user.Instituicoes)
	
	carregado, _, erroSync := servicos.GlobalCache.ObterStatus()
	opcoes.Sincronizado = carregado
	opcoes.ErroSync = erroSync

	c.JSON(http.StatusOK, opcoes)
}

// RecarregarCache força a atualização do cache a partir do banco de dados (Apenas Admin)
func RecarregarCache(dbSQLServer *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userObj, _ := c.Get(ChaveUsuarioInfo)
		user := userObj.(modelos.Usuario)
		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"erro": "Apenas administradores podem forçar a atualização do cache"})
			return
		}

		err := servicos.GlobalCache.Carregar(dbSQLServer)
		if err != nil {
			log.Printf("[ERRO CACHE] Falha ao recarregar: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao sincronizar dados do banco", "detalhe": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"mensagem": "Dados sincronizados com sucesso"})
	}
}

func ExportarDadosExcel(dbSQLServer *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params modelos.ParametrosBusca
		if err := c.ShouldBindQuery(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"erro": "Parâmetros inválidos"})
			return
		}
		params.Normalizar()

		userObj, _ := c.Get(ChaveUsuarioInfo)
		user := userObj.(modelos.Usuario)
		isAdmin := user.Role == "admin"
		
		tmpFile, err := os.CreateTemp("", "taxas_sejusp_*.xlsx")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao criar arquivo no servidor"})
			return
		}
		tmpPath := tmpFile.Name()
		tmpFile.Close() 
		defer os.Remove(tmpPath)

		err = servicos.GlobalCache.ExportarExcel(params, tmpPath, isAdmin, user.Instituicoes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha na geração do arquivo Excel"})
			return
		}

		c.FileAttachment(tmpPath, "Taxas_SEJUSP.xlsx")
	}
}

func ExportarDadosCSV(dbSQLServer *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params modelos.ParametrosBusca
		if err := c.ShouldBindQuery(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"erro": "Parâmetros inválidos"})
			return
		}
		params.Normalizar()

		userObj, _ := c.Get(ChaveUsuarioInfo)
		user := userObj.(modelos.Usuario)
		isAdmin := user.Role == "admin"
		
		tmpFile, err := os.CreateTemp("", "taxas_sejusp_*.csv")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao criar arquivo temporário"})
			return
		}
		tmpPath := tmpFile.Name()
		tmpFile.Close()
		defer os.Remove(tmpPath)

		err = servicos.GlobalCache.ExportarCSV(params, tmpPath, isAdmin, user.Instituicoes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha na geração do arquivo CSV"})
			return
		}

		c.FileAttachment(tmpPath, "Taxas_SEJUSP.csv")
	}
}
