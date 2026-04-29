package servicos

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"taxas-sejusp/backend/config"
	"taxas-sejusp/backend/modelos"
	"time"

	"github.com/xuri/excelize/v2"
)

// CacheManager gerencia o cache em memória dos dados das taxas
type CacheManager struct {
	mu            sync.RWMutex
	dados         []modelos.TaxaSejusp
	opcoesFiltros modelos.OpcoesFiltros
	carregado     bool
	ultimaAtu     string // Armazena a string formatada do horário
	erroSync      string // Armazena a mensagem do último erro de sincronização
}

// GlobalCache é o singleton do cache da aplicação
var GlobalCache = &CacheManager{}

// Carregar busca todos os dados do SQL Server e armazena em memória
func (cm *CacheManager) Carregar(db *sql.DB) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	log.Println("[CACHE] Iniciando carregamento total da View no SQL Server...")

	query := fmt.Sprintf(`
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
	`, config.ViewTaxas)

	rows, err := db.Query(query)
	if err != nil {
		cm.mu.Lock()
		cm.erroSync = err.Error()
		cm.mu.Unlock()
		return fmt.Errorf("erro query cache: %w", err)
	}
	defer rows.Close()

	var novosDados []modelos.TaxaSejusp
	setItens := make(map[string]struct{})
	setTributos := make(map[string]struct{})
	setMunicipios := make(map[string]struct{})
	setRefs := make(map[string]struct{})
	setInstituicoes := make(map[string]struct{})

	for rows.Next() {
		var d modelos.TaxaSejusp
		err := rows.Scan(
			&d.ItemSubItem, &d.Descricao, &d.Tributo, &d.DataPagamento,
			&d.Referencia, &d.Municipio, &d.ValorPrincipal,
			&d.QuantidadeUferms, &d.ValorTotal,
		)
		if err != nil {
			log.Printf("[ERRO CACHE] Erro ao ler registro: %v", err)
			continue
		}

		// Remover espaços residuais (comum em campos CHAR do SQL Server)
		d.ItemSubItem = strings.TrimSpace(d.ItemSubItem)
		d.Descricao = strings.TrimSpace(d.Descricao)
		d.Tributo = strings.TrimSpace(d.Tributo)
		d.Municipio = strings.TrimSpace(d.Municipio)
		d.Referencia = strings.TrimSpace(d.Referencia)

		// Lógica de Derivação da Instituição
		descUpper := strings.ToUpper(d.Descricao)
		if strings.Contains(descUpper, "POLICIA CIVIL") || strings.Contains(descUpper, "COORDENADORIA GERAL DE PERICIAS") {
			d.Instituicao = "POLICIA CIVIL"
		} else if strings.Contains(descUpper, "POLICIA MILITAR") && d.Tributo == "510" {
			d.Instituicao = "POLICIA MILITAR"
		} else if strings.Contains(descUpper, "CORPO DE BOMBEIROS MILITAR") {
			d.Instituicao = "BOMBEIRO MILITAR"
		} else {
			d.Instituicao = "outros"
		}

		novosDados = append(novosDados, d)

		// Coleta de filtros (Valores Distintos)
		// Para Item/SubItem, vamos concatenar com a descrição para facilitar a busca do usuário
		if strings.TrimSpace(d.ItemSubItem) != "" {
			descCurta := d.Descricao
			if len(descCurta) > 100 {
				descCurta = descCurta[:100] + "..."
			}
			labelItem := fmt.Sprintf("%s - %s", d.ItemSubItem, descCurta)
			setItens[labelItem] = struct{}{}
		}
		if strings.TrimSpace(d.Tributo) != "" {
			setTributos[d.Tributo] = struct{}{}
		}
		if strings.TrimSpace(d.Municipio) != "" {
			setMunicipios[d.Municipio] = struct{}{}
		}
		if strings.TrimSpace(d.Referencia) != "" {
			setRefs[d.Referencia] = struct{}{}
		}
		if d.Instituicao != "" {
			setInstituicoes[d.Instituicao] = struct{}{}
		}
	}

	cm.dados = novosDados
	cm.ultimaAtu = time.Now().Format("02/01/2006 15:04")
	cm.opcoesFiltros = modelos.OpcoesFiltros{
		Instituicoes:      keysToSlice(setInstituicoes),
		Itens:             keysToSlice(setItens),
		Tributos:          keysToSlice(setTributos),
		Municipios:        keysToSlice(setMunicipios),
		Referencias:       keysToSlice(setRefs),
		UltimaAtualizacao: cm.ultimaAtu,
	}
	cm.carregado = true
	cm.erroSync = "" // Limpa erro se tiver sucesso

	log.Printf("[CACHE] Sucesso! %d registros carregados em RAM.", len(novosDados))
	return nil
}

