<template>
  <div>
    <canvas ref="chartRef" :style="{ height: height + 'px' }"></canvas>
  </div>
</template>

<script setup lang="ts">
import { Chart, registerables } from 'chart.js'
import { ref, onMounted, onBeforeUnmount } from 'vue'

// 注册 Chart.js 所有组件
Chart.register(...registerables)

// 定义 props
const props = defineProps<{
  data: any
  height: number
}>()

// chart 实例引用
const chartRef = ref<HTMLCanvasElement | null>(null)
let chartInstance: Chart | null = null

onMounted(() => {
  if (chartRef.value) {
    const ctx = chartRef.value.getContext('2d')
    if (!ctx) return

    chartInstance = new Chart(ctx, {
      type: 'bar',
      data: props.data,
      options: {
        layout: {
          padding: { top: 30, right: 15, left: 10, bottom: 5 }
        },
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { display: false }
        },
        // Chart.js v3+ 的 tooltip 配置在 plugins.tooltip
        // 如果你用的是 v4，tooltips 配置需要改为 plugins.tooltip
        // 这里保留原写法
        tooltips: {
          enabled: true,
          mode: 'index',
          intersect: false
        },
        scales: {
          y: {
            grid: {
              display: true,
              color: 'rgba(255, 255, 255, .2)',
              borderDash: [6],
              borderDashOffset: 6
            },
            ticks: {
              suggestedMin: 0,
              suggestedMax: 1000,
              display: true,
              color: '#fff',
              font: {
                size: 14,
                lineHeight: 1.5,
                weight: '600',
                family: 'Open Sans'
              }
            }
          },
          x: {
            grid: { display: false },
            ticks: {
              display: true,
              color: '#fff',
              font: {
                size: 14,
                lineHeight: 1.5,
                weight: '600',
                family: 'Open Sans'
              }
            }
          }
        }
      }
    })
  }
})

onBeforeUnmount(() => {
  if (chartInstance) {
    chartInstance.destroy()
    chartInstance = null
  }
})
</script>

<style scoped lang="scss">
canvas {
  background-image: linear-gradient(to right, #00369E, #005CFD, #A18DFF);
}
</style>
