<template>
  <div class="create-post-container">
    <el-card class="create-post-card">
      <h2>发表新帖子</h2>
      <el-form 
        :model="postForm" 
        :rules="rules" 
        ref="postFormRef"
        label-position="top"
      >
        <el-form-item label="标题" prop="title">
          <el-input 
            v-model="postForm.title" 
            placeholder="请输入标题（至少3个字符）"
          />
        </el-form-item>

        <el-form-item label="内容" prop="content">
          <el-input 
            v-model="postForm.content" 
            type="textarea" 
            :rows="6"
            placeholder="请输入内容（至少10个字符）"
          />
        </el-form-item>

        <el-form-item label="选择板块" prop="board_id">
          <el-select 
            v-model="postForm.board_id" 
            placeholder="请选择板块"
            style="width: 100%"
          >
            <el-option 
              v-for="board in boards"
              :key="board.ID"
              :label="board.Name"
              :value="board.ID"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="标签" prop="tags">
          <el-select
            v-model="postForm.tags"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="请输入标签（可选）"
            style="width: 100%"
          >
          </el-select>
        </el-form-item>

        <el-form-item label="上传图片">
          <el-upload
            class="upload-demo"
            :action="`${baseURL}/api/posts/image`"
            :headers="uploadHeaders"
            :on-success="handleImageSuccess"
            :on-remove="handleImageRemove"
            :file-list="uploadedImages"
            list-type="picture"
            name="image"
            :data="{ type: 'post' }"
          >
            <el-button type="primary">上传图片</el-button>
            <template #tip>
              <div class="el-upload__tip">支持jpg/png文件</div>
            </template>
          </el-upload>
        </el-form-item>

        <el-form-item>
          <el-button 
            type="primary" 
            @click="submitPost"
            :loading="loading"
          >
            发表帖子
          </el-button>
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
import { getAllBoards } from '../api/board'
import { createPost } from '../api/post'

export default {
  name: 'CreatePost',
  setup() {
    const router = useRouter()
    const postFormRef = ref(null)
    const loading = ref(false)
    const boards = ref([])
    const baseURL = 'http://127.0.0.1:8989'
    const uploadedImages = ref([])

    const postForm = ref({
      title: '',
      content: '',
      board_id: '',
      tags: [],
      image_ids: []
    })

    const rules = {
      title: [
        { required: true, message: '请输入标题', trigger: 'blur' },
        { min: 3, message: '标题至少3个字符', trigger: 'blur' }
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
      loading.value = true
      try {
        const res = await getAllBoards()
        if (Array.isArray(res.data)) {
          boards.value = res.data
        } else if (res.data && Array.isArray(res.data.data)) {
          boards.value = res.data.data
        }
      } catch (error) {
        console.error('获取板块列表错误:', error)
        ElMessage.error('获取板块列表失败')
      } finally {
        loading.value = false
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

    const submitPost = async () => {
      if (!postFormRef.value) return
      
      await postFormRef.value.validate(async (valid) => {
        if (valid) {
          loading.value = true
          try {
            console.log('提交的数据:', postForm.value)
            const res = await createPost(postForm.value)
            console.log('服务器响应:', res)
            if (res.data.code === 1000) {
              ElMessage.success('发表成功')
              router.push('/')
            } else {
              ElMessage.error(res.data.msg || '发表失败')
            }
          } catch (error) {
            console.error('发表帖子失败:', error)
            ElMessage.error('发表失败')
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
      rules,
      loading,
      boards,
      baseURL,
      uploadedImages,
      uploadHeaders,
      handleImageSuccess,
      handleImageRemove,
      submitPost
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

.create-post-card {
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
  padding: 20px;
}

h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #303133;
}

.el-form-item {
  margin-bottom: 25px;
}

.el-upload__tip {
  color: #909399;
  font-size: 12px;
  margin-top: 5px;
}
</style>