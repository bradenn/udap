import {defineConfig, loadEnv} from 'vite';
import vue from '@vitejs/plugin-vue'
import {fileURLToPath, URL} from 'url'
import {VitePWA} from "vite-plugin-pwa";
// import {VitePWA} from "vite-plugin-pwa";

// https://vitejs.dev/config/
export default defineConfig(({mode}) => {

    let env = loadEnv(mode, process.cwd(), "UDAP");

    return {
        resolve: {
            alias: {
                '@': fileURLToPath(new URL('./src', import.meta.url)),
                'udap-ui': fileURLToPath(new URL('./common', import.meta.url))
            },
            dedupe: ["moment", "vue", "vue-router"]
        },
        build: {
            rollupOptions: {
                external: ["moment"]
            },

            // sourcemap: process.env.SOURCE_MAP === 'true',
        },
        envDir: "./",
        envPrefix: "UDAP",
        plugins: [vue(),
            VitePWA({
                mode: 'development',
                registerType: 'autoUpdate',
                injectRegister: "auto",
                srcDir: 'src',
                filename: 'pwa.ts',
                strategies: "injectManifest",
                injectManifest: {
                    globPatterns: ['**/*.{js,css,html,svg,png,woff2,ttf,otf}'],
                    maximumFileSizeToCacheInBytes: 5000000,
                },

                manifest: {
                    name: 'UDAP',
                    protocol_handlers: [],
                    orientation: "portrait-primary",
                    short_name: 'UDAP',
                    description: 'Universal Device Aggregation Platform Interface',
                    start_url: '/',
                    display: 'standalone',
                    background_color: '#000000',
                    theme_color: '#000000',
                    icons: [
                        {
                            src: `/udap-mobile-${env.UDAP_APP_ICON}-192x192.png`,
                            sizes: '192x192',
                            type: 'image/png',
                        },
                        {
                            src: `/udap-mobile-${env.UDAP_APP_ICON}-512x512.png`,
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
    }
})
