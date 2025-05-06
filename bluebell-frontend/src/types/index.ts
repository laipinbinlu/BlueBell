// 用户相关类型
export interface User {
  user_id: number;
  username: string;
  token: string;
}

// 登录响应类型
export interface LoginResponse {
  user_id: string;  // 注意这里是字符串类型
  user_name: string;
  token: string;
}

// 社区相关类型
export interface Community {
  id: number;
  name: string;
  introduction: string;
  create_time: string;
  update_time: string;
}

// 帖子相关类型
export interface Post {
  id: number;
  author_id: number;
  community_id: number;
  title: string;
  content: string;
  status: number;
  create_time: string;
  update_time: string;
  vote_num: number;
  author_name: string;
  community_name: string;
}

// 评论相关类型
export interface Comment {
  id: number;
  post_id: number;
  author_id: number;
  parent_id: number;
  content: string;
  create_time: string;
  update_time: string;
  author_name: string;
}

// 投票相关类型
export interface VoteData {
  postId: string;
  direction: number; // 1: 赞成, -1: 反对, 0: 取消
}

// 消息相关类型
export interface Message {
  id: number;
  fromUserId: number;
  toUserId: number;
  content: string;
  createTime: string;
  isRead: boolean;
}

// API 响应类型
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

// 分页参数类型
export interface PageParams {
  page: number;
  size: number;
}

// WebSocket 消息类型
export interface WsMessage {
  type: 'chat' | 'system';
  data: any;
} 