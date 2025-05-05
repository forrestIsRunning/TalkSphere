<template>
  <div class="edit-post-container">
    <div class="edit-header">
      <h2>编辑帖子</h2>
      <el-button @click="$router.back()">返回</el-button>
    </div>

    <div class="edit-form" v-loading="loading">
      <el-form :model="postForm" :rules="rules" ref="postFormRef">
        <el-form-item label="标题" prop="title">
          <el-input v-model="postForm.title" placeholder="请输入标题" />
        </el-form-item>

        <el-form-item label="内容" prop="content">
          <div class="editor-container">
            <Toolbar
              style="border-bottom: 1px solid #ccc"
              :editor="editorRef"
              :defaultConfig="toolbarConfig"
              :mode="mode"
            />
            <Editor
              style="height: 500px; overflow-y: hidden;"
              v-model="postForm.content"
              :defaultConfig="editorConfig"
              :mode="mode"
              @onCreated="handleCreated"
            />
          </div>
        </el-form-item>

        <el-form-item label="标签" prop="tags">
          <el-select
            v-model="postForm.tags"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="请选择或创建标签"
          >
            <el-option
              v-for="tag in existingTags"
              :key="tag"
              :label="tag"
              :value="tag"
            />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            保存修改
          </el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, shallowRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getPostDetail, updatePost } from '@/api/post'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const submitting = ref(false)
const editorRef = shallowRef()
const existingTags = ref([])
const mode = ref('default')

const postForm = ref({
  title: '',
  content: '',
  tags: [],
  image_ids: []
})

const rules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { min: 3, max: 100, message: '标题长度在 3 到 100 个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入内容', trigger: 'blur' },
    { min: 10, message: '内容至少 10 个字符', trigger: 'blur' }
  ]
}

// 编辑器配置
const toolbarConfig = {
  excludeKeys: [
    'insertVideo',
    'uploadVideo',
    'group-video',
    'insertTable'
  ]
}

const editorConfig = {
  placeholder: '请输入内容...',
  MENU_CONF: {
    uploadImage: {
      server: '/api/posts/image',
      fieldName: 'image',
      headers: {
        Authorization: localStorage.getItem('token')
      },
      maxFileSize: 20 * 1024 * 1024,
      maxNumberOfFiles: 10,
      allowedFileTypes: ['image/jpeg', 'image/png', 'image/gif'],
      metaWithUrl: true,
      customInsert: (res, insertFn) => {
        if (res.code === 1000) {
          const { image_url, image_id } = res.data
          postForm.value.image_ids.push(image_id)
          insertFn(image_url)
        } else {
          ElMessage.error('图片上传失败')
        }
      },
      onError: (err) => {
        console.error('图片上传错误:', err)
        ElMessage.error('图片上传失败')
      }
    }
  }
}

const handleCreated = (editor) => {
  editorRef.value = editor
}

// 加载帖子详情
const loadPost = async () => {
  loading.value = true
  try {
    const postId = route.params.id
    const res = await getPostDetail(postId)
    if (res.data.code === 1000) {
      const post = res.data.data
      postForm.value.title = post.title
      postForm.value.content = post.content
      postForm.value.tags = post.tags?.map(tag => tag.Name) || []
      postForm.value.image_ids = post.images?.map(img => img.ID) || []
      
      // 收集已有的标签
      if (post.tags) {
        existingTags.value = [...new Set([...existingTags.value, ...post.tags.map(tag => tag.Name)])]
      }
    }
  } catch (error) {
    console.error('加载帖子失败:', error)
    ElMessage.error('加载帖子失败')
  } finally {
    loading.value = false
  }
}

// 提交表单
const submitForm = async () => {
  if (!postForm.value.title.trim()) {
    ElMessage.warning('请输入标题')
    return
  }

  if (!postForm.value.content.trim()) {
    ElMessage.warning('请输入内容')
    return
  }

  submitting.value = true
  try {
    const postId = route.params.id
    const res = await updatePost(postId, {
      title: postForm.value.title,
      content: postForm.value.content,
      tags: postForm.value.tags,
      image_ids: postForm.value.image_ids
    })

    if (res.data.code === 1000) {
      ElMessage.success('更新成功')
      router.push('/user/posts')
    } else {
      ElMessage.error(res.data.msg || '更新失败')
    }
  } catch (error) {
    console.error('更新帖子失败:', error)
    ElMessage.error('更新帖子失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadPost()
})

onBeforeUnmount(() => {
  if (editorRef.value) {
    editorRef.value.destroy()
  }
})
</script>

<style scoped>
.edit-post-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.edit-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.edit-header h2 {
  margin: 0;
  color: #1d2129;
}

.editor-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

:deep(.w-e-toolbar) {
  border-top-left-radius: 4px;
  border-top-right-radius: 4px;
}

:deep(.w-e-text-container) {
  border-bottom-left-radius: 4px;
  border-bottom-right-radius: 4px;
}

.el-form-item {
  margin-bottom: 20px;
}
</style> 