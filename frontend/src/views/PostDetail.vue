<template>
  <div class="post-detail-container">
    <div class="post-content" v-if="post">
      <div class="back-button">
        <el-button @click="$router.back()">返回</el-button>
      </div>
      <h1 class="title">{{ post.title }}</h1>
      <div class="meta">
        <span class="author">作者：{{ authorName }}</span>
        <span class="time">发布时间：{{ formatDate(post.created_at) }}</span>
        <span class="views">阅读：{{ post.view_count || 0 }}</span>
      </div>
      <div class="content">{{ post.content }}</div>
      
      <!-- 图片展示 -->
      <div class="images" v-if="post.images?.length">
        <el-image 
          v-for="img in post.images"
          :key="img.id"
          :src="img.url"
          :preview-src-list="post.images.map(img => img.url)"
        />
      </div>

      <!-- 帖子统计信息 -->
      <div class="post-stats">
        <span class="stat-item" @click="handleLike">
          <span :class="{ 'liked': isLiked }">👍</span>
          点赞 {{ post.like_count || 0 }}
        </span>
        <span class="stat-item">
          <el-icon><ChatLineRound /></el-icon>
          评论 {{ post.comment_count || 0 }}
        </span>
        <span class="stat-item" @click="handleFavorite">
          <span :class="{ 'favorited': isFavorited }">📚</span>
          收藏 {{ post.favorite_count || 0 }}
        </span>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-else class="loading">
      <el-skeleton :rows="10" animated />
    </div>

    <!-- 评论区域 -->
    <div class="comments-section" v-if="post">
      <div class="comments-header">
        <h3>评论 ({{ comments.length }})</h3>
        <div class="comment-input">
          <div class="comment-user">
            <el-avatar 
              :size="40" 
              :src="userInfo?.avatar || defaultAvatar"
            />
          </div>
          <el-input
            v-model="newComment"
            type="textarea"
            :rows="3"
            placeholder="写下你的评论..."
          />
          <el-button 
            type="primary" 
            @click="submitComment"
            :loading="submitting"
          >发表评论</el-button>
        </div>
      </div>

      <!-- 评论列表 -->
      <div class="comments-list" v-if="comments.length > 0">
        <div v-for="comment in comments" 
          :key="comment.id" 
          class="comment-item"
        >
          <div class="comment-user">
            <el-avatar 
              :size="40" 
              :src="comment.user?.avatar_url || defaultAvatar"
            />
            <div class="comment-info">
              <div class="comment-username">{{ comment.user?.username || '未知用户' }}</div>
              <div class="comment-time">{{ formatDate(comment.created_at) }}</div>
            </div>
          </div>
          <div class="comment-content">{{ comment.content }}</div>
          
          <!-- 回复按钮 -->
          <div class="comment-actions">
            <span class="reply-btn" @click="showReplyInput(comment.id)">
              回复
            </span>
            <span class="like-btn" @click="handleCommentLike(comment.id)">
              <el-icon :class="{ 'liked': comment.isLiked }"><Star /></el-icon>
              {{ comment.like_count }}
            </span>
            <span class="reply-count">
              <el-icon><ChatLineRound /></el-icon>
              {{ comment.reply_count }}
            </span>
          </div>

          <!-- 回复输入框 -->
          <div class="reply-input" v-if="replyingTo === comment.id">
            <el-input
              v-model="replyContent"
              type="textarea"
              :rows="2"
              placeholder="回复评论..."
            />
            <div class="reply-actions">
              <el-button size="small" @click="cancelReply">取消</el-button>
              <el-button 
                type="primary" 
                size="small"
                @click="submitReply(comment.id)"
                :loading="submitting"
              >回复</el-button>
            </div>
          </div>

          <!-- 回复列表 -->
          <div class="reply-list" v-if="comment.children?.length">
            <div v-for="reply in comment.children" 
              :key="reply.id" 
              class="reply-item"
            >
              <div class="comment-user">
                <el-avatar 
                  :size="32" 
                  :src="reply.user?.avatar_url || defaultAvatar"
                />
                <div class="comment-info">
                  <div class="comment-username">{{ reply.user?.username || '未知用户' }}</div>
                  <div class="comment-time">{{ formatDate(reply.created_at) }}</div>
                </div>
              </div>
              <div class="comment-content">{{ reply.content }}</div>
              
              <!-- 回复的回复按钮 -->
              <div class="comment-actions">
                <span class="reply-btn" @click="showReplyInput(reply.id)">
                  回复
                </span>
              </div>

              <!-- 回复的回复输入框 -->
              <div class="reply-input" v-if="replyingTo === reply.id">
                <el-input
                  v-model="replyContent"
                  type="textarea"
                  :rows="2"
                  placeholder="回复评论..."
                />
                <div class="reply-actions">
                  <el-button size="small" @click="cancelReply">取消</el-button>
                  <el-button 
                    type="primary" 
                    size="small"
                    @click="submitReply(reply.id)"
                    :loading="submitting"
                  >回复</el-button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 无评论时的提示 -->
      <div v-else class="no-comments">
        暂无评论，快来抢沙发吧！
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useStore } from 'vuex'
import { getPostDetail } from '../api/post'
import { getPostComments, createComment } from '../api/comment'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { Star, ChatLineRound } from '@element-plus/icons-vue'
import { getUserProfile } from '../api/user'
import { toggleLike, getLikeStatus } from '../api/like'
import { toggleFavorite } from '../api/favorite'

