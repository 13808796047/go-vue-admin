import { computed } from 'vue'
import { useAuth } from '@/composables/useAuth.js'
import {
  User,
  UserFilled,
  Menu,
  Lock,
  Avatar
} from '@element-plus/icons-vue'

// 图标映射：菜单配置用字符串名，侧边栏渲染需要组件对象
const ICONS = {
  User,
  UserFilled,
  Menu,
  Lock,
  Avatar
}

// 侧边栏菜单项：perm 为后端 menus.permKey（type=2 菜单权限标识）
// super_admin 经 hasPermission 自动全量放行；null 表示不限制（人人可见）
const MENU_ITEMS = [
  { path: '/user-manage', title: '用户管理', icon: 'User', perm: 'sys:user' },
  { path: '/role-manage', title: '角色管理', icon: 'UserFilled', perm: 'sys:role' },
  { path: '/menu-manage', title: '菜单管理', icon: 'Menu', perm: 'sys:menu' },
  { path: '/permission', title: '权限分配', icon: 'Lock', perm: 'sys:role' },
  { path: '/profile', title: '个人中心', icon: 'Avatar', perm: null },
]

export function usePermission() {
  const { user, permKeys, hasPermission } = useAuth()

  const visibleMenus = computed(() =>
    MENU_ITEMS.filter((m) => !m.perm || hasPermission(m.perm)).map((m) => ({
      ...m,
      icon: ICONS[m.icon]
    }))
  )

  return { user, permKeys, hasPermission, visibleMenus }
}
