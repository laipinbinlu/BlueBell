import request from './request'
import type { Comment } from '@/types'

interface CreateCommentData {
  post_id: number
  content: string
  parent_id?: number
}

export const getCommentList = (postId: number) => {
  return request.get<Comment[]>(`/api/v1/post/${postId}/comment`)
}

export const createComment = (data: CreateCommentData) => {
  return request.post<Comment>('/api/v1/comment', data)
}

export const updateComment = (id: number, data: Pick<CreateCommentData, 'content'>) => {
  return request.put<Comment>(`/api/v1/comment/${id}`, data)
}

export const deleteComment = (id: number) => {
  return request.delete(`/api/v1/comment/${id}`)
}

export const voteComment = (id: string | number, direction: 'up' | 'down') => {
  return request.post<ApiResponse<Comment>>(`/comments/${id}/vote`, { direction })
} 