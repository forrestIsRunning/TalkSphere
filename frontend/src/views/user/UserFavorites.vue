<template>
  <div class="user-favorites">
    <h2>我的收藏</h2>
    <ReadOnlyPostList :fetch-posts="getUserFavoritePosts" />
  </div>
</template>

<script setup>
import ReadOnlyPostList from '@/components/ReadOnlyPostList.vue'
import { getUserFavoritePosts } from '@/api/post'
</script> 