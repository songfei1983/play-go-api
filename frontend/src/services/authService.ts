import apiClient from './api';
import type { LoginRequest, RegisterRequest, LoginResponse, RegisterResponse, User } from '@/types';

export const authService = {
  // 用户登录
  login: (credentials: LoginRequest) =>
    apiClient.post<LoginResponse>('/login', credentials),

  // 用户注册
  register: (userData: RegisterRequest) =>
    apiClient.post<RegisterResponse>('/register', userData),

  // 获取当前用户信息
  getCurrentUser: () =>
    apiClient.get<User>('/users/current'),

  // 刷新token
  refreshToken: () =>
    apiClient.post<{ token: string }>('/auth/refresh'),

  // 退出登录
  logout: () =>
    apiClient.post('/auth/logout'),
};
