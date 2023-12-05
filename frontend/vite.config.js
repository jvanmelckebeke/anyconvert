import {fileURLToPath, URL} from 'node:url'

import {defineConfig, loadEnv} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig(
    ({mode}) => {
        const {PORT} = loadEnv(mode, process.cwd(), '')

        console.log(PORT)
        return {
            plugins: [
                vue(),
            ],
            resolve: {
                alias: {
                    '@': fileURLToPath(new URL('./src', import.meta.url))
                }
            },
            server: {
                host: true,
                port: PORT
            }
        }
    })
