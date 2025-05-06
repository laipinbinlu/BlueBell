import axios from 'axios';
import type { ApiResponse, User, Community, Post, VoteData, Message, PageParams } from '../types';

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 5000,
});

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    if (error.response?.status === 401) {
      // 未授权，清除token并跳转到登录页
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// 用户相关API
export const userApi = {
  signup: (data: { username: string; password: string; rePassword: string }) =>
    api.post<ApiResponse<User>>('/signup', data),
  login: (data: { username: string; password: string }) =>
    api.post<ApiResponse<User>>('/login', data),
};

// 社区相关API
export const communityApi = {
  getList: () => api.get<ApiResponse<Community[]>>('/community'),
  getDetail: (id: number) => api.get<ApiResponse<Community>>(`/community/${id}`),
};

// 帖子相关API
export const postApi = {
  create: (data: { title: string; content: string; communityId: number }) =>
    api.post<ApiResponse<Post>>('/post', data),
  getDetail: (id: number) => api.get<ApiResponse<Post>>(`/post/${id}`),
  getList: (params: PageParams) => api.get<ApiResponse<Post[]>>('/posts', { params }),
  getList2: (params: PageParams & { order: string }) =>
    api.get<ApiResponse<Post[]>>('/posts2', { params }),
  vote: (data: VoteData) => api.post<ApiResponse<null>>('/vote', data),
};

// 消息相关API
export const messageApi = {
  send: (data: { toUserId: number; content: string }) =>
    api.post<ApiResponse<null>>('/message', data),
  getList: () => api.get<ApiResponse<Message[]>>('/messages'),
  getUnreadCount: () => api.get<ApiResponse<number>>('/messages/unread/count'),
  getChatHistory: (toUserId: number) =>
    api.get<ApiResponse<Message[]>>('/messages/history', { params: { toUserId } }),
}; 