<template>
  <div class="post-growth">
    <h2>帖子增长分析</h2>
    
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>每日帖子增长</span>
              <span class="subtitle">最近7天</span>
            </div>
          </template>
          <div class="chart-container" ref="dailyChartRef"></div>
        </el-card>
      </el-col>
      
      <el-col :span="24" style="margin-top: 20px;">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>每周帖子增长</span>
              <span class="subtitle">最近7周</span>
            </div>
          </template>
          <div class="chart-container" ref="weeklyChartRef"></div>
        </el-card>
      </el-col>
      
      <el-col :span="24" style="margin-top: 20px;">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>每月帖子增长</span>
              <span class="subtitle">最近6个月</span>
            </div>
          </template>
          <div class="chart-container" ref="monthlyChartRef"></div>
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

const dailyChartRef = ref(null)
const weeklyChartRef = ref(null)
const monthlyChartRef = ref(null)
let dailyChart = null
let weeklyChart = null
let monthlyChart = null

// 初始化图表
const initChart = (el, title) => {
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
      name: '新增帖子数'
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
          color: '#67C23A' // 使用绿色来区分用户增长图表
        },
        itemStyle: {
          borderWidth: 2,
          color: '#67C23A'
        }
      }
    ]
  }
  chart.setOption(option)
  return chart
}

// 更新图表数据
const updateChart = (chart, data, title) => {
  // 反转数组，使早期数据在左侧
  const reversedData = [...data].reverse()
  const dates = reversedData.map(item => item.date)
  const counts = reversedData.map(item => item.count)
  
  chart.setOption({
    xAxis: {
      data: dates
    },
    series: [{
      name: title,
      data: counts
    }]
  })
}

// 获取数据
const fetchData = async () => {
  try {
    const res = await request({
      url: '/api/analysis/posts/growth',
      method: 'get'
    })
    
    if (res.data.code === 1000) {
      const { daily_growth, weekly_growth, monthly_growth } = res.data.data
      
      updateChart(dailyChart, daily_growth, '日增长')
      updateChart(weeklyChart, weekly_growth, '周增长')
      updateChart(monthlyChart, monthly_growth, '月增长')
    } else {
      ElMessage.error(res.data.msg || '获取数据失败')
    }
  } catch (error) {
    console.error('获取帖子增长数据失败:', error)
    ElMessage.error('获取数据失败')
  }
}

// 处理窗口大小变化
const handleResize = () => {
  dailyChart?.resize()
  weeklyChart?.resize()
  monthlyChart?.resize()
}

onMounted(() => {
  dailyChart = initChart(dailyChartRef.value, '日增长')
  weeklyChart = initChart(weeklyChartRef.value, '周增长')
  monthlyChart = initChart(monthlyChartRef.value, '月增长')
  
  fetchData()
  
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  dailyChart?.dispose()
  weeklyChart?.dispose()
  monthlyChart?.dispose()
})
</script>

<style scoped>
.post-growth {
  padding: 20px;
}

.chart-card {
  background: #fff;
  border-radius: 8px;
}

.card-header {
  display: flex;
  align-items: center;
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
</style> 