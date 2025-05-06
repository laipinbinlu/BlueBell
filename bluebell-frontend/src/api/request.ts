import axios from 'axios'
import type { AxiosResponse } from 'axios'
import type { ApiResponse } from '@/types'
import { useUserStore } from '@/stores/user'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 5000
})

request.interceptors.request.use(
  config => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse<unknown>>) => {
    const res = response.data
    if (res.code !== 0) {
      return Promise.reject(new Error(res.message || 'Error'))
    }
    return res.data
  },
  error => {
    return Promise.reject(error)
  }
)

export default request 