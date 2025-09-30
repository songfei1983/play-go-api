import dayjs from 'dayjs';
import type { TableColumn } from '@/types';

// 格式化日期
export const formatDate = (date: string | Date, format = 'YYYY-MM-DD HH:mm:ss') => {
  return dayjs(date).format(format);
};

// 格式化相对时间
export const formatRelativeTime = (date: string | Date) => {
  return dayjs(date).fromNow();
};

// 延迟函数
export const delay = (ms: number) => {
  return new Promise(resolve => setTimeout(resolve, ms));
};

// 防抖函数
export const debounce = <T extends (...args: unknown[]) => unknown>(
  func: T,
  wait: number
): ((...args: Parameters<T>) => void) => {
  let timeout: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
};

// 节流函数
export const throttle = <T extends (...args: unknown[]) => unknown>(
  func: T,
  limit: number
): ((...args: Parameters<T>) => void) => {
  let inThrottle: boolean;
  return (...args: Parameters<T>) => {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
};

// 深拷贝
export const deepClone = <T>(obj: T): T => {
  if (obj === null || typeof obj !== 'object') {
    return obj;
  }

  if (obj instanceof Date) {
    return new Date(obj.getTime()) as T;
  }

  if (obj instanceof Array) {
    return obj.map(item => deepClone(item)) as T;
  }

  if (typeof obj === 'object') {
    const clonedObj = {} as T;
    for (const key in obj) {
      if (Object.prototype.hasOwnProperty.call(obj, key)) {
        clonedObj[key] = deepClone(obj[key]);
      }
    }
    return clonedObj;
  }

  return obj;
};

// 生成随机ID
export const generateId = (length = 8) => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
};

// 格式化文件大小
export const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

// 验证邮箱格式
export const isValidEmail = (email: string) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

// 验证手机号格式（中国大陆）
export const isValidPhone = (phone: string) => {
  const phoneRegex = /^1[3-9]\d{9}$/;
  return phoneRegex.test(phone);
};

// 获取文件扩展名
export const getFileExtension = (filename: string) => {
  return filename.slice(((filename.lastIndexOf('.') - 1) >>> 0) + 2);
};

// 首字母大写
export const capitalize = (str: string) => {
  return str.charAt(0).toUpperCase() + str.slice(1);
};

// 驼峰转下划线
export const camelToSnake = (str: string) => {
  return str.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`);
};

// 下划线转驼峰
export const snakeToCamel = (str: string) => {
  return str.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());
};

// 获取URL参数
export const getUrlParams = (url?: string) => {
  const urlStr = url || window.location.href;
  const params = new URLSearchParams(urlStr.split('?')[1]);
  const result: Record<string, string> = {};

  params.forEach((value, key) => {
    result[key] = value;
  });

  return result;
};

// 设置URL参数
export const setUrlParams = (params: Record<string, string | number>) => {
  const url = new URL(window.location.href);

  Object.entries(params).forEach(([key, value]) => {
    url.searchParams.set(key, String(value));
  });

  window.history.replaceState({}, '', url.toString());
};

// 移除URL参数
export const removeUrlParams = (keys: string[]) => {
  const url = new URL(window.location.href);

  keys.forEach(key => {
    url.searchParams.delete(key);
  });

  window.history.replaceState({}, '', url.toString());
};

// 表格列工具函数
export const createTableColumn = <T = Record<string, unknown>>(
  title: string,
  dataIndex: keyof T,
  options?: Partial<TableColumn<T>>
): TableColumn<T> => {
  return {
    title,
    dataIndex: dataIndex as string,
    key: dataIndex as string,
    ...options,
  };
};

// 处理API错误
export const handleApiError = (error: unknown) => {
  if (error.response) {
    // 服务器响应了错误状态码
    const { status, data } = error.response;
    return {
      status,
      message: data?.message || `HTTP Error ${status}`,
      data: data?.data,
    };
  } else if (error.request) {
    // 请求已发出但没有收到响应
    return {
      status: 0,
      message: '网络连接失败，请检查网络设置',
      data: null,
    };
  } else {
    // 其他错误
    return {
      status: -1,
      message: error.message || '未知错误',
      data: null,
    };
  }
};

// 本地存储工具
export const storage = {
  get: <T = unknown>(key: string, defaultValue?: T): T | null => {
    try {
      const item = localStorage.getItem(key);
      return item ? JSON.parse(item) : defaultValue || null;
    } catch {
      return defaultValue || null;
    }
  },

  set: (key: string, value: unknown) => {
    try {
      localStorage.setItem(key, JSON.stringify(value));
    } catch (error) {
      console.error('Failed to save to localStorage:', error);
    }
  },

  remove: (key: string) => {
    try {
      localStorage.removeItem(key);
    } catch (error) {
      console.error('Failed to remove from localStorage:', error);
    }
  },

  clear: () => {
    try {
      localStorage.clear();
    } catch (error) {
      console.error('Failed to clear localStorage:', error);
    }
  },
};
