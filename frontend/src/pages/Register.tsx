import React from 'react';
import { Form, Input, Button, Card, Typography, message, Divider } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons';
import { Link, useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/stores/authStore';
import { ROUTES, VALIDATION_RULES, ERROR_MESSAGES } from '@/constants';
import type { RegisterRequest } from '@/types';

const { Title, Text } = Typography;

const Register: React.FC = () => {
  const navigate = useNavigate();
  const { register, isLoading, error, clearError } = useAuthStore();
  const [form] = Form.useForm();

  const handleSubmit = async (values: RegisterRequest) => {
    try {
      clearError();
      await register(values);
      message.success('注册成功，正在为您登录...');
      navigate(ROUTES.DASHBOARD);
    } catch (err: unknown) {
      const error = err as Error;
      message.error(error.message || ERROR_MESSAGES.REGISTER_FAILED);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-green-50 to-emerald-100 flex items-center justify-center p-4">
      <Card className="w-full max-w-md shadow-lg">
        <div className="text-center mb-8">
          <Title level={2} className="text-gray-800 mb-2">
            创建账户
          </Title>
          <Text type="secondary">
            请填写以下信息完成注册
          </Text>
        </div>

        <Form
          form={form}
          name="register"
          onFinish={handleSubmit}
          layout="vertical"
          size="large"
        >
          <Form.Item
            name="username"
            label="用户名"
            rules={[
              { required: true, message: '请输入用户名' },
              {
                min: VALIDATION_RULES.USERNAME.MIN_LENGTH,
                message: `用户名至少${VALIDATION_RULES.USERNAME.MIN_LENGTH}个字符`,
              },
              {
                max: VALIDATION_RULES.USERNAME.MAX_LENGTH,
                message: `用户名最多${VALIDATION_RULES.USERNAME.MAX_LENGTH}个字符`,
              },
              {
                pattern: VALIDATION_RULES.USERNAME.PATTERN,
                message: '用户名只能包含字母、数字和下划线',
              },
            ]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="请输入用户名"
              autoComplete="username"
            />
          </Form.Item>

          <Form.Item
            name="name"
            label="姓名"
            rules={[
              { required: true, message: '请输入姓名' },
              { min: 2, message: '姓名至少2个字符' },
              { max: 50, message: '姓名最多50个字符' },
            ]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="请输入真实姓名"
              autoComplete="name"
            />
          </Form.Item>

          <Form.Item
            name="email"
            label="邮箱"
            rules={[
              { required: true, message: '请输入邮箱地址' },
              {
                pattern: VALIDATION_RULES.EMAIL.PATTERN,
                message: '请输入有效的邮箱地址',
              },
            ]}
          >
            <Input
              prefix={<MailOutlined />}
              placeholder="请输入邮箱地址"
              autoComplete="email"
            />
          </Form.Item>

          <Form.Item
            name="phone"
            label="手机号（可选）"
            rules={[
              {
                pattern: /^1[3-9]\d{9}$/,
                message: '请输入有效的手机号码',
              },
            ]}
          >
            <Input
              prefix={<PhoneOutlined />}
              placeholder="请输入手机号码"
              autoComplete="tel"
            />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[
              { required: true, message: '请输入密码' },
              {
                min: VALIDATION_RULES.PASSWORD.MIN_LENGTH,
                message: `密码至少${VALIDATION_RULES.PASSWORD.MIN_LENGTH}个字符`,
              },
              {
                max: VALIDATION_RULES.PASSWORD.MAX_LENGTH,
                message: `密码最多${VALIDATION_RULES.PASSWORD.MAX_LENGTH}个字符`,
              },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="请输入密码"
              autoComplete="new-password"
            />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            label="确认密码"
            dependencies={['password']}
            rules={[
              { required: true, message: '请确认密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('password') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'));
                },
              }),
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="请再次输入密码"
              autoComplete="new-password"
            />
          </Form.Item>

          {error && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
              <Text type="danger">{error}</Text>
            </div>
          )}

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={isLoading}
              className="w-full h-12 text-base font-medium"
            >
              {isLoading ? '注册中...' : '注册'}
            </Button>
          </Form.Item>
        </Form>

        <Divider plain>
          <Text type="secondary">已有账户？</Text>
        </Divider>

        <div className="text-center">
          <Link to={ROUTES.LOGIN}>
            <Button type="link" size="large" className="font-medium">
              立即登录
            </Button>
          </Link>
        </div>
      </Card>
    </div>
  );
};

export default Register;
