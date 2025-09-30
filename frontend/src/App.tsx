import React, { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { ConfigProvider, App as AntdApp } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { Helmet, HelmetProvider } from 'react-helmet-async';
import { useAuthStore } from '@/stores/authStore';
import { useAppStore } from '@/stores/appStore';
import Layout from '@/components/layout/MainLayout';
import Login from '@/pages/Login';
import Register from '@/pages/Register';
import Dashboard from '@/pages/Dashboard';
import Users from '@/pages/Users';
import Products from '@/pages/Products';
import UserDetail from '@/pages/UserDetail';
import ProductDetail from '@/pages/ProductDetail';
import ProductCreate from '@/pages/ProductCreate';
import ProductEdit from '@/pages/ProductEdit';
import { ROUTES } from '@/constants';
import './App.css';

// 创建 React Query 客户端
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 5 * 60 * 1000, // 5分钟
    },
    mutations: {
      retry: 1,
    },
  },
});

// 受保护的路由组件
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuthStore();

  if (!isAuthenticated) {
    return <Navigate to={ROUTES.LOGIN} replace />;
  }

  return <>{children}</>;
};

// 公共路由组件（已登录用户不能访问）
const PublicRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuthStore();

  if (isAuthenticated) {
    return <Navigate to={ROUTES.DASHBOARD} replace />;
  }

  return <>{children}</>;
};

const App: React.FC = () => {
  const { getCurrentUser } = useAuthStore();
  const { theme } = useAppStore();

  // 应用初始化时检查用户登录状态
  useEffect(() => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      getCurrentUser();
    }
  }, [getCurrentUser]);

  return (
    <HelmetProvider>
      <QueryClientProvider client={queryClient}>
        <ConfigProvider
          locale={zhCN}
          theme={{
            token: {
              colorPrimary: '#1890ff',
            },
            algorithm: theme === 'dark' ? undefined : undefined, // 可以后续添加暗色主题
          }}
        >
          <AntdApp>
            <Helmet>
              <title>{import.meta.env.VITE_APP_TITLE || 'React Admin'}</title>
              <meta name="description" content="React Admin Dashboard" />
            </Helmet>

            <Router>
              <Routes>
                {/* 公共路由 */}
                <Route
                  path={ROUTES.LOGIN}
                  element={
                    <PublicRoute>
                      <Login />
                    </PublicRoute>
                  }
                />
                <Route
                  path={ROUTES.REGISTER}
                  element={
                    <PublicRoute>
                      <Register />
                    </PublicRoute>
                  }
                />

                {/* 受保护的路由 */}
                <Route
                  path="/*"
                  element={
                    <ProtectedRoute>
                      <Layout>
                        <Routes>
                          <Route path={ROUTES.DASHBOARD} element={<Dashboard />} />
                          <Route path={ROUTES.USERS} element={<Users />} />
                          <Route path={`${ROUTES.USERS}/:id`} element={<UserDetail />} />
                          <Route path={ROUTES.PRODUCTS} element={<Products />} />
                          <Route path={`${ROUTES.PRODUCTS}/create`} element={<ProductCreate />} />
                          <Route path={`${ROUTES.PRODUCTS}/:id`} element={<ProductDetail />} />
                          <Route path={`${ROUTES.PRODUCTS}/:id/edit`} element={<ProductEdit />} />
                          <Route path={ROUTES.HOME} element={<Navigate to={ROUTES.DASHBOARD} replace />} />
                          <Route path="*" element={<Navigate to={ROUTES.DASHBOARD} replace />} />
                        </Routes>
                      </Layout>
                    </ProtectedRoute>
                  }
                />
              </Routes>
            </Router>
          </AntdApp>

          {/* React Query 开发工具 - 暂时禁用以避免hook问题 */}
          {/* {import.meta.env.DEV && <ReactQueryDevtools initialIsOpen={false} />} */}
        </ConfigProvider>
      </QueryClientProvider>
    </HelmetProvider>
  );
};

export default App;
