export interface User {
  id: number
  username: string
  created_at: string
}

export interface Community {
  id: number
  name: string
  introduction: string
  created_at: string
  member_count: number
}

export interface Post {
  id: number
  title: string
  content: string
  author: string
  community: string
  community_id: number
  created_at: string
  up_vote_count: number
  down_vote_count: number
  vote_status: number // 1: 已点赞, -1: 已踩, 0: 未投票
}

export interface Comment {
  id: number
  content: string
  author: string
  post_id: number
  created_at: string
  up_vote_count: number
  down_vote_count: number
  vote_status: number // 1: 已点赞, -1: 已踩, 0: 未投票
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
}

export interface CreatePostRequest {
  title: string
  content: string
  community_id: number
}

export interface CreateCommentRequest {
  post_id: number
  content: string
}

export interface VoteRequest {
  id: number
  direction: 'up' | 'down'
} 