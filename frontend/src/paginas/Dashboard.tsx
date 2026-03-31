import React, { useState } from 'react';
import { Layout, Card, Form, Button, Typography, Space, Input } from 'antd';
import { SearchOutlined, FileExcelOutlined } from '@ant-design/icons';
import TabelaTaxas from '../componentes/TabelaTaxas';

const { Content } = Layout;
const { Text, Title } = Typography;

const Dashboard: React.FC = () => {
  const [form] = Form.useForm();
  const [filtros, setFiltros] = useState<{ 
    itemSubItem: string; 
    descricao: string; 
    tributo: string; 
    referencia: string;
  }>({ 
    itemSubItem: '',
    descricao: '',
    tributo: '',
    referencia: ''
  });

  const handleFiltrar = (values: any) => {
    setFiltros({
      itemSubItem: values.itemSubItem || '',
      descricao: values.descricao || '',
      tributo: values.tributo || '',
      referencia: values.referencia || ''
    });
  };

  const handleLimpar = () => {
    form.resetFields();
    setFiltros({
      itemSubItem: '',
      descricao: '',
      tributo: '',
      referencia: ''
    });
  };

  const handleExportarExcel = async () => {
    try {
      const params = new URLSearchParams();
      if (filtros.itemSubItem) params.append('itemSubItem', filtros.itemSubItem);
      if (filtros.descricao) params.append('descricao', filtros.descricao);
      if (filtros.tributo) params.append('tributo', filtros.tributo);
      if (filtros.referencia) params.append('referenciaFlt', filtros.referencia);

      // Usar a URL base da API (ajustar conforme necessário se houver env)
      const urlBase = import.meta.env.VITE_API_URL || ''; 
      const urlFinal = `${urlBase}/api/taxas/exportar?${params.toString()}`;
      
      // Disparar o download de forma silenciosa e direta
      const link = document.createElement('a');
      link.href = urlFinal;
      link.setAttribute('download', 'Taxas_SEJUSP.xlsx');
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);

    } catch (error) {
      console.error("Erro ao exportar Excel:", error);
    }
  };

  return (
    <Layout style={{ minHeight: '80vh' }}>
      <div style={{ marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0, color: '#004f9f' }}>Controle de Arrecadação - SEJUSP</Title>
        <Text type="secondary">Utilize os campos abaixo para filtrar dados específicos da base de dados</Text>
      </div>

      <Card style={{ marginBottom: 16, borderRadius: 8, boxShadow: '0 2px 8px rgba(0,0,0,0.05)' }}>
        <Form 
          form={form} 
          layout="inline" 
          onFinish={handleFiltrar} 
          style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}
        >
          <Form.Item name="itemSubItem" label="Item/SubItem">
            <Input placeholder="Ex: 01.09" style={{ width: 120 }} allowClear />
          </Form.Item>
          
          <Form.Item name="descricao" label="Descrição">
            <Input placeholder="Palavra-chave na descrição" style={{ width: 250 }} allowClear />
          </Form.Item>

          <Form.Item name="tributo" label="Tributo">
            <Input placeholder="Ex: 512" style={{ width: 100 }} allowClear />
          </Form.Item>

          <Form.Item 
            name="referencia" 
            label="Referência"
            getValueFromEvent={(e) => {
              const val = e.target.value.replace(/\D/g, '').substring(0, 6);
              if (val.length <= 2) return val;
              return `${val.substring(0, 2)}/${val.substring(2)}`;
            }}
          >
            <Input placeholder="MM/YYYY" style={{ width: 120 }} allowClear maxLength={7} />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" icon={<SearchOutlined />} htmlType="submit">
                Filtrar
              </Button>
              <Button onClick={handleLimpar}>
                Limpar
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>

      <Content>
        <Card styles={{ body: { padding: 0 } }} style={{ borderRadius: 8, overflow: 'hidden' }}>
          <div style={{ padding: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', backgroundColor: '#fafafa', borderBottom: '1px solid #f0f0f0' }}>
            <Title level={5} style={{ margin: 0 }}>Listagem de Taxas</Title>
            <Space>
              <Button 
                icon={<FileExcelOutlined />} 
                onClick={handleExportarExcel}
                style={{ backgroundColor: '#217346', color: '#fff', border: 'none' }}
              >
                Excel
              </Button>
            </Space>
          </div>
          <div style={{ padding: '0 16px 16px 16px' }}>
            <TabelaTaxas filtros={filtros} />
          </div>
        </Card>
      </Content>
    </Layout>
  );
};

export default Dashboard;