<template>
  <div class="min-h-screen bg-gray-100">
    <!-- 导航栏 -->
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex">
            <div class="flex-shrink-0 flex items-center">
              <router-link to="/" class="text-xl font-bold text-primary-600">
                Blue Bell
              </router-link>
            </div>
          </div>
          <div class="flex items-center" v-if="userStore.isLoggedIn">
            <router-link
              to="/post/create"
              class="btn btn-primary mr-4"
            >
              发帖
            </router-link>
            <button
              class="btn btn-ghost"
              @click="handleLogout"
            >
              退出
            </button>
          </div>
          <div class="flex items-center" v-else>
            <router-link
              to="/login"
              class="btn btn-ghost mr-4"
            >
              登录
            </router-link>
            <router-link
              to="/register"
              class="btn btn-primary"
            >
              注册
            </router-link>
          </div>
        </div>
      </div>
    </nav>

    <!-- 主要内容 -->
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <router-view></router-view>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from './stores/user'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style>
.app {
  min-height: 100vh;
}
</style>
