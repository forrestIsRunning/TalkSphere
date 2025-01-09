<template>
  <div class="home-container">
    <div class="header">
      <div class="left-section">
        <h1>欢迎来到 TalkSphere</h1>
        <el-button type="primary" @click="$router.push('/post/create')">
          <el-icon><Plus /></el-icon>发表帖子
        </el-button>
      </div>
      <div class="user-info" v-if="userInfo && userInfo.username">
        <el-dropdown trigger="click">
          <el-avatar 
            :size="40" 
            :src="userInfo.avatar || userInfo.avatar_url || defaultAvatar"
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

    <div class="main-content">
      <!-- 这里可以添加帖子列表等内容 -->
    </div>
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { getUserProfile } from '../api/user'
import { Plus } from '@element-plus/icons-vue'

export default {
  name: 'HomePage',
  components: {
    Plus
  },
  setup() {
    const store = useStore()
    const router = useRouter()
    const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
    
    const userInfo = computed(() => store.state.userInfo)

    const fetchUserInfo = async () => {
      try {
        const token = localStorage.getItem('token')
        if (token && (!userInfo.value || !userInfo.value.username)) {
          const res = await getUserProfile()
          if (res.data.code === 1000) {
            store.commit('SET_USERINFO', res.data.data)
          }
        }
      } catch (error) {
        console.error('获取用户信息失败:', error)
      }
    }

    const handleLogout = () => {
      store.dispatch('logout')
      router.push('/login')
    }

    onMounted(() => {
      fetchUserInfo()
    })

    return {
      userInfo,
      defaultAvatar,
      handleLogout
    }
  }
}
</script>

<style scoped>
.home-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 10px 0;
  border-bottom: 1px solid #ebeef5;
}

.left-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.left-section h1 {
  margin: 0;
  font-size: 24px;
  color: #303133;
}

.user-info {
  cursor: pointer;
}

.el-button {
  display: flex;
  align-items: center;
  gap: 5px;
}

.main-content {
  background-color: #fff;
  border-radius: 8px;
  padding: 20px;
  min-height: calc(100vh - 140px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .home-container {
    padding: 10px;
  }

  .header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }

  .left-section {
    width: 100%;
    justify-content: space-between;
  }

  .left-section h1 {
    font-size: 20px;
  }

  .user-info {
    align-self: flex-end;
  }
}
</style> 