// ObterStatus retorna o estado atual da sincronização com o banco
func (cm *CacheManager) ObterStatus() (bool, string, string) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.carregado, cm.ultimaAtu, cm.erroSync
}

// Listar realiza filtros e paginação diretamente na memória com suporte a RBAC
func (cm *CacheManager) Listar(params modelos.ParametrosBusca, isAdmin bool, allowed []string) ([]modelos.TaxaSejusp, int64) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if !cm.carregado {
		return []modelos.TaxaSejusp{}, 0
	}

	var filtrados []modelos.TaxaSejusp

	// 1. Aplicar Filtros
	for _, d := range cm.dados {
		if matchFiltros(d, params, isAdmin, allowed) {
			filtrados = append(filtrados, d)
		}
	}

	total := int64(len(filtrados))

	// 2. Aplicar Paginação
	offset := params.CalcularOffset()
	fim := offset + params.Limite

	if offset >= len(filtrados) {
		return []modelos.TaxaSejusp{}, total
	}
	if fim > len(filtrados) {
		fim = len(filtrados)
	}

	return filtrados[offset:fim], total
}

// ObterOpcoes retorna os valores distintos para os selects do frontend de forma dinâmica (Filtros Relevantes) respeitando o RBAC
func (cm *CacheManager) ObterOpcoes(p modelos.ParametrosBusca, isAdmin bool, allowed []string) modelos.OpcoesFiltros {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if !cm.carregado {
		return cm.opcoesFiltros
	}

	// Caso contrário, calculamos as opções relevantes baseadas na seleção atual (Filtro Cruzado)
	log.Printf("[OPCOES] Processando filtros cruzados. Parâmetros recebidos: ItemSubItem: %d, Descricao: %s, Tributo: %d, Referencia: %d, Municipio: %d, Instituicao: %d",
		len(p.ItemSubItem), p.Descricao, len(p.Tributo), len(p.ReferenciaFlt), len(p.Municipio), len(p.Instituicao))
	setItens := make(map[string]struct{})
	setTributos := make(map[string]struct{})
	setMunicipios := make(map[string]struct{})
	setRefs := make(map[string]struct{})
	setInstituicoes := make(map[string]struct{})

	for _, d := range cm.dados {
		// Para cada seletor, calculamos os valores possíveis aplicando os filtros de TODOS os outros campos.
		// Isso permite que o seletor atual mostre todas as opções válidas para o contexto dos demais filtros.

		// 1. Itens Relevantes: Ignora o filtro do próprio campo ItemSubItem
		pItens := p
		pItens.ItemSubItem = []string{}
		if matchFiltros(d, pItens, isAdmin, allowed) {
			if strings.TrimSpace(d.ItemSubItem) != "" {
				descCurta := d.Descricao
				if len(descCurta) > 100 {
					descCurta = descCurta[:100] + "..."
				}
				setItens[fmt.Sprintf("%s - %s", d.ItemSubItem, descCurta)] = struct{}{}
			}
		}

		// 2. Tributos Relevantes: Ignora o próprio filtro Tributo
		pTrib := p
		pTrib.Tributo = []string{}
		if matchFiltros(d, pTrib, isAdmin, allowed) {
			if strings.TrimSpace(d.Tributo) != "" {
				setTributos[d.Tributo] = struct{}{}
			}
		}

		// 3. Municípios Relevantes: Ignora o próprio filtro Município
		pMun := p
		pMun.Municipio = []string{}
		if matchFiltros(d, pMun, isAdmin, allowed) {
			if strings.TrimSpace(d.Municipio) != "" {
				setMunicipios[d.Municipio] = struct{}{}
			}
		}

		// 4. Referências Relevantes: Ignora o próprio filtro de Referência individual
		pRef := p
		pRef.ReferenciaFlt = []string{}
		if matchFiltros(d, pRef, isAdmin, allowed) {
			if strings.TrimSpace(d.Referencia) != "" {
				setRefs[d.Referencia] = struct{}{}
			}
		}

		// 5. Instituições Relevantes: Ignora o próprio filtro Instituição
		pInst := p
		pInst.Instituicao = []string{}
		if matchFiltros(d, pInst, isAdmin, allowed) {
			if d.Instituicao != "" {
				setInstituicoes[d.Instituicao] = struct{}{}
			}
		}
	}

	return modelos.OpcoesFiltros{
		Instituicoes:      keysToSlice(setInstituicoes),
		Itens:             keysToSlice(setItens),
		Tributos:          keysToSlice(setTributos),
		Municipios:        keysToSlice(setMunicipios),
		Referencias:       keysToSlice(setRefs),
		UltimaAtualizacao: cm.ultimaAtu,
	}
}

