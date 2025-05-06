<template>
  <div class="space-y-6">
    <!-- 社区信息 -->
    <div class="bg-white rounded-lg shadow p-6">
      <div v-if="loading" class="text-center py-8">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
      </div>
      <div v-else-if="!community" class="text-center py-8 text-gray-500">
        社区不存在
      </div>
      <div v-else>
        <h1 class="text-2xl font-bold mb-4">{{ community.name }}</h1>
        <p class="text-gray-600">{{ community.introduction }}</p>
        <div class="mt-4 text-sm text-gray-500">
          创建时间：{{ new Date(community.create_time).toLocaleString() }}
        </div>
      </div>
    </div>

    <!-- 帖子列表 -->
    <div class="bg-white rounded-lg shadow">
      <div class="p-4 border-b flex justify-between items-center">
        <h2 class="text-lg font-semibold">社区帖子</h2>
        <router-link
          :to="{ name: 'create-post', query: { community_id: communityId }}"
          class="btn btn-primary"
        >
          发帖
        </router-link>
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
              <span>{{ new Date(post.create_time).toLocaleString() }}</span>
            </div>
          </router-link>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import type { Community, Post } from '@/types'
import { getCommunityDetail } from '@/api/community'
import { getPostList } from '@/api/post'

const route = useRoute()
const communityId = Number(route.params.id)

const loading = ref(false)
const community = ref<Community | null>(null)
const posts = ref<Post[]>([])

const loadCommunity = async () => {
  try {
    loading.value = true
    community.value = await getCommunityDetail(communityId)
  } catch (error) {
    console.error('Failed to load community:', error)
  } finally {
    loading.value = false
  }
}

const loadPosts = async () => {
  try {
    loading.value = true
    posts.value = await getPostList({ community_id: communityId })
  } catch (error) {
    console.error('Failed to load posts:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCommunity()
  loadPosts()
})
</script> 