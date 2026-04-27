import React from 'react';
import { Table, Card, Button, Modal, Form, Input, Select, Tag, Space, message, Popconfirm, Typography } from 'antd';
import { UserAddOutlined, EditOutlined, DeleteOutlined, SafetyCertificateOutlined } from '@ant-design/icons';
import api from '../servicos/api';

const { Title } = Typography;
const { Option } = Select;

interface Usuario {
  id: number;
  nome: string;
  cpf: string;
  role: string;
  instituicoes: string[];
}

const Usuarios: React.FC = () => {
  const [usuarios, setUsuarios] = React.useState<Usuario[]>([]);
  const [loading, setLoading] = React.useState(false);
  const [modalVisible, setModalVisible] = React.useState(false);
  const [form] = Form.useForm();
  const [editingId, setEditingId] = React.useState<number | null>(null);

  const carregarUsuarios = async () => {
    setLoading(true);
    try {
      const response = await api.get('/usuarios/');
      setUsuarios(response.data);
    } catch (error) {
      message.error('Falha ao carregar usuários');
    } finally {
      setLoading(false);
    }
  };

  React.useEffect(() => {
    carregarUsuarios();
  }, []);

  const handleAdd = () => {
    setEditingId(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (record: Usuario) => {
    setEditingId(record.id);
    form.setFieldsValue(record);
    setModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/usuarios/${id}`);
      message.success('Usuário removido');
      carregarUsuarios();
    } catch (error) {
      message.error('Falha ao remover usuário');
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      const payload = editingId ? { ...values, id: editingId } : values;

      await api.post('/usuarios/', payload);
      message.success(editingId ? 'Usuário atualizado' : 'Usuário cadastrado');
      setModalVisible(false);
      carregarUsuarios();
    } catch (error) {
      message.error('Falha ao salvar usuário');
    }
  };

  const columns = [
    {
      title: 'Nome',
      dataIndex: 'nome',
      key: 'nome',
      render: (text: string, record: Usuario) => (
        <Space>
          {text}
          {record.role === 'admin' && <Tag color="gold" icon={<SafetyCertificateOutlined />}>ADMIN</Tag>}
        </Space>
      )
    },

    {
      title: 'CPF',
      dataIndex: 'cpf',
      key: 'cpf',
    },
    {
      title: 'Instituições Permitidas',
      dataIndex: 'instituicoes',
      key: 'instituicoes',
      render: (insts: string[]) => (
        <>
          {insts && insts.map(inst => (
            <Tag color="blue" key={inst}>{inst}</Tag>
          ))}
          {(!insts || insts.length === 0) && <Tag color="default">NENHUMA</Tag>}
        </>
      )
    },
    {
      title: 'Ações',
      key: 'acoes',
      width: 150,
      render: (_: any, record: Usuario) => (
        <Space size="middle">
          <Button icon={<EditOutlined />} onClick={() => handleEdit(record)} />
          <Popconfirm title="Tem certeza?" onConfirm={() => handleDelete(record.id)}>
            <Button danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Card
        title={
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Title level={4} style={{ margin: 0 }}>Gestão de Acessos</Title>
            <Button type="primary" icon={<UserAddOutlined />} onClick={handleAdd}>
              Novo Usuário
            </Button>
          </div>
        }
      >
        <Table
          columns={columns}
          dataSource={usuarios}
          rowKey="id"
          loading={loading}
          pagination={{ pageSize: 10 }}
        />
      </Card>

      <Modal
        title={editingId ? 'Editar Usuário' : 'Cadastrar Usuário'}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
        okText="Salvar"
        cancelText="Cancelar"
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="nome" label="Nome Completo" rules={[{ required: true }]}>
            <Input placeholder="Digite o nome" />
          </Form.Item>
          <Form.Item name="cpf" label="CPF" rules={[{ required: true }]}>
            <Input placeholder="Apenas números" />
          </Form.Item>
          <Form.Item name="role" label="Perfil de Acesso" rules={[{ required: true }]} initialValue="user">
            <Select>
              <Option value="user">Usuário</Option>
              <Option value="admin">Administrador</Option>
            </Select>
          </Form.Item>
          <Form.Item name="instituicoes" label="Vertentes Permitidas">
            <Select mode="multiple" placeholder="Selecione as polícias/instituições">
              <Option value="POLICIA CIVIL">POLICIA CIVIL</Option>
              <Option value="BOMBEIRO MILITAR">BOMBEIRO MILITAR</Option>
              <Option value="POLICIA MILITAR">POLICIA MILITAR</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Usuarios;
