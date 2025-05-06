<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold">帖子列表</h2>
      <div class="join">
        <button
          class="join-item btn"
          :class="{ 'btn-active': order === 'time' }"
          @click="changeOrder('time')"
        >
          最新
        </button>
        <button
          class="join-item btn"
          :class="{ 'btn-active': order === 'score' }"
          @click="changeOrder('score')"
        >
          最热
        </button>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="posts.length === 0" class="text-center py-8">
      <p class="text-gray-500">暂无帖子</p>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="post in posts"
        :key="post.id"
        class="card bg-base-100 shadow-xl"
      >
        <div class="card-body">
          <h3 class="card-title">
            <router-link :to="`/post/${post.id}`" class="hover:text-primary">
              {{ post.title }}
            </router-link>
          </h3>
          <p class="text-gray-500">
            作者：{{ post.authorName }} | 社区：{{ post.communityName }} |
            时间：{{ formatTime(post.createTime) }}
          </p>
          <div class="card-actions justify-end">
            <div class="join">
              <button
                class="join-item btn btn-sm"
                :class="{ 'btn-active': post.voteCount > 0 }"
                @click="vote(post.id, 1)"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-4 w-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M5 15l7-7 7 7"
                  />
                </svg>
                {{ post.voteCount }}
              </button>
              <button
                class="join-item btn btn-sm"
                :class="{ 'btn-active': post.voteCount < 0 }"
                @click="vote(post.id, -1)"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-4 w-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M19 9l-7 7-7-7"
                  />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="join flex justify-center">
        <button
          class="join-item btn"
          :class="{ 'btn-disabled': page === 1 }"
          @click="prevPage"
        >
          上一页
        </button>
        <button class="join-item btn">第 {{ page }} 页</button>
        <button
          class="join-item btn"
          :class="{ 'btn-disabled': !hasMore }"
          @click="nextPage"
        >
          下一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import type { Post } from "../types";
import { postApi } from "../api";
import { useUserStore } from "../stores/user";
import { formatTime } from "../utils/time";

const userStore = useUserStore();
const posts = ref<Post[]>([]);
const loading = ref(false);
const page = ref(1);
const size = ref(10);
const hasMore = ref(true);
const order = ref<"time" | "score">("time");

const fetchPosts = async () => {
  loading.value = true;
  try {
    const res = await postApi.getList2({
      page: page.value,
      size: size.value,
      order: order.value,
    });
    posts.value = res.data;
    hasMore.value = res.data.length === size.value;
  } catch (error) {
    console.error("获取帖子列表失败:", error);
  } finally {
    loading.value = false;
  }
};

const changeOrder = (newOrder: "time" | "score") => {
  order.value = newOrder;
  page.value = 1;
  fetchPosts();
};

const prevPage = () => {
  if (page.value > 1) {
    page.value--;
    fetchPosts();
  }
};

const nextPage = () => {
  if (hasMore.value) {
    page.value++;
    fetchPosts();
  }
};

const vote = async (postId: number, direction: number) => {
  if (!userStore.token) {
    // 提示用户登录
    return;
  }
  try {
    await postApi.vote({
      postId: postId.toString(),
      direction,
    });
    // 刷新帖子列表
    fetchPosts();
  } catch (error) {
    console.error("投票失败:", error);
  }
};

onMounted(() => {
  fetchPosts();
});

watch([page, order], () => {
  fetchPosts();
});
</script>
