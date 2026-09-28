import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'node:path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        // 让每个组件可直接使用 tokens，无需重复 @use
        additionalData: `@use "@/styles/tokens.scss" as *;\n`
      }
    }
  },
  server: {
    port: 5173,
    proxy: {
      // 后端就绪后把 /api 转发到 Go 服务；当前 mock 在开发态兜底
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
