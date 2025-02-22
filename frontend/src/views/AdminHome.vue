<template>
  <el-container class="admin-container">
    <!-- 左侧导航栏 -->
    <el-aside :width="isCollapse ? '64px' : '200px'" class="aside">
      <div class="logo">
        <h2>TalkSphere Admin</h2>
      </div>
      <el-menu
        :default-openeds="['1', '2', '3']"
        :collapse="isCollapse"
        :collapse-transition="false"
        :router="true"
        :default-active="activeMenu"
        class="el-menu-vertical"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
      >
        <el-menu-item index="/admin/dashboard">
          <el-icon><DataLine /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>

        <el-sub-menu index="1">
          <template #title>
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </template>
          <el-menu-item index="/admin/users">用户列表</el-menu-item>
          <el-menu-item index="/admin/permissions">权限管理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="2">
          <template #title>
            <el-icon><Grid /></el-icon>
            <span>板块管理</span>
          </template>
          <el-menu-item index="/admin/boards">板块列表</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="3">
          <template #title>
            <el-icon><TrendCharts /></el-icon>
            <span>数据分析</span>
          </template>
          <el-menu-item index="/admin/analytics/users/growth">用户增长</el-menu-item>
          <el-menu-item index="/admin/analytics/active/users">最近活跃用户</el-menu-item>

          <el-menu-item index="/admin/analytics/posts/growth">帖子增长</el-menu-item>
          <el-menu-item index="/admin/analytics/hot/posts">最近活跃帖子</el-menu-item>
          <el-menu-item index="/admin/analytics/wordcloud">词云图</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <!-- 右侧内容区 -->
    <el-container>
      <!-- 顶部导航栏 -->
      <el-header class="header">
        <div class="header-left">
          <el-icon class="toggle-sidebar" @click="toggleSidebar">
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/admin/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>{{ currentPath }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="el-dropdown-link">
              {{ userInfo.username }}
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主要内容区 -->
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter, useRoute } from 'vue-router'
import {
  DataLine,
  User,
  Grid,
  TrendCharts,
  Fold,
  Expand,
  ArrowDown
} from '@element-plus/icons-vue'

const store = useStore()
const router = useRouter()
const route = useRoute()
const isCollapse = ref(false)
const activeMenu = ref('/admin/dashboard')

// 计算当前路径名称
const currentPath = computed(() => {
  const pathMap = {
    '/admin/dashboard': '仪表盘',
    '/admin/users': '用户列表',
    '/admin/permissions': '权限管理',
    '/admin/boards': '板块管理',
    '/admin/analytics/users/growth': '用户增长分析',
    '/admin/analytics/active/users': '最近活跃用户',

    '/admin/analytics/active/posts/growth': '帖子增长分析',
    '/admin/analytics/hot/posts': '最近活跃帖子',

    '/admin/analytics/wordcloud': '词云图'
  }
  return pathMap[route.path] || '仪表盘'
})

const userInfo = computed(() => store.state.userInfo)

const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
}

const handleLogout = () => {
  store.commit('CLEAR_USER_INFO')
  router.push('/admin/login')
}
</script>

<style scoped>
.admin-container {
  height: 100vh;
}

.aside {
  background-color: #304156;
  transition: width 0.3s;
  overflow-x: hidden;
}

.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  color: #fff;
}

.logo h2 {
  margin: 0;
  font-size: 18px;
}

.el-menu-vertical {
  border-right: none;
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 60px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.toggle-sidebar {
  font-size: 20px;
  cursor: pointer;
}

.header-right {
  display: flex;
  align-items: center;
}

.el-dropdown-link {
  cursor: pointer;
  display: flex;
  align-items: center;
  color: #606266;
}

.main {
  background-color: #f0f2f5;
  padding: 20px;
  height: calc(100vh - 60px);
  overflow-y: auto;
}

:deep(.el-menu) {
  border-right: none;
}

:deep(.el-menu-item.is-active) {
  background-color: #263445;
}

:deep(.el-sub-menu__title:hover) {
  background-color: #263445 !important;
}

:deep(.el-menu-item:hover) {
  background-color: #263445 !important;
}
</style> 