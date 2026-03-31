import { useEffect, useState, useCallback } from 'react';
import { Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { SorterResult } from 'antd/es/table/interface';
import api from '../servicos/api';

interface IRespostaAPI<T> {
    total_registros?: number;
    pagina_atual?: number;
    limite_por_pagina?: number;
    dados?: T[];
}

interface FiltrosState {
    periodoInicio?: string;
    periodoFim?: string;
    nomeAdmin?: string;
}

interface TabelaGenericaProps<T> {
    filtros: FiltrosState;
    url: string;
    colunas: ColumnsType<T>;
    onCarregamentoFinalizado?: (sucesso: boolean, total: number) => void;
}

interface TabelaState<T> {
    dados: T[];
    totalRegistros: number;
    loading: boolean;
    paginaAtual: number;
    limitePorPagina: number;
}

type Sorter<T> = SorterResult<T>;

function useTabelaDados<T>(
    url: string,
    filtros: FiltrosState,
    onCarregamentoFinalizado?: (sucesso: boolean, total: number) => void
) {
    const [estado, setEstado] = useState<TabelaState<T>>({
        dados: [],
        totalRegistros: 0,
        loading: false,
        paginaAtual: 1,
        limitePorPagina: 15,
    });

    const obterParametroOrdenacao = (sorter: Sorter<T>): string => {
        if (sorter.field && sorter.order) {
            const campo = String(sorter.field);
            const ordem = sorter.order === 'ascend' ? 'ASC' : 'DESC';
            return `${campo} ${ordem}`;
        }
        return '';
    };

    const buscarDados = useCallback(async (
        pagina: number = estado.paginaAtual,
        limite: number | undefined,
        ordenacao: string = ''
    ) => {
        if (!filtros.periodoInicio || !filtros.periodoFim) return;

        setEstado(prev => ({ ...prev, loading: true }));

        try {
            const params = new URLSearchParams();
            params.append('pagina', String(pagina));
            const limiteFinal = limite || estado.limitePorPagina;
            params.append('limite', String(limiteFinal));
            if (ordenacao) params.append('ordenarPor', ordenacao);

            params.append('inicio', filtros.periodoInicio);
            params.append('fim', filtros.periodoFim);
            if (filtros.nomeAdmin) params.append('nomeAdmin', filtros.nomeAdmin);

            const resposta = await api.get<IRespostaAPI<T>>(`${url}?${params.toString()}`);
            const data = resposta.data;

            setEstado(prev => ({
                ...prev,
                dados: data.dados || [],
                totalRegistros: data.total_registros || 0,
                paginaAtual: data.pagina_atual || 1,
                limitePorPagina: data.limite_por_pagina || limiteFinal,
                loading: false,
            }));

            if (onCarregamentoFinalizado) {
                onCarregamentoFinalizado(true, data.total_registros || 0);
            }

        } catch (error) {
            console.error("Erro API:", error);
            setEstado(prev => ({ ...prev, loading: false }));
            if (onCarregamentoFinalizado) {
                onCarregamentoFinalizado(false, 0);
            }
        }
    }, [url, filtros.periodoInicio, filtros.periodoFim, filtros.nomeAdmin, estado.paginaAtual, estado.limitePorPagina, onCarregamentoFinalizado]);

    const handleTabelaChange = useCallback((
        pagination: any,
        sorter: Sorter<T> | Sorter<T>[]
    ) => {
        const sorterArray = Array.isArray(sorter) ? sorter : [sorter];
        const ordenacao = sorterArray.map(s => obterParametroOrdenacao(s)).filter(s => s).join(', ');
        buscarDados(pagination.current, pagination.pageSize, ordenacao);
    }, [buscarDados]);

    return { ...estado, buscarDados, handleTabelaChange };
}

function TabelaGenerica<T extends object>({ filtros, url, colunas, onCarregamentoFinalizado }: TabelaGenericaProps<T>) {
    const { dados, loading, buscarDados, handleTabelaChange, totalRegistros, paginaAtual, limitePorPagina } = useTabelaDados<T>(url, filtros, onCarregamentoFinalizado);

    useEffect(() => {
        if (filtros.periodoInicio && filtros.periodoFim) {
            buscarDados(1, limitePorPagina, '');
        }
    }, [filtros, buscarDados, limitePorPagina]);

    return (
        <Table
            columns={colunas}
            dataSource={dados}
            loading={loading}
            onChange={handleTabelaChange}
            pagination={{
                current: paginaAtual,
                pageSize: limitePorPagina,
                total: totalRegistros,
                showSizeChanger: true,
                pageSizeOptions: ['15', '20', '50', '100']
            }}
            scroll={{ x: 'max-content', y: 500 }}
            rowKey={(record: any) => record.id || Math.random().toString()} // Fallback key
            size="small"
            bordered
            className="tabela-customizada"
        />
    );
}

export default TabelaGenerica;