export default {
  name: 'PostDetail',
  components: {
    Star,
    ChatLineRound
  },
  setup() {
    const store = useStore()
    const route = useRoute()
    const post = ref(null)
    const comments = ref([])
    const newComment = ref('')
    const replyContent = ref('')
    const replyingTo = ref(null)
    const submitting = ref(false)
    const authorName = ref('加载中...')
    const isLiked = ref(false)
    const isFavorited = ref(false)
    
    const defaultAvatar = computed(() => {
      return 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
    })

    const loadPost = async () => {
      try {
        const res = await getPostDetail(route.params.id)
        console.log('帖子详情响应:', res)
        if (res.data.code === 1000) {
          post.value = res.data.data
          console.log('处理后的帖子详情:', post.value)
          getAuthorInfo(post.value.author_id)
          await checkLikeStatus(post.value.id)
        }
      } catch (error) {
        console.error('获取帖子详情失败:', error)
        ElMessage.error('获取帖子详情失败')
      }
    }

    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }

    // 加载评论列表
    const loadComments = async () => {
      try {
        const res = await getPostComments(route.params.id)
        console.log('评论列表响应:', res)
        if (res.data.code === 1000) {
          comments.value = res.data.data.comments || []
          console.log('处理后的评论数据:', comments.value)
        }
      } catch (error) {
        console.error('获取评论失败:', error)
        ElMessage.error('获取评论失败')
      }
    }

    // 提交评论
    const submitComment = async () => {
      if (!newComment.value.trim()) {
        ElMessage.warning('请输入评论内容')
        return
      }

      submitting.value = true
      try {
        const res = await createComment({
          post_id: parseInt(route.params.id),
          content: newComment.value.trim()
        })
        console.log('评论创建响应:', res)
        if (res.data.code === 1000) {
          ElMessage.success('评论成功')
          newComment.value = ''
          await loadComments()
          await loadPost()
        }
      } catch (error) {
        console.error('提交评论失败:', error)
        ElMessage.error('提交评论失败')
      } finally {
        submitting.value = false
      }
    }

    // 显示回复输入框
    const showReplyInput = (commentId) => {
      replyingTo.value = commentId
      replyContent.value = ''
    }

    // 取消回复
    const cancelReply = () => {
      replyingTo.value = null
      replyContent.value = ''
    }

    // 提交回复
    const submitReply = async (commentId) => {
      if (!replyContent.value.trim()) {
        ElMessage.warning('请输入回复内容')
        return
      }

      submitting.value = true
      try {
        const res = await createComment({
          post_id: parseInt(route.params.id),
          parent_id: commentId,
          content: replyContent.value.trim()
        })
        if (res.data.code === 1000) {
          ElMessage.success('回复成功')
          replyContent.value = ''
          replyingTo.value = null
          await loadComments()
        }
      } catch (error) {
        console.error('提交回复失败:', error)
        ElMessage.error('提交回复失败')
      } finally {
        submitting.value = false
      }
    }

    const getAuthorInfo = async () => {
      try {
        const { data } = await getUserProfile()
        if (data.code === 1000) {
          authorName.value = data.data.username
        }
      } catch (error) {
        console.error('获取作者信息失败:', error)
        authorName.value = '未知用户'
      }
    }

    // 获取点赞状态
    const checkLikeStatus = async (postId) => {
      try {
        const res = await getLikeStatus(postId, 1) // 1 表示帖子
        if (res.data.code === 1000) {
          isLiked.value = res.data.data.status === 'liked'
        }
      } catch (error) {
        console.error('获取点赞状态失败:', error)
      }
    }

    // 处理点赞
    const handleLike = async () => {
      if (!post.value?.id) return

      try {
        const res = await toggleLike({
          target_id: post.value.id,
          target_type: 1 // 1 表示帖子
        })

        if (res.data.code === 1000) {
          isLiked.value = res.data.data.status === 'liked'
          // 更新点赞数
          await loadPost() // 重新加载帖子信息以更新点赞数
          ElMessage.success(isLiked.value ? '点赞成功' : '取消点赞成功')
        }
      } catch (error) {
        console.error('点赞操作失败:', error)
        ElMessage.error('点赞操作失败')
      }
    }

    // 处理评论点赞
    const handleCommentLike = async (commentId) => {
      try {
        const res = await toggleLike({
          target_id: commentId,
          target_type: 2 // 2 表示评论
        })

        if (res.data.code === 1000) {
          // 更新评论列表以刷新点赞状态和数量
          await loadComments()
          ElMessage.success(res.data.data.status === 'liked' ? '点赞成功' : '取消点赞成功')
        }
      } catch (error) {
        console.error('评论点赞失败:', error)
        ElMessage.error('评论点赞失败')
      }
    }

    // 处理收藏
    const handleFavorite = async () => {
      if (!post.value?.id) return

      try {
        const res = await toggleFavorite(post.value.id)
        if (res.data.code === 1000) {
          isFavorited.value = res.data.data.status === 'favorited'
          // 更新帖子信息以刷新收藏数
          await loadPost()
          ElMessage.success(isFavorited.value ? '收藏成功' : '取消收藏成功')
        }
      } catch (error) {
        console.error('收藏操作失败:', error)
        ElMessage.error('收藏操作失败')
      }
    }

    onMounted(async () => {
      await loadPost()
      await loadComments()
    })

    return {
      post,
      comments,
      newComment,
      replyContent,
      replyingTo,
      submitting,
      defaultAvatar,
      formatDate,
      submitComment,
      showReplyInput,
      cancelReply,
      submitReply,
      userInfo: computed(() => store.state.userInfo),
      authorName,
      isLiked,
      handleLike,
      handleCommentLike,
      isFavorited,
      handleFavorite
    }
  }
}
</script>

