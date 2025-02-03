<template>
  <div class="home-container">
    <!-- 顶部导航条 -->
    <div class="top-header">
      <div class="header-content">
        <div class="left">
          <h1 class="logo">TalkSphere</h1>
        </div>
        <div class="right">
          <el-button 
            type="primary" 
            @click="$router.push('/post/create')"
          >
            发表帖子
          </el-button>
          <div class="user-info" v-if="userInfo">
            <el-dropdown trigger="click">
              <el-avatar 
                :size="40" 
                :src="userInfo.avatar || defaultAvatar"
              />
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="$router.push('/profile')">个人资料</el-dropdown-item>
                  <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>
    </div>

    <!-- 主体内容区 -->
    <div class="main-wrapper">
      <!-- 左侧板块列表 -->
      <div class="side-nav">
        <div class="nav-item" 
          v-for="board in boards" 
          :key="board?.ID"
          :class="{ active: selectedBoardId === board?.ID }"
          @click="selectBoard(board?.ID)"
        >
          <span class="board-name">{{ board?.Name }}</span>
          <span class="post-count">({{ board?.post_count || 0 }})</span>
        </div>
      </div>

      <!-- 右侧内容区 -->
      <div class="content-area">
        <!-- 欢迎信息 -->
        <div v-if="!selectedBoardId" class="welcome-message">
          <h2>欢迎来到 TalkSphere</h2>
          <p>请从左侧选择一个板块开始浏览</p>
        </div>

        <!-- 当前板块信息 -->
        <template v-else>
          <div class="board-header" v-if="currentBoardName">
            <h2>{{ currentBoardName }}</h2>
          </div>
        </template>

        <!-- 帖子列表 -->
        <div class="post-list">
          <div v-for="post in posts" 
            :key="post.id" 
            class="post-item"
            @click="viewPost(post.id)"
          >
            <div class="post-main">
              <div class="post-user">
                <el-avatar 
                  :size="40" 
                  :src="post.author?.avatar_url || defaultAvatar"
                />
                <span class="username">{{ post.author?.username || '未知用户' }}</span>
                <span class="post-time">{{ formatDate(post.created_at) }}</span>
              </div>
              
              <div class="post-content">
                <h3 class="post-title">{{ post.title }}</h3>
                <p class="post-excerpt">{{ post.content?.slice(0, 100) }}...</p>
              </div>

              <div class="post-meta">
                <span class="meta-item">
                  <el-icon><View /></el-icon>
                  {{ post.view_count }}
                </span>
                <span class="meta-item">
                  👍
                  {{ post.like_count }}
                </span>
                <span class="meta-item">
                  📚
                  {{ post.favorite_count }}
                </span>
                <span class="meta-item">
                  <el-icon><ChatLineRound /></el-icon>
                  {{ post.comment_count }}
                </span>
              </div>
            </div>

            <div class="post-cover" v-if="post.images?.length">
              <el-image :src="post.images[0].ImageURL" fit="cover" />
            </div>
            <div class="post-content" v-else>
              <p class="post-excerpt">{{ post.content?.slice(0, 100) }}...</p>
            </div>
          </div>
        </div>

        <!-- 分页器 -->
        <div class="pagination-wrapper" v-if="total > 0">
          <span class="total-text">共 {{ total }} 条</span>
          <el-pagination
            v-model:current-page="currentPage"
            :page-size="pageSize"
            :total="total"
            layout="total, prev, pager, next"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getAllBoards } from '../api/board'
import { getBoardPosts } from '../api/post'
import { getUserProfile } from '../api/user'
import dayjs from 'dayjs'
import { View, ChatLineRound } from '@element-plus/icons-vue'

