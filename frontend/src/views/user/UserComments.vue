<template>
  <div class="user-comments">
    <h2>我的评论</h2>
    <ReadOnlyPostList :fetch-posts="getUserCommentedPosts" />
  </div>
</template>

<script setup>
import ReadOnlyPostList from '@/components/ReadOnlyPostList.vue'
import { getUserCommentedPosts } from '@/api/post'
</script>

<style scoped>
.user-comments {
  padding: 20px;
}

h2 {
  margin-bottom: 20px;
  font-size: 24px;
  color: #333;
}
</style> 