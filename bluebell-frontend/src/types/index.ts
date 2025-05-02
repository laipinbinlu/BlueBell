// 用户相关类型
export interface User {
  user_id: number;
  username: string;
  token: string;
}

// 社区相关类型
export interface Community {
  id: number;
  name: string;
  introduction: string;
  create_time: string;
}

// 帖子相关类型
export interface Post {
  post_id: number;
  title: string;
  content: string;
  author_id: number;
  community_id: number;
  create_time: string;
  vote_num?: number;
}

// API 响应类型
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 帖子详情类型
export interface PostDetail extends Post {
  author_name: string;
  community_detail: Community;
}

// 登录响应类型
export interface LoginResponse {
  code: number;
  message: string;
  data: {
    token: string;
    user_id: number;
    username: string;
  };
}

// 注册响应类型
export interface SignupResponse {
  code: number;
  message: string;
  data: null | {
    user_id: number;
    username: string;
  };
} 