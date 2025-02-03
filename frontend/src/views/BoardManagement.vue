<template>
  <div class="board-management">
    <h2>板块管理</h2>
    <div class="action-bar">
      <el-button type="primary" @click="handleAdd">新增板块</el-button>
    </div>
    <el-table :data="boards" style="width: 100%">
      <el-table-column prop="ID" label="ID" width="180" />
      <el-table-column prop="Name" label="名称" width="180" />
      <el-table-column prop="Description" label="描述" />
      <el-table-column prop="CreatedAt" label="创建时间" width="180">
        <template #default="scope">
          {{ formatDate(scope.row.CreatedAt) }}
        </template>
      </el-table-column>
      <el-table-column prop="Status" label="状态" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.Status === 1 ? 'success' : 'danger'">
            {{ scope.row.Status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="scope">
          <el-button-group>
            <el-button
              size="small"
              type="primary"
              @click="handleEdit(scope.row)"
            >编辑</el-button>
            <el-button
              size="small"
              type="danger"
              @click="handleDelete(scope.row)"
            >删除</el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form :model="boardForm" :rules="rules" ref="boardFormRef">
        <el-form-item label="板块名称" prop="Name">
          <el-input v-model="boardForm.Name" />
        </el-form-item>
        <el-form-item label="板块描述" prop="Description">
          <el-input
            type="textarea"
            v-model="boardForm.Description"
            :rows="3"
          />
        </el-form-item>
        <el-form-item label="状态" prop="Status">
          <el-switch
            v-model="boardForm.Status"
            :active-value="1"
            :inactive-value="0"
            active-text="正常"
            inactive-text="禁用"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="loading">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import dayjs from 'dayjs'

export default {
  name: 'BoardManagement',
  setup() {
    const boards = ref([])
    const dialogVisible = ref(false)
    const dialogTitle = ref('')
    const loading = ref(false)
    const boardFormRef = ref(null)
    const isEdit = ref(false)

    const boardForm = ref({
      ID: '',
      Name: '',
      Description: '',
      Status: 1
    })

    const rules = {
      Name: [
        { required: true, message: '请输入板块名称', trigger: 'blur' },
        { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
      ],
      Description: [
        { required: true, message: '请输入板块描述', trigger: 'blur' },
        { max: 200, message: '最多200个字符', trigger: 'blur' }
      ]
    }

    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
    }

    // 获取所有板块
    const getBoards = async () => {
      try {
        const res = await request.get('/api/boards')
        if (res.data.code === 1000) {
          boards.value = res.data.data
        }
      } catch (error) {
        ElMessage.error('获取板块列表失败')
      }
    }

    // 新增板块
    const handleAdd = () => {
      isEdit.value = false
      dialogTitle.value = '新增板块'
      boardForm.value = {
        ID: '',
        Name: '',
        Description: '',
        Status: 1
      }
      dialogVisible.value = true
    }

    // 编辑板块
    const handleEdit = (row) => {
      isEdit.value = true
      dialogTitle.value = '编辑板块'
      boardForm.value = { ...row }
      dialogVisible.value = true
    }

    // 删除板块
    const handleDelete = async (row) => {
      try {
        await ElMessageBox.confirm('确定要删除该板块吗？', '提示', {
          type: 'warning'
        })
        const res = await request.delete(`/api/boards/${row.ID}`)
        if (res.data.code === 1000) {
          ElMessage.success('删除成功')
          getBoards()
        } else {
          ElMessage.error(res.data.msg || '删除失败')
        }
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error('删除失败')
        }
      }
    }

    // 提交表单
    const submitForm = async () => {
      if (!boardFormRef.value) return
      
      try {
        await boardFormRef.value.validate()
        loading.value = true

        const requestData = {
          name: boardForm.value.Name,
          description: boardForm.value.Description
        }
        
        let res
        if (isEdit.value) {
          res = await request.put(`/api/boards/${boardForm.value.ID}`, requestData)
        } else {
          res = await request.post('/api/boards', requestData)
        }
        
        if (res.data.code === 1000) {
          ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
          dialogVisible.value = false
          getBoards()
        } else {
          ElMessage.error(res.data.msg || (isEdit.value ? '更新失败' : '创建失败'))
        }
      } catch (error) {
        ElMessage.error(error.message || '操作失败')
      } finally {
        loading.value = false
      }
    }

    onMounted(() => {
      getBoards()
    })

    return {
      boards,
      dialogVisible,
      dialogTitle,
      loading,
      boardForm,
      boardFormRef,
      rules,
      handleAdd,
      handleEdit,
      handleDelete,
      submitForm,
      formatDate
    }
  }
}
</script>

<style scoped>
.board-management {
  padding: 20px;
}
.action-bar {
  margin-bottom: 20px;
}
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>