import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import {fileURLToPath, URL} from 'url'
import {VitePWA} from "vite-plugin-pwa";
// import {VitePWA} from "vite-plugin-pwa";
// https://vitejs.dev/config/
export default defineConfig({
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
            'udap-ui': fileURLToPath(new URL('../common', import.meta.url))
        },
        dedupe: ["moment", "vue", "vue-router"]
    },
    build: {
        rollupOptions: {
            external: ["moment"]
        }
        // sourcemap: process.env.SOURCE_MAP === 'true',
    },
    plugins: [vue(),
        VitePWA({
            mode: 'development',
            registerType: 'autoUpdate',
            injectRegister: 'auto',
            manifest: {
                name: 'UDAP',
                short_name: 'UDAP',
                description: 'Universal Device Aggregation Platform Interface',
                start_url: '/home/dashboard',
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
            devOptions: {
                enabled: true
                /* other options */
            },
            // add this to cache all the imports
        }),
    ],
    server: {
        port: 5045,
    }
})