// matchFiltros contém a lógica de comparação (LIKE, Case-Insensitive) e Segurança
func matchFiltros(d modelos.TaxaSejusp, p modelos.ParametrosBusca, isAdmin bool, allowed []string) bool {
	// 1. Filtro de Segurança (RBAC)
	if !isAdmin {
		permitido := false
		for _, a := range allowed {
			if strings.EqualFold(strings.TrimSpace(d.Instituicao), strings.TrimSpace(a)) {
				permitido = true
				break
			}
		}
		if !permitido {
			return false
		}
	}

	// 2. Filtro de Referência Global (periodoInicio)
	if p.PeriodoInicio != nil && *p.PeriodoInicio != "" {
		if d.Referencia != *p.PeriodoInicio {
			return false
		}
	}

	// Filtro ItemSubItem (Suporta seleção múltipla de "CÓDIGO - DESCRIÇÃO")
	if len(p.ItemSubItem) > 0 {
		found := false
		for _, sel := range p.ItemSubItem {
			if sel == "" {
				continue
			}
			// O front envia "CÓDIGO - DESCRIÇÃO", então vamos extrair o código para uma busca exata
			parts := strings.SplitN(sel, " - ", 2)
			if len(parts) > 0 && strings.EqualFold(parts[0], d.ItemSubItem) {
				found = true
				break
			}
		}
		// Se o filtro estava preenchido mas nenhum valor bateu
		if len(p.ItemSubItem) > 0 && !found {
			// Mas só bloqueamos se houver pelo menos um critério REAL (não vazio)
			temCriterioReal := false
			for _, s := range p.ItemSubItem {
				if s != "" {
					temCriterioReal = true
					break
				}
			}
			if temCriterioReal {
				return false
			}
		}
	}

	// Filtro Descricao (LIKE)
	if p.Descricao != "" {
		if !strings.Contains(strings.ToLower(d.Descricao), strings.ToLower(p.Descricao)) {
			return false
		}
	}

	// Filtro Tributo (Seleção Múltipla e Case-Insensitive)
	if len(p.Tributo) > 0 {
		found := false
		temCriterio := false
		for _, t := range p.Tributo {
			if t == "" {
				continue
			}
			temCriterio = true
			if strings.EqualFold(d.Tributo, t) {
				found = true
				break
			}
		}
		if temCriterio && !found {
			return false
		}
	}

	// Filtro Município (Seleção Múltipla)
	if len(p.Municipio) > 0 {
		found := false
		temCriterio := false
		for _, m := range p.Municipio {
			if m == "" {
				continue
			}
			temCriterio = true
			if strings.EqualFold(d.Municipio, m) {
				found = true
				break
			}
		}
		if temCriterio && !found {
			return false
		}
	}

	// Filtro Instituição (Seleção Múltipla)
	if len(p.Instituicao) > 0 {
		found := false
		temCriterio := false
		for _, i := range p.Instituicao {
			if i == "" {
				continue
			}
			temCriterio = true
			if strings.EqualFold(d.Instituicao, i) {
				found = true
				break
			}
		}
		if temCriterio && !found {
			return false
		}
	}

	// Filtro Referencia Individual (Seleção Múltipla)
	if len(p.ReferenciaFlt) > 0 {
		found := false
		temCriterio := false
		for _, r := range p.ReferenciaFlt {
			if r == "" {
				continue
			}
			temCriterio = true
			if d.Referencia == r {
				found = true
				break
			}
		}
		if temCriterio && !found {
			return false
		}
	}

	return true
}

