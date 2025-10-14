<template>
<div class="header-outside">
   <a-layout-header class="header">

    <div class="header-menu">
      <div></div>
      <div>
      
      </div>
    </div>
  </a-layout-header>
</div>
 

</template>
<script lang="ts" setup>
import { DownOutlined } from '@ant-design/icons-vue';

</script>
<style scoped>
.header-menu{
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-outside {
  position: sticky;
  top: 0;
  left: 0;
  right: 0;
  /* z-index: 9999; */
}

</style>
