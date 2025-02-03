import request from '../utils/request'

// 创建帖子
export function createPost(data) {
  return request({
    url: 'api/posts',
    method: 'post',
    data
  })
}

// 获取帖子详情
export function getPostDetail(id) {
  return request({
    url: `api/posts/${id}`,
    method: 'get'
  })
}

// 获取用户帖子列表
export function getUserPosts(userId, params) {
  return request({
    url: `api/posts/user/${userId}`,
    method: 'get',
    params
  })
}

// 获取板块下的帖子列表
export const getBoardPosts = (boardId, params) => {
  return request({
    url: `/api/posts/board/${boardId}`,
    method: 'get',
    params // 直接传递params对象
  })
}