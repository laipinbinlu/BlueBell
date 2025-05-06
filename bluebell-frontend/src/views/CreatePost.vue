<template>
  <div class="max-w-4xl mx-auto">
    <div class="bg-white rounded-lg shadow">
      <div class="p-6">
        <h1 class="text-2xl font-bold mb-6">发布新帖子</h1>
        <form @submit.prevent="handleSubmit">
          <div class="space-y-6">
            <!-- 社区选择 -->
            <div>
              <label for="community" class="block text-sm font-medium text-gray-700">
                选择社区
              </label>
              <select
                id="community"
                v-model="form.community_id"
                required
                class="mt-1 input"
              >
                <option value="" disabled>请选择社区</option>
                <option
                  v-for="community in communities"
                  :key="community.id"
                  :value="community.id"
                >
                  {{ community.name }}
                </option>
              </select>
            </div>

            <!-- 标题 -->
            <div>
              <label for="title" class="block text-sm font-medium text-gray-700">
                标题
              </label>
              <input
                id="title"
                v-model="form.title"
                type="text"
                required
                class="mt-1 input"
                placeholder="请输入标题"
              >
            </div>

            <!-- 内容 -->
            <div>
              <label for="content" class="block text-sm font-medium text-gray-700">
                内容
              </label>
              <textarea
                id="content"
                v-model="form.content"
                rows="10"
                required
                class="mt-1 textarea"
                placeholder="请输入内容"
              ></textarea>
            </div>

            <!-- 提交按钮 -->
            <div class="flex justify-end">
              <button
                type="submit"
                class="btn btn-primary"
                :disabled="submitting"
              >
                <span v-if="submitting" class="flex items-center">
                  <span class="animate-spin h-4 w-4 mr-2 border-b-2 border-white rounded-full"></span>
                  发布中...
                </span>
                <span v-else>发布</span>
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import type { Community } from '@/types'
import { getCommunityList } from '@/api/community'
import { createPost } from '@/api/post'

const router = useRouter()
const route = useRoute()

const submitting = ref(false)
const communities = ref<Community[]>([])
const form = ref({
  community_id: route.query.community_id ? Number(route.query.community_id) : '',
  title: '',
  content: ''
})

const loadCommunities = async () => {
  try {
    communities.value = await getCommunityList()
  } catch (error) {
    console.error('Failed to load communities:', error)
  }
}

const handleSubmit = async () => {
  try {
    submitting.value = true
    const post = await createPost(form.value)
    router.push({
      name: 'post-detail',
      params: { id: post.id }
    })
  } catch (error: any) {
    alert(error.message || '发布失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadCommunities()
})
</script> 