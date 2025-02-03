<template>
  <div class="login-container">
    <el-card class="login-card">
      <h2>{{ isAdminLogin ? '管理员登录' : '登录 TalkSphere' }}</h2>
      <p class="subtitle">欢迎回来！请登录您的账号</p>
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef" class="login-form">
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
          <div class="btn-row" :class="{ 'admin-login': isAdminLogin }">
            <el-button 
              type="primary" 
              @click="handleLogin" 
              :loading="loading"
              size="large"
              :class="{ 'admin-btn': isAdminLogin }"
            >登录</el-button>
            <template v-if="!isAdminLogin">
              <el-button 
                type="success" 
                @click="goToAdminLogin" 
                size="large"
              >管理员登录</el-button>
              <el-button 
                @click="$router.push('/register')"
                size="large"
              >注册账号</el-button>
            </template>
          </div>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useStore } from 'vuex'
import { login } from '../api/user'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { isAdmin } from '@/utils/permission'

export default {
  name: 'LoginPage',
  props: {
    isAdminLogin: {
      type: Boolean,
      default: false
    }
  },
  components: {
    User,
    Lock
  },
  setup(props) {
    const router = useRouter()
    const route = useRoute()
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
          // 直接从 data 中获取数据
          const { token, userID, username } = res.data.data
          
          store.commit('SET_TOKEN', token)
          store.commit('SET_USERINFO', { userID, username })
          ElMessage.success('登录成功')
          
          // 检查是否是管理员账号
          if (await isAdmin(userID)) { // 使用正确的 userID
            router.push('/board-manage')
            return
          }
          
          // 如果是管理员登录页面但不是管理员账号
          if (props.isAdminLogin && !(await isAdmin(userID))) { // 使用正确的 userID
            ElMessage.error('非管理员账号，请使用管理员账号登录')
            return
          }
          
          // 如果有重定向地址则跳转到重定向地址
          const redirect = route.query.redirect
          router.push(redirect || '/')
        } else {
          ElMessage.error(res.data.msg || '登录失败') // 使用 msg 而不是 message
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

    // 跳转到管理员登录页面
    const goToAdminLogin = () => {
      router.push('/admin/login')
    }

    return {
      loginForm,
      loginFormRef,
      rules,
      loading,
      handleLogin,
      goToAdminLogin
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

.login-form {
  display: flex;
  flex-direction: column;
}

.btn-container {
  margin-top: 35px;
  display: flex;
  justify-content: center;
}

.btn-row {
  display: flex;
  gap: 12px;
  justify-content: center;
  width: 100%;
}

.btn-row:not(.admin-login) {
  justify-content: space-between;
}

.el-button {
  height: 44px;
  font-size: 16px;
  border-radius: 8px;
  font-weight: 500;
}

.btn-row:not(.admin-login) .el-button {
  flex: 1;
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

.el-button--success {
  background-color: #67C23A;
  border: none;
  transition: all 0.3s ease;
}

.el-button--success:hover {
  background-color: #85ce61;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(103, 194, 58, 0.2);
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
  
  .btn-row {
    flex-direction: column;
    gap: 12px;
    align-items: center;
  }
  
  .el-button {
    width: 100%;
  }

  .admin-btn {
    width: 80% !important;
  }
}

/* 修改管理员登录页面的按钮样式 */
.btn-row.admin-login {
  justify-content: center;
}

.admin-btn {
  width: 200px !important;
}
</style> 