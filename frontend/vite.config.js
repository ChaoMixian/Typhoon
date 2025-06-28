import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import { VitePWA } from 'vite-plugin-pwa'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['favicon.ico', 'apple-touch-icon.png', 'masked-icon.svg'], // Customize as needed
      manifest: {
        name: 'Typhoon Panel',
        short_name: 'Typhoon',
        description: 'Typhoon Panel for managing Mihomo',
        theme_color: '#ffffff',
        background_color: '#ffffff',
        display: 'standalone',
        scope: '/',
        start_url: '/',
        icons: [
          {
            src: 'img/icons/android-chrome-192x192.png',
            sizes: '192x192',
            type: 'image/png',
          },
          {
            src: 'img/icons/android-chrome-512x512.png',
            sizes: '512x512',
            type: 'image/png',
          },
          {
            src: 'img/icons/apple-touch-icon.png',
            sizes: '180x180',
            type: 'image/png',
            purpose: 'apple touch icon',
          },
          {
            src: 'img/icons/masked-icon.svg',
            sizes: '512x512',
            type: 'image/svg+xml',
            purpose: 'maskable',
          }
        ],
      },
      // Optional: Service worker strategies
      // workbox: {
      //   globPatterns: ['**/*.{js,css,html,ico,png,svg}'],
      //   runtimeCaching: [
      //     {
      //       urlPattern: /^https:\/\/api\.example\.com\/.*/, // Example API caching
      //       handler: 'NetworkFirst',
      //       options: {
      //         cacheName: 'api-cache',
      //         expiration: {
      //           maxEntries: 10,
      //           maxAgeSeconds: 60 * 60 * 24 // 1 day
      //         },
      //         cacheableResponse: {
      //           statuses: [0, 200]
      //         }
      //       }
      //     }
      //   ]
      // }
    })
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
})
