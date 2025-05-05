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

// 获取板块下的帖子列表
export const getBoardPosts = (boardId, params) => {
  return request({
    url: `/api/posts/board/${boardId}`,
    method: 'get',
    params // 直接传递params对象
  })
}

// 获取用户帖子列表
export function getUserPosts(params) {
  return request({
    url: '/api/posts/user',
    method: 'get',
    params: {
      page: params.page,
      size: params.size  // 使用 size 而不是 page_size
    }
  })
}

// 获取用户点赞的帖子
export function getUserLikedPosts(params) {
  return request({
    url: '/api/posts/user/likes',
    method: 'get',
    params: {
      page: params.page,
      size: params.size
    }
  })
}

// 获取用户收藏的帖子
export function getUserFavoritePosts(params) {
  return request({
    url: '/api/posts/user/favorites',
    method: 'get',
    params: {
      page: params.page,
      size: params.size
    }
  })
}

// 删除帖子
export function deletePost(id) {
  return request({
    url: `/api/posts/${id}`,
    method: 'delete'
  })
}

// 更新帖子
export function updatePost(id, data) {
  return request({
    url: `/api/posts/${id}`,
    method: 'put',
    data
  })
}