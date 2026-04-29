import React from 'react';
import { useSearchParams } from 'react-router-dom';
import api from '../servicos/api';

const PortalLogin: React.FC = () => {
  const [searchParams] = useSearchParams();

  React.useEffect(() => {
    const id = searchParams.get('id');
    if (id) {
      validarPortalLogin(id);
    } else {
      // Se cair aqui sem ID, volta para o portal oficial imediatamente
      window.location.href = 'https://it.fazenda.ms.gov.br/';
    }
  }, [searchParams]);

  const validarPortalLogin = async (id: string) => {
    try {
      await api.get(`/auth/portal-login?id=${id}`);
      // Login silencioso concluído, vai para a home
      window.location.href = '/';
    } catch (error) {
      // Em caso de erro, volta para o portal para o usuário tentar novamente
      window.location.href = 'https://it.fazenda.ms.gov.br/';
    }
  };

  // Não renderiza nada (silencioso)
  return null;
};

export default PortalLogin;
