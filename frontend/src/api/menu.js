import http from '@/utils/http.js'

// 菜单树（带 children）：GET /menus
export function menuTree() {
  return http.get('/menus')
}

// 全量平铺菜单（用于权限树构建）：GET /menus/all
export function allMenus() {
  return http.get('/menus/all')
}

// 创建菜单：{ parentId, name, permKey, type, icon, url, component, sort, status }
export function createMenu(data) {
  return http.post('/menus', data)
}

// 更新菜单
export function updateMenu(id, data) {
  return http.put(`/menus/${id}`, data)
}

// 删除菜单
export function removeMenu(id) {
  return http.delete(`/menus/${id}`)
}
