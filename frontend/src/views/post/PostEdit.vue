<template>
  <div class="edit-post-container">
    <div class="edit-header">
      <h2>编辑帖子</h2>
      <el-button @click="$router.back()">返回</el-button>
    </div>

    <div class="edit-form" v-loading="loading">
      <el-form :model="postForm" label-width="80px">
        <el-form-item label="标题">
          <el-input v-model="postForm.title" placeholder="请输入标题" />
        </el-form-item>

        <el-form-item label="内容">
          <div class="editor-container">
            <div ref="editorRef"></div>
          </div>
        </el-form-item>

        <el-form-item label="标签">
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
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getPostDetail, updatePost } from '@/api/post'
import '@wangeditor/editor/dist/css/style.css'
import { createEditor } from '@wangeditor/editor'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const submitting = ref(false)
const editorRef = ref(null)
const editor = ref(null)
const existingTags = ref([])

const postForm = ref({
  title: '',
  content: '',
  tags: [],
  image_ids: []
})

// 初始化编辑器
const initEditor = () => {
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

  editor.value = createEditor({
    selector: editorRef.value,
    config: editorConfig,
    html: postForm.value.content
  })
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
      
      // 初始化编辑器并设置内容
      initEditor()
      editor.value.setHtml(post.content)
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

  const content = editor.value.getHtml()
  if (!content.trim()) {
    ElMessage.warning('请输入内容')
    return
  }

  submitting.value = true
  try {
    const postId = route.params.id
    const res = await updatePost(postId, {
      title: postForm.value.title,
      content: content,
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
  if (editor.value) {
    editor.value.destroy()
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
}

.editor-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  min-height: 400px;
}

:deep(.w-e-text-container) {
  min-height: 300px !important;
}
</style> 