// IniciarSincronizacaoAutomatica inicia um worker em background para atualizar o cache
// nos horários agendados: 06:00, 12:00 e 16:00
func IniciarSincronizacaoAutomatica(db *sql.DB) {
	go func() {
		log.Println("[CRON] Agendador de sincronização iniciado (06:00, 12:00, 16:00)")

		// Ticker de 1 minuto para checar o horário agendado
		tickerAgendado := time.NewTicker(1 * time.Minute)
		// Ticker de 2 minutos para retentativa se estiver com erro
		tickerRetry := time.NewTicker(2 * time.Minute)

		defer tickerAgendado.Stop()
		defer tickerRetry.Stop()

		for {
			select {
			case <-tickerAgendado.C:
				agora := time.Now()
				horaMinuto := agora.Format("15:04")

				// Verifica se coincide com os horários desejados
				if horaMinuto == "06:00" || horaMinuto == "12:00" || horaMinuto == "16:00" {
					log.Printf("[CRON] Horário agendado atingido (%s). Iniciando atualização...", horaMinuto)
					err := GlobalCache.Carregar(db)
					if err != nil {
						log.Printf("[CRON ERRO] Falha na atualização agendada: %v", err)
					}
				}

			case <-tickerRetry.C:
				// Se o cache nunca foi carregado ou está com erro recente, tenta novamente
				cm_cache := GlobalCache
				cm_cache.mu.RLock()
				semDados := !cm_cache.carregado
				temErro := cm_cache.erroSync != ""
				cm_cache.mu.RUnlock()

				if semDados || temErro {
					log.Println("[RETRY] Sistema ainda sem conexão estável com Banco. Tentando reconectar...")
					err := GlobalCache.Carregar(db)
					if err == nil {
						log.Println("[RETRY] Conexão com Banco de Dados restabelecida com sucesso!")
					}
				}
			}
		}
	}()
}

