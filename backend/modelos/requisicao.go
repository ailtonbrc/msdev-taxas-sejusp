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
	ItemSubItem   []string `form:"itemSubItem"`
	Descricao     string  `form:"descricao" binding:"omitempty"`
	Tributo       []string `form:"tributo"`
	ReferenciaFlt []string `form:"referenciaFlt"`
	Municipio     []string `form:"municipio"`
	Instituicao   []string `form:"instituicao"`
}

func (p *ParametrosBusca) Normalizar() {
	if p.Pagina <= 0 {
		p.Pagina = 1
	}
	if p.Limite <= 0 {
		p.Limite = 20
	}

	// Limpar strings vazias de slices (evita filtros fantasmas como ?tributo=)
	p.ItemSubItem = limparSlice(p.ItemSubItem)
	p.Tributo = limparSlice(p.Tributo)
	p.ReferenciaFlt = limparSlice(p.ReferenciaFlt)
	p.Municipio = limparSlice(p.Municipio)
	p.Instituicao = limparSlice(p.Instituicao)
}

func limparSlice(s []string) []string {
	var r []string
	for _, v := range s {
		if v != "" {
			r = append(r, v)
		}
	}
	return r
}

func (p *ParametrosBusca) CalcularOffset() int {
	return (p.Pagina - 1) * p.Limite
}

