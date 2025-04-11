import request from './request'

// 检查是否是管理员
// 调用后端的权限校验，看用户是否为admin或super_admin
export async function isAdmin() {
  try {
    const res = await request({
      url: '/api/user/role',
      method: 'get'
    })
    console.log('权限检查响应:', res)
    return res.data && res.data.code === 1000 && 
           (res.data.data.role === 'admin' || res.data.data.role === 'super_admin')
  } catch (error) {
    console.error('管理员检查失败:', error)
    return false
  }
} 