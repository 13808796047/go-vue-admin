// Dev 服务器（后台）+ Playwright 截图 的驱动脚本
// 用法: node scripts/run-visual-check.cjs
const { spawn } = require('node:child_process')
const path = require('node:path')
const http = require('node:http')

const DIR = path.resolve(__dirname, '..')
const VITE = path.join(DIR, 'node_modules', 'vite', 'bin', 'vite.js')
const PORT = 5173
const BASE = `http://localhost:${PORT}`

function waitUntil(url, timeoutMs = 30000) {
  const start = Date.now()
  return new Promise((resolve, reject) => {
    const attempt = () => {
      const req = http.get(url, (res) => {
        res.resume()
        res.on('end', () => resolve())
      })
      req.on('error', () => {
        if (Date.now() - start > timeoutMs) reject(new Error('dev server not reachable'))
        else setTimeout(attempt, 500)
      })
    }
    attempt()
  })
}

async function main() {
  const vite = spawn(process.execPath, [VITE, '--port', String(PORT), '--strictPort'], {
    cwd: DIR,
    stdio: ['ignore', 'pipe', 'pipe'],
  })
  vite.stdout.on('data', (d) => console.log('[vite]', d.toString().trim()))
  vite.stderr.on('data', (d) => process.stderr.write('[vite:err] ' + d.toString()))

  let ok = false
  try {
    await waitUntil(`${BASE}/`, 30000)
    console.log('dev server ready')
    const script = spawn(process.execPath, [path.join(DIR, 'scripts', 'visual-check.cjs')], {
      cwd: DIR,
      stdio: 'inherit',
      env: { ...process.env, BASE },
    })
    const code = await new Promise((r) => script.on('exit', r))
    ok = code === 0
  } finally {
    vite.kill()
  }
  process.exit(ok ? 0 : 1)
}

main()
