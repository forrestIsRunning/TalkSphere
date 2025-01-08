import axios from 'axios'

const baseURL = 'http://127.0.0.1:8989'

const request = axios.create({
  baseURL,
  timeout: 5000
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = token
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

export const login = (data) => {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

export const register = (data) => {
  return request({
    url: '/register', 
    method: 'post',
    data
  })
}

export const getUserProfile = () => {
  return request({
    url: '/profile',
    method: 'get'
  })
}

export const updateBio = (data) => {
  return request({
    url: '/bio',
    method: 'post',
    data
  })
} 