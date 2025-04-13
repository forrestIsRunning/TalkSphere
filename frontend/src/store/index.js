import { createStore } from 'vuex'
import { getUserProfile } from '../api/user'

export default createStore({
  state: {
    token: localStorage.getItem('token') || '',
    userInfo: JSON.parse(localStorage.getItem('userInfo') || '{}')
  },
  mutations: {
    SET_TOKEN(state, token) {
      state.token = token
      localStorage.setItem('token', token)
    },
    SET_USERINFO(state, userInfo) {
      state.userInfo = {
        ...userInfo,
        userID: String(userInfo.userID)
      }
      localStorage.setItem('userInfo', JSON.stringify(state.userInfo))
    },
    CLEAR_USER(state) {
      state.token = ''
      state.userInfo = {}
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
    }
  },
  actions: {
    // 添加登出动作
    logout({ commit }) {
      commit('CLEAR_USER')
    },
    // 添加更新用户信息的动作
    async updateUserInfo({ commit }) {
      try {
        const res = await getUserProfile()
        if (res.data.code === 1000) {
          commit('SET_USERINFO', res.data.data)
          return true
        }
        return false
      } catch (error) {
        console.error('更新用户信息失败:', error)
        return false
      }
    }
  },
  getters: {
    userAvatar: state => {
      return state.userInfo?.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
    }
  }
}) 