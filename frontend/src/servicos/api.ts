import axios from 'axios';

// A URL base será pega da variável de ambiente VITE_API_BASE_URL (definida no .env.production)
// Se não estiver definida (dev), usa '/api' que será capturado pelo proxy do Vite
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 120000, // Aumentado para 120 segundos (2 minutos)
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor para tratar erros de autenticação (401)
// Interceptor removido temporariamente para produção rápida sem autenticação
/*
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      window.location.href = '/api/auth/login';
    }
    return Promise.reject(error);
  }
);
*/

export default api;