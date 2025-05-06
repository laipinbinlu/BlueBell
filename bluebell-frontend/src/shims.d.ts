declare module '@/stores/user' {
  import { StoreDefinition } from 'pinia'
  import { User } from '@/types'
  
  export const useUserStore: StoreDefinition<'user', {
    user: User | null
    token: string | null
  }, {
    isLoggedIn: boolean
  }, {
    login(username: string, password: string): Promise<void>
    register(username: string, password: string, rePassword: string): Promise<void>
    logout(): void
  }>
} 