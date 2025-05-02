import axios, { AxiosResponse } from 'axios';
import { message } from 'antd';
import { User, Community, Post, PostDetail, ApiResponse, LoginResponse, SignupResponse } from '../types';

const api = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json'
  }
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
  (response: AxiosResponse) => {
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      // 未授权，清除token并刷新页面
      localStorage.removeItem('token');
      window.location.reload();
      return Promise.reject(new Error('请重新登录'));
    }
    message.error(error.response?.data?.message || '请求失败');
    return Promise.reject(error);
  }
);

// 用户相关 API
export const userApi = {
  login: (username: string, password: string) =>
    api.post<LoginResponse>('/login', { username, password }),
  
  signup: (username: string, password: string) =>
    api.post<SignupResponse>('/signup', { username, password }),
};

// 社区相关 API
export const communityApi = {
  getCommunityList: (): Promise<ApiResponse<Community[]>> =>
    api.get('/community'),
  
  getCommunityDetail: (id: number): Promise<ApiResponse<Community>> =>
    api.get(`/community/${id}`),
};

// 帖子相关 API
export const postApi = {
  createPost: (post: Partial<Post>): Promise<ApiResponse<null>> =>
    api.post('/post', post),
  
  getPostDetail: (id: number): Promise<ApiResponse<PostDetail>> =>
    api.get(`/post/${id}`),
  
  getPostList: (page: number = 1, size: number = 10): Promise<ApiResponse<PostDetail[]>> =>
    api.get('/posts2', {
      params: { page, size, order: 'time' }
    }),
  
  votePost: (post_id: number, direction: 1 | -1): Promise<ApiResponse<null>> =>
    api.post('/vote', { post_id, direction }),
}; 