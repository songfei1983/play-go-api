import apiClient from './api';
import type { Product, GetProductsParams, CreateProductRequest, UpdateProductRequest, PaginatedResponse } from '@/types';

export const productService = {
  // 获取产品列表
  getProducts: (params?: GetProductsParams) =>
    apiClient.get<PaginatedResponse<Product>>('/products', { params }),

  // 获取单个产品
  getProduct: (id: string) =>
    apiClient.get<Product>(`/products/${id}`),

  // 创建产品
  createProduct: (data: CreateProductRequest) =>
    apiClient.post<Product>('/products', data),

  // 更新产品信息
  updateProduct: (id: string, data: UpdateProductRequest) =>
    apiClient.put<Product>(`/products/${id}`, data),

  // 删除产品
  deleteProduct: (id: string) =>
    apiClient.delete(`/products/${id}`),

  // 软删除产品
  softDeleteProduct: (id: string) =>
    apiClient.delete(`/products/${id}/soft`),

  // 恢复产品
  restoreProduct: (id: string) =>
    apiClient.post(`/products/${id}/restore`),
};
