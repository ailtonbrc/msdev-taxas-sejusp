import React from 'react';
import { Layout, Typography, ConfigProvider } from 'antd';
import ptBR from 'antd/locale/pt_BR';
import dayjs from 'dayjs';
import 'dayjs/locale/pt-br';

import Dashboard from './paginas/Dashboard';

dayjs.locale('pt-br');

const { Header, Content, Footer } = Layout;
const { Title } = Typography;

const headerStyle: React.CSSProperties = {
  textAlign: 'left',
  color: '#fff',
  height: 64, // Reduzido em ~25% (de 85px)
  padding: '0 50px',
  backgroundColor: '#004F9F',
  display: 'flex',
  alignItems: 'center',
};

const footerStyle: React.CSSProperties = {
  backgroundColor: '#004F9F',
  height: 56,
  padding: '5px 50px',
  color: 'white',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
};

const contentStyle: React.CSSProperties = {
  padding: '4px 50px', // Reduzido drasticamente (de 12px)
  minHeight: 'calc(100vh - 64px - 56px)',
  backgroundColor: '#f5f5f5',
};

const layoutStyle: React.CSSProperties = {
  minHeight: '100vh',
};

function App() {
  return (
    <ConfigProvider locale={ptBR}>
      <Layout style={layoutStyle}>
        <Header style={headerStyle}>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', width: '100%' }}>

            {/* Lado Esquerdo */}
            <div style={{ display: 'flex', alignItems: 'center' }}>
              <img
                src="logo_cotin.png"
                alt="Logo Cotin"
                style={{ height: 50, objectFit: 'contain', marginRight: 20 }}
              />
              <Title level={3} style={{ color: 'white', margin: 0, fontSize: '1.5rem' }}>
                Consulta de Taxas SEJUSP
              </Title>
            </div>

            {/* Lado Direito (LIMPO CONFORME SOLICITADO) */}
            <div style={{ display: 'flex', alignItems: 'center', height: 65, color: 'white', fontSize: '12px' }}>
              {/* Futuro local para Nome do Usuário / Logout */}
            </div>
          </div>
        </Header>

        <Content style={contentStyle}>
          <div style={{ padding: '12px 24px', minHeight: '100%', background: '#fff' }}>
            <Dashboard />
          </div>
        </Content>

        <Footer style={footerStyle}>
          <img
            src="logo-cotin-roda-pe.png"
            alt="Logo COTIN Rodapé"
            style={{ height: 40, objectFit: 'contain' }}
          />
        </Footer>

      </Layout>
    </ConfigProvider>
  );
}

export default App;