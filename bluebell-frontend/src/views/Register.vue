<template>
  <div class="min-h-[calc(100vh-4rem)] flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          注册账号
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          或者
          <router-link to="/login" class="font-medium text-primary-600 hover:text-primary-500">
            登录已有账号
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
              class="input rounded-t-none rounded-b-none"
              placeholder="密码"
            >
          </div>
          <div>
            <label for="rePassword" class="sr-only">确认密码</label>
            <input
              id="rePassword"
              v-model="form.rePassword"
              type="password"
              required
              class="input rounded-t-none"
              placeholder="确认密码"
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
              注册中...
            </span>
            <span v-else>注册</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const form = ref({
  username: '',
  password: '',
  rePassword: ''
})

const handleSubmit = async () => {
  if (form.value.password !== form.value.rePassword) {
    alert('两次输入的密码不一致')
    return
  }

  try {
    loading.value = true
    await userStore.register(
      form.value.username,
      form.value.password,
      form.value.rePassword
    )
    router.push('/')
  } catch (error: any) {
    alert(error.message || '注册失败')
  } finally {
    loading.value = false
  }
}
</script>
