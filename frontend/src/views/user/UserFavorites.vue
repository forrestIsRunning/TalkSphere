<template>
  <div class="user-favorites">
    <h2>我的收藏</h2>
    <UserPostList :fetch-posts="getUserFavoritePosts" />
  </div>
</template>

<script setup>
import UserPostList from '@/components/UserPostList.vue'
import { getUserFavoritePosts } from '@/api/post'
</script> 