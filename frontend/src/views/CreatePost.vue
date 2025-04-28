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
          <div class="editor-container">
            <Toolbar
              style="border-bottom: 1px solid #ccc"
              :editor="editorRef"
              :defaultConfig="toolbarConfig"
              :mode="mode"
            />
            <Editor
              style="height: 400px; overflow-y: hidden;"
              v-model="postForm.content"
              :defaultConfig="editorConfig"
              :mode="mode"
              @onCreated="handleCreated"
            />
          </div>
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
import { ref, onMounted, shallowRef, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getAllBoards } from '../api/board'
import { createPost } from '../api/post'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'

export default {
  name: 'CreatePost',
  components: { Editor, Toolbar },
  setup() {
    const router = useRouter()
    const postFormRef = ref(null)
    const loading = ref(false)
    const boards = ref([])
    const baseURL = 'http://127.0.0.1:8989'

    const postForm = ref({
      title: '',
      content: '',
      board_id: '',
      tags: [],
      image_ids: []
    })

    // 编辑器实例，必须用 shallowRef
    const editorRef = shallowRef()
    
    // 模式
    const mode = ref('default')
    
    // 工具栏配置
    const toolbarConfig = {
      excludeKeys: [
        'insertVideo',
        'uploadVideo',
        'group-video',
        'insertTable'
      ]
    }
    
    // 编辑器配置
    const editorConfig = {
      placeholder: '请输入内容...',
      MENU_CONF: {
        uploadImage: {
          server: `${baseURL}/api/posts/image`,
          fieldName: 'image',
          headers: {
            Authorization: localStorage.getItem('token')
          },
          maxFileSize: 20 * 1024 * 1024, // 20MB
          maxNumberOfFiles: 10,
          allowedFileTypes: ['image/jpeg', 'image/png', 'image/gif'],
          metaWithUrl: true,
          customInsert: (res, insertFn) => {
            // 从服务器响应中获取图片url
            if (res.code === 1000) {
              const { image_url, image_id } = res.data
              // 将图片ID添加到表单数据中
              postForm.value.image_ids.push(image_id)
              // 插入图片到编辑器
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

    // 组件销毁时，也及时销毁编辑器
    onBeforeUnmount(() => {
      const editor = editorRef.value
      if (editor == null) return
      editor.destroy()
    })

    const handleCreated = (editor) => {
      editorRef.value = editor // 记录 editor 实例
    }

    const rules = {
      title: [
        { required: true, message: '请输入标题', trigger: 'blur' },
        { min: 3, message: '标题至少3个字符', trigger: 'blur' }
      ],
      content: [
        { required: true, message: '请输入内容', trigger: 'blur' },
        { 
          validator: (rule, value, callback) => {
            if (value && value.length < 10) {
              callback(new Error('内容至少10个字符'))
            } else {
              callback()
            }
          },
          trigger: 'blur'
        }
      ],
      board_id: [
        { required: true, message: '请选择板块', trigger: 'change' }
      ]
    }

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

    const submitPost = async () => {
      if (!postFormRef.value) return
      
      await postFormRef.value.validate(async (valid) => {
        if (valid) {
          loading.value = true
          try {
            // 获取编辑器 HTML 内容
            const content = postForm.value.content
            
            const postData = {
              ...postForm.value,
              content: content
            }
            
            console.log('提交的数据:', postData)
            const res = await createPost(postData)
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
      submitPost,
      editorRef,
      mode,
      toolbarConfig,
      editorConfig,
      handleCreated
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

.editor-container {
  border: 1px solid #ccc;
  z-index: 100;
  border-radius: 4px;
}

:deep(.w-e-text-container) {
  min-height: 300px !important;
}

:deep(.w-e-toolbar) {
  border-top-left-radius: 4px;
  border-top-right-radius: 4px;
}

:deep(.w-e-text-container) {
  border-bottom-left-radius: 4px;
  border-bottom-right-radius: 4px;
}
</style>