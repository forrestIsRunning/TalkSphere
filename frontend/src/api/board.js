import request from '../utils/request'

// 获取所有板块
export function getAllBoards() {
  return request({
    url: '/boards',
    method: 'get'
  }).catch(error => {
    console.error('获取板块列表失败:', error)
    throw error
  })
} 