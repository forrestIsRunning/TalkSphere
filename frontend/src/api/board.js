import request from '../utils/request'

// 获取所有板块
export function getAllBoards() {
  return request({
    url: 'api/boards',
    method: 'get'
  })
}

// 创建板块
export function createBoard(data) {
  return request({
    url: 'api/boards',
    method: 'post',
    data
  })
}

// 更新板块
export function updateBoard(id, data) {
  return request({
    url: `/boards/${id}`,
    method: 'put',
    data
  })
}

// 删除板块
export function deleteBoard(id) {
  return request({
    url: `/boards/${id}`,
    method: 'delete'
  })
} 