<template>
    <div class="create-post-container">
      <el-card class="post-card">
        <template #header>
          <div class="card-header">
            <h2>发表新帖子</h2>
          </div>
        </template>
  
        <el-form :model="postForm" :rules="rules" ref="postFormRef">
          <el-form-item prop="title">
            <el-input 
              v-model="postForm.title" 
              placeholder="请输入标题"
              maxlength="100"
              show-word-limit
            />
          </el-form-item>
  
          <el-form-item prop="board_id">
            <el-select 
              v-model="postForm.board_id" 
              placeholder="请选择板块"
              style="width: 100%"
              :loading="loadingBoards"
            >
              <el-option 
                v-for="board in boards" 
                :key="board.ID" 
                :label="board.Name" 
                :value="board.ID"
              />
            </el-select>
          </el-form-item>
  
          <div v-if="false">
            <pre>{{ JSON.stringify(boards, null, 2) }}</pre>
          </div>
  
          <el-form-item prop="content">
            <el-input 
              v-model="postForm.content" 
              type="textarea" 
              :rows="6"
              placeholder="请输入内容"
              maxlength="10000"
              show-word-limit
            />
          </el-form-item>
  
          <el-form-item label="标签">
            <el-select
              v-model="postForm.tags"
              multiple
              filterable
              allow-create
              default-first-option
              placeholder="请选择或创建标签"
              style="width: 100%"
            >
              <el-option
                v-for="tag in commonTags"
                :key="tag"
                :label="tag"
                :value="tag"
              />
            </el-select>
          </el-form-item>
  
          <el-form-item label="图片">
            <el-upload
              :action="`${baseURL}/posts/image`"
              :headers="uploadHeaders"
              list-type="picture-card"
              :on-success="handleImageSuccess"
              :on-remove="handleImageRemove"
              :before-upload="beforeImageUpload"
            >
              <el-icon><Plus /></el-icon>
            </el-upload>
          </el-form-item>
  
          <el-form-item>
            <el-button 
              type="primary" 
              :loading="loading"
              @click="submitPost"
            >发布帖子</el-button>
            <el-button @click="$router.push('/')">取消</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </template>
  
  <script>
  import { ref, onMounted, computed } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { Plus } from '@element-plus/icons-vue'
  import { getAllBoards } from '../api/board'
  import { createPost } from '../api/post'
  
  export default {
    name: 'CreatePost',
    components: { Plus },
    setup() {
      const router = useRouter()
      const postFormRef = ref(null)
      const loading = ref(false)
      const boards = ref([])
      const baseURL = 'http://127.0.0.1:8989'
      const uploadedImages = ref([])
      const loadingBoards = ref(false)
  
      const postForm = ref({
        title: '',
        content: '',
        board_id: '',
        tags: [],
        image_ids: []
      })
  
      const commonTags = ref(['讨论', '提问', '分享', '建议'])
  
      const rules = {
        title: [
          { required: true, message: '请输入标题', trigger: 'blur' },
          { min: 2, max: 100, message: '标题长度在2到100个字符之间', trigger: 'blur' }
        ],
        content: [
          { required: true, message: '请输入内容', trigger: 'blur' },
          { min: 10, message: '内容至少10个字符', trigger: 'blur' }
        ],
        board_id: [
          { required: true, message: '请选择板块', trigger: 'change' }
        ]
      }
  
      const uploadHeaders = computed(() => ({
        Authorization: localStorage.getItem('token')
      }))
  
      const loadBoards = async () => {
        loadingBoards.value = true
        try {
          const res = await getAllBoards()
          console.log('板块响应详细信息:', {
            完整响应: res,
            数据: res.data,
            数据类型: typeof res.data,
            是否数组: Array.isArray(res.data)
          })
          
          if (Array.isArray(res.data)) {
            boards.value = res.data
          } else if (res.data && Array.isArray(res.data.data)) {
            boards.value = res.data.data
          } else {
            console.error('意外的数据格式:', res.data)
            ElMessage.error('获取板块列表失败：数据格式错误')
          }

          if (boards.value.length === 0) {
            ElMessage.warning('暂无可用板块')
          } else {
            console.log('处理后的板块数据:', boards.value)
          }
        } catch (error) {
          console.error('获取板块列表错误:', error)
          ElMessage.error('获取板块列表失败，请稍后重试')
        } finally {
          loadingBoards.value = false
        }
      }
  
      const handleImageSuccess = (res) => {
        if (res.code === 1000) {
          uploadedImages.value.push({
            id: res.data.image_id,
            url: res.data.image_url
          })
          postForm.value.image_ids.push(res.data.image_id)
        } else {
          ElMessage.error('图片上传失败')
        }
      }
  
      const handleImageRemove = (file) => {
        const index = uploadedImages.value.findIndex(img => img.url === file.url)
        if (index !== -1) {
          uploadedImages.value.splice(index, 1)
          postForm.value.image_ids.splice(index, 1)
        }
      }
  
      const beforeImageUpload = (file) => {
        const isValidType = ['image/jpeg', 'image/png', 'image/gif'].includes(file.type)
        const isLt20M = file.size / 1024 / 1024 < 20
  
        if (!isValidType) {
          ElMessage.error('只能上传 JPG/PNG/GIF 格式的图片!')
        }
        if (!isLt20M) {
          ElMessage.error('图片大小不能超过 20MB!')
        }
        return isValidType && isLt20M
      }
  
      const submitPost = async () => {
        if (!postFormRef.value) return
        
        await postFormRef.value.validate(async (valid) => {
          if (valid) {
            loading.value = true
            try {
              const res = await createPost(postForm.value)
              if (res.data.code === 1000) {
                ElMessage.success('发布成功')
                router.push('/')
              } else {
                ElMessage.error(res.data.msg || '发布失败')
              }
            } catch (error) {
              ElMessage.error('发布失败')
            } finally {
              loading.value = false
            }
          }
        })
      }
  
      onMounted(() => {
        loadBoards()
      })
  
      return {
        postForm,
        postFormRef,
        loading,
        boards,
        rules,
        commonTags,
        baseURL,
        uploadHeaders,
        uploadedImages,
        handleImageSuccess,
        handleImageRemove,
        beforeImageUpload,
        submitPost,
        loadingBoards,
        loadBoards
      }
    }
  }
  </script>
  
  <style scoped>
  .create-post-container {
    max-width: 800px;
    margin: 20px auto;
    padding: 0 20px;
  }
  
  .post-card {
    margin-bottom: 20px;
  }
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .card-header h2 {
    margin: 0;
    font-size: 18px;
    color: #303133;
  }
  
  :deep(.el-upload--picture-card) {
    width: 100px;
    height: 100px;
    line-height: 100px;
  }
  
  :deep(.el-form-item__label) {
    font-weight: 500;
  }
  </style>