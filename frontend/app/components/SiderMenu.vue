<template>
  <div>
    <div class="menu-out">
      <a-layout-sider v-model:collapsed="collapsed" collapsible style="height: 100%;">
        <div class="logo" />
        <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="menuClick">
          <a-menu-item key="">
            <pie-chart-outlined />
            <span>控制台</span>
          </a-menu-item>
          <a-sub-menu key="site">
            <template #title>
              <span>
                <user-outlined />
                <span>仿站</span>
              </span>
            </template>
            <a-menu-item key="download">整站下载</a-menu-item>
            <a-menu-item key="record">仿站记录</a-menu-item>
          </a-sub-menu>
        </a-menu>
      </a-layout-sider>
    </div>
    <div>
      <a-layout-sider v-model:collapsed="collapsed" collapsible style="height: 100%;"></a-layout-sider>
    </div>
  </div>
  
</template>
<script lang="ts" setup>
import {
  PieChartOutlined,
  UserOutlined,
} from '@ant-design/icons-vue';
import { ref } from 'vue';
const collapsed = ref<boolean>(false);
const selectedKeys = ref<string[]>(['']);
function menuClick(e: any) {
  console.log(e);
  const keyList = e.keyPath
  let path = ''
  for (let i = 0; i < keyList.length; i++) {
    const element = keyList[i];
    path += "/"+element
  }
  navigateTo(path)
}


</script>
<style scoped>
.logo {
  height: 32px;
  margin: 16px;
  background: rgba(255, 255, 255, 0.3);
}

.site-layout .site-layout-background {
  background: #fff;
}
[data-theme='dark'] .site-layout .site-layout-background {
  background: #141414;
}

.menu-out{
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  z-index: 9999;
  min-height: 100vh;
}
</style>