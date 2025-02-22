<template>
  <div class="hot-analysis">
    <h2>热门帖子分析</h2>
    
    <el-card class="filter-card">
      <el-radio-group v-model="timeRange" @change="fetchData">
        <el-radio-button label="daily">每日</el-radio-button>
        <el-radio-button label="weekly">每周</el-radio-button>
        <el-radio-button label="monthly">每月</el-radio-button>
      </el-radio-group>
    </el-card>

    <el-card v-loading="loading" class="post-list">
      <el-table 
        :data="activePosts" 
        stripe
        height="600"
      >
        <el-table-column label="排名" type="index" width="80" align="center" />
        
        <el-table-column label="帖子信息" min-width="300">
          <template #default="{ row }">
            <div class="post-info">
              <div class="post-title">{{ row.title }}</div>
              <div class="post-content">{{ row.content }}</div>
              <div class="post-author">
                <el-avatar :src="row.author_avatar" :size="24" />
                <span>{{ row.author_name }}</span>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="互动指标" width="300">
          <template #default="{ row }">
            <div class="metrics">
              <el-tooltip content="获得点赞" placement="top">
                <div class="metric-item">
                  <el-icon><Star /></el-icon>
                  {{ row.like_count }}
                </div>
              </el-tooltip>
              
              <el-tooltip content="获得收藏" placement="top">
                <div class="metric-item">
                  <el-icon><Collection /></el-icon>
                  {{ row.favorite_count }}
                </div>
              </el-tooltip>
              
              <el-tooltip content="评论数量" placement="top">
                <div class="metric-item">
                  <el-icon><ChatLineRound /></el-icon>
                  {{ row.comment_count }}
                </div>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="活跃度" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getScoreType(row.activity_score)">
              {{ row.activity_score.toFixed(1) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column 
          prop="created_at" 
          label="发布时间" 
          width="180" 
          align="center"
        />
      </el-table>

      <div class="score-info">
        <h4>活跃度计算规则：</h4>
        <ul>
          <li>发布时间 (40%):
            <ul>
              <li>24小时内：100分</li>
              <li>72小时内：70分</li>
              <li>7天内：40分</li>
              <li>其他：10分</li>
            </ul>
          </li>
          <li>获赞数量 (30%)：每个赞1分</li>
          <li>收藏数量 (20%)：每个收藏1分</li>
          <li>评论数量 (10%)：每个评论1分</li>
        </ul>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Star, Collection, ChatLineRound } from '@element-plus/icons-vue'
import request from '@/utils/request'

const timeRange = ref('daily')
const activePosts = ref([])
const loading = ref(false)

// 获取活跃帖子数据
const fetchData = async () => {
  loading.value = true
  try {
    const res = await request({
      url: '/api/analysis/posts/active',
      method: 'get',
      params: {
        time_range: timeRange.value
      }
    })
    
    if (res.data.code === 1000) {
      activePosts.value = res.data.data.active_posts || []
    } else {
      ElMessage.error(res.data.msg || '获取数据失败')
    }
  } catch (error) {
    console.error('获取活跃帖子数据失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

// 根据活跃度得分返回不同的标签类型
const getScoreType = (score) => {
  if (score >= 80) return 'success'
  if (score >= 60) return 'warning'
  return 'info'
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.hot-analysis {
  padding: 20px;
}

.filter-card {
  margin-bottom: 20px;
}

.post-list {
  min-height: 400px;
}

.post-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.post-title {
  font-weight: 500;
  font-size: 16px;
}

.post-content {
  color: #606266;
  font-size: 14px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #909399;
  font-size: 14px;
}

.metrics {
  display: flex;
  gap: 24px;
}

.metric-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #606266;
}

.score-info {
  margin-top: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.score-info h4 {
  margin: 0 0 10px 0;
  color: #303133;
}

.score-info ul {
  margin: 0;
  padding-left: 20px;
  color: #606266;
}

.score-info ul ul {
  margin: 5px 0;
}
</style> 