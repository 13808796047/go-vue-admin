import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

/**
 * 通用表格数据管理：分页、加载、搜索、增删改后的刷新。
 * @param fetcher (params) => { list, total }  params = { page, pageSize, keyword }
 */
export function useTable(fetcher) {
  const data = ref([])
  const total = ref(0)
  const loading = ref(false)
  const keyword = ref('')
  const pagination = reactive({ page: 1, pageSize: 10 })

  async function reload() {
    loading.value = true
    try {
      const res = await fetcher({ ...pagination, keyword: keyword.value })
      data.value = res.list || []
      total.value = res.total || 0
    } catch (e) {
      ElMessage.error(e.message || '加载失败')
    } finally {
      loading.value = false
    }
  }

  async function onSearch(kw = '') {
    keyword.value = kw
    pagination.page = 1
    await reload()
  }

  async function onPage(page) {
    pagination.page = page
    await reload()
  }

  async function onPageSize(size) {
    pagination.pageSize = size
    pagination.page = 1
    await reload()
  }

  async function confirmRemove(label) {
    try {
      await ElMessageBox.confirm(`确定删除「${label}」吗？删除后不可恢复。`, '删除确认', {
        type: 'warning',
        confirmButtonText: '删除',
        cancelButtonText: '取消'
      })
      return true
    } catch {
      return false
    }
  }

  return {
    data,
    total,
    loading,
    keyword,
    pagination,
    reload,
    onSearch,
    onPage,
    onPageSize,
    confirmRemove
  }
}
