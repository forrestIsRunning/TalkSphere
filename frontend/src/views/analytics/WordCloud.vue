<template>
  <div class="word-cloud">
    <h2>词云分析</h2>
    <div class="word-cloud-container">
      <div v-if="loading" class="loading">
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const wordCloudUrl = ref('')

const fetchWordCloud = async () => {
  loading.value = true
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
    loading.value = false
  }
}

onMounted(() => {
  fetchWordCloud()
})
</script>

<style scoped>
.word-cloud {
  padding: 20px;
}

.word-cloud-container {
  margin-top: 20px;
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  min-height: 400px;
}

.word-cloud-image {
  width: 100%;
  height: 100%;
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

:deep(.el-image) {
  max-width: 100%;
  max-height: 600px;
}
</style> 