export default {
  name: 'HomePage',
  components: {
    View,
    ChatLineRound
  },
  setup() {
    const store = useStore()
    const router = useRouter()
    const posts = ref([])
    const boards = ref([])
    const currentPage = ref(1)
    const pageSize = ref(20)
    const total = ref(0)
    const selectedBoardId = ref(null)
    const currentBoardName = ref('')
    const defaultAvatar = ref('https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png')

    const loadPosts = async () => {
      try {
        if (!selectedBoardId.value) {
          posts.value = []
          total.value = 0
          return
        }

        const params = {
          page: currentPage.value,
          size: pageSize.value
        }
        
        const res = await getBoardPosts(selectedBoardId.value, params)
        console.log('API返回数据:', res.data)
        
        if (res.data.code === 1000) {
          posts.value = res.data.data.posts
          total.value = res.data.data.total
          console.log('设置的总数:', total.value)
        }
      } catch (error) {
        console.error('获取帖子列表失败:', error)
        ElMessage.error('获取帖子列表失败')
      }
    }

    const loadBoards = async () => {
      try {
        const res = await getAllBoards()
        if (res.data.code === 1000) {
          // 获取所有板块
          const boardsData = res.data.data || []
          
          // 为每个板块获取其帖子数量
          for (let board of boardsData) {
            try {
              const postsRes = await getBoardPosts(board.ID, { page: 1, size: 1 })
              if (postsRes.data.code === 1000) {
                board.post_count = postsRes.data.data.total || 0
              }
            } catch (error) {
              console.error(`获取板块 ${board.ID} 帖子数量失败:`, error)
              board.post_count = 0
            }
          }
          
          boards.value = boardsData
        }
      } catch (error) {
        console.error('获取板块列表失败:', error)
        boards.value = []
      }
    }

    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }

    const viewPost = (postId) => {
      if (postId) {
        router.push(`/post/${postId}`)
      }
    }

    const selectBoard = (boardId) => {
      if (boardId) {
        selectedBoardId.value = boardId
        const selectedBoard = boards.value.find(b => b.ID === boardId)
        currentBoardName.value = selectedBoard?.Name || '未知板块'
        currentPage.value = 1
        loadPosts()
      }
    }

    const handleSizeChange = (val) => {
      pageSize.value = val
      loadPosts()
    }

    const handleCurrentChange = (val) => {
      currentPage.value = val
      loadPosts()
    }

    const handleLogout = () => {
      store.dispatch('logout')
      router.push('/login')
    }

    const loadUserInfo = async () => {
      try {
        const res = await getUserProfile()
        if (res.data.code === 1000) {
          store.commit('SET_USERINFO', res.data.data)
        }
      } catch (error) {
        console.error('获取用户信息失败:', error)
        ElMessage.error('获取用户信息失败')
      }
    }

    onMounted(async () => {
      await loadBoards()
      // 如果有板块数据，自动选择第一个板块
      if (boards.value.length > 0) {
        selectBoard(boards.value[0].ID)
      }
      loadUserInfo()
    })

    return {
      posts,
      boards,
      currentPage,
      pageSize,
      total,
      selectedBoardId,
      currentBoardName,
      defaultAvatar,
      userInfo: computed(() => store.state.userInfo),
      formatDate,
      viewPost,
      selectBoard,
      handleSizeChange,
      handleCurrentChange,
      handleLogout
    }
  }
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background-color: #fff;
}

.top-header {
  height: 60px;
  background-color: #fff;
  border-bottom: 1px solid #e4e6eb;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
}

.header-content {
  max-width: 1440px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
}

.logo {
  color: #1e80ff;
  font-size: 24px;
  font-weight: 600;
  margin: 0;
}

.right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.main-wrapper {
  display: flex;
  padding-top: 60px;
  max-width: 1440px;
  margin: 0 auto;
}

.side-nav {
  width: 240px;
  padding: 20px 0;
  border-right: 1px solid #e4e6eb;
  height: calc(100vh - 60px);
  position: fixed;
  background: #fff;
}

.nav-item {
  padding: 12px 24px;
  cursor: pointer;
  color: #333;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nav-item:hover, .nav-item.active {
  color: #1e80ff;
  background-color: #f4f5f5;
}

.content-area {
  flex: 1;
  margin-left: 240px;
  padding: 20px;
}

.operation-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #fff;
  border-radius: 4px;
  margin-bottom: 16px;
}

.tabs {
  display: flex;
  gap: 24px;
}

.tab-item {
  font-size: 14px;
  color: #86909c;
  cursor: pointer;
  padding: 4px 12px;
}

.tab-item.active {
  color: #1e80ff;
  font-weight: 500;
}

.post-list {
  background: #fff;
  border-radius: 4px;
}

.post-item {
  padding: 16px;
  border-bottom: 1px solid #e4e6eb;
  cursor: pointer;
  display: flex;
  gap: 16px;
}

.post-item:hover {
  background-color: #fafafa;
}

.post-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.post-user {
  display: flex;
  align-items: center;
  gap: 12px;
}

.username {
  font-size: 14px;
  font-weight: 500;
  color: #1d2129;
}

.post-time {
  font-size: 12px;
  color: #86909c;
}

.post-title {
  font-size: 16px;
  font-weight: 500;
  color: #1d2129;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.post-excerpt {
  font-size: 14px;
  color: #86909c;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.post-meta {
  display: flex;
  align-items: center;
  gap: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #86909c;
  font-size: 13px;
}

.meta-item:hover {
  color: #1e80ff;
}

.post-cover {
  width: 200px;
  height: 150px;
  border-radius: 4px;
  overflow: hidden;
}

.post-cover .el-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
}

.total-text {
  color: #606266;
  font-size: 14px;
}

.page-size-select {
  width: 120px;
}

:deep(.el-pagination) {
  justify-content: flex-start;
}

.create-post-btn {
  border-radius: 4px;
}

.board-name {
  font-size: 14px;
}

.post-count {
  font-size: 12px;
  color: #909399;
}

.board-header {
  margin-bottom: 20px;
  padding: 16px;
  background: #fff;
  border-radius: 4px;
}

.board-header h2 {
  margin: 0;
  font-size: 20px;
  color: #1d2129;
}

.welcome-message {
  text-align: center;
  padding: 100px 0;
}

.welcome-message h2 {
  color: #1e80ff;
  font-size: 28px;
  margin-bottom: 16px;
}

.welcome-message p {
  color: #86909c;
  font-size: 16px;
}
</style> 