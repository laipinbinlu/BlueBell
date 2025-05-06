<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">社区列表</h1>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="community in communities"
        :key="community.id"
        class="card bg-base-100 shadow-xl"
      >
        <div class="card-body">
          <h2 class="card-title">{{ community.name }}</h2>
          <p>{{ community.introduction }}</p>
          <div class="card-actions justify-end">
            <router-link
              :to="'/community/' + community.id"
              class="btn btn-primary"
            >
              进入社区
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getCommunityList } from '@/api/community'

const communities = ref([])

onMounted(async () => {
  try {
    const response = await getCommunityList()
    communities.value = response.data
  } catch (error) {
    console.error('获取社区列表失败:', error)
  }
})
</script> 