import request from '../utils/request'

// 获取帖子评论列表
export function getPostComments(postId) {
  return request({
    url: `/comments/post/${postId}`,
    method: 'get'
  })
}

// 创建评论
export function createComment(data) {
  return request({
    url: '/comments',
    method: 'post',
    data
  })
} 