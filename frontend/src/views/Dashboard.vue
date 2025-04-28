<template>
  <div class="dashboard">
    <h2>系统概览</h2>
    
    <!-- 基础数据统计 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stats-card">
          <template #header>
            <div class="stats-header">
              <el-icon><User /></el-icon>
              <span>用户总数</span>
            </div>
          </template>
          <div class="stats-content">
            <span class="stats-number">{{ stats.userCount || 0 }}</span>
            <span v-if="stats.userGrowth" class="stats-trend" :class="{ 'up': stats.userGrowth > 0, 'down': stats.userGrowth < 0 }">
              {{ (stats.userGrowth > 0 ? '+' : '') + stats.userGrowth.toFixed(1) }}
            </span>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stats-card">
          <template #header>
            <div class="stats-header">
              <el-icon><Document /></el-icon>
              <span>帖子总数</span>
            </div>
          </template>
          <div class="stats-content">
            <span class="stats-number">{{ stats.postCount || 0 }}</span>
            <span v-if="stats.postGrowth" class="stats-trend" :class="{ 'up': stats.postGrowth > 0, 'down': stats.postGrowth < 0 }">
              {{ (stats.postGrowth > 0 ? '+' : '') + stats.postGrowth.toFixed(1) }}
            </span>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stats-card">
          <template #header>
            <div class="stats-header">
              <el-icon><Grid /></el-icon>
              <span>板块总数</span>
            </div>
          </template>
          <div class="stats-content">
            <span class="stats-number">{{ stats.boardCount || 0 }}</span>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stats-card">
          <template #header>
            <div class="stats-header">
              <el-icon><ChatLineRound /></el-icon>
              <span>评论总数</span>
            </div>
          </template>
          <div class="stats-content">
            <span class="stats-number">{{ stats.commentCount || 0 }}</span>
            <span v-if="stats.commentGrowth" class="stats-trend" :class="{ 'up': stats.commentGrowth > 0, 'down': stats.commentGrowth < 0 }">
              {{ (stats.commentGrowth > 0 ? '+' : '') + stats.commentGrowth.toFixed(1) }}
            </span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 增长趋势图表 -->
    <el-row :gutter="20" class="charts-row">
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>用户增长趋势</span>
              <el-radio-group v-model="userTimeRange" size="small" @change="fetchUserGrowthData">
                <el-radio-button label="7">最近7天</el-radio-button>
                <el-radio-button label="30">最近30天</el-radio-button>
                <el-radio-button label="180">最近6个月</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container" ref="userGrowthChartRef"></div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>帖子增长趋势</span>
              <el-radio-group v-model="postTimeRange" size="small" @change="fetchPostGrowthData">
                <el-radio-button label="7">最近7天</el-radio-button>
                <el-radio-button label="30">最近30天</el-radio-button>
                <el-radio-button label="180">最近6个月</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container" ref="postGrowthChartRef"></div>
        </el-card>
      </el-col>
      
      <el-col :span="24" style="margin-top: 20px;">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>热门话题</span>
              <span class="subtitle">基于帖子内容分析</span>
            </div>
          </template>
          <div class="word-cloud-container">
            <div v-if="wordCloudLoading" class="loading">
              <el-skeleton :rows="10" animated />
            </div>
            <div v-else-if="wordCloudUrl" class="word-cloud-image">
              <el-image 
                :src="wordCloudUrl" 
                fit="contain"
                :preview-src-list="[wordCloudUrl]"
              />
            </div>
            <div v-else class="no-data">
              暂无数据
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'
import { User, Document, Grid, ChatLineRound } from '@element-plus/icons-vue'

// 基础统计数据
const stats = ref({
  userCount: 0,
  postCount: 0,
  boardCount: 0,
  commentCount: 0,
  userGrowth: 0,
  postGrowth: 0,
  commentGrowth: 0
})

// 图表引用
const userGrowthChartRef = ref(null)
const postGrowthChartRef = ref(null)
let userGrowthChart = null
let postGrowthChart = null

// 词云图相关
const wordCloudLoading = ref(false)
const wordCloudUrl = ref('')

// 时间范围选择
const userTimeRange = ref('7')
const postTimeRange = ref('7')

// 初始化图表
const initChart = (el, title, color) => {
  const chart = echarts.init(el)
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'line'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: []
    },
    yAxis: {
      type: 'value',
      name: title === '用户增长' ? '新增用户数' : '新增帖子数'
    },
    series: [
      {
        name: title,
        type: 'line',
        smooth: true,
        data: [],
        areaStyle: {
          opacity: 0.3
        },
        lineStyle: {
          width: 2,
          color: color
        },
        itemStyle: {
          borderWidth: 2,
          color: color
        }
      }
    ]
  }
  chart.setOption(option)
  return chart
}

