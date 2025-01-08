<template>
  <div class="login-container">
    <el-card class="login-card">
      <h2>登录 TalkSphere</h2>
      <p class="subtitle">欢迎回来！请登录您的账号</p>
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef">
        <el-form-item prop="username">
          <el-input 
            v-model="loginForm.username" 
            placeholder="请输入用户名"
            size="large"
          >
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input 
            v-model="loginForm.password" 
            type="password" 
            placeholder="请输入密码"
            size="large"
          >
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item class="btn-container">
          <el-button 
            type="primary" 
            @click="handleLogin" 
            :loading="loading"
            size="large"
          >登录</el-button>
          <div class="btn-spacer"></div>
          <el-button 
            @click="$router.push('/register')"
            size="large"
          >注册账号</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { login } from '../api/user'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'

export default {
  name: 'LoginPage',
  components: {
    User,
    Lock
  },
  setup() {
    const router = useRouter()
    const store = useStore()
    const loginFormRef = ref(null)
    const loading = ref(false)
    const loginForm = ref({
      username: '',
      password: ''
    })

    const rules = {
      username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
      password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
    }

    const handleLogin = async () => {
      if (!loginFormRef.value) return
      
      await loginFormRef.value.validate(async (valid) => {
        if (valid) {
          loading.value = true
          try {
            const res = await login(loginForm.value)
            if (res.data.code === 1000) {
              const token = res.data.data.token
              store.commit('SET_TOKEN', token)
              
              store.commit('SET_USERINFO', {
                userID: res.data.data.userID,
                username: res.data.data.username,
                avatar_url: res.data.data.avatar_url,
                bio: res.data.data.bio
              })

              ElMessage.success('登录成功')
              await router.push('/')
            } else {
              ElMessage.error(res.data.msg || '登录失败')
            }
          } catch (error) {
            console.error('登录错误:', error)
            ElMessage.error('登录失败')
          } finally {
            loading.value = false
          }
        }
      })
    }

    return {
      loginForm,
      loginFormRef,
      rules,
      loading,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f7fa;
  padding: 20px;
}

.login-card {
  width: 440px;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  background: #ffffff;
  padding: 40px;
}

h2 {
  text-align: center;
  margin-bottom: 12px;
  color: #303133;
  font-size: 28px;
  font-weight: 600;
}

.subtitle {
  text-align: center;
  color: #909399;
  font-size: 14px;
  margin-bottom: 35px;
}

.el-form-item {
  margin-bottom: 25px;
}

:deep(.el-input__wrapper) {
  padding: 4px 11px;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #409EFF inset;
}

:deep(.el-input__inner) {
  height: 42px;
  font-size: 14px;
}

:deep(.el-input__prefix-inner) {
  font-size: 18px;
  color: #909399;
}

.btn-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-top: 35px;
}

.btn-spacer {
  height: 10px;
}

.el-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
  border-radius: 8px;
  font-weight: 500;
}

.el-button--primary {
  background-color: #409EFF;
  border: none;
  transition: all 0.3s ease;
}

.el-button--primary:hover {
  background-color: #66b1ff;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
}

.el-button--default {
  border: 1px solid #dcdfe6;
  color: #606266;
}

.el-button--default:hover {
  color: #409EFF;
  border-color: #c6e2ff;
  background-color: #ecf5ff;
}

:deep(.el-input__prefix) {
  display: flex;
  align-items: center;
}

:deep(.el-icon) {
  font-size: 18px;
  color: #909399;
}

@media (max-width: 480px) {
  .login-card {
    width: 100%;
    padding: 30px 20px;
  }
  
  h2 {
    font-size: 24px;
  }
}
</style> 