import { defineStore } from 'pinia';
import type { User, ApiResponse, LoginResponse } from '@/types/index';
import { login, register, logout, getUserInfo } from '@/api/user';

interface UserState {
  user: User | null;
  token: string | null;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    user: null,
    token: localStorage.getItem('token')
  }),

  getters: {
    isLoggedIn: (state) => !!state.token
  },

  actions: {
    async login(username: string, password: string): Promise<LoginResponse> {
      const response = await login({ username, password });
      const data = response.data.data;
      this.setUser({
        user: {
          user_id: parseInt(data.user_id),
          username: data.user_name,
          token: data.token
        },
        token: data.token
      });
      return data;
    },

    async register(username: string, password: string, rePassword: string): Promise<LoginResponse> {
      const response = await register({ username, password, re_password: rePassword });
      const data = response.data.data;
      this.setUser({
        user: {
          user_id: parseInt(data.user_id),
          username: data.user_name,
          token: data.token
        },
        token: data.token
      });
      return data;
    },

    setUser(data: { user: User; token: string }) {
      this.user = data.user;
      this.token = data.token;
      localStorage.setItem('token', data.token);
    },

    async logout() {
      try {
        await logout();
      } finally {
        this.user = null;
        this.token = null;
        localStorage.removeItem('token');
      }
    },

    async fetchUserInfo() {
      const response = await getUserInfo();
      this.user = response.data.data;
    }
  }
}); 