// 更新图表数据
const updateChart = (chart, data, title) => {
  if (!chart || !data) {
    console.warn('图表或数据为空')
    return
  }

  // 确保数据是数组且不为空
  if (!Array.isArray(data) || data.length === 0) {
    console.warn('数据格式不正确或为空')
    return
  }

  // 按日期升序排序
  const sortedData = [...data].sort((a, b) => {
    // 处理不同格式的日期
    const dateA = a.date.includes('W') ? 
      // 周数据格式：2025-W17
      a.date.replace('W', '-') :
      // 月数据格式：2025-04 或 日数据格式：2025-04-25
      a.date
    
    const dateB = b.date.includes('W') ?
      b.date.replace('W', '-') :
      b.date
    
    return new Date(dateA) - new Date(dateB)
  })

  // 提取日期和数量
  const dates = sortedData.map(item => item.date)
  const counts = sortedData.map(item => item.count)
  
  const option = {
    xAxis: {
      data: dates,
      axisLabel: {
        formatter: function(value) {
          // 根据不同的时间范围格式化日期
          if (value.includes('W')) {
            // 周数据格式：2025-W17
            return value
          } else if (value.includes('-') && value.length === 7) {
            // 月数据格式：2025-04
            return value
          } else {
            // 日数据格式：2025-04-25
            return value.split('-').slice(1).join('-') // 只显示月-日
          }
        }
      }
    },
    yAxis: {
      type: 'value',
      name: title === '用户增长' ? '新增用户数' : '新增帖子数'
    },
    series: [{
      name: title,
      type: 'line',
      smooth: true,
      data: counts,
      areaStyle: {
        opacity: 0.3
      },
      lineStyle: {
        width: 2,
        color: title === '用户增长' ? '#409EFF' : '#67C23A'
      },
      itemStyle: {
        borderWidth: 2,
        color: title === '用户增长' ? '#409EFF' : '#67C23A'
      }
    }]
  }
  
  chart.setOption(option)
}

// 获取用户增长数据
const fetchUserGrowthData = async () => {
  try {
    const res = await request({
      url: '/api/analysis/users/growth',
      method: 'get'
    })
    
    if (res.data.code === 1000) {
      let growthData = []
      switch (userTimeRange.value) {
        case '7':
          growthData = res.data.data.daily_growth
          break
        case '30':
          growthData = res.data.data.weekly_growth
          break
        case '180':
          growthData = res.data.data.monthly_growth
          break
      }
      
      if (growthData && growthData.length > 0) {
        // 确保数据按日期升序排列
        const sortedData = [...growthData].sort((a, b) => 
          new Date(a.date) - new Date(b.date)
        )
        updateChart(userGrowthChart, sortedData, '用户增长')
      } else {
        console.warn('用户增长数据为空')
      }
    } else {
      ElMessage.error(res.data.msg || '获取用户增长数据失败')
    }
  } catch (error) {
    console.error('获取用户增长数据失败:', error)
    ElMessage.error('获取用户增长数据失败')
  }
}

// 获取帖子增长数据
const fetchPostGrowthData = async () => {
  try {
    const res = await request({
      url: '/api/analysis/posts/growth',
      method: 'get'
    })
    
    if (res.data.code === 1000) {
      let growthData = []
      switch (postTimeRange.value) {
        case '7':
          growthData = res.data.data.daily_growth
          break
        case '30':
          growthData = res.data.data.weekly_growth
          break
        case '180':
          growthData = res.data.data.monthly_growth
          break
      }
      
      if (growthData && growthData.length > 0) {
        // 确保数据按日期升序排列
        const sortedData = [...growthData].sort((a, b) => 
          new Date(a.date) - new Date(b.date)
        )
        updateChart(postGrowthChart, sortedData, '帖子增长')
      } else {
        console.warn('帖子增长数据为空')
      }
    } else {
      ElMessage.error(res.data.msg || '获取帖子增长数据失败')
    }
  } catch (error) {
    console.error('获取帖子增长数据失败:', error)
    ElMessage.error('获取帖子增长数据失败')
  }
}

// 获取词云图
const fetchWordCloud = async () => {
  wordCloudLoading.value = true
  try {
    const res = await request({
      url: '/api/analysis/posts/wordcloud',
      method: 'get'
    })
    
    if (res.data.code === 1000) {
      wordCloudUrl.value = res.data.data.word_cloud_url
    } else {
      ElMessage.error(res.data.msg || '获取词云图失败')
    }
  } catch (error) {
    console.error('获取词云图失败:', error)
    ElMessage.error('获取词云图失败')
  } finally {
    wordCloudLoading.value = false
  }
}

// 获取基础统计数据
const fetchStats = async () => {
  try {
    const res = await request({
      url: '/api/admin/stats',
      method: 'get'
    })
    
    if (res.data.code === 1000) {
      stats.value = res.data.data
    } else {
      ElMessage.error(res.data.msg || '获取统计数据失败')
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
    ElMessage.error('获取统计数据失败')
  }
}

// 处理窗口大小变化
const handleResize = () => {
  userGrowthChart?.resize()
  postGrowthChart?.resize()
}

onMounted(() => {
  fetchStats()
  userGrowthChart = initChart(userGrowthChartRef.value, '用户增长', '#409EFF')
  postGrowthChart = initChart(postGrowthChartRef.value, '帖子增长', '#67C23A')
  
  fetchUserGrowthData()
  fetchPostGrowthData()
  fetchWordCloud()
  
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  userGrowthChart?.dispose()
  postGrowthChart?.dispose()
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.stats-row {
  margin-bottom: 20px;
}

.stats-card {
  background: #fff;
  border-radius: 8px;
  height: 100%;
}

.stats-header {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
}

.stats-content {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.stats-number {
  font-size: 24px;
  font-weight: 500;
  color: #303133;
}

.stats-trend {
  font-size: 14px;
}

.stats-trend.up {
  color: #67C23A;
}

.stats-trend.down {
  color: #F56C6C;
}

.chart-card {
  background: #fff;
  border-radius: 8px;
  height: 100%;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.subtitle {
  color: #909399;
  font-size: 14px;
}

.chart-container {
  height: 300px;
  width: 100%;
}

.word-cloud-container {
  min-height: 400px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.loading, .no-data {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 400px;
}

.no-data {
  color: #909399;
  font-size: 14px;
}

.word-cloud-image {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

:deep(.el-image) {
  max-width: 100%;
  max-height: 400px;
}
</style>