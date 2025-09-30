import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { authService } from '@/services/authService';
import type { User, LoginRequest, RegisterRequest } from '@/types';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;

  // Actions
  login: (credentials: LoginRequest) => Promise<void>;
  register: (userData: RegisterRequest) => Promise<void>;
  logout: () => void;
  clearError: () => void;
  setLoading: (loading: boolean) => void;
  getCurrentUser: () => Promise<void>;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,

      login: async (credentials: LoginRequest) => {
        try {
          set({ isLoading: true, error: null });

          const response = await authService.login(credentials);
          const { token, user } = response.data;

          localStorage.setItem('auth_token', token);
          set({
            user,
            token,
            isAuthenticated: true,
            isLoading: false
          });
        } catch (error: unknown) {
          const err = error as Error;
          set({
            error: err.message || '登录失败',
            isLoading: false
          });
          throw error;
        }
      },

      register: async (userData: RegisterRequest) => {
        try {
          set({ isLoading: true, error: null });

          await authService.register(userData);

          // 注册成功后自动登录
          await get().login({
            username: userData.username,
            password: userData.password
          });
        } catch (error: unknown) {
          const err = error as Error;
          set({
            error: err.message || '注册失败',
            isLoading: false
          });
          throw error;
        }
      },

      logout: () => {
        localStorage.removeItem('auth_token');
        set({
          user: null,
          token: null,
          isAuthenticated: false,
          error: null
        });
      },

      clearError: () => {
        set({ error: null });
      },

      setLoading: (loading: boolean) => {
        set({ isLoading: loading });
      },

      getCurrentUser: async () => {
        try {
          const token = localStorage.getItem('auth_token');
          if (!token) {
            return;
          }

          set({ isLoading: true });
          const response = await authService.getCurrentUser();

          set({
            user: response.data,
            token,
            isAuthenticated: true,
            isLoading: false
          });
        } catch {
          // Token可能已过期，清除认证信息
          get().logout();
          set({ isLoading: false });
        }
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        token: state.token,
        user: state.user,
        isAuthenticated: state.isAuthenticated
      }),
    }
  )
);
