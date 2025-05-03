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
  msg: string;
  data: T;
}

// 帖子详情类型
export interface PostDetail extends Post {
  author_name: string;
  community_detail: Community;
}

// 登录响应类型
export interface LoginResponse extends ApiResponse<{
  user_id: string;
  user_name: string;
  token: string;
}> {}

// 注册响应类型
export interface SignupResponse extends ApiResponse<null> {} 