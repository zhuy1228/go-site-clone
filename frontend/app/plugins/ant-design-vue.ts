import { defineNuxtPlugin } from '#app'
import Antd, { notification }  from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css' // v4 使用新样式重置文件

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.use(Antd)
  nuxtApp.provide('notification', notification)
})