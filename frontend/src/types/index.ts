// 用户相关类型
export interface User {
  id: string;
  username: string;
  email: string;
  name: string;
  phone?: string;
  status: 'active' | 'inactive';
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  name: string;
  phone?: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface RegisterResponse {
  message: string;
}

export interface UpdateUserRequest {
  name?: string;
  email?: string;
  phone?: string;
  status?: 'active' | 'inactive';
}

export interface GetUsersParams {
  page?: number;
  limit?: number;
  search?: string;
  status?: 'active' | 'inactive';
}

// 产品相关类型
export interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  stock: number;
  status: 'active' | 'inactive';
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface CreateProductRequest {
  name: string;
  description: string;
  price: number;
  stock: number;
  status?: 'active' | 'inactive';
}

export interface UpdateProductRequest {
  name?: string;
  description?: string;
  price?: number;
  stock?: number;
  status?: 'active' | 'inactive';
}

export interface GetProductsParams {
  page?: number;
  limit?: number;
  search?: string;
  status?: 'active' | 'inactive';
}

// API响应类型
export interface ApiResponse<T> {
  data: T;
  message?: string;
  success: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

// 通用类型
export interface PaginationConfig {
  current: number;
  pageSize: number;
  total: number;
  showSizeChanger?: boolean;
  showQuickJumper?: boolean;
}

export interface ColumnDef<T> {
  key: string;
  title: string;
  dataIndex?: keyof T;
  render?: (value: unknown, record: T, index: number) => React.ReactNode;
  width?: number;
  align?: 'left' | 'center' | 'right';
  sorter?: boolean;
  filters?: Array<{ text: string; value: string | number | boolean }>;
}

export interface TableColumn<T> {
  title: string;
  dataIndex: string;
  key: string;
  render?: (value: unknown, record: T, index: number) => React.ReactNode;
  width?: number;
  align?: 'left' | 'center' | 'right';
  sorter?: boolean;
  filters?: Array<{ text: string; value: string | number | boolean }>;
}

// 表单相关类型
export interface FormFieldProps {
  label: string;
  name: string;
  type?: 'text' | 'email' | 'password' | 'number' | 'textarea';
  required?: boolean;
  placeholder?: string;
  rules?: Array<Record<string, unknown>>;
}

// 状态类型
export type LoadingState = 'idle' | 'loading' | 'success' | 'error';

export type Theme = 'light' | 'dark';

export type StatusType = 'active' | 'inactive' | 'pending' | 'deleted';
