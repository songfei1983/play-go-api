import React from 'react';
import { Form, Input, Button, Card, Typography, message, Divider } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { Link, useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/stores/authStore';
import { ROUTES, VALIDATION_RULES, ERROR_MESSAGES } from '@/constants';
import type { LoginRequest } from '@/types';

const { Title, Text } = Typography;

const Login: React.FC = () => {
  const navigate = useNavigate();
  const { login, isLoading, error, clearError } = useAuthStore();
  const [form] = Form.useForm();

  const handleSubmit = async (values: LoginRequest) => {
    try {
      clearError();
      await login(values);
      message.success('登录成功');
      navigate(ROUTES.DASHBOARD);
    } catch (err: unknown) {
      const error = err as Error;
      message.error(error.message || ERROR_MESSAGES.LOGIN_FAILED);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <Card className="w-full max-w-md shadow-lg">
        <div className="text-center mb-8">
          <Title level={2} className="text-gray-800 mb-2">
            欢迎回来
          </Title>
          <Text type="secondary">
            请登录您的账户
          </Text>
        </div>

        <Form
          form={form}
          name="login"
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
            name="password"
            label="密码"
            rules={[
              { required: true, message: '请输入密码' },
              {
                min: VALIDATION_RULES.PASSWORD.MIN_LENGTH,
                message: `密码至少${VALIDATION_RULES.PASSWORD.MIN_LENGTH}个字符`,
              },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="请输入密码"
              autoComplete="current-password"
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
              {isLoading ? '登录中...' : '登录'}
            </Button>
          </Form.Item>
        </Form>

        <Divider plain>
          <Text type="secondary">还没有账户？</Text>
        </Divider>

        <div className="text-center">
          <Link to={ROUTES.REGISTER}>
            <Button type="link" size="large" className="font-medium">
              立即注册
            </Button>
          </Link>
        </div>
      </Card>
    </div>
  );
};

export default Login;
