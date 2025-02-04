<template>
  <div class="user-likes">
    <h2>我的点赞</h2>
    <UserPostList :fetch-posts="getUserLikedPosts" />
  </div>
</template>

<script setup>
import UserPostList from '@/components/UserPostList.vue'
import { getUserLikedPosts } from '@/api/post'
</script> 