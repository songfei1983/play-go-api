import React, { useState } from 'react';
import {
  Card,
  Form,
  Input,
  Button,
  Select,
  Space,
  Avatar,
  Typography,
  Row,
  Col,
  Spin,
  message,
} from 'antd';
import {
  UserOutlined,
  SaveOutlined,
  ArrowLeftOutlined,
  EditOutlined,
} from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate, useParams } from 'react-router-dom';
import { userService } from '@/services/userService';
import { ROUTES, USER_STATUS, FORM_RULES } from '@/constants';
import type { UpdateUserRequest } from '@/types';

const { Title, Text } = Typography;
const { Option } = Select;

const UserDetail: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const queryClient = useQueryClient();
  const [form] = Form.useForm();
  const [isEditing, setIsEditing] = useState(false);

  // 获取用户详情
  const { data: user, isLoading } = useQuery({
    queryKey: ['user', id],
    queryFn: () => userService.getUser(id!),
    enabled: !!id,
  });

  // 更新用户
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateUserRequest }) =>
      userService.updateUser(id, data),
    onSuccess: () => {
      message.success('用户信息更新成功');
      setIsEditing(false);
      queryClient.invalidateQueries({ queryKey: ['user', id] });
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
    onError: (error: unknown) => {
      const err = error as Error;
      message.error(err.message || '更新失败');
    },
  });

  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      updateMutation.mutate({ id: id!, data: values });
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  const handleCancel = () => {
    form.setFieldsValue(user?.data);
    setIsEditing(false);
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spin size="large" />
      </div>
    );
  }

  if (!user?.data) {
    return (
      <div className="text-center py-8">
        <Text type="secondary">用户不存在</Text>
      </div>
    );
  }

  const userData = user.data;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <Button
            icon={<ArrowLeftOutlined />}
            onClick={() => navigate(ROUTES.USERS)}
          >
            返回
          </Button>
          <Title level={2} className="mb-0">
            用户详情
          </Title>
        </div>
        <Space>
          {!isEditing ? (
            <Button
              type="primary"
              icon={<EditOutlined />}
              onClick={() => setIsEditing(true)}
            >
              编辑
            </Button>
          ) : (
            <>
              <Button onClick={handleCancel}>取消</Button>
              <Button
                type="primary"
                icon={<SaveOutlined />}
                loading={updateMutation.isPending}
                onClick={handleSave}
              >
                保存
              </Button>
            </>
          )}
        </Space>
      </div>

      <Row gutter={24}>
        <Col xs={24} lg={8}>
          <Card>
            <div className="text-center space-y-4">
              <Avatar size={80} icon={<UserOutlined />} />
              <div>
                <Title level={4} className="mb-1">
                  {userData.name || userData.username}
                </Title>
                <Text type="secondary">{userData.email}</Text>
              </div>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <Text type="secondary">状态:</Text>
                  <Text
                    className={`font-medium ${
                      userData.status === 'active'
                        ? 'text-green-600'
                        : 'text-red-600'
                    }`}
                  >
                    {userData.status === 'active' ? '活跃' : '非活跃'}
                  </Text>
                </div>
                <div className="flex justify-between">
                  <Text type="secondary">创建时间:</Text>
                  <Text>{new Date(userData.created_at).toLocaleDateString('zh-CN')}</Text>
                </div>
                <div className="flex justify-between">
                  <Text type="secondary">更新时间:</Text>
                  <Text>{new Date(userData.updated_at).toLocaleDateString('zh-CN')}</Text>
                </div>
              </div>
            </div>
          </Card>
        </Col>

        <Col xs={24} lg={16}>
          <Card title="基本信息">
            <Form
              form={form}
              layout="vertical"
              initialValues={userData}
              disabled={!isEditing}
            >
              <Row gutter={16}>
                <Col xs={24} sm={12}>
                  <Form.Item
                    label="用户名"
                    name="username"
                    rules={FORM_RULES.USERNAME}
                  >
                    <Input placeholder="请输入用户名" />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12}>
                  <Form.Item
                    label="姓名"
                    name="name"
                    rules={FORM_RULES.NAME}
                  >
                    <Input placeholder="请输入姓名" />
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16}>
                <Col xs={24} sm={12}>
                  <Form.Item
                    label="邮箱"
                    name="email"
                    rules={FORM_RULES.EMAIL}
                  >
                    <Input placeholder="请输入邮箱" />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12}>
                  <Form.Item
                    label="手机号"
                    name="phone"
                    rules={FORM_RULES.PHONE}
                  >
                    <Input placeholder="请输入手机号" />
                  </Form.Item>
                </Col>
              </Row>

              <Form.Item label="状态" name="status">
                <Select placeholder="请选择状态">
                  <Option value={USER_STATUS.ACTIVE}>活跃</Option>
                  <Option value={USER_STATUS.INACTIVE}>非活跃</Option>
                </Select>
              </Form.Item>
            </Form>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default UserDetail;
