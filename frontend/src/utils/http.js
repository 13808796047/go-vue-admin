import axios from 'axios'
import { ElMessage } from 'element-plus'
import { CODE_SUCCESS } from '@/constants/index.js'
import { getToken, clearAuth } from '@/utils/storage.js'

const http = axios.create({
  baseURL: '/api',
  timeout: 15000
})

// 请求拦截：自动带 token
http.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 后端统一返回体 { code, msg, data }；HTTP 恒为 200，业务码在 body.code
http.interceptors.response.use(
  (res) => {
    const body = res.data
    if (body && typeof body === 'object' && 'code' in body) {
      if (body.code === CODE_SUCCESS) return body.data
      if (body.code === 401) {
        clearAuth()
        ElMessage.error(body.msg || '登录已过期，请重新登录')
        window.location.href = '/login'
        return Promise.reject(new Error(body.msg || '未认证'))
      }
      ElMessage.error(body.msg || '请求失败')
      return Promise.reject(new Error(body.msg || '业务错误'))
    }
    // 非约定体（文件流等）原样返回
    return body
  },
  (err) => {
    // 网络/HTTP 层错误（CORS、502 等）
    const status = err.response?.status
    if (status === 401) {
      clearAuth()
      ElMessage.error('登录已过期，请重新登录')
      window.location.href = '/login'
    } else {
      ElMessage.error(err.response?.data?.msg || err.message || '网络异常')
    }
    return Promise.reject(err)
  }
)

export { http }
export default http
