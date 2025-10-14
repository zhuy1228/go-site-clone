<template>
  <a-card :bordered="false" class="dashboard-bar-line header-solid">
    <template #title>
      <h2>Sales Overview</h2>
      <p>than last year <span class="text-success">+20%</span></p>
    </template>

    <template #extra>
      <a-badge color="primary" class="badge-dot-primary" text="Traffic" />
      <a-badge color="primary" class="badge-dot-secondary" text="Sales" />
    </template>

    <chart-line :height="310" :data="lineChartData" />
  </a-card>
</template>

<script setup lang="ts">
// 定义折线图数据
const lineChartData = {
  labels: ["Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
  datasets: [
    {
      label: "Mobile apps",
      tension: 0.4,
      borderWidth: 0,
      pointRadius: 0,
      borderColor: "#1890FF",
      borderWidth: 3,
      data: [50, 40, 300, 220, 500, 250, 400, 230, 500],
      maxBarThickness: 6
    },
    {
      label: "Websites",
      tension: 0.4,
      borderWidth: 0,
      pointRadius: 0,
      borderColor: "#B37FEB",
      borderWidth: 3,
      data: [30, 90, 40, 140, 290, 290, 340, 230, 400],
      maxBarThickness: 6
    }
  ]
}
</script>
