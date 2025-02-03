<template>
  <div class="user-management">
    <div class="search-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索用户名/ID/邮箱/简介"
        clearable
        @clear="handleSearch"
        @keyup.enter="handleSearch"
        style="width: 300px"
      >
        <template #append>
          <el-button @click="handleSearch">
            <el-icon><Search /></el-icon>
          </el-button>
        </template>
      </el-input>
    </div>

    <h2>用户管理</h2>
    <el-table :data="users" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="180" />
      <el-table-column prop="username" label="用户名" width="180" />
      <el-table-column prop="email" label="邮箱" width="180" />
      <el-table-column label="头像" width="100">
        <template #default="scope">
          <el-avatar :src="scope.row.avatar" />
        </template>
      </el-table-column>
      <el-table-column prop="bio" label="个人简介" />
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="scope">
          {{ formatDate(scope.row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
        </template>
      </el-table-column>
    </el-table>

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
</template>

<script>
import request from '@/utils/request'
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'

export default {
  name: 'UserManagement',
  components: { Search },
  setup() {
    const searchKeyword = ref('')
    const loading = ref(false)
    const users = ref([])
    const total = ref(0)
    const currentPage = ref(1)
    const pageSize = ref(10)

    const fetchUsers = async () => {
      try {
        loading.value = true
        const res = await request({
          url: '/api/users',
          method: 'get',
          params: {
            page: currentPage.value,
            size: pageSize.value,
            keyword: searchKeyword.value
          }
        })
        
        if (res.data.code === 1000) {
          users.value = res.data.data.users
          total.value = res.data.data.total
        } else {
          ElMessage.error(res.data.msg || '获取用户列表失败')
        }
      } catch (error) {
        console.error('获取用户列表失败:', error)
        ElMessage.error('获取用户列表失败')
      } finally {
        loading.value = false
      }
    }

    const handleSizeChange = (val) => {
      pageSize.value = val
      fetchUsers()
    }

    const handleCurrentChange = (val) => {
      currentPage.value = val
      fetchUsers()
    }

    const handleEdit = (row) => {
      console.log('编辑用户:', row)
    }

    const formatDate = (date) => {
      return new Date(date).toLocaleString()
    }

    const handleSearch = () => {
      currentPage.value = 1
      fetchUsers()
    }

    onMounted(() => {
      fetchUsers()
    })

    return {
      searchKeyword,
      loading,
      users,
      total,
      currentPage,
      pageSize,
      handleEdit,
      handleSizeChange,
      handleCurrentChange,
      formatDate,
      handleSearch,
      fetchUsers
    }
  }
}
</script>

<style scoped>
.user-management {
  padding: 20px;
}

.search-bar {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>