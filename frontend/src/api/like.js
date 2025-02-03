import request from '../utils/request'

// 点赞/取消点赞
export function toggleLike(data) {
  return request({
    url: 'api/likes',
    method: 'post',
    data
  })
}

// 获取点赞状态
export function getLikeStatus(targetId, targetType) {
  return request({
    url: `api/likes/status`,
    method: 'get',
    params: {
      target_id: targetId,
      target_type: targetType
    }
  })
} 