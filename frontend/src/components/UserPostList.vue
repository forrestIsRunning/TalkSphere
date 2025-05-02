<template>
  <div class="user-post-list">
    <div v-if="loading" class="loading">
      <el-skeleton :rows="3" animated />
    </div>
    <div v-else>
      <div v-if="!posts || posts.length === 0" class="empty">
        <el-empty description="暂无数据" />
      </div>
      <div v-else class="post-list">
        <div v-for="post in posts" :key="post.id" class="post-item">
          <el-card>
            <div class="post-header">
              <el-avatar 
                :src="post.author?.avatar_url || defaultAvatar" 
                :size="40" 
              />
              <div class="post-info">
                <div class="author">{{ post.author?.username || '未知用户' }}</div>
                <div class="time">{{ formatTime(post.created_at) }}</div>
              </div>
            </div>
            <div class="post-content" @click="goToPost(post.id)">
              <h3>{{ post.title }}</h3>
              <div class="rich-content" v-html="truncateContent(post.content || '')"></div>
              <div class="post-images" v-if="post.images?.length">
                <div v-for="image in post.images" :key="image.id" class="image-container">
                  <el-image 
                    :src="image.ImageURL"
                    :preview-src-list="post.images.map(img => img.ImageURL)"
                    fit="cover"
                    class="post-image"
                  >
                    <template #error>
                      <div class="image-error">
                        <el-icon><picture-filled /></el-icon>
                        <span>加载失败</span>
                      </div>
                    </template>
                    <template #placeholder>
                      <div class="image-placeholder">
                        <el-icon class="is-loading"><loading /></el-icon>
                        <span>加载中</span>
                      </div>
                    </template>
                  </el-image>
                </div>
              </div>
            </div>
            <div class="post-footer">
              <div class="stats">
                <span><i class="el-icon-view"></i> {{ post.view_count || 0 }}</span>
                <span><i class="el-icon-star-on"></i> {{ post.like_count || 0 }}</span>
                <span><i class="el-icon-collection-tag"></i> {{ post.favorite_count || 0 }}</span>
                <span><i class="el-icon-chat-dot-round"></i> {{ post.comment_count || 0 }}</span>
              </div>
              <div class="actions">
                <el-button 
                  type="primary" 
                  size="small" 
                  @click.stop="editPost(post)"
                >
                  编辑
                </el-button>
                <el-button 
                  type="danger" 
                  size="small" 
                  @click.stop="confirmDelete(post)"
                >
                  删除
                </el-button>
              </div>
              <div class="tags" v-if="post.tags?.length">
                <el-tag v-for="tag in post.tags" :key="tag.ID" size="small">
                  {{ tag.Name }}
                </el-tag>
              </div>
            </div>
          </el-card>
        </div>
      </div>
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 30, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, defineProps } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatTime } from '@/utils/time'
import { PictureFilled, Loading } from '@element-plus/icons-vue'
import { deletePost } from '@/api/post'

const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

const props = defineProps({
  fetchPosts: {
    type: Function,
    required: true
  }
})

const router = useRouter()
const loading = ref(false)
const posts = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const loadPosts = async () => {
  loading.value = true
  try {
    const res = await props.fetchPosts({
      page: currentPage.value,
      size: pageSize.value
    })
    console.log('API Response:', res)
    if (res.data?.code === 1000 && res.data?.data) {
      posts.value = res.data.data.posts || []
      total.value = res.data.data.total || 0
      console.log('Processed posts:', posts.value)
    } else {
      posts.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('加载帖子失败:', error)
    posts.value = []
    total.value = 0
    ElMessage.error('加载帖子失败')
  } finally {
    loading.value = false
  }
}

const truncateContent = (content) => {
  if (!content) return ''
  // 创建一个临时的 div 来解析 HTML 内容
  const tempDiv = document.createElement('div')
  tempDiv.innerHTML = content
  const textContent = tempDiv.textContent || tempDiv.innerText
  return textContent.length > 200 ? textContent.slice(0, 200) + '...' : textContent
}

const goToPost = (postId) => {
  if (postId) {
    router.push(`/post/${postId}`)
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1  // 重置页码
  loadPosts()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  loadPosts()
}

// 确认删除
const confirmDelete = (post) => {
  ElMessageBox.confirm(
    '确定要删除这篇帖子吗？此操作不可恢复。',
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      const res = await deletePost(post.id)
      if (res.data.code === 1000) {
        ElMessage.success('删除成功')
        // 重新加载帖子列表
        loadPosts()
      } else {
        ElMessage.error(res.data.msg || '删除失败')
      }
    } catch (error) {
      console.error('删除帖子失败:', error)
      ElMessage.error('删除帖子失败')
    }
  }).catch(() => {
    // 用户取消删除
  })
}

// 编辑帖子
const editPost = (post) => {
  router.push(`/post/edit/${post.id}`)
}

onMounted(() => {
  loadPosts()
})
</script>

<script>
export default {
  name: 'UserPostList',
  components: {
    PictureFilled,
    Loading
  }
}
</script>

<style scoped>
.user-post-list {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.post-item {
  margin-bottom: 20px;
}

.post-header {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.post-info {
  margin-left: 10px;
}

.author {
  font-weight: bold;
}

.time {
  font-size: 12px;
  color: #999;
}

.post-content {
  cursor: pointer;
}

.post-content h3 {
  margin: 0 0 10px 0;
}

.post-images {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin: 10px 0;
}

.image-container {
  width: 200px;
  height: 200px;
  position: relative;
  border-radius: 8px;
  overflow: hidden;
}

.post-image {
  width: 100%;
  height: 100%;
}

.image-error, .image-placeholder {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  background-color: #f5f7fa;
}

.image-error .el-icon, .image-placeholder .el-icon {
  font-size: 24px;
  margin-bottom: 8px;
}

.image-error span, .image-placeholder span {
  font-size: 12px;
  color: #909399;
}

.post-footer {
  margin-top: 15px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stats {
  display: flex;
  gap: 15px;
  color: #666;
}

.stats span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.actions {
  display: flex;
  gap: 8px;
  margin-right: 16px;
}

.tags {
  display: flex;
  gap: 5px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.loading {
  padding: 20px;
}

.empty {
  padding: 40px 0;
}

.rich-content {
  overflow: hidden;
  line-height: 1.6;
  color: #333;
}

.rich-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
  margin: 10px 0;
}

.rich-content :deep(p) {
  margin: 8px 0;
}

.rich-content :deep(a) {
  color: #409eff;
  text-decoration: none;
}

.rich-content :deep(pre) {
  background-color: #f6f8fa;
  padding: 12px;
  border-radius: 4px;
  overflow-x: auto;
}

.rich-content :deep(blockquote) {
  border-left: 4px solid #dcdfe6;
  margin: 10px 0;
  padding-left: 10px;
  color: #666;
}
</style> 