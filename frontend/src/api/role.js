import http from '@/utils/http.js'

// 分页角色列表：params = { keyword, page, pageSize }
// 返回 PageResult { total, page, pageSize, list }（list 内 role.menus / role.apis 为关联数组）
export function listRoles(params) {
  return http.get('/roles', { params })
}

// 全量角色（用于用户分配角色下拉）：返回 role 数组
export function allRoles() {
  return http.get('/roles/all')
}

// 角色详情（含 menus / apis）
export function getRole(id) {
  return http.get(`/roles/${id}`)
}

// 创建角色：{ name, code, sort, status, remark }
export function createRole(data) {
  return http.post('/roles', data)
}

// 更新角色：{ name, sort, status, remark }
export function updateRole(id, data) {
  return http.put(`/roles/${id}`, data)
}

// 删除角色
export function removeRole(id) {
  return http.delete(`/roles/${id}`)
}

// 角色分配菜单：{ menuIds }
export function assignRoleMenus(id, menuIds) {
  return http.put(`/roles/${id}/menus`, { menuIds })
}

// 角色分配接口：{ apiIds }
export function assignRoleApis(id, apiIds) {
  return http.put(`/roles/${id}/apis`, { apiIds })
}
