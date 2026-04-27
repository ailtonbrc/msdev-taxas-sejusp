import React from 'react';
import { Card, Form, Input, Button, Typography, message, Layout } from 'antd';
import { UserOutlined, LoginOutlined } from '@ant-design/icons';
import api from '../servicos/api';

const { Title, Text } = Typography;
const { Content } = Layout;

const Login: React.FC = () => {
  const [loading, setLoading] = React.useState(false);

  const onFinish = async (values: { cpf: string }) => {
    setLoading(true);
    try {
      await api.post('/auth/mock-login', { cpf: values.cpf });
      message.success('Login realizado com sucesso!');
      // Redireciona para o dashboard
      window.location.href = '/';
    } catch (error: any) {
      const msg = error.response?.data?.erro || 'Falha ao realizar login';
      message.error(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout style={{ minHeight: '100vh', background: 'linear-gradient(135deg, #001529 0%, #004F9F 100%)' }}>
      <Content style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
        <Card
          style={{
            width: 400,
            borderRadius: 12,
            boxShadow: '0 8px 24px rgba(0,0,0,0.2)',
            border: 'none'
          }}
        >
          <div style={{ textAlign: 'center', marginBottom: 30 }}>
            <img
              src="/logo_cotin.png"
              alt="Logo"
              style={{ height: 60, marginBottom: 16 }}
            />
            <Title level={3} style={{ margin: 0, color: '#004F9F' }}>Taxas SEJUSP</Title>
            <Text type="secondary">Portal de Arrecadação - Homologação</Text>
          </div>

          <Form
            name="login"
            onFinish={onFinish}
            layout="vertical"
            size="large"
          >
            <Form.Item
              name="cpf"
              rules={[
                { required: true, message: 'Por favor, insira seu CPF!' },
              ]}
            >
              <Input
                prefix={<UserOutlined style={{ color: 'rgba(0,0,0,.25)' }} />}
                placeholder="Digite seu CPF (apenas números)"
              />
            </Form.Item>

            <Form.Item style={{ marginBottom: 0 }}>
              <Button
                type="primary"
                htmlType="submit"
                loading={loading}
                icon={<LoginOutlined />}
                style={{ width: '100%', height: 45, borderRadius: 6, backgroundColor: '#004F9F' }}
              >
                ENTRAR
              </Button>
            </Form.Item>
          </Form>

          <div style={{ marginTop: 24, textAlign: 'center' }}>
            <Text style={{ fontSize: 12, color: '#999' }}>
              Utilize o CPF cadastrado pelo administrador para acessar o simulador de e-Fazenda.
            </Text>
          </div>
        </Card>
      </Content>
    </Layout>
  );
};

export default Login;
