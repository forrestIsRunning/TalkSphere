import request from '../utils/request'

// 创建帖子
export function createPost(data) {
  return request({
    url: '/posts',
    method: 'post',
    data
  })
}

// 获取帖子详情
export function getPostDetail(id) {
  return request({
    url: `/posts/${id}`,
    method: 'get'
  })
}

// 获取用户帖子列表
export function getUserPosts(userId, params) {
  return request({
    url: `/posts/user/${userId}`,
    method: 'get',
    params
  })
}

// 获取板块帖子列表
export function getBoardPosts(boardId, params) {
  return request({
    url: `/posts/board/${boardId}`,
    method: 'get',
    params
  })
}