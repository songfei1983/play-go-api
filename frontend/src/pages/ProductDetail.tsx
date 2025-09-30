import React from 'react';
import {
  Card,
  Button,
  Typography,
  Row,
  Col,
  Spin,
  Tag,
  Image,
  Descriptions,
  Space,
} from 'antd';
import {
  ArrowLeftOutlined,
  EditOutlined,
} from '@ant-design/icons';
import { useQuery } from '@tanstack/react-query';
import { useNavigate, useParams } from 'react-router-dom';
import { productService } from '@/services/productService';
import { ROUTES } from '@/constants';

const { Title, Text, Paragraph } = Typography;

const ProductDetail: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();

  // 获取产品详情
  const { data: product, isLoading } = useQuery({
    queryKey: ['product', id],
    queryFn: () => productService.getProduct(id!),
    enabled: !!id,
  });

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('zh-CN', {
      style: 'currency',
      currency: 'CNY',
    }).format(price);
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spin size="large" />
      </div>
    );
  }

  if (!product?.data) {
    return (
      <div className="text-center py-8">
        <Text type="secondary">产品不存在</Text>
      </div>
    );
  }

  const productData = product.data;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <Button
            icon={<ArrowLeftOutlined />}
            onClick={() => navigate(ROUTES.PRODUCTS)}
          >
            返回
          </Button>
          <Title level={2} className="mb-0">
            产品详情
          </Title>
        </div>
        <Space>
          <Button
            type="primary"
            icon={<EditOutlined />}
            onClick={() => navigate(ROUTES.PRODUCT_EDIT.replace(':id', productData.id))}
          >
            编辑产品
          </Button>
        </Space>
      </div>

      <Row gutter={24}>
        <Col xs={24} lg={10}>
          <Card>
            <div className="text-center space-y-4">
              <Image
                width="100%"
                height={300}
                src={productData.image_url}
                fallback="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAMIAAADDCAYAAADQvc6UAAABRWlDQ1BJQ0MgUHJvZmlsZQAAKJFjYGASSSwoyGFhYGDIzSspCnJ3UoiIjFJgf8LAwSDCIMogwMCcmFxc4BgQ4ANUwgCjUcG3awyMIPqyLsis7PPOq3QdDFcvjV3jOD1boQVTPQrgSkktTgbSf4A4LbmgqISBgTEFyFYuLykAsTuAbJEioKOA7DkgdjqEvQHEToKwj4DVhAQ5A9k3gGyB5IxEoBmML4BsnSQk8XQkNtReEOBxcfXxUQg1Mjc0dyHgXNJBSWpFCYh2zi+oLMpMzyhRcASGUqqCZ16yno6CkYGRAQMDKMwhqj/fAIcloxgHQqxAjIHBEugw5sUIsSQpBobtQPdLciLEVJYzMPBHMDBsayhILEqEO4DxG0txmrERhM29nYGBddr//5/DGRjYNRkY/l7////39v///y4Dmn+LgeHANwDrkl1AuO+pmgAAADhlWElmTU0AKgAAAAgAAYdpAAQAAAABAAAAGgAAAAAAAqACAAQAAAABAAAAwqADAAQAAAABAAAAwwAAAAD9b/HnAAAHlklEQVR4Ae3dP3Ik1RnG4W+FgYxN"
                className="rounded object-cover"
              />
              <div className="space-y-2">
                <div className="flex items-center justify-center space-x-2">
                  <Text className="text-3xl font-bold text-green-600">
                    {formatPrice(productData.price)}
                  </Text>
                  <Tag color={productData.status === 'active' ? 'green' : 'red'}>
                    {productData.status === 'active' ? '上架' : '下架'}
                  </Tag>
                </div>
                <div className="flex items-center justify-center space-x-4">
                  <Text type="secondary">库存: </Text>
                  <Text
                    className={`font-medium ${
                      productData.stock > 10 ? 'text-green-600' : 'text-red-600'
                    }`}
                  >
                    {productData.stock} 件
                  </Text>
                </div>
              </div>
            </div>
          </Card>
        </Col>

        <Col xs={24} lg={14}>
          <Card title="产品信息">
            <div className="space-y-6">
              <div>
                <Title level={3} className="mb-2">
                  {productData.name}
                </Title>
                <Paragraph className="text-gray-600">
                  {productData.description}
                </Paragraph>
              </div>

              <Descriptions column={1} bordered>
                <Descriptions.Item label="产品ID">
                  {productData.id}
                </Descriptions.Item>
                <Descriptions.Item label="产品名称">
                  {productData.name}
                </Descriptions.Item>
                <Descriptions.Item label="价格">
                  <span className="text-green-600 font-medium">
                    {formatPrice(productData.price)}
                  </span>
                </Descriptions.Item>
                <Descriptions.Item label="库存数量">
                  <span
                    className={`font-medium ${
                      productData.stock > 10 ? 'text-green-600' : 'text-red-600'
                    }`}
                  >
                    {productData.stock} 件
                  </span>
                </Descriptions.Item>
                <Descriptions.Item label="状态">
                  <Tag color={productData.status === 'active' ? 'green' : 'red'}>
                    {productData.status === 'active' ? '上架' : '下架'}
                  </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="创建时间">
                  {new Date(productData.created_at).toLocaleString('zh-CN')}
                </Descriptions.Item>
                <Descriptions.Item label="更新时间">
                  {new Date(productData.updated_at).toLocaleString('zh-CN')}
                </Descriptions.Item>
              </Descriptions>
            </div>
          </Card>
        </Col>
      </Row>

      {productData.description && (
        <Card title="产品描述">
          <Paragraph>
            {productData.description}
          </Paragraph>
        </Card>
      )}
    </div>
  );
};

export default ProductDetail;
