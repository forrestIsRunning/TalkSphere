<template>
  <div class="home-container">
    <!-- 顶部导航条 -->
    <div class="top-header">
      <div class="header-content">
        <div class="left">
          <h1 class="logo">TalkSphere</h1>
        </div>
        <div class="center">
          <div class="search-section">
            <el-input
              v-model="searchQuery"
              placeholder="搜索用户名或帖子内容..."
              class="search-input"
              clearable
              style="min-width: 400px"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-select v-model="searchType" class="search-type">
              <el-option label="全部" value="all" />
              <el-option label="用户名" value="username" />
              <el-option label="帖子内容" value="content" />
            </el-select>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
          </div>
        </div>
        <div class="right">
          <el-button 
            type="primary" 
            @click="$router.push('/post/create')"
          >
            发表帖子
          </el-button>
          <div class="user-info" v-if="userInfo">
            <el-dropdown trigger="click" @command="handleCommand">
              <div class="avatar-container">
                <el-avatar 
                  :size="40" 
                  :src="userInfo.avatar_url || defaultAvatar"
                />
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人资料</el-dropdown-item>
                  <el-dropdown-item command="posts">我的帖子</el-dropdown-item>
                  <el-dropdown-item command="likes">我的点赞</el-dropdown-item>
                  <el-dropdown-item command="favorites">我的收藏</el-dropdown-item>
                  <el-dropdown-item command="comments">我的评论</el-dropdown-item>
                  <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>
    </div>

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
                <p class="post-excerpt">{{ post.excerpt || formatContent(post.content) }}</p>
                <div class="post-tags" v-if="post.tags?.length">
                  <el-tag
                    v-for="tag in post.tags"
                    :key="tag.ID"
                    size="small"
                    type="info"
                    effect="plain"
                    class="post-tag"
                  >
                    {{ tag.Name }}
                  </el-tag>
                </div>
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
              <el-image 
                :src="post.images[0].url" 
                fit="cover"
                :preview-src-list="post.images.map(img => img.url)"
                :hide-on-click-modal="true"
              >
                <template #error>
                  <div class="image-error">
                    <el-icon><PictureFilled /></el-icon>
                    <span>加载失败</span>
                  </div>
                </template>
                <template #placeholder>
                  <div class="image-placeholder">
                    <el-icon><Loading /></el-icon>
                    <span>加载中</span>
                  </div>
                </template>
              </el-image>
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
import dayjs from 'dayjs'
import { View, ChatLineRound, Search, PictureFilled, Loading } from '@element-plus/icons-vue'
import { getUserProfile } from '../api/user'

export default {
  name: 'HomePage',
  components: {
    View,
    ChatLineRound,
    Search,
    PictureFilled,
    Loading
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
    const searchQuery = ref('')
    const searchType = ref('all')
    const loading = ref(false)

    const loadPosts = async () => {
      loading.value = true
      try {
        if (!selectedBoardId.value) {
          posts.value = []
          total.value = 0
          return
        }

        // 构建查询参数
        const params = {
          page: currentPage.value,
          size: pageSize.value
        }
        
        // 只有当搜索条件不为空时才添加搜索参数
        if (searchQuery.value) {
          params.search_query = searchQuery.value
          params.search_type = searchType.value
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
      } finally {
        loading.value = false
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

    const handleCommand = (command) => {
      switch (command) {
        case 'profile':
          router.push('/profile')
          break
        case 'posts':
          router.push('/user/posts')
          break
        case 'likes':
          router.push('/user/likes')
          break
        case 'favorites':
          router.push('/user/favorites')
          break
        case 'comments':
          router.push('/user/comments')
          break
        case 'logout':
          store.dispatch('logout')
          router.push('/login')
          break
      }
    }

    const loadUserInfo = async () => {
      try {
        const res = await getUserProfile()
        if (res.data.code === 1000) {
          store.commit('SET_USERINFO', res.data.data)
        } else {
          // 如果返回的不是成功码，设置默认的 guest 用户信息
          setDefaultGuestInfo()
        }
      } catch (error) {
        console.error('获取用户信息失败:', error)
        // 任何错误情况下都设置默认的 guest 用户信息
        setDefaultGuestInfo()
      }
    }

    const setDefaultGuestInfo = () => {
      store.commit('SET_USERINFO', {
        id: 0,
        username: 'momo',
        email: '',
        avatar_url: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
        bio: '我是大帅哥',
        role: 'guest',
        status: 1
      })
    }

    const handleSearch = () => {
      currentPage.value = 1
      loadPosts()
    }

    const formatContent = (content) => {
      if (!content) return ''
      
      // 创建一个临时的 div 来解析 HTML
      const div = document.createElement('div')
      div.innerHTML = content
      
      // 获取所有图片标签
      const images = div.getElementsByTagName('img')
      const imageCount = images.length
      
      // 获取纯文本内容
      let text = div.textContent || div.innerText
      text = text.trim()
      
      // 限制文本长度
      const maxLength = 100
      if (text.length > maxLength) {
        text = text.substring(0, maxLength) + '...'
      }
      
      // 如果有图片，添加图片提示
      if (imageCount > 0) {
        text += ` [${imageCount}张图片]`
      }
      
      return text
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
      handleCommand,
      searchQuery,
      searchType,
      loading,
      handleSearch,
      formatContent
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

.left {
  width: 200px;
}

.center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.right {
  width: 200px;
  display: flex;
  align-items: center;
  gap: 16px;
  justify-content: flex-end;
}

.main-wrapper {
  display: flex;
  padding-top: 60px;
  max-width: 1200px;
  margin: 0 auto;
  height: calc(100vh - 60px);
}

.side-nav {
  width: 200px;
  padding: 20px 0;
  border-right: 1px solid #e4e6eb;
  height: 100%;
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
  margin-left: 200px;
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
  max-height: 180px;
  overflow: hidden;
}

.post-item:hover {
  background-color: #fafafa;
}

.post-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.post-user {
  display: flex;
  align-items: center;
  gap: 8px;
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
  color: #666;
  font-size: 14px;
  line-height: 1.5;
  margin: 8px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.post-cover {
  width: 160px;
  height: 120px;
  border-radius: 4px;
  overflow: hidden;
  flex-shrink: 0;
}

.post-cover .el-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.post-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: auto;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #86909c;
  font-size: 13px;
}

.username {
  font-size: 13px;
  font-weight: 500;
  color: #1d2129;
}

.post-time {
  font-size: 12px;
  color: #86909c;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  padding: 16px 0;
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

.search-section {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  padding: 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.search-input {
  flex: 1;
}

.search-type {
  width: 120px;
}

.avatar-container {
  cursor: pointer;
}

.post-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 8px;
}

.post-tag {
  font-size: 12px;
  padding: 0 8px;
  height: 22px;
  line-height: 20px;
  border-radius: 4px;
  background-color: #f2f3f5;
  color: #86909c;
  border: none;
}

.post-tag:hover {
  background-color: #e4e6eb;
  color: #1d2129;
}
</style> 