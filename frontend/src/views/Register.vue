<template>
  <div class="register-container">
    <el-card class="register-card">
      <h2>注册 TalkSphere</h2>
      <p class="subtitle">创建您的账号，开始畅聊</p>
      <el-form :model="registerForm" :rules="rules" ref="registerFormRef">
        <el-form-item prop="username">
          <el-input 
            v-model="registerForm.username" 
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
            v-model="registerForm.password" 
            type="password" 
            placeholder="请输入密码"
            size="large"
          >
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="email">
          <el-input 
            v-model="registerForm.email" 
            placeholder="请输入邮箱"
            size="large"
          >
            <template #prefix>
              <el-icon><Message /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item class="btn-container">
          <el-button 
            type="primary" 
            @click="handleRegister" 
            :loading="loading"
            size="large"
          >注册</el-button>
          <el-button 
            @click="$router.push('/login')"
            size="large"
          >返回登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '../api/user'
import { ElMessage } from 'element-plus'
import { User, Lock, Message } from '@element-plus/icons-vue'

export default {
  name: 'RegisterPage',
  components: {
    User,
    Lock,
    Message
  },
  setup() {
    const router = useRouter()
    const registerFormRef = ref(null)
    const loading = ref(false)

    const registerForm = reactive({
      username: '',
      password: '',
      email: ''
    })

    const rules = {
      username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
      password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
      email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
      ]
    }

    const handleRegister = async () => {
      if (!registerFormRef.value) return
      
      await registerFormRef.value.validate(async (valid) => {
        if (valid) {
          loading.value = true
          try {
            const res = await register(registerForm)
            if (res.data.code === 1000) {
              ElMessage.success('注册成功')
              router.push('/login')
            } else {
              ElMessage.error(res.data.msg)
            }
          } catch (error) {
            ElMessage.error('注册失败')
          } finally {
            loading.value = false
          }
        }
      })
    }

    return {
      registerForm,
      registerFormRef,
      rules,
      loading,
      handleRegister
    }
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f7fa;
  padding: 20px;
}

.register-card {
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
  gap: 16px;
  margin-top: 35px;
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

@media (max-width: 480px) {
  .register-card {
    width: 100%;
    padding: 30px 20px;
  }
  
  h2 {
    font-size: 24px;
  }
}

/* 添加图标样式 */
:deep(.el-input__prefix) {
  display: flex;
  align-items: center;
}

:deep(.el-icon) {
  font-size: 18px;
  color: #909399;
}
</style> 