import { ref, computed } from 'vue'
import * as authApi from '@/api/auth.js'
import { setAuth, getToken, getStoredUser, clearAuth } from '@/utils/storage.js'

// 模块级共享状态：登录态、当前用户、权限标识集合
const user = ref(getStoredUser())
const permKeys = ref([])
const roleCodes = ref([])
let restoring = null

const isSuper = computed(() => roleCodes.value.includes('super_admin'))

export function useAuth() {
  const isAuthenticated = computed(() => !!getToken())

  async function login(credentials) {
    const data = await authApi.login(credentials)
    setAuth(data.accessToken, data.user)
    user.value = data.user
    await loadPermissions()
    return data.user
  }

  // 拉取 /me + /permissions，填充 user / permKeys / roleCodes
  async function loadPermissions() {
    try {
      const [me, perms] = await Promise.all([
        authApi.getProfile(),
        authApi.getPermissions()
      ])
      user.value = me
      permKeys.value = perms.permKeys || []
      roleCodes.value = perms.roleCodes || []
    } catch {
      // 拉取失败不清空已有权限，避免刷新闪烁
      permKeys.value = []
      roleCodes.value = []
    }
  }

  // 刷新恢复：已登录但权限为空时去重地拉一次
  function ensurePermissions() {
    if (permKeys.value.length) return Promise.resolve()
    if (!getToken()) return Promise.resolve()
    if (!restoring) restoring = loadPermissions().finally(() => { restoring = null })
    return restoring
  }

  function logout() {
    clearAuth()
    user.value = null
    permKeys.value = []
    roleCodes.value = []
  }

  // super_admin 持有全部权限
  function hasPermission(code) {
    if (isSuper.value) return true
    return permKeys.value.includes(code)
  }

  return {
    user,
    permKeys,
    roleCodes,
    isSuper,
    isAuthenticated,
    login,
    logout,
    hasPermission,
    loadPermissions,
    ensurePermissions,
  }
}
