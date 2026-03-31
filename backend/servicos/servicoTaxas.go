package servicos

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
	"taxas-sejusp/backend/config"
	"taxas-sejusp/backend/modelos"
)

// ListarDadosTaxas busca os dados diretamente da View no SQL Server
func ListarDadosTaxas(db *sql.DB, params modelos.ParametrosBusca) ([]modelos.TaxaSejusp, int64, error) {
	where := []string{"1=1"}
	var args []interface{}
	paramCounter := 1

	// Filtro de Período (Mês/Ano)
	if params.PeriodoInicio != nil && *params.PeriodoInicio != "" {
		where = append(where, fmt.Sprintf("[Referência] = @p%d", paramCounter))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramCounter), *params.PeriodoInicio))
		paramCounter++
	}

	// --- Busca Global Removida para Ganho de Performance ---
	// (O sistema agora utiliza filtros individuais por coluna)


	// Filtros Específicos por Coluna
	if params.ItemSubItem != "" {
		where = append(where, fmt.Sprintf("CONVERT(VARCHAR(MAX), [Item.SubItem]) LIKE @f%d", paramCounter))
		args = append(args, sql.Named(fmt.Sprintf("f%d", paramCounter), "%"+params.ItemSubItem+"%"))
		paramCounter++
	}
	if params.Descricao != "" {
		where = append(where, fmt.Sprintf("[Descricao] LIKE @f%d", paramCounter))
		args = append(args, sql.Named(fmt.Sprintf("f%d", paramCounter), "%"+params.Descricao+"%"))
		paramCounter++
	}
	if params.Tributo != "" {
		where = append(where, fmt.Sprintf("CONVERT(VARCHAR(MAX), [Tributo]) LIKE @f%d", paramCounter))
		args = append(args, sql.Named(fmt.Sprintf("f%d", paramCounter), "%"+params.Tributo+"%"))
		paramCounter++
	}
	if params.ReferenciaFlt != "" {
		flt := params.ReferenciaFlt
		// Se o usuário digitar no formato MM/YYYY (ex: 01/2021), convertemos para YYYYMM (202101)
		if strings.Contains(flt, "/") {
			partes := strings.Split(flt, "/")
			if len(partes) == 2 && len(partes[0]) == 2 && len(partes[1]) == 4 {
				flt = partes[1] + partes[0]
			}
		}
		where = append(where, fmt.Sprintf("CONVERT(VARCHAR(MAX), [Referência]) LIKE @f%d", paramCounter))
		args = append(args, sql.Named(fmt.Sprintf("f%d", paramCounter), "%"+flt+"%"))
		paramCounter++
	}

	whereClause := strings.Join(where, " AND ")

	queryCount := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", config.ViewTaxas, whereClause)

	var total int64
	err := db.QueryRow(queryCount, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("erro count taxas: %w", err)
	}

	// Consulta com conversão explícita para evitar o erro "varchar to numeric"
	queryData := fmt.Sprintf(`
		SELECT 
			ISNULL(CONVERT(VARCHAR(MAX), [Item.SubItem]), '') as ItemSubItem, 
			ISNULL([Descricao], '') as Descricao, 
			ISNULL(CONVERT(VARCHAR(MAX), [Tributo]), '') as Tributo, 
			ISNULL(CONVERT(VARCHAR(MAX), [DataPagamento], 120), '') as DataPagamento, 
			ISNULL(CONVERT(VARCHAR(MAX), [Referência]), '') as Referencia, 
			ISNULL(CONVERT(VARCHAR(MAX), [Municipio]), 'NÃO INFORMADO') as Municipio, 
			ISNULL(TRY_CAST([Valor Principal] AS FLOAT), 0) as ValorPrincipal, 
			ISNULL(TRY_CAST([Quantidade de UFERMS] AS FLOAT), 0) as QuantidadeUferms, 
			ISNULL(TRY_CAST([Valor Total] AS FLOAT), 0) as ValorTotal
		FROM %s
		WHERE %s
		ORDER BY (SELECT NULL)
		OFFSET %d ROWS FETCH NEXT %d ROWS ONLY
	`, config.ViewTaxas, whereClause, params.CalcularOffset(), params.Limite)


	rows, err := db.Query(queryData, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("erro query taxas: %w", err)
	}
	defer rows.Close()

	var dados []modelos.TaxaSejusp
	for rows.Next() {
		var d modelos.TaxaSejusp
		err := rows.Scan(
			&d.ItemSubItem,
			&d.Descricao,
			&d.Tributo,
			&d.DataPagamento,
			&d.Referencia,
			&d.Municipio,
			&d.ValorPrincipal,
			&d.QuantidadeUferms,
			&d.ValorTotal,
		)
		if err != nil {
			log.Printf("[ERRO SCAN] Erro ao ler linha de taxas: %v", err)
			continue
		}
		dados = append(dados, d)
	}

	// NOVO: Verifica se houve erro geral durante a iteração das linhas
	if err = rows.Err(); err != nil {
		log.Printf("[ERRO ROWS] Erro durante a iteração das linhas: %v", err)
		return nil, 0, fmt.Errorf("erro durante leitura: %w", err)
	}

	return dados, total, nil
}

