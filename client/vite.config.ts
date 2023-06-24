import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import {fileURLToPath, URL} from 'url'

process.env.VUE_APP_VERSION = process.env.npm_package_version
// https://vitejs.dev/config/
export default defineConfig({
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
            'udap-ui': '../common'
        }
    },
    plugins: [vue()],
    server: {
        port: 5002,
    }
})
