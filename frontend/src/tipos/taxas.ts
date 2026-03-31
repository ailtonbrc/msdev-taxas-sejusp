export interface ITaxaSejusp {
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
