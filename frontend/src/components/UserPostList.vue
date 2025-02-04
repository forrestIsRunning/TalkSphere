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
              <p>{{ truncateContent(post.content || '') }}</p>
              <div v-if="post.images?.length" class="post-images">
                <el-image 
                  v-for="img in post.images.slice(0, 3)" 
                  :key="img.ID"
                  :src="img.ImageURL"
                  :preview-src-list="[img.ImageURL]"
                  fit="cover"
                />
              </div>
            </div>
            <div class="post-footer">
              <div class="stats">
                <span><i class="el-icon-view"></i> {{ post.view_count || 0 }}</span>
                <span><i class="el-icon-star-on"></i> {{ post.like_count || 0 }}</span>
                <span><i class="el-icon-collection-tag"></i> {{ post.favorite_count || 0 }}</span>
                <span><i class="el-icon-chat-dot-round"></i> {{ post.comment_count || 0 }}</span>
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
import { ElMessage } from 'element-plus'
import { formatTime } from '@/utils/time'

const defaultAvatar = '/defaultAvatar.jpg'

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
  return content.length > 200 ? content.slice(0, 200) + '...' : content
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

onMounted(() => {
  loadPosts()
})
</script>

<script>
export default {
  name: 'UserPostList'
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
  gap: 10px;
  margin-top: 10px;
}

.post-images .el-image {
  width: 150px;
  height: 150px;
  border-radius: 4px;
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
</style> 