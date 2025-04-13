import request from '../utils/request'

export const login = (data) => {
  return request({
    url: 'api/login',
    method: 'post',
    data
  })
}

export const register = (data) => {
  return request({
    url: 'api/register',
    method: 'post',
    data
  })
}

export const getUserProfile = () => {
  return request({
    url: 'api/profile',
    method: 'get'
  })
}

export const updateBio = (data) => {
  return request({
    url: 'api/bio',
    method: 'post',
    data
  })
}