<template>
  <div class="user-posts">
    <h2>我的帖子</h2>
    <UserPostList :fetch-posts="getUserPosts" />
  </div>
</template>

<script setup>
import UserPostList from '@/components/UserPostList.vue'
import { getUserPosts } from '@/api/post'
</script> 