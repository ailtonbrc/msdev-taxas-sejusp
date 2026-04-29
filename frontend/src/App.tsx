import React from 'react';
import { Layout, Typography, ConfigProvider, Button, Space, Avatar, Menu } from 'antd';
import { LogoutOutlined, UserOutlined, DashboardOutlined, TeamOutlined } from '@ant-design/icons';
import { BrowserRouter, Routes, Route, Link, Navigate, useLocation, Outlet } from 'react-router-dom';
import ptBR from 'antd/locale/pt_BR';
import dayjs from 'dayjs';
import 'dayjs/locale/pt-br';

import Dashboard from './paginas/Dashboard';
import Usuarios from './paginas/Usuarios';
import Login from './paginas/PortalLogin';
import api from './servicos/api';

dayjs.locale('pt-br');

const { Header, Content, Footer } = Layout;
const { Title } = Typography;

interface Usuario {
  id: number;
  nome: string;
  cpf: string;
  role: string;
  instituicoes: string[];
  autenticado: boolean;
}

const App: React.FC = () => {
  const [usuario, setUsuario] = React.useState<Usuario | null>(null);
  const [loading, setLoading] = React.useState(true);

  const carregarUsuario = async () => {
    try {
      const response = await api.get('/auth/user');
      if (response.data && response.data.autenticado) {
        setUsuario(response.data);
      } else {
        setUsuario(null);
      }
    } catch (error: any) {
      // 401 é esperado se não estiver logado, não será tratado como erro crítico no console
      if (error.response?.status !== 401) {
        console.error('Erro ao carregar usuário:', error);
      }
      setUsuario(null);
    } finally {
      setLoading(false);
    }
  };

  React.useEffect(() => {
    carregarUsuario();
  }, []);

  const handleLogout = async () => {
    try {
      await api.get('/auth/logout');
      setUsuario(null);
      // Redireciona o usuário de volta para o Portal e-Fazenda oficial
      window.location.href = 'https://it.fazenda.ms.gov.br/';
    } catch (error) {
      console.error('Erro ao sair:', error);
      window.location.href = 'https://it.fazenda.ms.gov.br/';
    }
  };

  if (loading) return null;

  return (
    <ConfigProvider locale={ptBR}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          
          {/* Rotas Protegidas */}
          <Route element={usuario ? <MainLayout usuario={usuario} onLogout={handleLogout} /> : <Navigate to="/login" replace />}>
            <Route path="/" element={<Dashboard />} />
            {usuario?.role === 'admin' && <Route path="/usuarios" element={<Usuarios />} />}
          </Route>

          {/* Fallback */}
          <Route path="*" element={<Navigate to={usuario ? "/" : "/login"} replace />} />
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  );
};

const MainLayout: React.FC<{ usuario: Usuario; onLogout: () => void }> = ({ usuario, onLogout }) => {
  const location = useLocation();

  const menuItems = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: <Link to="/">Dashboard</Link>,
    },
    ...(usuario.role === 'admin' ? [{
      key: '/usuarios',
      icon: <TeamOutlined />,
      label: <Link to="/usuarios">Usuários</Link>,
    }] : []),
  ];

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center', backgroundColor: '#004F9F', padding: '0 24px', height: 64 }}>
        <div style={{ display: 'flex', alignItems: 'center', flex: 1 }}>
          <img src="/logo_cotin.png" alt="Logo" style={{ height: 40, marginRight: 16 }} />
          <Title level={4} style={{ color: 'white', margin: 0, marginRight: 48 }}>Taxas SEJUSP</Title>
          
          {usuario.role === 'admin' && (
            <Menu 
              theme="dark" 
              mode="horizontal" 
              selectedKeys={[location.pathname]} 
              items={menuItems} 
              style={{ backgroundColor: 'transparent', flex: 1, borderBottom: 'none' }}
            />
          )}
        </div>

        <Space size="middle" style={{ color: 'white' }}>
          <div style={{ textAlign: 'right', lineHeight: '1.3' }}>
            <div style={{ fontWeight: 600, fontSize: '14px', letterSpacing: '0.5px' }}>
              {usuario.nome.toUpperCase()}
            </div>
            <div style={{ fontSize: '11px', opacity: 0.85, fontStyle: 'italic' }}>
              {usuario.role === 'admin' 
                ? 'Administrador (Acesso Total)' 
                : `Permissões: ${usuario.instituicoes.length > 0 ? usuario.instituicoes.join(', ') : 'Nenhuma'}`
              }
            </div>
          </div>
          <Avatar 
            size="large"
            icon={<UserOutlined />} 
            style={{ backgroundColor: usuario.role === 'admin' ? '#f5222d' : '#1890ff', border: '2px solid rgba(255,255,255,0.2)' }} 
          />
          <Button 
            type="text" 
            icon={<LogoutOutlined />} 
            style={{ color: 'white', marginLeft: 8 }} 
            onClick={onLogout}
            title="Sair do Sistema"
          />
        </Space>
      </Header>

      <Content style={{ padding: '24px', backgroundColor: '#f5f5f5' }}>
        <div style={{ background: '#fff', padding: '24px', borderRadius: 8, minHeight: '100%' }}>
          <Outlet context={{ usuario }} />
        </div>
      </Content>

      <Footer style={{ textAlign: 'center', padding: '12px', backgroundColor: '#004F9F' }}>
        <img src="/logo-cotin-roda-pe.png" alt="Logo Rodapé" style={{ height: 30 }} />
      </Footer>
    </Layout>
  );
};

export default App;