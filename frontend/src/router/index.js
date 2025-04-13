import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Home from '../views/Home.vue'
import UserProfile from '../views/UserProfile.vue'
import store from '../store'
import { ElMessage } from 'element-plus'

import CreatePost from '../views/CreatePost.vue'
import PostDetail from '../views/PostDetail.vue'
import AdminHome from '../views/AdminHome.vue'

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
        path: 'permissions',
        name: 'Permissions',
        component: () => import('../views/admin/Permissions.vue'),
        meta: { requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'boards',
        name: 'BoardManagement',
        component: () => import('../views/BoardManagement.vue'),
        meta: { requiresAuth: true, requiresAdmin: true }
      },
      //DONE
      {
        path: 'analytics/users/growth', 
        name: 'UserGrowth',
        component: () => import('@/views/analytics/UserGrowth.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: '用户增长'
        }
      },
      //DONE
      {
        path: 'analytics/active/users',
        name: 'ActiveUsers',
        component: () => import('@/views/analytics/ActiveUsers.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: '最近活跃用户'
        }
      },
      //DONE
      {
        path: 'analytics/posts/growth', 
        name: 'PostsGrowth',
        component: () => import('@/views/analytics/PostGrowth.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: '帖子增长'
        }
      },
      {
        path: 'analytics/hot/posts',
        name: 'HotPosts',
        component: () => import('@/views/analytics/HotPosts.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: '最近活跃帖子'
        }
      },
      {
        path: 'analytics/wordcloud',
        name: 'WordCloud',
        component: () => import('@/views/analytics/WordCloud.vue'),
        meta: {
          requiresAuth: true,
          requiresAdmin: true,
          title: '词云图'
        }
      },
    ]
  },
  {
    path: '/user/likes',
    name: 'UserLikes',
    component: () => import('@/views/user/UserLikes.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/user/favorites',
    name: 'UserFavorites',
    component: () => import('@/views/user/UserFavorites.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/user/posts',
    name: 'UserPosts',
    component: () => import('@/views/user/UserPosts.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 改进的路由守卫
router.beforeEach(async (to, from, next) => {
  // 检查是否需要认证
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // 检查是否已登录
    if (!store.state.userInfo) {
      ElMessage.warning('请先登录')
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
      return
    }

    // 检查是否需要管理员权限
    if (to.matched.some(record => record.meta.requiresAdmin)) {
      // 检查用户角色
      if (store.state.userInfo.role !== 'admin' && store.state.userInfo.role !== 'super_admin') {
        ElMessage.error('需要管理员权限')
        next('/')
        return
      }
    }
  }
  next()
})

export default router 