import axios from 'axios';

import type { IOpcoesFiltros } from '../tipos/taxas';

// A URL base será pega da variável de ambiente VITE_API_BASE_URL (definida no .env.production)
// Se não estiver definida (dev), usa '/api' que será capturado pelo proxy do Vite
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 120000, // Aumentado para 120 segundos (2 minutos)
  headers: {
    'Content-Type': 'application/json',
  },
  paramsSerializer: {
    indexes: null // Previne o envio de colchetes [] em parâmetros de array, compatível com o backend Go (Gin)
  }
});

// Funções de API
export const getOpcoesFiltros = async (filtros?: any): Promise<IOpcoesFiltros> => {
  const response = await api.get<IOpcoesFiltros>('/taxas/opcoes-filtros', { params: filtros });
  return response.data;
};

export const sincronizarBanco = async (): Promise<void> => {
  await api.post('/taxas/refresh');
};

// Interceptor para tratar erros de autenticação (401)
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Apenas rejeitamos o erro. O App.tsx ou as páginas decidirão o que fazer.
    return Promise.reject(error);
  }
);

export default api;