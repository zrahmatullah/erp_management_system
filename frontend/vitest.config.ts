/// <reference types="vitest" />
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./tests/setup.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'lcov', 'json', 'html'],
      reportsDirectory: './coverage',
      include: [
        'src/utils/**/*.ts',
        'src/stores/**/*.ts',
        'src/composables/**/*.ts',
        'src/services/**/*.ts',
        'src/router/guards.ts',
        'src/constants/**/*.ts',
        'src/components/common/**/*.vue'
      ],
      exclude: [
        'src/main.ts',
        'src/vite-env.d.ts',
        'src/types/**',
        'src/**/*.d.ts',
        'node_modules/**',
        'dist/**',
        'src/components/common/GlobalDialogContainer.vue',
        'src/components/common/ToastContainer.vue'
      ],
    },
  },
})
