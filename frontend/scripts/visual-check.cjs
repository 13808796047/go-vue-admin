// Playwright 截图（对照《设计规范》做视觉校准的取证工具）
const { chromium } = require('playwright')
const path = require('node:path')
const fs = require('node:fs')

const BASE = process.env.BASE || 'http://localhost:5173'
const OUT = path.join(__dirname, '..', 'screenshots')

const PAGES = [
  { name: '02-user', url: '/user-manage' },
  { name: '03-role', url: '/role-manage' },
  { name: '04-menu', url: '/menu-manage' },
  { name: '05-permission', url: '/permission' },
  { name: '06-profile', url: '/profile' },
]

async function shoot(page, p, errors) {
  await page.goto(BASE + p.url, { waitUntil: 'networkidle' })
  await page.waitForTimeout(500)
  await page.screenshot({ path: path.join(OUT, `${p.name}.png`), fullPage: true })
  console.log(`[${p.name}] ${p.url} errors=${errors.length}`, errors.slice(0, 5))
  return { page: p.name, url: p.url, consoleErrors: errors }
}

async function main() {
  const reports = []
  const browser = await chromium.launch()

  // —— 登录页：干净上下文（未登录）——
  {
    const ctx = await browser.newContext({ viewport: { width: 1440, height: 900 } })
    const page = await ctx.newPage()
    const errors = []
    page.on('console', (m) => m.type() === 'error' && errors.push(m.text()))
    page.on('pageerror', (e) => errors.push('pageerror: ' + e.message))
    reports.push(await shoot(page, { name: '01-login', url: '/login' }, errors))
    await ctx.close()
  }

  // —— 管理页：登录一次，token 存 localStorage，循环内不再重复登录 ——
  {
    const ctx = await browser.newContext({ viewport: { width: 1440, height: 900 } })
    const page = await ctx.newPage()
    await page.goto(BASE + '/login', { waitUntil: 'networkidle' })
    await page.getByPlaceholder('用户名').fill('admin')
    await page.getByPlaceholder('密码').fill('123456')
    await page.getByRole('button', { name: '登录' }).click()
    await page.waitForURL('**/profile', { timeout: 10000 })

    for (const p of PAGES) {
      const errors = []
      page.removeAllListeners('console')
      page.removeAllListeners('pageerror')
      page.on('console', (m) => m.type() === 'error' && errors.push(m.text()))
      page.on('pageerror', (e) => errors.push('pageerror: ' + e.message))
      reports.push(await shoot(page, p, errors))
    }
    await ctx.close()
  }

  await browser.close()
  fs.mkdirSync(OUT, { recursive: true })
  fs.writeFileSync(
    path.join(OUT, 'console-report.json'),
    JSON.stringify(reports, null, 2),
  )
  const bad = reports.filter((r) => r.consoleErrors.length)
  console.log('done →', OUT, bad.length ? `⚠ ${bad.length} page(s) with errors` : '· 0 console errors')
  process.exitCode = bad.length ? 1 : 0
}

main().catch((e) => {
  console.error('visual-check failed:', e)
  process.exit(1)
})
