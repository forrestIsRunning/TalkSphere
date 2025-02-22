<template>
    <div class="hot-analysis">
      <h2>热门用户分析</h2>
      
      <el-card class="filter-card">
        <el-radio-group v-model="timeRange" @change="fetchData">
          <el-radio-button label="daily">每日</el-radio-button>
          <el-radio-button label="weekly">每周</el-radio-button>
          <el-radio-button label="monthly">每月</el-radio-button>
        </el-radio-group>
      </el-card>
  
      <el-card class="user-list">
        <el-table :data="activeUsers" stripe>
          <el-table-column label="排名" type="index" width="80" align="center" />
          
          <el-table-column label="用户信息" min-width="200">
            <template #default="{ row }">
              <div class="user-info">
                <el-avatar :src="row.avatar_url" :size="40" />
                <span class="username">{{ row.username }}</span>
              </div>
            </template>
          </el-table-column>
  
          <el-table-column label="活跃度指标" min-width="400">
            <template #default="{ row }">
              <div class="metrics">
                <el-tooltip content="发帖数量" placement="top">
                  <div class="metric-item">
                    <el-icon><Document /></el-icon>
                    {{ row.post_count }}
                  </div>
                </el-tooltip>
                
                <el-tooltip content="获得点赞" placement="top">
                  <div class="metric-item">
                    <el-icon><Star /></el-icon>
                    {{ row.like_received_count }}
                  </div>
                </el-tooltip>
                
                <el-tooltip content="获得收藏" placement="top">
                  <div class="metric-item">
                    <el-icon><Collection /></el-icon>
                    {{ row.favorite_received_count }}
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
            prop="last_login_at" 
            label="最后登录" 
            width="180" 
            align="center"
            :formatter="formatDate"
          />
        </el-table>
  
        <div class="score-info">
          <h4>活跃度计算规则：</h4>
          <ul>
            <li>登录时间 (30%):
              <ul>
                <li>24小时内：100分</li>
                <li>72小时内：70分</li>
                <li>7天内：40分</li>
                <li>其他：10分</li>
              </ul>
            </li>
            <li>发帖数量 (30%)：每篇帖子10分</li>
            <li>获赞数量 (20%)：每个赞1分</li>
            <li>收藏数量 (20%)：每个收藏1分</li>
          </ul>
        </div>
      </el-card>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { Document, Star, Collection } from '@element-plus/icons-vue'
  import request from '@/utils/request'
  
  const timeRange = ref('daily')
  const activeUsers = ref([])
  
  // 获取活跃用户数据
  const fetchData = async () => {
    try {
      const res = await request({
        url: '/api/analysis/users/active',
        method: 'get',
        params: {
          time_range: timeRange.value
        }
      })
      
      if (res.data.code === 1000) {
        activeUsers.value = res.data.data.active_users
      } else {
        ElMessage.error(res.data.msg || '获取数据失败')
      }
    } catch (error) {
      console.error('获取活跃用户数据失败:', error)
      ElMessage.error('获取数据失败')
    }
  }
  
  // 格式化日期
  const formatDate = (row) => {
    if (!row.last_login_at) return '暂无记录'
    return row.last_login_at
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
  
  .user-info {
    display: flex;
    align-items: center;
    gap: 12px;
  }
  
  .username {
    font-weight: 500;
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