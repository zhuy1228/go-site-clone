

<template>
  <div>
    <a-layout style="min-height: 100vh">
      <SiderMenu />
      <a-layout>
        <Header></Header>
        <a-layout-content style="margin: 0 16px">
          <div style="min-height: calc(100vh - 155px)">
            <slot />
          </div>
        </a-layout-content>
        <a-layout-footer style="text-align: center">
          <!-- Ant Design ©2018 Created by Ant UED -->
        </a-layout-footer>
      </a-layout>
    </a-layout>
  </div>
</template>

<script setup lang="ts">
onMounted(() => {
  
});
</script>