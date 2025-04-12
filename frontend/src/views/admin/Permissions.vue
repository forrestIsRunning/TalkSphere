<template>
  <div class="permissions-manager">
    <h2>权限管理</h2>
    
    <!-- 搜索和过滤 -->
    <el-card class="filter-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="searchQuery"
            placeholder="搜索用户名/ID"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="6">
          <el-select
            v-model="roleFilter"
            placeholder="按角色筛选"
            clearable
            @change="handleSearch"
          >
            <el-option label="普通用户" value="user" />
            <el-option label="管理员" value="admin" />
            <el-option label="超级管理员" value="super_admin" />
          </el-select>
        </el-col>
      </el-row>
    </el-card>

    <!-- 用户列表 -->
    <el-card class="user-list-card">
      <el-table
        :data="filteredUsers"
        stripe
        v-loading="loading"
      >
        <!-- 基本信息列 -->
        <el-table-column label="用户信息" min-width="200">
          <template #default="{ row }">
            <div class="user-info">
              <el-avatar :src="row.avatar_url" :size="40" />
              <div class="user-details">
                <div class="username">{{ row.username }}</div>
                <div class="user-id">ID: {{ row.user_id }}</div>
              </div>
            </div>
          </template>
        </el-table-column>

        <!-- 当前角色列 -->
        <el-table-column label="当前角色" width="150">
          <template #default="{ row }">
            <el-tag :type="getRoleTagType(row.role)">
              {{ getRoleDisplayName(row.role) }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 角色管理列 -->
        <el-table-column label="角色管理" width="200">
          <template #default="{ row }">
            <el-select
              v-model="row.role"
              :disabled="row.role === 'super_admin' && row.user_id !== currentUser.userID"
              @change="(value) => handleRoleChange(row, value)"
            >
              <el-option label="普通用户" value="user" />
              <el-option label="管理员" value="admin" />
              <el-option 
                label="超级管理员" 
                value="super_admin"
                :disabled="!isSuperAdmin"
              />
            </el-select>
          </template>
        </el-table-column>

        <!-- 特殊权限列 -->
        <el-table-column label="特殊权限" min-width="300">
          <template #default="{ row }">
            <el-space wrap>
              <el-tooltip 
                v-for="permission in availablePermissions"
                :key="permission.key"
                :content="permission.description"
                placement="top"
              >
                <el-checkbox
                  v-model="row.permissions[permission.key]"
                  :label="permission.name"
                  :disabled="!canManagePermission(row)"
                  @change="(value) => handlePermissionChange(row, permission.key, value)"
                />
              </el-tooltip>
            </el-space>
          </template>
        </el-table-column>

        <!-- 操作列 -->
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              @click="showPermissionDetails(row)"
              :disabled="!canManagePermission(row)"
            >
              详细权限
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 详细权限对话框 -->
    <el-dialog
      v-model="permissionDialog.visible"
      :title="'详细权限设置 - ' + (permissionDialog.user?.username || '')"
      width="60%"
    >
      <el-tabs v-model="permissionDialog.activeTab">
        <el-tab-pane label="API权限" name="api">
          <el-tree
            ref="permissionTree"
            :data="apiPermissions"
            show-checkbox
            node-key="path"
            :default-checked-keys="permissionDialog.checkedPermissions"
            :props="{ label: 'name', children: 'children' }"
          />
        </el-tab-pane>
        <el-tab-pane label="功能权限" name="feature">
          <el-checkbox-group v-model="permissionDialog.checkedFeatures">
            <el-row :gutter="20">
              <el-col :span="8" v-for="feature in featurePermissions" :key="feature.key">
                <el-checkbox :label="feature.key">
                  {{ feature.name }}
                  <el-tooltip :content="feature.description" placement="top">
                    <el-icon class="info-icon"><InfoFilled /></el-icon>
                  </el-tooltip>
                </el-checkbox>
              </el-col>
            </el-row>
          </el-checkbox-group>
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="permissionDialog.visible = false">取消</el-button>
          <el-button type="primary" @click="saveDetailedPermissions">
            保存
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { ElMessage } from 'element-plus'
import { Search, InfoFilled } from '@element-plus/icons-vue'
import request from '@/utils/request'

// 状态管理
const store = useStore()
const loading = ref(false)
const searchQuery = ref('')
const roleFilter = ref('')
const users = ref([])

// 当前用户信息
const currentUser = computed(() => store.state.userInfo)
const isSuperAdmin = computed(() => currentUser.value?.role === 'super_admin')

// 可用的特殊权限列表
const availablePermissions = [
  { key: 'manage_boards', name: '板块管理', description: '允许创建、编辑和删除板块' },
  { key: 'manage_posts', name: '帖子管理', description: '允许管理所有用户的帖子' },
  { key: 'view_analytics', name: '查看分析', description: '允许查看数据分析页面' },
  { key: 'manage_users', name: '用户管理', description: '允许管理普通用户' },
]

// API权限树形数据
const apiPermissions = [
  {
    name: '用户管理',
    path: '/api/users',
    children: [
      { name: '查看用户列表', path: '/api/users/list', actions: ['read'] },
      { name: '编辑用户信息', path: '/api/users/edit', actions: ['write'] },
      { name: '删除用户', path: '/api/users/delete', actions: ['delete'] }
    ]
  },
  {
    name: '内容管理',
    path: '/api/content',
    children: [
      { 
        name: '帖子管理', 
        path: '/api/posts',
        children: [
          { name: '查看帖子', path: '/api/posts/view', actions: ['read'] },
          { name: '创建帖子', path: '/api/posts/create', actions: ['write'] },
          { name: '编辑帖子', path: '/api/posts/edit', actions: ['write'] },
          { name: '删除帖子', path: '/api/posts/delete', actions: ['delete'] },
          { name: '置顶帖子', path: '/api/posts/pin', actions: ['write'] }
        ]
      },
      { 
        name: '评论管理', 
        path: '/api/comments',
        children: [
          { name: '查看评论', path: '/api/comments/view', actions: ['read'] },
          { name: '发表评论', path: '/api/comments/create', actions: ['write'] },
          { name: '删除评论', path: '/api/comments/delete', actions: ['delete'] }
        ]
      }
    ]
  },
  {
    name: '系统管理',
    path: '/api/system',
    children: [
      { 
        name: '系统配置', 
        path: '/api/system/config',
        children: [
          { name: '查看配置', path: '/api/system/config/view', actions: ['read'] },
          { name: '修改配置', path: '/api/system/config/edit', actions: ['write'] }
        ]
      },
      { 
        name: '数据分析', 
        path: '/api/analysis',
        children: [
          { name: '查看统计', path: '/api/analysis/view', actions: ['read'] },
          { name: '导出数据', path: '/api/analysis/export', actions: ['write'] }
        ]
      }
    ]
  }
]

// 功能权限列表
const featurePermissions = [
  { key: 'create_board', name: '创建板块', description: '允许创建新的讨论板块' },
  { key: 'delete_board', name: '删除板块', description: '允许删除现有板块' },
  { key: 'pin_post', name: '置顶帖子', description: '允许将帖子置顶' },
  { key: 'ban_user', name: '封禁用户', description: '允许临时封禁违规用户' },
]

// 权限详情对话框
const permissionDialog = ref({
  visible: false,
  activeTab: 'api',
  user: null,
  checkedPermissions: [],
  checkedFeatures: []
})

// 添加 permissionTree ref
const permissionTree = ref(null)

// 过滤后的用户列表
const filteredUsers = computed(() => {
  return users.value.filter(user => {
    const matchQuery = searchQuery.value ? 
      (user.username.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
       user.user_id.toString().includes(searchQuery.value)) : true
    
    const matchRole = roleFilter.value ? 
      user.role === roleFilter.value : true
    
    return matchQuery && matchRole
  })
})

// 获取用户列表
const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await request({
      url: '/api/users',
      method: 'get'
    })
    
    if (res.data.code === 1000) {
      // 直接使用后端返回的用户数据，包括角色信息
      users.value = res.data.data.users.map(user => ({
        user_id: user.id,
        id: user.id,
        username: user.username,
        avatar_url: user.avatar,
        bio: user.bio,
        email: user.email,
        created_at: user.created_at,
        role: user.role || 'user',  // 使用后端返回的角色，如果没有则默认为 user
        permissions: {},
        originalRole: user.role || 'user'  // 同样使用后端返回的角色
      }))
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

// 处理角色变更
const handleRoleChange = async (user, newRole) => {
  try {
    const res = await request({
      url: `/api/permission/user/role/${user.user_id}`,
      method: 'post',
      data: {
        role: newRole
      }
    })
    
    if (res.data.code === 1000) {
      ElMessage.success('角色更新成功')
      user.role = newRole
      user.originalRole = newRole  // 更新原始角色
    } else {
      ElMessage.error(res.data.msg || '角色更新失败')
      // 回滚 UI 状态
      user.role = user.originalRole
    }
  } catch (error) {
    console.error('更新角色失败:', error)
    ElMessage.error('角色更新失败')
    // 回滚 UI 状态
    user.role = user.originalRole
  }
}

// 处理权限变更
const handlePermissionChange = async (user, permission, value) => {
  try {
    const res = await request({
      url: `/api/permission/user/${user.user_id}`,
      method: 'post',
      data: {
        user_id: user.user_id,
        permissions: [{
          path: permission,
          actions: [value ? 'write' : 'read']
        }]
      }
    })
    
    if (res.data.code === 1000) {
      ElMessage.success('权限更新成功')
    } else {
      ElMessage.error(res.data.msg || '权限更新失败')
      // 回滚UI状态
      user.permissions[permission] = !value
    }
  } catch (error) {
    console.error('更新权限失败:', error)
    ElMessage.error('权限更新失败')
    // 回滚UI状态
    user.permissions[permission] = !value
  }
}

// 显示详细权限设置
const showPermissionDetails = (user) => {
  permissionDialog.value = {
    visible: true,
    activeTab: 'api',
    user,
    checkedPermissions: user.detailed_permissions?.api || [],
    checkedFeatures: user.detailed_permissions?.features || []
  }
}

// 保存详细权限设置
const saveDetailedPermissions = async () => {
  const user = permissionDialog.value.user
  if (!user || !permissionTree.value) return  // 添加空值检查

  try {
    // 获取所有选中的节点
    const checkedNodes = permissionTree.value.getCheckedNodes()
    
    // 转换为权限策略格式
    const permissions = checkedNodes
      .filter(node => node.actions) // 只处理有 actions 的叶子节点
      .flatMap(node => 
        node.actions.map(action => [
          user.user_id,
          node.path,
          action
        ])
      )

    const res = await request({
      url: '/api/user/detailed-permissions',
      method: 'put',
      data: {
        user_id: user.user_id,
        permissions: permissions
      }
    })
    
    if (res.data.code === 1000) {
      ElMessage.success('权限更新成功')
      permissionDialog.value.visible = false
      // 更新用户的详细权限
      user.detailed_permissions = permissions
    } else {
      ElMessage.error(res.data.msg || '权限更新失败')
    }
  } catch (error) {
    console.error('更新详细权限失败:', error)
    ElMessage.error('权限更新失败')
  }
}

// 工具函数
const getRoleTagType = (role) => {
  switch (role) {
    case 'super_admin': return 'danger'
    case 'admin': return 'warning'
    default: return 'info'
  }
}

const getRoleDisplayName = (role) => {
  switch (role) {
    case 'super_admin': return '超级管理员'
    case 'admin': return '管理员'
    case 'user': return '普通用户'
    default: return '未知角色'
  }
}

const canManagePermission = (user) => {
  // 超级管理员可以管理所有人的权限，除了其他超级管理员
  if (isSuperAdmin.value) {
    return user.role !== 'super_admin' || user.user_id === currentUser.value?.userID
  }
  // 管理员只能管理普通用户的权限
  return currentUser.value?.role === 'admin' && user.role === 'user'
}

// 搜索处理
const handleSearch = () => {
  // 已通过计算属性 filteredUsers 处理
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.permissions-manager {
  padding: 20px;
}

.filter-card {
  margin-bottom: 20px;
}

.user-list-card {
  margin-top: 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.username {
  font-weight: 500;
  font-size: 16px;
}

.user-id {
  color: #909399;
  font-size: 12px;
}

.info-icon {
  margin-left: 4px;
  color: #909399;
  cursor: help;
}

:deep(.el-checkbox) {
  margin-right: 20px;
  margin-bottom: 10px;
}

:deep(.el-tree) {
  margin: 10px 0;
}

.dialog-footer {
  margin-top: 20px;
}
</style> 