<template>
  <a-card :bordered="false" class="widget-1" style="margin-bottom: 20px;">
    <a-statistic
      :title="title"
      :value="value"
      :prefix="prefix"
      :suffix="suffix"
      :precision="0"
      class="text-success"
      :class="'text-' + status"
    />
    <div class="icon" v-html="icon"></div>
  </a-card>
</template>

<script setup lang="ts">
const props = defineProps<{
  title?: string
  value?: number
  prefix?: string
  suffix?: string
  icon?: string
  status?: string
}>()

// 默认值处理（可选）
const {
  title = '',
  value = 0,
  prefix = '',
  suffix = '',
  icon = '',
  status = 'success'
} = props
</script>
