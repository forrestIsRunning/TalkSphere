import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Home from '../views/Home.vue'
import UserProfile from '../views/UserProfile.vue'
import store from '../store'
import { ElMessage } from 'element-plus'
import { isAdmin } from '@/utils/permission'

import CreatePost from '../views/CreatePost.vue'
import PostDetail from '../views/PostDetail.vue'
import AdminHome from '../views/AdminHome.vue'
import {getUserProfile} from "@/api/user";

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: Login,
    props: { isAdminLogin: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  },
  {
    path: '/profile',
    name: 'UserProfile',
    component: UserProfile,
    meta: { requiresAuth: true }
  },
  {
    path: '/post/create',
    name: 'CreatePost',
    component: CreatePost,
    meta: { requiresAuth: true }
  },
  {
    path: '/post/:id',
    name: 'PostDetail',
    component: PostDetail,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin',
    component: AdminHome,
    meta: { 
      requiresAuth: true,
      requiresAdmin: true 
    },
    children: [
      {
        path: '',  // 默认子路由
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'users',
        name: 'UserManagement',
        component: () => import('../views/UserManagement.vue'),
        meta: { requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'boards',
        name: 'BoardManagement',
        component: () => import('../views/BoardManagement.vue'),
        meta: { requiresAuth: true, requiresAdmin: true }
      },
      // ... 其他管理页面路由
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 改进的路由守卫
router.beforeEach(async (to, from, next) => {
  const token = localStorage.getItem('token')
  const userInfo = store.state.userInfo

  if (to.meta.requiresAuth) {
    if (!token) {
      next({ path: '/login', query: { redirect: to.fullPath } })
      return
    }
    
    // 如果有token但没有用户信息，尝试获取用户信息
    if (token && (!userInfo || !userInfo.userID)) {
      try {
        const res = await getUserProfile()
        if (res.data.code === 1000) {
          store.commit('SET_USERINFO', res.data.data)
          // 检查管理员权限
          if (to.meta.requiresAdmin && !(await isAdmin(res.data.data.userID))) {
            ElMessage.error('需要管理员权限')
            next({ path: '/' })
            return
          }
          next()
          return
        } else {
          // 获取用户信息失败，清除token并跳转到登录页
          store.dispatch('logout')
          next({ path: '/login', query: { redirect: to.fullPath } })
          return
        }
      } catch (error) {
        store.dispatch('logout')
        next({ path: '/login', query: { redirect: to.fullPath } })
        return
      }
    }
  }
  
  // 已有用户信息，检查管理员权限
  if (to.meta.requiresAdmin && !(await isAdmin(userInfo.userID))) {
    ElMessage.error('需要管理员权限')
    next({ path: '/' })
    return
  }
  
  next()
})

export default router 