// ExportarExcel gera um arquivo .xlsx com os dados filtrados respeitando RBAC
func (cm *CacheManager) ExportarExcel(params modelos.ParametrosBusca, targetPath string, isAdmin bool, allowed []string) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if !cm.carregado {
		return fmt.Errorf("dados não estão carregados na memória")
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Taxas SEJUSP"
	f.SetSheetName("Sheet1", sheet)

	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		return err
	}

	// Estilos do Cabeçalho
	styleHeader, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"E0E0E0"}, Pattern: 1},
	})

	headers := []interface{}{"Instituição", "Item/SubItem", "Descrição", "Tributo", "Data Pagamento", "Referência", "Município", "Qtd UFERMS", "Valor Principal", "Valor Total"}
	if err := sw.SetRow("A1", headers); err != nil {
		return err
	}

	rowIdx := 2
	for _, d := range cm.dados {
		if !matchFiltros(d, params, isAdmin, allowed) {
			continue
		}

		// Formatações (Stream)
		dtStr := d.DataPagamento
		if len(dtStr) >= 10 {
			dtStr = dtStr[8:10] + "/" + dtStr[5:7] + "/" + dtStr[0:4]
		}

		// Formatar Referência como Texto Padrão "MM/YYYY"
		refStr := d.Referencia
		if len(refStr) == 6 {
			refStr = fmt.Sprintf("%s/%s", refStr[4:], refStr[:4])
		}

		// Gravar linha
		row := []interface{}{
			d.Instituicao,
			d.ItemSubItem,
			d.Descricao,
			d.Tributo,
			dtStr,
			refStr,
			d.Municipio,
			d.QuantidadeUferms,
			d.ValorPrincipal,
			d.ValorTotal,
		}

		cell, _ := excelize.CoordinatesToCellName(1, rowIdx)
		if err := sw.SetRow(cell, row); err != nil {
			log.Printf("[ERRO EXCEL] Falha no fluxo na linha %d: %v", rowIdx, err)
			continue
		}
		rowIdx++
	}

	if err := sw.Flush(); err != nil {
		return err
	}

	// Ajustar Estilos
	styleMoney, _ := f.NewStyle(&excelize.Style{NumFmt: 2})
	styleFloat, _ := f.NewStyle(&excelize.Style{NumFmt: 2})
	customFmt := "@"
	styleText, _ := f.NewStyle(&excelize.Style{CustomNumFmt: &customFmt})

	f.SetColStyle(sheet, "F:F", styleText)
	f.SetColStyle(sheet, "H:H", styleFloat)
	f.SetColStyle(sheet, "I:I", styleMoney)
	f.SetColStyle(sheet, "J:J", styleMoney)
	f.SetCellStyle(sheet, "A1", "J1", styleHeader)
	f.SetColWidth(sheet, "C", "J", 15)
	f.SetColWidth(sheet, "F", "F", 12)

	return f.SaveAs(targetPath)
}

// ExportarCSV gera um arquivo CSV separado por ; respeitando RBAC
func (cm *CacheManager) ExportarCSV(params modelos.ParametrosBusca, targetPath string, isAdmin bool, allowed []string) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if !cm.carregado {
		return fmt.Errorf("cache não carregado")
	}

	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Escrever BOM para UTF-8
	file.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	headers := []string{"Instituição", "Item/SubItem", "Descrição", "Tributo", "Data Pagamento", "Referência", "Município", "Qtd UFERMS", "Valor Principal", "Valor Total"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, d := range cm.dados {
		if !matchFiltros(d, params, isAdmin, allowed) {
			continue
		}

		dtStr := d.DataPagamento
		if len(dtStr) >= 10 {
			dtStr = dtStr[8:10] + "/" + dtStr[5:7] + "/" + dtStr[0:4]
		}

		refStr := d.Referencia
		if len(refStr) == 6 {
			refStr = fmt.Sprintf("%s/%s", refStr[4:], refStr[:4])
		}

		qtdStr := strings.Replace(fmt.Sprintf("%.2f", d.QuantidadeUferms), ".", ",", 1)
		valPrincipalStr := strings.Replace(fmt.Sprintf("%.2f", d.ValorPrincipal), ".", ",", 1)
		valTotalStr := strings.Replace(fmt.Sprintf("%.2f", d.ValorTotal), ".", ",", 1)

		row := []string{
			d.Instituicao,
			d.ItemSubItem,
			d.Descricao,
			d.Tributo,
			dtStr,
			refStr,
			d.Municipio,
			qtdStr,
			valPrincipalStr,
			valTotalStr,
		}

		if err := writer.Write(row); err != nil {
			log.Printf("[ERRO CSV] Falha ao gravar record: %v", err)
			continue
		}
	}
	return nil
}

func keysToSlice(m map[string]struct{}) []string {
	res := make([]string, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}
