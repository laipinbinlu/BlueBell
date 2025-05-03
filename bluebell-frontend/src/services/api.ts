import axios, { AxiosResponse } from 'axios';
import { message } from 'antd';
import { Community, Post, PostDetail, ApiResponse, LoginResponse, SignupResponse } from '../types';

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
      console.log('添加token到请求头:', token); // 添加调试日志
    } else {
      console.log('未找到token'); // 添加调试日志
    }
    return config;
  },
  (error) => {
    console.error('请求拦截器错误:', error); // 添加调试日志
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    // 检查响应中的业务状态码
    const data = response.data;
    if (data && data.code !== 0) {
      message.error(data.message || '操作失败');
      return Promise.reject(new Error(data.message || '操作失败'));
    }
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
  
  signup: (username: string, password: string, re_password: string) =>
    api.post<SignupResponse>('/signup', { username, password, re_password }),
};

// 社区相关 API
export const communityApi = {
  getCommunityList: (): Promise<ApiResponse<Community[]>> =>
    api.get<ApiResponse<Community[]>>('/community')
      .then(response => {
        console.log('获取社区列表响应:', response.data); // 添加调试日志
        return response.data;
      })
      .catch(error => {
        console.error('获取社区列表失败:', error);
        message.error('获取社区列表失败');
        throw error;
      }),
  
  getCommunityDetail: (id: number): Promise<ApiResponse<Community>> =>
    api.get<ApiResponse<Community>>(`/community/${id}`)
      .then(response => response.data),
};

// 帖子相关 API
export const postApi = {
  createPost: (post: Partial<Post>): Promise<ApiResponse<null>> =>
    api.post<ApiResponse<null>>('/post', post)
      .then(response => response.data),
  
  getPostDetail: (id: number): Promise<ApiResponse<PostDetail>> =>
    api.get<ApiResponse<PostDetail>>(`/post/${id}`)
      .then(response => response.data),
  
  getPostList: (page: number = 1, size: number = 10): Promise<ApiResponse<PostDetail[]>> =>
    api.get<ApiResponse<PostDetail[]>>('/posts2', {
      params: { 
        page, 
        size, 
        order: 'time' 
      }
    })
    .then(response => {
      console.log('获取帖子列表响应:', response.data); // 添加调试日志
      return response.data;
    })
    .catch(error => {
      console.error('获取帖子列表失败:', error);
      if (error.response?.status === 401) {
        message.error('请先登录');
      } else {
        message.error('获取帖子列表失败，请重试');
      }
      throw error;
    }),
  
  votePost: (post_id: number, direction: 1 | -1): Promise<ApiResponse<null>> =>
    api.post<ApiResponse<null>>('/vote', { post_id, direction })
      .then(response => response.data),
}; 