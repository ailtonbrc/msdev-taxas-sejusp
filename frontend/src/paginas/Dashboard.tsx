import React, { useState, useEffect } from 'react';
import { Layout, Card, Form, Button, Typography, Space, Input, Modal } from 'antd';
import { FileExcelOutlined, MailOutlined, WarningOutlined } from '@ant-design/icons';
import { useOutletContext } from 'react-router-dom';
import TabelaTaxas from '../componentes/TabelaTaxas';
import SelectTableau from '../componentes/SelectTableau';
import api, { getOpcoesFiltros } from '../servicos/api';
import type { IOpcoesFiltros, IFiltrosState } from '../tipos/taxas';

const { Content } = Layout;
const { Text, Title, Paragraph } = Typography;

interface UsuarioContext {
  usuario: {
    nome: string;
    role: string;
    instituicoes: string[];
  };
}

const Dashboard: React.FC = () => {
  const { usuario } = useOutletContext<UsuarioContext>();
  const isGuest = usuario.role === 'guest';
  const [form] = Form.useForm();
  
  // Modal de solicitação de acesso
  const [showAccessModal, setShowAccessModal] = useState(isGuest);

  const [opcoes, setOpcoes] = useState<IOpcoesFiltros>({
    instituicoes: [],
    itens: [],
    tributos: [],
    municipios: [],
    referencias: [],
    ultimaAtualizacao: ''
  });

  const [filtros, setFiltros] = useState<IFiltrosState>({ 
    itemSubItem: [],
    descricao: '',
    tributo: '',
    referencia: [],
    municipio: [],
    instituicao: []
  });

  const [loadingFiltros, setLoadingFiltros] = useState(false);

  const carregarOpcoes = async (filtrosAtuais: IFiltrosState) => {
    if (isGuest) return;
    setLoadingFiltros(true);
    try {
      const data = await getOpcoesFiltros(filtrosAtuais);
      setOpcoes(data);
    } catch (error) {
      console.error('Erro ao carregar opções de filtros:', error);
    } finally {
      setLoadingFiltros(false);
    }
  };

  useEffect(() => {
    carregarOpcoes(filtros);
  }, [filtros]);

  const handleFiltroChange = (values: any) => {
    if (isGuest) return;
    setFiltros({
      itemSubItem: values.itemSubItem || [],
      descricao: values.descricao || '',
      tributo: values.tributo || '',
      referencia: values.referencia || [],
      municipio: values.municipio || [],
      instituicao: values.instituicao || []
    });
  };

  const handleExportExcel = () => {
    if (isGuest) return;
    const urlFinal = api.getUri({ url: '/taxas/exportar', params: filtros as any });
    window.open(urlFinal, '_blank');
  };

  const handleExportCSV = () => {
    if (isGuest) return;
    const urlFinal = api.getUri({ url: '/taxas/exportar-csv', params: filtros as any });
    window.open(urlFinal, '_blank');
  };

  return (
    <Layout style={{ background: '#fff' }}>
      <Content>
        {/* Modal Informativo de Acesso */}
        <Modal
          title={
            <Space>
              <WarningOutlined style={{ color: '#faad14' }} />
              <span>Acesso Restrito - Solicitação de Permissão</span>
            </Space>
          }
          
          open={showAccessModal}
          onCancel={() => setShowAccessModal(false)}
          footer={[
            <Button key="ok" type="primary" onClick={() => setShowAccessModal(false)}>
              Entendi
            </Button>
          ]}
          centered
          width={650}
        >
          <div style={{ textAlign: 'center', padding: '10px 0' }}>
            <Title level={4}>Olá! Bem-vindo ao Sistema de Taxas SEJUSP.</Title>
            <Paragraph style={{ fontSize: '16px', color: '#555' }}>
              Identificamos sua autenticação via e-Fazenda, porém seu CPF ainda não possui 
              permissões cadastradas para visualizar os dados de arrecadação.
            </Paragraph>
            
            <Card style={{ backgroundColor: '#f0f5ff', border: '1px solid #adc6ff', borderRadius: 12, marginTop: 24 }}>
              <Paragraph strong style={{ fontSize: '16px', marginBottom: 16 }}>
                Como solicitar a liberação do acesso?
              </Paragraph>
              <Paragraph>
                Para habilitar a visualização dos dados da sua vertente (PM, Bombeiros ou Polícia Civil), 
                por favor encaminhe uma solicitação para o e-mail:
              </Paragraph>
              <Title level={4} style={{ color: '#004f9f', margin: '20px 0' }}>
                <MailOutlined style={{ marginRight: 10 }} />
                parle_tableau@fazenda.ms.gov.br
              </Title>
              <Paragraph type="secondary" style={{ fontSize: '13px', fontStyle: 'italic' }}>
                No e-mail, informe seu Nome Completo, CPF e a Instituição à qual você pertence.
              </Paragraph>
            </Card>
          </div>
        </Modal>

        <div style={{ marginBottom: 16 }}>
          <Title level={4} style={{ margin: 0, color: '#004f9f' }}>Dashboard de Arrecadação</Title>
          <Text type="secondary">Monitoramento de taxas e custas estaduais</Text>
        </div>

        <Card style={{ marginBottom: 16, borderRadius: 8, boxShadow: '0 2px 8px rgba(0,0,0,0.05)' }}>
          <Form 
            form={form} 
            layout="vertical" 
            onValuesChange={(_, values) => handleFiltroChange(values)}
            style={{ display: 'flex', flexWrap: 'wrap', gap: '12px', alignItems: 'flex-start' }}
          >
            <Form.Item name="itemSubItem" label="Item/SubItem" style={{ marginBottom: 0, flex: '1 1 400px' }}>
              <SelectTableau 
                mode="multiple"
                placeholder={isGuest ? "Bloqueado - Sem permissão" : "Selecione os itens..."} 
                loading={loadingFiltros}
                options={opcoes.itens}
                disabled={isGuest}
              />
            </Form.Item>

            <Form.Item name="instituicao" label="Instituição" style={{ marginBottom: 0, flex: '0 0 180px' }}>
              <SelectTableau 
                mode="multiple"
                placeholder="Instituição" 
                loading={loadingFiltros}
                options={opcoes.instituicoes}
                disabled={isGuest}
              />
            </Form.Item>

            <Form.Item name="descricao" label="Descrição" style={{ marginBottom: 0, flex: '0 0 200px' }}>
              <Input 
                placeholder="Busca livre..." 
                disabled={isGuest}
              />
            </Form.Item>

            <Form.Item name="tributo" label="Tributo" style={{ marginBottom: 0, flex: '0 0 120px' }}>
              <SelectTableau 
                mode="multiple"
                placeholder="Código" 
                loading={loadingFiltros}
                options={opcoes.tributos}
                disabled={isGuest}
              />
            </Form.Item>

            <Form.Item name="referencia" label="Referência" style={{ marginBottom: 0, flex: '0 0 120px' }}>
              <SelectTableau 
                mode="multiple"
                placeholder="Mês/Ano" 
                loading={loadingFiltros}
                options={opcoes.referencias}
                disabled={isGuest}
              />
            </Form.Item>

            <Form.Item label=" " style={{ marginBottom: 0, flex: '0 0 auto', alignSelf: 'flex-end' }}>
              <Button onClick={() => { form.resetFields(); handleFiltroChange({}); }} disabled={isGuest}>Limpar</Button>
            </Form.Item>
          </Form>
        </Card>

        <Card styles={{ body: { padding: 0 } }} style={{ borderRadius: 8, overflow: 'hidden' }}>
          <div style={{ padding: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', backgroundColor: '#fafafa', borderBottom: '1px solid #f0f0f0' }}>
            <Title level={5} style={{ margin: 0 }}>Listagem de Taxas</Title>
            <Space>
              <Button 
                icon={<FileExcelOutlined />} 
                onClick={handleExportExcel}
                disabled={isGuest}
                style={{ backgroundColor: isGuest ? '#ccc' : '#1d6f42', borderColor: isGuest ? '#ccc' : '#1d6f42', color: '#fff' }}
              >
                Excel
              </Button>
              <Button 
                icon={<FileExcelOutlined />} 
                onClick={handleExportCSV}
                disabled={isGuest}
                style={{ backgroundColor: isGuest ? '#ccc' : '#555', borderColor: isGuest ? '#ccc' : '#555', color: '#fff' }}
              >
                CSV
              </Button>
            </Space>
          </div>
          <div style={{ padding: '0 16px 16px 16px' }}>
            <TabelaTaxas filtros={filtros} isGuest={isGuest} />
          </div>
        </Card>
        
        {opcoes.ultimaAtualizacao && (
          <div style={{ textAlign: 'right', marginTop: 16 }}>
            <Text type="secondary" style={{ fontSize: 11 }}>
              Última sincronização: {opcoes.ultimaAtualizacao}
            </Text>
          </div>
        )}
      </Content>
    </Layout>
  );
};

export default Dashboard;