// 时间格式化工具
const pad = (n) => String(n).padStart(2, '0')

/** ISO 字符串 → YYYY-MM-DD HH:mm */
export function formatDateTime(input) {
  if (!input) return '-'
  const d = new Date(input)
  if (Number.isNaN(d.getTime())) return input
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}
