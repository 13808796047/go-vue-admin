import http from '@/utils/http.js'

// 登录：返回 { accessToken, refreshToken, expiresIn, user }
export function login(payload) {
  return http.post('/login', payload)
}

// 当前用户（含角色）
export function getProfile() {
  return http.get('/me')
}

// 当前用户权限：{ menus, apis, permKeys, roleCodes }
export function getPermissions() {
  return http.get('/permissions')
}

// 更新资料：返回最新 user
export function updateProfile(payload) {
  return http.put('/profile', payload)
}

// 修改密码：{ oldPassword, newPassword }
export function changePassword(payload) {
  return http.put('/me/password', payload)
}
