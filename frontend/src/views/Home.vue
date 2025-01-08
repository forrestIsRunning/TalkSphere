<template>
  <div class="home-container">
    <div class="header">
      <h1>欢迎来到 TalkSphere</h1>
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
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { getUserProfile } from '../api/user'

export default {
  name: 'HomePage',
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
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.user-info {
  cursor: pointer;
}
</style> 