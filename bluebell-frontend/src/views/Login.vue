<template>
  <div class="min-h-[calc(100vh-4rem)] flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          登录账号
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          或者
          <router-link to="/register" class="font-medium text-primary-600 hover:text-primary-500">
            注册新账号
          </router-link>
        </p>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleSubmit">
        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <label for="username" class="sr-only">用户名</label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              required
              class="input rounded-t-md rounded-b-none"
              placeholder="用户名"
            >
          </div>
          <div>
            <label for="password" class="sr-only">密码</label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              required
              class="input rounded-t-none"
              placeholder="密码"
            >
          </div>
        </div>

        <div>
          <button
            type="submit"
            class="btn btn-primary w-full"
            :disabled="loading"
          >
            <span v-if="loading" class="flex items-center justify-center">
              <span class="animate-spin h-5 w-5 mr-3 border-b-2 border-white rounded-full"></span>
              登录中...
            </span>
            <span v-else>登录</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import type { LoginResponse } from '@/types/index'

interface FormData {
  username: string
  password: string
}

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(false)
const form = ref<FormData>({
  username: '',
  password: ''
})

const handleSubmit = async () => {
  try {
    loading.value = true
    const response = await userStore.login(form.value.username, form.value.password)
    if (response && 'user_id' in response && 'token' in response) {
      alert('登录成功！')
      const redirectPath = route.query.redirect as string || '/'
      router.replace(redirectPath)
    } else {
      throw new Error('登录失败：服务器返回数据格式错误')
    }
  } catch (error: any) {
    console.error('登录失败:', error)
    alert(error.message || '登录失败，请检查用户名和密码')
  } finally {
    loading.value = false
  }
}
</script>
