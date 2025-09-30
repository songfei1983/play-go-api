import React, { useState } from 'react';
import {
  Card,
  Form,
  Input,
  Button,
  Select,
  InputNumber,
  Upload,
  Typography,
  Row,
  Col,
  message,
  Space,
} from 'antd';
import {
  ArrowLeftOutlined,
  SaveOutlined,
  UploadOutlined,
} from '@ant-design/icons';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { productService } from '@/services/productService';
import { ROUTES } from '@/constants';
import type { CreateProductRequest } from '@/types';
import type { UploadFile } from 'antd/es/upload/interface';

const { Title } = Typography;
const { TextArea } = Input;
const { Option } = Select;

const ProductCreate: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [form] = Form.useForm();
  const [fileList, setFileList] = useState<UploadFile[]>([]);

  // 创建产品
  const createMutation = useMutation({
    mutationFn: productService.createProduct,
    onSuccess: () => {
      message.success('产品创建成功');
      queryClient.invalidateQueries({ queryKey: ['products'] });
      navigate(ROUTES.PRODUCTS);
    },
    onError: (error: unknown) => {
      const err = error as Error;
      message.error(err.message || '创建失败');
    },
  });

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const productData: CreateProductRequest = {
        ...values,
        image_url: fileList.length > 0 ? fileList[0].url || fileList[0].thumbUrl : undefined,
      };
      createMutation.mutate(productData);
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  const handleUploadChange = ({ fileList: newFileList }: { fileList: UploadFile[] }) => {
    setFileList(newFileList);
  };

  const uploadProps = {
    name: 'file',
    listType: 'picture-card' as const,
    fileList,
    onChange: handleUploadChange,
    beforeUpload: (file: File) => {
      const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
      if (!isJpgOrPng) {
        message.error('只能上传 JPG/PNG 格式的图片!');
        return false;
      }
      const isLt2M = file.size / 1024 / 1024 < 2;
      if (!isLt2M) {
        message.error('图片大小不能超过 2MB!');
        return false;
      }

      // 创建预览URL
      const reader = new FileReader();
      reader.onload = (e) => {
        const newFile = {
          uid: file.name,
          name: file.name,
          status: 'done' as const,
          url: e.target?.result as string,
          thumbUrl: e.target?.result as string,
        };
        setFileList([newFile]);
      };
      reader.readAsDataURL(file);

      return false; // 阻止自动上传
    },
    onRemove: () => {
      setFileList([]);
    },
  };

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
            创建产品
          </Title>
        </div>
        <Space>
          <Button onClick={() => navigate(ROUTES.PRODUCTS)}>取消</Button>
          <Button
            type="primary"
            icon={<SaveOutlined />}
            loading={createMutation.isPending}
            onClick={handleSubmit}
          >
            保存
          </Button>
        </Space>
      </div>

      <Row gutter={24}>
        <Col xs={24} lg={16}>
          <Card title="基本信息">
            <Form form={form} layout="vertical">
              <Form.Item
                label="产品名称"
                name="name"
                rules={[
                  { required: true, message: '请输入产品名称' },
                  { min: 2, max: 100, message: '产品名称长度应在2-100个字符之间' },
                ]}
              >
                <Input placeholder="请输入产品名称" />
              </Form.Item>

              <Form.Item
                label="产品描述"
                name="description"
                rules={[
                  { required: true, message: '请输入产品描述' },
                  { min: 10, max: 500, message: '产品描述长度应在10-500个字符之间' },
                ]}
              >
                <TextArea
                  rows={4}
                  placeholder="请输入产品描述"
                  showCount
                  maxLength={500}
                />
              </Form.Item>

              <Row gutter={16}>
                <Col xs={24} sm={12}>
                  <Form.Item
                    label="价格 (元)"
                    name="price"
                    rules={[
                      { required: true, message: '请输入产品价格' },
                      { type: 'number', min: 0.01, message: '价格必须大于0' },
                    ]}
                  >
                    <InputNumber
                      style={{ width: '100%' }}
                      placeholder="请输入价格"
                      min={0.01}
                      precision={2}
                      formatter={(value) => `¥ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
                      parser={(value) => value!.replace(/¥\s?|(,*)/g, '')}
                    />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12}>
                  <Form.Item
                    label="库存数量"
                    name="stock"
                    rules={[
                      { required: true, message: '请输入库存数量' },
                      { type: 'number', min: 0, message: '库存数量不能小于0' },
                    ]}
                  >
                    <InputNumber
                      style={{ width: '100%' }}
                      placeholder="请输入库存数量"
                      min={0}
                      precision={0}
                    />
                  </Form.Item>
                </Col>
              </Row>

              <Form.Item
                label="状态"
                name="status"
                initialValue="active"
                rules={[{ required: true, message: '请选择产品状态' }]}
              >
                <Select placeholder="请选择产品状态">
                  <Option value="active">上架</Option>
                  <Option value="inactive">下架</Option>
                </Select>
              </Form.Item>
            </Form>
          </Card>
        </Col>

        <Col xs={24} lg={8}>
          <Card title="产品图片">
            <Form.Item label="上传图片">
              <Upload {...uploadProps}>
                {fileList.length >= 1 ? null : (
                  <div className="text-center p-4">
                    <UploadOutlined className="text-2xl mb-2" />
                    <div>点击上传图片</div>
                    <div className="text-gray-500 text-xs mt-1">
                      支持 JPG、PNG 格式，大小不超过 2MB
                    </div>
                  </div>
                )}
              </Upload>
            </Form.Item>

            <div className="mt-4 p-4 bg-gray-50 rounded">
              <div className="text-sm text-gray-600 space-y-1">
                <div>• 建议图片尺寸：800x800px</div>
                <div>• 支持格式：JPG、PNG</div>
                <div>• 文件大小：不超过2MB</div>
                <div>• 图片将用于产品展示</div>
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default ProductCreate;
