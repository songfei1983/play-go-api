import apiClient from './api';
import type { User, GetUsersParams, UpdateUserRequest, PaginatedResponse } from '@/types';

export const userService = {
  // 获取用户列表
  getUsers: (params?: GetUsersParams) =>
    apiClient.get<PaginatedResponse<User>>('/users', { params }),

  // 获取单个用户
  getUser: (id: string) =>
    apiClient.get<User>(`/users/${id}`),

  // 更新用户信息
  updateUser: (id: string, data: UpdateUserRequest) =>
    apiClient.put<User>(`/users/${id}`, data),

  // 删除用户
  deleteUser: (id: string) =>
    apiClient.delete(`/users/${id}`),

  // 软删除用户
  softDeleteUser: (id: string) =>
    apiClient.delete(`/users/${id}/soft`),

  // 恢复用户
  restoreUser: (id: string) =>
    apiClient.post(`/users/${id}/restore`),
};
