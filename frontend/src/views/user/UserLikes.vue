<template>
  <div class="user-likes">
    <h2>我的点赞</h2>
    <ReadOnlyPostList :fetch-posts="getUserLikedPosts" />
  </div>
</template>

<script setup>
import ReadOnlyPostList from '@/components/ReadOnlyPostList.vue'
import { getUserLikedPosts } from '@/api/post'
</script> 