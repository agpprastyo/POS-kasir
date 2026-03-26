import { defineConfig } from 'vite'
import { TanStackRouterVite } from '@tanstack/router-plugin/vite'
import viteReact from '@vitejs/plugin-react'
import viteTsConfigPaths from 'vite-tsconfig-paths'
import tailwindcss from '@tailwindcss/vite'

const config = defineConfig({
  base: '/',
  plugins: [
    viteTsConfigPaths({
      projects: ['./tsconfig.json'],
    }),
    TanStackRouterVite(),
    tailwindcss(),
    viteReact(),
  ],
  build: {
    rollupOptions: {
      output: {
        manualChunks: (id) => {
          if (id.includes('node_modules')) {
            if (id.includes('react') || id.includes('react-dom') || id.includes('axios') || id.includes('i18next')) {
              return 'vendor-framework'
            }
            if (id.includes('@tanstack')) {
              return 'vendor-tanstack'
            }
            if (id.includes('@radix-ui') || id.includes('lucide-react')) {
              return 'vendor-ui'
            }
            if (id.includes('recharts') || id.includes('victory')) {
              return 'vendor-charts'
            }
          }
        },
      },
    },
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
      '/swagger': 'http://localhost:8080',
      '/healthz': 'http://localhost:8080',
    },
  },
})

export default config
