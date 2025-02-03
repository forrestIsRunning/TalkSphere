<template>
  <div class="dashboard">
    <h2>仪表盘</h2>
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card>
          <template #header>用户总数</template>
          <div class="card-content">{{ stats.userCount || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>帖子总数</template>
          <div class="card-content">{{ stats.postCount || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>板块总数</template>
          <div class="card-content">{{ stats.boardCount || 0 }}</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: 'AdminDashboard',
  data() {
    return {
      stats: {
        userCount: 0,
        postCount: 0,
        boardCount: 0
      }
    }
  }
}
</script>

<style scoped>
.dashboard {
  padding: 20px;
}
.card-content {
  font-size: 24px;
  text-align: center;
  padding: 20px;
}
</style> 