import request from './request'

// 检查用户是否有特定权限
export async function checkPermission(userId) {
  try {
    const res = await request({
      url: `/user/check/${userId}`,
      method: 'get'
    })
    return res.data.code === 1000 && res.data.data.is_admin
  } catch (error) {
    console.error('权限检查失败:', error)
    return false
  }
}

// 检查是否是管理员
export async function isAdmin(userId) {
  try {
    const res = await request({
      url: `/user/check/${userId}`,
      method: 'get'
    })
    console.log('权限检查响应:', res) // 添加调试日志
    return res.data && res.data.code === 1000 && res.data.data.is_admin
  } catch (error) {
    console.error('管理员检查失败:', error)
    return false
  }
} 