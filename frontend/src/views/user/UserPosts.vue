<template>
  <div class="user-posts">
    <div class="page-header">
      <h2>我的帖子</h2>
      <el-button type="primary" @click="$router.push('/post/create')">
        发表新帖子
      </el-button>
    </div>
    <UserPostList :fetch-posts="getUserPosts" />
  </div>
</template>

<script setup>
import UserPostList from '@/components/UserPostList.vue'
import { getUserPosts } from '@/api/post'
</script>

<style scoped>
.user-posts {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

h2 {
  margin: 0;
  color: #1d2129;
}
</style> 