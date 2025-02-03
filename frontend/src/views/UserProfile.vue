<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <div class="user-info">
        <h2>{{ userProfile.username }}</h2>
        <p>{{ userProfile.email }}</p>
      </div>
      <div class="avatar-container">
        <el-avatar 
          :size="100" 
          :src="userProfile.avatar || defaultAvatar"
          @error="() => true"
        />
        <el-upload
          class="avatar-uploader"
          :action="`${baseURL}/api/avatar`"
          :headers="uploadHeaders"
          :show-file-list="false"
          :on-success="handleAvatarSuccess"
          :before-upload="beforeAvatarUpload"
          name="avatar"
        >
          <el-button size="small" type="primary">更换头像</el-button>
        </el-upload>
      </div>
      <div class="bio-container">
        <div class="bio-header">
          <h3>个人简介</h3>
          <el-button 
            v-if="!isEditingBio" 
            type="primary" 
            link
            @click="editBio"
          >编辑</el-button>
        </div>
        <div v-if="!isEditingBio" class="bio-content">
          {{ userProfile.bio || '暂无简介' }}
        </div>
        <div v-else class="bio-edit">
          <el-input
            v-model="editingBio"
            type="textarea"
            :rows="3"
            placeholder="请输入个人简介"
          />
          <div class="bio-actions">
            <el-button type="primary" @click="saveBio">保存</el-button>
            <el-button @click="cancelEditBio">取消</el-button>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { getUserProfile, updateBio } from '../api/user'
import { ElMessage } from 'element-plus'
import { useStore } from 'vuex'

export default {
  name: 'UserProfilePage',
  setup() {
    const store = useStore()
    const baseURL = 'http://127.0.0.1:8989'
    const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
    const userProfile = ref({})
    const isEditingBio = ref(false)
    const editingBio = ref('')

    const uploadHeaders = computed(() => ({
      Authorization: localStorage.getItem('token')
    }))

    const uploadData = ref({})

    const loadUserProfile = async () => {
      try {
        const res = await getUserProfile()
        console.log('Profile response:', res)  // 添加日志
        if (res.data.code === 1000) {
          userProfile.value = res.data.data
          editingBio.value = res.data.data.bio || ''
          // 更新 store 中的用户信息
          store.commit('SET_USERINFO', res.data.data)
        } else {
          ElMessage.error(res.data.msg || '获取用户信息失败')
        }
      } catch (error) {
        console.error('获取用户信息错误:', error)
        ElMessage.error('获取用户信息失败')
      }
    }

    const handleAvatarSuccess = (res) => {
      console.log('Avatar upload response:', res)  // 添加日志
      if (res.code === 1000) {
        userProfile.value.avatar_url = res.data.avatar_url
        // 更新 store 中的头像
        store.commit('SET_USERINFO', {
          ...store.state.userInfo,
          avatar_url: res.data.avatar_url
        })
        ElMessage.success('头像更新成功')
      } else {
        ElMessage.error(res.msg || '上传失败')
      }
    }

    const beforeAvatarUpload = (file) => {
      const isValidType = ['image/jpeg', 'image/png'].includes(file.type)
      const isLt2M = file.size / 1024 / 1024 < 20

      if (!isValidType) {
        ElMessage.error('头像只能是 JPG 或 PNG 格式!')
      }
      if (!isLt2M) {
        ElMessage.error('头像大小不能超过 20MB!')
      }
      return isValidType && isLt2M
    }

    const editBio = () => {
      editingBio.value = userProfile.value.bio
      isEditingBio.value = true
    }

    const saveBio = async () => {
      try {
        const res = await updateBio({ bio: editingBio.value })
        if (res.data.code === 1000) {
          userProfile.value.bio = editingBio.value
          // 更新 store 中的用户信息
          store.commit('SET_USERINFO', {
            ...store.state.userInfo,
            bio: editingBio.value
          })
          isEditingBio.value = false
          ElMessage.success('个人简介更新成功')
        } else {
          ElMessage.error(res.data.msg || '更新失败')
        }
      } catch (error) {
        console.error('更新失败:', error)
        ElMessage.error('更新失败')
      }
    }

    const cancelEditBio = () => {
      editingBio.value = userProfile.value.bio || ''
      isEditingBio.value = false
    }

    onMounted(() => {
      loadUserProfile()
    })

    return {
      baseURL,
      defaultAvatar,
      userProfile,
      isEditingBio,
      editingBio,
      uploadHeaders,
      uploadData,
      handleAvatarSuccess,
      beforeAvatarUpload,
      editBio,
      saveBio,
      cancelEditBio
    }
  }
}
</script>

<style scoped>
.profile-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.profile-card {
  padding: 20px;
}

.avatar-container {
  text-align: center;
  margin: 20px 0;
}

.bio-container {
  margin-top: 20px;
}

.bio-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.bio-content {
  color: #666;
  line-height: 1.6;
}

.bio-edit {
  margin-top: 10px;
}

.bio-actions {
  margin-top: 10px;
  text-align: right;
}

.user-info {
  text-align: center;
}

.user-info h2 {
  margin: 0;
  color: #303133;
}

.user-info p {
  margin: 5px 0;
  color: #909399;
}
</style> 