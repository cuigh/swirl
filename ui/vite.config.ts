import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

const config = loadEnv('development', './')

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      vue: "vue/dist/vue.esm-bundler.js",
      'vue-i18n': "vue-i18n/dist/vue-i18n.cjs.js",
      '@': path.resolve(__dirname, './src'),
    }
  },
  build: {
    cssCodeSplit: false,
    // rollupOptions: {
    //   output: {
    //     manualChunks(id) {
    //       if (id.includes('node_modules')) {
    //         return id.toString().split('node_modules/')[1].split('/')[0].toString();
    //       }
    //     }
    //   }
    // },
  },
  server: {
    port: 3002,
    proxy: {
      '/api': {
        target: 'http://' + config.VITE_PROXY_URL,
        changeOrigin: true,
      },
    }
  },
})
