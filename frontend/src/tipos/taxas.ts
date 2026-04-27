export interface ITaxaSejusp {
  instituicao: string;
  itemSubItem: string;
  descricao: string;
  tributo: string;
  dataPagamento: string;
  referencia: string;
  municipio: string;
  valorPrincipal: number;
  quantidadeUferms: number;
  valorTotal: number;
}

export interface IRespostaDadosTaxas {
  total_registros: number;
  pagina_atual: number;
  limite_por_pagina: number;
  dados: ITaxaSejusp[];
}

export interface IOpcoesFiltros {
  instituicoes: string[];
  itens: string[];
  tributos: string[];
  municipios: string[];
  referencias: string[];
  ultimaAtualizacao: string;
}

export interface IFiltrosState {
  itemSubItem: string[];
  descricao: string;
  tributo: string;
  referencia: string[];
  municipio: string[];
  instituicao: string[];
}
