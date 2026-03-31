package modelos

type ParametrosBusca struct {
	Pagina     int    `form:"pagina" binding:"omitempty,gte=1"`
	Limite     int    `form:"limite" binding:"omitempty,gte=1"`
	OrdenarPor string `form:"ordenarPor" binding:"omitempty"`

	// --- ALTERADO: Filtro por Período (YYYYMM) ---
	// Tags corrigidas para coincidir com o frontend (periodoInicio, periodoFim)
	PeriodoInicio *string `form:"periodoInicio" binding:"omitempty,len=6"` // Ex: "202501"
	PeriodoFim    *string `form:"periodoFim" binding:"omitempty,len=6"`    // Ex: "202503"

	NomeAdmin *string `form:"nomeAdmin" binding:"omitempty"`
	Busca         string  `form:"busca" binding:"omitempty"`
	ItemSubItem   string  `form:"itemSubItem" binding:"omitempty"`
	Descricao     string  `form:"descricao" binding:"omitempty"`
	Tributo       string  `form:"tributo" binding:"omitempty"`
	ReferenciaFlt string  `form:"referenciaFlt" binding:"omitempty"`
}

func (p *ParametrosBusca) Normalizar() {
	if p.Pagina <= 0 {
		p.Pagina = 1
	}
	if p.Limite <= 0 {
		p.Limite = 20
	}
}

func (p *ParametrosBusca) CalcularOffset() int {
	return (p.Pagina - 1) * p.Limite
}

