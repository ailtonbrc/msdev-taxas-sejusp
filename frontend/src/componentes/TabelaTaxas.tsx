import React, { useEffect, useState, useCallback } from 'react';
import { Table, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table'; 
import './TabelaTaxas.css'; 
import api from '../servicos/api';
import type { ITaxaSejusp, IRespostaDadosTaxas, IFiltrosState } from '../tipos/taxas';

const { Text } = Typography;

interface TabelaTaxasProps {
    filtros: IFiltrosState;
    isGuest?: boolean;
}

interface TabelaState {
    dados: ITaxaSejusp[];
    totalRegistros: number;
    loading: boolean;
    paginaAtual: number;
    limitePorPagina: number;
}


const TabelaTaxas: React.FC<TabelaTaxasProps> = ({ filtros, isGuest }) => {
  const [estado, setEstado] = useState<TabelaState>({
    dados: [],
    totalRegistros: 0,
    loading: false,
    paginaAtual: 1,
    limitePorPagina: 10, 
  });

  const buscarDados = useCallback(async (
    pagina: number = estado.paginaAtual, 
    limite: number = estado.limitePorPagina
  ) => {
    if (isGuest) {
        setEstado(prev => ({ ...prev, dados: [], totalRegistros: 0, loading: false }));
        return;
    }

    setEstado(prev => ({ ...prev, loading: true }));
    
    try {
      const params = new URLSearchParams();
      params.append('pagina', String(pagina));
      params.append('limite', String(limite)); 
      
      // Filtros do Dashboard (Seleção Múltipla)
      if (filtros.itemSubItem && filtros.itemSubItem.length > 0) {
        filtros.itemSubItem.forEach(v => params.append('itemSubItem', v));
      }
      if (filtros.referencia && filtros.referencia.length > 0) {
        filtros.referencia.forEach(v => params.append('referenciaFlt', v));
      }
      if (filtros.municipio && filtros.municipio.length > 0) {
        filtros.municipio.forEach(v => params.append('municipio', v));
      }
      if (filtros.instituicao && filtros.instituicao.length > 0) {
        filtros.instituicao.forEach(v => params.append('instituicao', v));
      }

      // Filtros Simples
      if (filtros.descricao) params.append('descricao', filtros.descricao);
      if (filtros.tributo) params.append('tributo', filtros.tributo);

      const url = `/taxas/dados?${params.toString()}`;
      
      const resposta = await api.get<IRespostaDadosTaxas>(url);
      const data = resposta.data;

      setEstado({
        dados: data.dados || [],
        totalRegistros: data.total_registros || 0,
        paginaAtual: data.pagina_atual || 1,
        limitePorPagina: data.limite_por_pagina || limite, 
        loading: false,
      });

    } catch (error) {
      console.error("Erro API Taxas:", error);
      setEstado(prev => ({ ...prev, loading: false }));
    }
  }, [filtros, estado.paginaAtual, estado.limitePorPagina]);

  useEffect(() => {
    buscarDados(1, estado.limitePorPagina);
  }, [filtros]);

  const handleTabelaChange = (pagination: any) => {
    buscarDados(pagination.current, pagination.pageSize);
  };

  const formatarMoeda = (valor: number): string => {
    return new Intl.NumberFormat('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(valor);
  };

  const colunas: ColumnsType<ITaxaSejusp> = [
    { 
      title: 'Instituição', 
      dataIndex: 'instituicao', 
      key: 'instituicao', 
      width: 150,
      fixed: 'left',
      render: (text: string) => <Text strong>{text}</Text>,
      sorter: (a, b) => a.instituicao.localeCompare(b.instituicao)
    },
    { 
      title: <div style={{ lineHeight: '1.2' }}>Item/<br/>SubItem</div>, 
      dataIndex: 'itemSubItem', 
      key: 'itemSubItem', 
      width: 75,
      fixed: 'left',
      align: 'center',
      render: (text) => <Text strong style={{ fontSize: '12px' }}>{text}</Text>,
      sorter: (a, b) => a.itemSubItem.localeCompare(b.itemSubItem)
    },
    { 
      title: 'Descrição', 
      dataIndex: 'descricao', 
      key: 'descricao', 
      width: 450, 
      render: (text) => <div style={{ whiteSpace: 'normal', wordBreak: 'break-word', lineHeight: '1.4' }}>{text}</div>,
      sorter: (a, b) => a.descricao.localeCompare(b.descricao)
    },
    { 
      title: 'Tributo', 
      dataIndex: 'tributo', 
      key: 'tributo', 
      width: 80, 
      align: 'center',
      sorter: (a, b) => a.tributo.localeCompare(b.tributo)
    },
    { 
      title: 'Data Pagamento', 
      dataIndex: 'dataPagamento', 
      key: 'dataPagamento', 
      width: 110,
      align: 'center',
      render: (val) => val ? new Date(val).toLocaleDateString('pt-BR') : '-',
      sorter: (a, b) => new Date(a.dataPagamento).getTime() - new Date(b.dataPagamento).getTime()
    },
    { 
      title: 'Referência', 
      dataIndex: 'referencia', 
      key: 'referencia', 
      width: 95, 
      align: 'center',
      render: (val) => (val && val.length === 6) ? `${val.substring(4, 6)}/${val.substring(0, 4)}` : val,
      sorter: (a, b) => a.referencia.localeCompare(b.referencia)
    },
    { 
      title: 'Município', 
      dataIndex: 'municipio', 
      key: 'municipio', 
      width: 160,
      render: (val) => val === 'NULL' || !val ? <Text type="secondary">NÃO INFORMADO</Text> : val,
      sorter: (a, b) => a.municipio.localeCompare(b.municipio)
    },
    { 
      title: 'Qtd UFERMS', 
      dataIndex: 'quantidadeUferms', 
      key: 'quantidadeUferms', 
      width: 95, 
      align: 'right',
      render: (v) => v.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 }),
      sorter: (a, b) => a.quantidadeUferms - b.quantidadeUferms
    },
    { 
      title: 'Valor Principal', 
      dataIndex: 'valorPrincipal', 
      key: 'valorPrincipal', 
      width: 105, 
      align: 'right',
      render: (v) => formatarMoeda(v),
      sorter: (a, b) => a.valorPrincipal - b.valorPrincipal
    },
    { 
      title: 'Valor Total', 
      dataIndex: 'valorTotal', 
      key: 'valorTotal', 
      width: 105, 
      align: 'right',
      render: (v) => <Text>{formatarMoeda(v)}</Text>,
      sorter: (a, b) => a.valorTotal - b.valorTotal
    },
  ];

  return (
    <div className="tabela-taxas-container">
      <Table 
        columns={colunas} 
        dataSource={estado.dados} 
        loading={estado.loading} 
        onChange={handleTabelaChange} 
        pagination={{ 
          current: estado.paginaAtual, 
          pageSize: estado.limitePorPagina, 
          total: estado.totalRegistros, 
          showSizeChanger: true, 
          pageSizeOptions: ['10', '20', '50', '100'] 
        }} 
        scroll={{ x: 'max-content', y: 'calc(100vh - 400px)' }} 
        rowKey={(record) => `${record.itemSubItem}-${record.dataPagamento}-${record.valorTotal}-${record.municipio}-${record.descricao.substring(0, 10)}`}
        size="small" 
        bordered 
      />
    </div>
  );
};

export default TabelaTaxas;