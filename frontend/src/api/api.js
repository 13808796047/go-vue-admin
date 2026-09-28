import http from '@/utils/http.js'

// 分页接口权限列表：GET /apis
export function listApis(params) {
  return http.get('/apis', { params })
}

// 全量接口权限：GET /apis/all
export function allApis() {
  return http.get('/apis/all')
}
