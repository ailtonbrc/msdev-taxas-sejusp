package modelos

type TaxaSejusp struct {
	ItemSubItem      string  `json:"itemSubItem"`
	Descricao        string  `json:"descricao"`
	Tributo          string  `json:"tributo"`
	DataPagamento    string  `json:"dataPagamento"`
	Referencia       string  `json:"referencia"`
	Municipio        string  `json:"municipio"`
	ValorPrincipal   float64 `json:"valorPrincipal"`
	QuantidadeUferms float64 `json:"quantidadeUferms"`
	ValorTotal       float64 `json:"valorTotal"`
}

type RespostaDadosTaxas struct {
	TotalRegistros  int64        `json:"total_registros"`
	PaginaAtual     int          `json:"pagina_atual"`
	LimitePorPagina int          `json:"limite_por_pagina"`
	Dados           []TaxaSejusp `json:"dados"`
}
