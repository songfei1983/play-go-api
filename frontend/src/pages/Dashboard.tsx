import React, { useEffect, useState } from 'react';
import {
  Card,
  Row,
  Col,
  Statistic,
  Typography,
  Button,
  List,
  Avatar,
  Tag,
  Space,
  Spin,
} from 'antd';
import {
  UserOutlined,
  ShoppingOutlined,
  DollarOutlined,
  ArrowUpOutlined,
  PlusOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { userService } from '@/services/userService';
import { productService } from '@/services/productService';
import { ROUTES } from '@/constants';
import type { User, Product } from '@/types';

const { Title, Text } = Typography;

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const [stats, setStats] = useState({
    totalUsers: 0,
    totalProducts: 0,
    totalRevenue: 0,
    activeUsers: 0,
  });

  // 获取用户列表
  const { data: usersData, isLoading: usersLoading } = useQuery({
    queryKey: ['users', { limit: 5 }],
    queryFn: () => userService.getUsers({ limit: 5 }),
  });

  // 获取产品列表
  const { data: productsData, isLoading: productsLoading } = useQuery({
    queryKey: ['products', { limit: 5 }],
    queryFn: () => productService.getProducts({ limit: 5 }),
  });

  useEffect(() => {
    // 计算统计数据
    if (usersData?.data) {
      const users = usersData.data.data;
      const activeUsers = users.filter(user => user.status === 'active').length;
      setStats(prev => ({
        ...prev,
        totalUsers: usersData.data.total,
        activeUsers,
      }));
    }

    if (productsData?.data) {
      const products = productsData.data.data;
      const totalRevenue = products.reduce((sum, product) => sum + (product.price * product.stock), 0);
      setStats(prev => ({
        ...prev,
        totalProducts: productsData.data.total,
        totalRevenue,
      }));
    }
  }, [usersData, productsData]);

  const recentUsers = usersData?.data?.data || [];
  const recentProducts = productsData?.data?.data || [];

  const quickActions = [
    {
      title: '创建产品',
      description: '添加新的产品到系统',
      icon: <PlusOutlined />,
      action: () => navigate(ROUTES.PRODUCT_CREATE),
      color: '#1890ff',
    },
    {
      title: '查看用户',
      description: '管理系统用户',
      icon: <UserOutlined />,
      action: () => navigate(ROUTES.USERS),
      color: '#52c41a',
    },
    {
      title: '产品管理',
      description: '查看和管理产品',
      icon: <ShoppingOutlined />,
      action: () => navigate(ROUTES.PRODUCTS),
      color: '#722ed1',
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <Title level={2} className="mb-0">
          仪表盘
        </Title>
        <Text type="secondary">
          欢迎回到 Play Go API 管理系统
        </Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="总用户数"
              value={stats.totalUsers}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="活跃用户"
              value={stats.activeUsers}
              prefix={<ArrowUpOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="总产品数"
              value={stats.totalProducts}
              prefix={<ShoppingOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="总价值"
              value={stats.totalRevenue}
              prefix={<DollarOutlined />}
              precision={2}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        {/* 快速操作 */}
        <Col xs={24} lg={8}>
          <Card title="快速操作" className="h-full">
            <div className="space-y-4">
              {quickActions.map((action, index) => (
                <Card
                  key={index}
                  size="small"
                  hoverable
                  onClick={action.action}
                  className="cursor-pointer border-l-4"
                  style={{ borderLeftColor: action.color }}
                >
                  <div className="flex items-center space-x-3">
                    <div
                      className="w-10 h-10 rounded-lg flex items-center justify-center text-white"
                      style={{ backgroundColor: action.color }}
                    >
                      {action.icon}
                    </div>
                    <div className="flex-1">
                      <Text className="font-medium block">{action.title}</Text>
                      <Text type="secondary" className="text-sm">
                        {action.description}
                      </Text>
                    </div>
                  </div>
                </Card>
              ))}
            </div>
          </Card>
        </Col>

        {/* 最近用户 */}
        <Col xs={24} lg={8}>
          <Card
            title="最近用户"
            extra={
              <Button
                type="link"
                size="small"
                icon={<EyeOutlined />}
                onClick={() => navigate(ROUTES.USERS)}
              >
                查看全部
              </Button>
            }
            className="h-full"
          >
            {usersLoading ? (
              <div className="flex justify-center py-8">
                <Spin />
              </div>
            ) : (
              <List
                dataSource={recentUsers}
                renderItem={(user: User) => (
                  <List.Item>
                    <List.Item.Meta
                      avatar={<Avatar icon={<UserOutlined />} />}
                      title={user.name || user.username}
                      description={
                        <Space>
                          <Text type="secondary">{user.email}</Text>
                          <Tag color={user.status === 'active' ? 'green' : 'red'}>
                            {user.status === 'active' ? '活跃' : '非活跃'}
                          </Tag>
                        </Space>
                      }
                    />
                  </List.Item>
                )}
              />
            )}
          </Card>
        </Col>

        {/* 最近产品 */}
        <Col xs={24} lg={8}>
          <Card
            title="最近产品"
            extra={
              <Button
                type="link"
                size="small"
                icon={<EyeOutlined />}
                onClick={() => navigate(ROUTES.PRODUCTS)}
              >
                查看全部
              </Button>
            }
            className="h-full"
          >
            {productsLoading ? (
              <div className="flex justify-center py-8">
                <Spin />
              </div>
            ) : (
              <List
                dataSource={recentProducts}
                renderItem={(product: Product) => (
                  <List.Item>
                    <List.Item.Meta
                      avatar={
                        <Avatar
                          shape="square"
                          size="large"
                          icon={<ShoppingOutlined />}
                        />
                      }
                      title={product.name}
                      description={
                        <Space>
                          <Text type="secondary">¥{product.price}</Text>
                          <Text type="secondary">库存: {product.stock}</Text>
                          <Tag color={product.status === 'active' ? 'green' : 'red'}>
                            {product.status === 'active' ? '上架' : '下架'}
                          </Tag>
                        </Space>
                      }
                    />
                  </List.Item>
                )}
              />
            )}
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
