import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from '@/utils/storage.js'

const AdminLayout = () => import('@/components/layout/AdminLayout.vue')

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/pages/login/LoginPage.vue'),
    meta: { public: true, title: '登录' }
  },
  {
    path: '/',
    redirect: '/profile',
    component: AdminLayout,
    children: [
      {
        path: 'profile',
        name: 'profile',
        component: () => import('@/pages/profile/ProfilePage.vue'),
        meta: { title: '个人中心', icon: 'Avatar', order: 5 }
      }
    ]
  },
  {
    path: '/user-manage',
    component: AdminLayout,
    children: [
      {
        path: '',
        name: 'user-manage',
        component: () => import('@/pages/user/UserManagePage.vue'),
        meta: { title: '用户管理', icon: 'User', order: 1, perm: 'sys:user' }
      }
    ]
  },
  {
    path: '/role-manage',
    component: AdminLayout,
    children: [
      {
        path: '',
        name: 'role-manage',
        component: () => import('@/pages/role/RoleManagePage.vue'),
        meta: { title: '角色管理', icon: 'UserFilled', order: 2, perm: 'sys:role' }
      }
    ]
  },
  {
    path: '/menu-manage',
    component: AdminLayout,
    children: [
      {
        path: '',
        name: 'menu-manage',
        component: () => import('@/pages/menu/MenuManagePage.vue'),
        meta: { title: '菜单管理', icon: 'Menu', order: 3, perm: 'sys:menu' }
      }
    ]
  },
  {
    path: '/permission',
    component: AdminLayout,
    children: [
      {
        path: '',
        name: 'permission',
        component: () => import('@/pages/permission/PermissionPage.vue'),
        meta: { title: '权限分配', icon: 'Lock', order: 4, perm: 'sys:role' }
      }
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 登录态 + 权限码守卫：未登录跳 /login；刷新后恢复权限；带 meta.perm 的页面需持有对应权限标识
import { useAuth } from '@/composables/useAuth.js'
const { hasPermission, ensurePermissions } = useAuth()
router.beforeEach(async (to) => {
  if (!to.meta.public && !getToken()) {
    return { name: 'login' }
  }
  if (to.meta.public && getToken()) {
    return { path: '/' }
  }
  if (to.meta.perm) {
    await ensurePermissions()
    if (!hasPermission(to.meta.perm)) {
      return { name: 'profile' }
    }
  }
  document.title = to.meta.title ? `${to.meta.title} · Go Vue Admin` : 'Go Vue Admin'
})

export default router
