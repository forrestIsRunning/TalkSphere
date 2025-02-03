import request from '../utils/request'

// 收藏/取消收藏帖子
export const toggleFavorite = (postId) => {
  return request({
    url: `/api/favorites/post/${postId}`,
    method: 'post'
  })
} 