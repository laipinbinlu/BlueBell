import axios, { AxiosResponse } from 'axios'
import type { ApiResponse, User, LoginResponse } from '@/types/index'

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 5000
})

// 用户注册
export const register = async (data: {
  username: string
  password: string
  re_password: string
}): Promise<AxiosResponse<ApiResponse<LoginResponse>>> => {
  return api.post<ApiResponse<LoginResponse>>('/signup', data)
}

// 用户登录
export const login = async (data: {
  username: string
  password: string
}): Promise<AxiosResponse<ApiResponse<LoginResponse>>> => {
  return api.post<ApiResponse<LoginResponse>>('/login', data)
}

// 获取用户信息
export const getUserInfo = async (): Promise<AxiosResponse<ApiResponse<User>>> => {
  return api.get<ApiResponse<User>>('/user/info')
}

// 用户登出
export const logout = async (): Promise<AxiosResponse<ApiResponse<null>>> => {
  return api.post<ApiResponse<null>>('/logout')
} 