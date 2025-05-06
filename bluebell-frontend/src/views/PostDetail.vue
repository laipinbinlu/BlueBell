<template>
  <div class="space-y-6">
    <!-- 帖子详情 -->
    <div class="bg-white rounded-lg shadow">
      <div v-if="loading" class="text-center py-8">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
      </div>
      <div v-else-if="!post" class="text-center py-8 text-gray-500">
        帖子不存在
      </div>
      <div v-else class="p-6">
        <h1 class="text-2xl font-bold mb-4">{{ post.title }}</h1>
        <div class="flex items-center text-sm text-gray-500 mb-4">
          <span>{{ post.author_name }}</span>
          <span class="mx-2">·</span>
          <router-link
            :to="{ name: 'community-detail', params: { id: post.community_id }}"
            class="text-primary-600 hover:text-primary-500"
          >
            {{ post.community_name }}
          </router-link>
          <span class="mx-2">·</span>
          <span>{{ new Date(post.create_time).toLocaleString() }}</span>
        </div>
        <div class="prose max-w-none">
          {{ post.content }}
        </div>
        <div class="mt-4 flex items-center space-x-4">
          <button
            class="btn btn-ghost"
            :class="{ 'text-primary-600': post.vote === 1 }"
            @click="handleVote(1)"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M3.293 9.707a1 1 0 010-1.414l6-6a1 1 0 011.414 0l6 6a1 1 0 01-1.414 1.414L11 5.414V17a1 1 0 11-2 0V5.414L4.707 9.707a1 1 0 01-1.414 0z" clip-rule="evenodd" />
            </svg>
            {{ post.vote_num }}
          </button>
          <button
            class="btn btn-ghost"
            :class="{ 'text-primary-600': post.vote === -1 }"
            @click="handleVote(-1)"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M16.707 10.293a1 1 0 010 1.414l-6 6a1 1 0 01-1.414 0l-6-6a1 1 0 111.414-1.414L9 14.586V3a1 1 0 012 0v11.586l4.293-4.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- 评论列表 -->
    <div class="bg-white rounded-lg shadow">
      <div class="p-4 border-b">
        <h2 class="text-lg font-semibold">评论</h2>
      </div>

      <!-- 评论输入框 -->
      <div class="p-4 border-b">
        <textarea
          v-model="commentContent"
          rows="3"
          class="textarea w-full"
          placeholder="写下你的评论..."
        ></textarea>
        <div class="mt-4 flex justify-end">
          <button
            class="btn btn-primary"
            :disabled="!commentContent.trim() || submitting"
            @click="handleComment"
          >
            <span v-if="submitting" class="flex items-center">
              <span class="animate-spin h-4 w-4 mr-2 border-b-2 border-white rounded-full"></span>
              发送中...
            </span>
            <span v-else>发送</span>
          </button>
        </div>
      </div>

      <div v-if="loading" class="text-center py-8">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
      </div>

      <div v-else-if="comments.length === 0" class="text-center py-8 text-gray-500">
        暂无评论
      </div>

      <ul v-else class="divide-y">
        <li v-for="comment in comments" :key="comment.id" class="p-4">
          <div class="flex items-center text-sm text-gray-500 mb-2">
            <span>{{ comment.author_name }}</span>
            <span class="mx-2">·</span>
            <span>{{ new Date(comment.create_time).toLocaleString() }}</span>
          </div>
          <p class="text-gray-900">{{ comment.content }}</p>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import type { Post, Comment } from '@/types'
import { getPostDetail, votePost } from '@/api/post'
import { getCommentList, createComment } from '@/api/comment'

const route = useRoute()
const postId = Number(route.params.id)

const loading = ref(false)
const submitting = ref(false)
const post = ref<Post | null>(null)
const comments = ref<Comment[]>([])
const commentContent = ref('')

const loadPost = async () => {
  try {
    loading.value = true
    post.value = await getPostDetail(postId)
  } catch (error) {
    console.error('Failed to load post:', error)
  } finally {
    loading.value = false
  }
}

const loadComments = async () => {
  try {
    loading.value = true
    comments.value = await getCommentList(postId)
  } catch (error) {
    console.error('Failed to load comments:', error)
  } finally {
    loading.value = false
  }
}

const handleVote = async (direction: 1 | -1) => {
  if (!post.value) return

  try {
    await votePost(postId, direction)
    post.value.vote = direction
    post.value.vote_num += direction === 1 ? 1 : -1
  } catch (error) {
    console.error('Failed to vote:', error)
  }
}

const handleComment = async () => {
  if (!commentContent.value.trim()) return

  try {
    submitting.value = true
    await createComment({
      post_id: postId,
      content: commentContent.value
    })
    commentContent.value = ''
    loadComments()
  } catch (error) {
    console.error('Failed to create comment:', error)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadPost()
  loadComments()
})
</script> 