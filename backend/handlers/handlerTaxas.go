package handlers

import (
	"database/sql"
	"log"
	"net/http"
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


		// Consulta diretamente do SQL Server
		dados, total, err := servicos.ListarDadosTaxas(dbSQLServer, params)
		if err != nil {
			log.Printf("[ERRO-HANDLER] Falha ao consultar banco: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha na consulta de dados", "detalhe": err.Error()})
			return
		}

		if total == 0 {
			c.JSON(http.StatusOK, modelos.RespostaDadosTaxas{
				TotalRegistros:  0,
				PaginaAtual:     params.Pagina,
				LimitePorPagina: params.Limite,
				Dados:           []modelos.TaxaSejusp{},
			})
			return
		}

		c.JSON(http.StatusOK, modelos.RespostaDadosTaxas{
			TotalRegistros:  total,
			PaginaAtual:     params.Pagina,
			LimitePorPagina: params.Limite,
			Dados:           dados,
		})
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


		buf, err := servicos.ExportarDadosExcel(dbSQLServer, params)
		if err != nil {
			log.Printf("[ERRO-EXPORT] Falha ao gerar Excel: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha na geração do Excel", "detalhe": err.Error()})
			return
		}

		nomeArquivo := "Taxas_SEJUSP.xlsx"
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+nomeArquivo)
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Expires", "0")
		c.Header("Cache-Control", "must-revalidate")
		c.Header("Pragma", "public")

		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
	}
}
