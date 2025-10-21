// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  app: {
    head: {
      title: '仿站小工具',
      meta: [
        { name: 'description', content: '这是我的 Nuxt4 应用默认描述' }
      ]
    },
    baseURL: './' // 确保相对路径，方便打包到 Wails
  },
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  ssr: false, // 关闭 SSR，改为 CSR
  nitro: {
    preset: 'static',
    output: {
      publicDir: './dist' // 相对 Nuxt3 项目根目录的路径
    }
  },
  plugins: [
    '~/plugins/ant-design-vue',
  ]
})
