import request from './request'
import type { Post } from '@/types'

interface PostQuery {
  page?: number
  size?: number
  order?: string
  community_id?: number
}

interface CreatePostData {
  community_id: number
  title: string
  content: string
}

export const getPostList = (params: PostQuery = {}) => {
  return request.get<Post[]>('/api/v1/posts', { params })
}

export const getPostDetail = (id: number) => {
  return request.get<Post>(`/api/v1/post/${id}`)
}

export const createPost = (data: CreatePostData) => {
  return request.post<Post>('/api/v1/post', data)
}

export const updatePost = (id: number, data: Partial<CreatePostData>) => {
  return request.put<Post>(`/api/v1/post/${id}`, data)
}

export const deletePost = (id: number) => {
  return request.delete(`/api/v1/post/${id}`)
}

export const votePost = (id: number, direction: 1 | -1 | 0) => {
  return request.post(`/api/v1/vote`, {
    post_id: id,
    direction
  })
} 