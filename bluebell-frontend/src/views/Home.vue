<template>
  <div class="grid grid-cols-12 gap-6">
    <!-- 社区列表 -->
    <div class="col-span-3">
      <div class="bg-white rounded-lg shadow p-4">
        <h2 class="text-lg font-semibold mb-4">社区列表</h2>
        <div v-if="loading" class="text-center py-4">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
        </div>
        <ul v-else class="space-y-2">
          <li v-for="community in communities" :key="community.id">
            <router-link
              :to="{ name: 'community-detail', params: { id: community.id }}"
              class="block p-2 rounded hover:bg-gray-50"
            >
              {{ community.name }}
            </router-link>
          </li>
        </ul>
      </div>
    </div>

    <!-- 帖子列表 -->
    <div class="col-span-9">
      <div class="bg-white rounded-lg shadow">
        <div class="p-4 border-b">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold">最新帖子</h2>
            <div class="space-x-2">
              <button
                v-for="order in ['time', 'score']"
                :key="order"
                class="btn btn-ghost"
                :class="{ 'text-primary-600': currentOrder === order }"
                @click="handleOrderChange(order)"
              >
                {{ order === 'time' ? '最新' : '最热' }}
              </button>
            </div>
          </div>
        </div>
        
        <div v-if="loading" class="text-center py-8">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
        </div>
        
        <div v-else-if="posts.length === 0" class="text-center py-8 text-gray-500">
          暂无帖子
        </div>
        
        <ul v-else class="divide-y">
          <li v-for="post in posts" :key="post.id" class="p-4 hover:bg-gray-50">
            <router-link :to="{ name: 'post-detail', params: { id: post.id }}">
              <h3 class="text-lg font-medium mb-2">{{ post.title }}</h3>
              <div class="text-sm text-gray-500">
                <span>{{ post.author_name }}</span>
                <span class="mx-2">·</span>
                <span>{{ post.community_name }}</span>
                <span class="mx-2">·</span>
                <span>{{ new Date(post.create_time).toLocaleString() }}</span>
              </div>
            </router-link>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Community, Post } from '@/types'
import { getCommunityList } from '@/api/community'
import { getPostList } from '@/api/post'

const loading = ref(false)
const communities = ref<Community[]>([])
const posts = ref<Post[]>([])
const currentOrder = ref('time')

const loadCommunities = async () => {
  try {
    loading.value = true
    communities.value = await getCommunityList()
  } catch (error) {
    console.error('Failed to load communities:', error)
  } finally {
    loading.value = false
  }
}

const loadPosts = async () => {
  try {
    loading.value = true
    posts.value = await getPostList({
      order: currentOrder.value
    })
  } catch (error) {
    console.error('Failed to load posts:', error)
  } finally {
    loading.value = false
  }
}

const handleOrderChange = (order: string) => {
  currentOrder.value = order
  loadPosts()
}

onMounted(() => {
  loadCommunities()
  loadPosts()
})
</script>
