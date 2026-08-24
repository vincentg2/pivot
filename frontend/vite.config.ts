import { fileURLToPath, URL } from 'node:url'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vitest/config'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  envDir: fileURLToPath(new URL('..', import.meta.url)),
  resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
  server: { proxy: { '/api': 'http://localhost:8080', '/health': 'http://localhost:8080' } },
  test: { environment: 'jsdom', globals: true, exclude: ['e2e/**', 'node_modules/**'] },
})
