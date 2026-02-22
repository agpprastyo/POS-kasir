import { defineConfig } from 'vite'
import { tanstackStart } from '@tanstack/react-start/plugin/vite'
import viteReact from '@vitejs/plugin-react'
import viteTsConfigPaths from 'vite-tsconfig-paths'
import tailwindcss from '@tailwindcss/vite'

const config = defineConfig({
  base: '/',
  plugins: [
    viteTsConfigPaths({
      projects: ['./tsconfig.json'],
    }),
    tailwindcss(),
    tanstackStart({
      spa: {
        enabled: true,
        prerender: {
          enabled: false,
        },
      },
    }),
    viteReact(),
  ],
  build: {
    rollupOptions: {
      output: {
        // manualChunks hanya berlaku untuk client build (bukan SSR)
        // Vite menjalankan 2 build: client (browser) dan server (SSR)
        // Pada SSR, react & dependencies di-external-kan sehingga tidak bisa di-chunk manual
        manualChunks: (id, { getModuleInfo }) => {
          const isSSR = getModuleInfo?.(id)?.isEntry === false && id.includes('node_modules')

          if (id.includes('node_modules/recharts') || id.includes('node_modules/victory')) {
            return 'vendor-charts'
          }
          if (id.includes('node_modules/react-dom') || id.includes('node_modules/react/')) {
            return 'vendor-react'
          }
          if (
            id.includes('node_modules/@tanstack/react-router') ||
            id.includes('node_modules/@tanstack/react-start')
          ) {
            return 'vendor-router'
          }
          if (id.includes('node_modules/@tanstack/react-query')) {
            return 'vendor-query'
          }
          if (
            id.includes('node_modules/i18next') ||
            id.includes('node_modules/react-i18next')
          ) {
            return 'vendor-i18n'
          }
          if (id.includes('node_modules/@radix-ui')) {
            return 'vendor-radix'
          }
        },
      },
    },
  },
})

export default config