// ExportarDadosExcel gera um buffer contendo o arquivo Excel com todos os dados filtrados
func ExportarDadosExcel(db *sql.DB, params modelos.ParametrosBusca) (*bytes.Buffer, error) {
	where := []string{"1=1"}
	var args []interface{}
	paramCounter := 1

	// Filtro de Período (Mês/Ano) - Mesma lógica da Listagem
	if params.PeriodoInicio != nil && *params.PeriodoInicio != "" {
		where = append(where, fmt.Sprintf("[Referência] = @p%d", paramCounter))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramCounter), *params.PeriodoInicio))
		paramCounter++
	}

	// Filtros Específicos por Coluna
	if params.ItemSubItem != "" {
		where = append(where, fmt.Sprintf("CONVERT(VARCHAR(MAX), [Item.SubItem]) LIKE @f%d", paramCounter))
		args = append(args, sql.Named(fmt.Sprintf("f%d", paramCounter), "%"+params.ItemSubItem+"%"))
		paramCounter++
	}
	if params.Descricao != "" {
		where = append(where, fmt.Sprintf("[Descricao] LIKE @f%d", paramCounter))
		args = append(args, sql.Named(fmt.Sprintf("f%d", paramCounter), "%"+params.Descricao+"%"))
		paramCounter++
	}
	if params.Tributo != "" {
		where = append(where, fmt.Sprintf("CONVERT(VARCHAR(MAX), [Tributo]) LIKE @f%d", paramCounter))
		args = append(args, sql.Named(fmt.Sprintf("f%d", paramCounter), "%"+params.Tributo+"%"))
		paramCounter++
	}
	if params.ReferenciaFlt != "" {
		flt := params.ReferenciaFlt
		if strings.Contains(flt, "/") {
			partes := strings.Split(flt, "/")
			if len(partes) == 2 && len(partes[0]) == 2 && len(partes[1]) == 4 {
				flt = partes[1] + partes[0]
			}
		}
		where = append(where, fmt.Sprintf("CONVERT(VARCHAR(MAX), [Referência]) LIKE @f%d", paramCounter))
		args = append(args, sql.Named(fmt.Sprintf("f%d", paramCounter), "%"+flt+"%"))
		paramCounter++
	}

	whereClause := strings.Join(where, " AND ")

	// Query sem paginação para exportação total
	queryData := fmt.Sprintf(`
		SELECT 
			ISNULL(CONVERT(VARCHAR(MAX), [Item.SubItem]), '') as ItemSubItem, 
			ISNULL([Descricao], '') as Descricao, 
			ISNULL(CONVERT(VARCHAR(MAX), [Tributo]), '') as Tributo, 
			ISNULL(CONVERT(VARCHAR(MAX), [DataPagamento], 103), '') as DataPagamento, -- Formato DD/MM/YYYY para Excel
			ISNULL(CONVERT(VARCHAR(MAX), [Referência]), '') as Referencia, 
			ISNULL(CONVERT(VARCHAR(MAX), [Municipio]), 'NÃO INFORMADO') as Municipio, 
			ISNULL(TRY_CAST([Valor Principal] AS FLOAT), 0) as ValorPrincipal, 
			ISNULL(TRY_CAST([Quantidade de UFERMS] AS FLOAT), 0) as QuantidadeUferms, 
			ISNULL(TRY_CAST([Valor Total] AS FLOAT), 0) as ValorTotal
		FROM %s
		WHERE %s
		ORDER BY [Descricao]
	`, config.ViewTaxas, whereClause)

	rows, err := db.Query(queryData, args...)
	if err != nil {
		return nil, fmt.Errorf("erro query exportacao: %w", err)
	}
	defer rows.Close()

	// Criar o arquivo Excel
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Taxas SEJUSP"
	index, _ := f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)

	// Cabeçalhos
	headers := []string{"Item/SubItem", "Descrição", "Tributo", "Data Pagamento", "Referência", "Município", "Valor Principal", "Qtd UFERMS", "Valor Total"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Estilo para o cabeçalho (Negrito e Fundo Cinza)
	styleHeader, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"E0E0E0"}, Pattern: 1},
	})
	f.SetCellStyle(sheet, "A1", "I1", styleHeader)

	// Inserir dados
	rowIdx := 2
	for rows.Next() {
		var d modelos.TaxaSejusp
		err := rows.Scan(
			&d.ItemSubItem, &d.Descricao, &d.Tributo, &d.DataPagamento, &d.Referencia,
			&d.Municipio, &d.ValorPrincipal, &d.QuantidadeUferms, &d.ValorTotal,
		)
		if err != nil {
			log.Printf("[ERRO EXPORT] Erro ao ler linha: %v", err)
			continue
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIdx), d.ItemSubItem)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIdx), d.Descricao)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIdx), d.Tributo)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIdx), d.DataPagamento)
		
		// Formatar Referência (YYYYMM -> MM/YYYY) para facilidade no Excel
		refStr := d.Referencia
		if len(refStr) == 6 {
			refStr = refStr[4:] + "/" + refStr[:4]
		}
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIdx), refStr)
		
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIdx), d.Municipio)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", rowIdx), d.ValorPrincipal)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", rowIdx), d.QuantidadeUferms)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", rowIdx), d.ValorTotal)
		
		rowIdx++
	}

	// Estilo para números (2 casas decimais - Valores)
	styleMoney, _ := f.NewStyle(&excelize.Style{
		NumFmt: 2, // Formato "0.00"
	})
	// Estilo para Inteiros (UFERMS)
	styleInt, _ := f.NewStyle(&excelize.Style{
		NumFmt: 1, // Formato "0"
	})
	
	f.SetColStyle(sheet, "G", styleMoney)
	f.SetColStyle(sheet, "H", styleInt)
	f.SetColStyle(sheet, "I", styleMoney)

	// Ajustar largura das colunas
	f.SetColWidth(sheet, "A", "A", 12)
	f.SetColWidth(sheet, "B", "B", 60)
	f.SetColWidth(sheet, "C", "I", 15)

	// Gravar em buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar buffer excel: %w", err)
	}

	return buf, nil
}

