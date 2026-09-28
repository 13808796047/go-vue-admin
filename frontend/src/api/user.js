import http from '@/utils/http.js'

// 分页用户列表：params = { keyword, roleId, status, page, pageSize }
// 返回 PageResult { total, page, pageSize, list }（list 内 user.roles 为角色数组）
export function listUsers(params) {
  return http.get('/users', { params })
}

// 创建用户：{ username, nickname, password, email, phone, avatar, status, roleIds }
export function createUser(data) {
  return http.post('/users', data)
}

// 更新用户：{ nickname, email, phone, avatar, status, roleIds }
export function updateUser(id, data) {
  return http.put(`/users/${id}`, data)
}

// 启用/停用用户
export function setUserStatus(id, status) {
  return http.put(`/users/${id}/status`, { status })
}

// 删除用户（软删除）
export function removeUser(id) {
  return http.delete(`/users/${id}`)
}
