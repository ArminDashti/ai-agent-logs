import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'node:path'

const basePath = process.env.VITE_BASE_PATH || '/ai-agent-logs/'
const apiProxyTarget = process.env.VITE_API_PROXY || 'http://127.0.0.1:8120'
const baseNoSlash = basePath.replace(/\/$/, '') || ''

export default defineConfig({
  base: basePath,
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    host: true,
    port: Number(process.env.VITE_DEV_PORT || 5173),
    allowedHosts: true,
    watch: {
      usePolling: process.env.CHOKIDAR_USEPOLLING === 'true',
    },
    hmr: {
      clientPort: Number(process.env.VITE_HMR_CLIENT_PORT || process.env.VITE_DEV_PORT || 5173),
    },
    proxy: {
      [`${baseNoSlash}/api`]: {
        target: apiProxyTarget,
        changeOrigin: true,
        rewrite: (p) => p.replace(new RegExp(`^${baseNoSlash}/api`), '/api'),
      },
      '/api': {
        target: apiProxyTarget,
        changeOrigin: true,
      },
      [`${baseNoSlash}/health`]: {
        target: apiProxyTarget,
        changeOrigin: true,
        rewrite: (p) => p.replace(new RegExp(`^${baseNoSlash}/health`), '/health'),
      },
      '/health': {
        target: apiProxyTarget,
        changeOrigin: true,
      },
    },
  },
})
