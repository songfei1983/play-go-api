// API相关常量
export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/auth/login',
    REGISTER: '/auth/register',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    PROFILE: '/auth/profile',
  },
  USERS: {
    LIST: '/users',
    DETAIL: (id: string) => `/users/${id}`,
    UPDATE: (id: string) => `/users/${id}`,
    DELETE: (id: string) => `/users/${id}`,
    SOFT_DELETE: (id: string) => `/users/${id}/soft`,
    RESTORE: (id: string) => `/users/${id}/restore`,
  },
  PRODUCTS: {
    LIST: '/products',
    DETAIL: (id: string) => `/products/${id}`,
    CREATE: '/products',
    UPDATE: (id: string) => `/products/${id}`,
    DELETE: (id: string) => `/products/${id}`,
    SOFT_DELETE: (id: string) => `/products/${id}/soft`,
    RESTORE: (id: string) => `/products/${id}/restore`,
  },
} as const;

// 分页相关常量
export const PAGINATION = {
  DEFAULT_PAGE: 1,
  DEFAULT_PAGE_SIZE: 10,
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100],
} as const;

// 用户角色
export const USER_ROLES = {
  ADMIN: 'admin',
  USER: 'user',
} as const;

// 用户状态
export const USER_STATUS = {
  ACTIVE: 'active',
  INACTIVE: 'inactive',
  DELETED: 'deleted',
} as const;

// 产品状态
export const PRODUCT_STATUS = {
  ACTIVE: 'active',
  INACTIVE: 'inactive',
  DELETED: 'deleted',
} as const;

// 主题
export const THEMES = {
  LIGHT: 'light',
  DARK: 'dark',
} as const;

// 本地存储键名
export const STORAGE_KEYS = {
  AUTH_TOKEN: 'auth_token',
  USER_INFO: 'user_info',
  THEME: 'theme',
  SIDEBAR_COLLAPSED: 'sidebar_collapsed',
} as const;

// 路由路径
export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  DASHBOARD: '/dashboard',
  USERS: '/users',
  USER_DETAIL: '/users/:id',
  PRODUCTS: '/products',
  PRODUCT_DETAIL: '/products/:id',
  PRODUCT_CREATE: '/products/create',
  PROFILE: '/profile',
} as const;

// 表单验证规则
export const VALIDATION_RULES = {
  USERNAME: {
    MIN_LENGTH: 3,
    MAX_LENGTH: 20,
    PATTERN: /^[a-zA-Z0-9_]+$/,
  },
  PASSWORD: {
    MIN_LENGTH: 6,
    MAX_LENGTH: 50,
  },
  EMAIL: {
    PATTERN: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  },
  PRODUCT_NAME: {
    MIN_LENGTH: 1,
    MAX_LENGTH: 100,
  },
  PRODUCT_DESCRIPTION: {
    MAX_LENGTH: 500,
  },
} as const;

// 错误消息
export const ERROR_MESSAGES = {
  NETWORK_ERROR: '网络连接失败，请检查网络设置',
  UNAUTHORIZED: '未授权访问，请重新登录',
  FORBIDDEN: '权限不足，无法访问该资源',
  NOT_FOUND: '请求的资源不存在',
  SERVER_ERROR: '服务器内部错误，请稍后重试',
  VALIDATION_ERROR: '数据验证失败',
  LOGIN_FAILED: '用户名或密码错误',
  REGISTER_FAILED: '注册失败，请检查输入信息',
} as const;

// 成功消息
export const SUCCESS_MESSAGES = {
  LOGIN_SUCCESS: '登录成功',
  REGISTER_SUCCESS: '注册成功',
  LOGOUT_SUCCESS: '退出登录成功',
  CREATE_SUCCESS: '创建成功',
  UPDATE_SUCCESS: '更新成功',
  DELETE_SUCCESS: '删除成功',
  RESTORE_SUCCESS: '恢复成功',
} as const;

// 表单验证规则
export const FORM_RULES = {
  USERNAME: [
    { required: true, message: '请输入用户名' },
    { min: VALIDATION_RULES.USERNAME.MIN_LENGTH, message: `用户名至少${VALIDATION_RULES.USERNAME.MIN_LENGTH}个字符` },
    { max: VALIDATION_RULES.USERNAME.MAX_LENGTH, message: `用户名最多${VALIDATION_RULES.USERNAME.MAX_LENGTH}个字符` },
    { pattern: VALIDATION_RULES.USERNAME.PATTERN, message: '用户名只能包含字母、数字和下划线' },
  ],
  NAME: [
    { required: true, message: '请输入姓名' },
    { max: 50, message: '姓名最多50个字符' },
  ],
  EMAIL: [
    { required: true, message: '请输入邮箱' },
    { pattern: VALIDATION_RULES.EMAIL.PATTERN, message: '请输入有效的邮箱地址' },
  ],
  PHONE: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码' },
  ],
  PASSWORD: [
    { required: true, message: '请输入密码' },
    { min: VALIDATION_RULES.PASSWORD.MIN_LENGTH, message: `密码至少${VALIDATION_RULES.PASSWORD.MIN_LENGTH}个字符` },
    { max: VALIDATION_RULES.PASSWORD.MAX_LENGTH, message: `密码最多${VALIDATION_RULES.PASSWORD.MAX_LENGTH}个字符` },
  ],
} as const;
