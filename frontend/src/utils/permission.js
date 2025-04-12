import request from './request'
import store from '../store'

// 检查是否是管理员
// 调用后端的权限校验，看用户是否为admin或super_admin
export async function isAdmin() {
  try {
    // 从 store 中获取用户ID，并确保它是字符串形式
    const userID = String(store.state.userInfo.userID)

    if (!userID) {
      console.error('未找到用户ID')
      return false
    }

    console.log('使用用户ID:', userID)
    const res = await request({
      url: `api/permission/user/role/${userID}`,
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