<style scoped>
.post-detail-container {
  max-width: 800px;
  margin: 20px auto;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.back-button {
  margin-bottom: 20px;
}

.title {
  font-size: 24px;
  font-weight: bold;
  color: #1d2129;
  margin-bottom: 16px;
}

.meta {
  display: flex;
  gap: 24px;
  color: #86909c;
  font-size: 14px;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e4e6eb;
}

.content {
  font-size: 16px;
  line-height: 1.8;
  color: #1d2129;
  margin: 24px 0;
  white-space: pre-wrap;
}

.images {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  margin-top: 24px;
}

.images .el-image {
  width: 100%;
  border-radius: 4px;
}

.loading {
  padding: 20px;
}

.post-stats {
  display: flex;
  gap: 24px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e4e6eb;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #86909c;
  font-size: 14px;
  cursor: pointer;
}

.stat-item .el-icon {
  font-size: 16px;
}

.comments-section {
  margin-top: 40px;
  padding-top: 24px;
  border-top: 1px solid #e4e6eb;
}

.comments-header {
  margin-bottom: 24px;
}

.comments-header h3 {
  margin-bottom: 16px;
  font-size: 18px;
  color: #1d2129;
}

.comment-input {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 24px;
}

.comment-user {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.comment-item {
  padding: 16px 0;
  border-bottom: 1px solid #e4e6eb;
}

.comment-user {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.comment-info {
  display: flex;
  flex-direction: column;
}

.comment-username {
  font-size: 14px;
  font-weight: 500;
  color: #1d2129;
}

.comment-time {
  font-size: 12px;
  color: #86909c;
}

.comment-content {
  margin: 8px 0;
  font-size: 14px;
  line-height: 1.6;
  color: #1d2129;
}

.comment-actions {
  display: flex;
  gap: 16px;
  margin-top: 8px;
}

.like-count,
.reply-count {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #86909c;
  font-size: 13px;
  cursor: pointer;
}

.like-count:hover,
.reply-count:hover {
  color: #1e80ff;
}

.reply-btn {
  font-size: 13px;
  color: #86909c;
  cursor: pointer;
}

.reply-btn:hover {
  color: #1e80ff;
}

.reply-input {
  margin: 12px 0;
  padding-left: 52px;
}

.reply-actions {
  margin-top: 8px;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.reply-list {
  margin-top: 12px;
  margin-left: 48px;
  border-left: 2px solid #e4e6eb;
  padding-left: 16px;
}

.reply-item {
  padding: 12px 0;
  border-bottom: 1px solid #f2f3f5;
}

.reply-item:last-child {
  border-bottom: none;
}

.no-comments {
  text-align: center;
  padding: 40px 0;
  color: #86909c;
  font-size: 14px;
}

.post-meta {
  color: #666;
  font-size: 14px;
}
.post-meta span {
  margin-right: 15px;
}

.stat-item span {
  font-size: 16px;
  line-height: 1;
  cursor: pointer;
  transition: transform 0.2s;
}

.stat-item span:hover {
  transform: scale(1.2);
}

.liked, .favorited {
  filter: brightness(1.2);
}
</style> 