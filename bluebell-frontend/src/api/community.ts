import request from './request'
import type { Community } from '@/types'

export const getCommunityList = () => {
  return request.get<Community[]>('/api/v1/community')
}

export const getCommunityDetail = (id: number) => {
  return request.get<Community>(`/api/v1/community/${id}`)
}

export const createCommunity = (data: Omit<Community, 'id' | 'create_time' | 'update_time'>) => {
  return request.post<Community>('/api/v1/community', data)
}

export const updateCommunity = (id: number, data: Partial<Omit<Community, 'id' | 'create_time' | 'update_time'>>) => {
  return request.put<Community>(`/api/v1/community/${id}`, data)
}

export const deleteCommunity = (id: number) => {
  return request.delete(`/api/v1/community/${id}`)
} 