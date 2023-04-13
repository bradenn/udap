import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import {fileURLToPath, URL} from 'url'
import {VitePWA} from "vite-plugin-pwa";
// https://vitejs.dev/config/
export default defineConfig({
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        }
    },
    build: {
        sourcemap: process.env.SOURCE_MAP === 'true',
    },
    plugins: [vue(),
        VitePWA({
            mode: 'development',
            manifest: {
                name: 'UDAP',
                short_name: 'UDAP',
                description: 'Universal Device Aggregation Platform Interface',
                start_url: '/',
                display: 'standalone',
                background_color: '#000000',
                theme_color: '#000000',
                icons: [
                    {
                        src: '/pwa-192x192.png',
                        sizes: '192x192',
                        type: 'image/png',
                    },
                    {
                        src: '/pwa-512x512.png',
                        sizes: '512x512',
                        type: 'image/png',
                    },
                ],
            },
            workbox: {
                // Workbox configuration options
            },
        }),],
    server: {
        port: 5045,
    }
})
