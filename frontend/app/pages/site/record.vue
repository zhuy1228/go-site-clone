<template>
  <a-breadcrumb style="margin: 10px 0">
    <a-breadcrumb-item>仿站</a-breadcrumb-item>
    <a-breadcrumb-item>仿站记录</a-breadcrumb-item>
  </a-breadcrumb>
  <div>
    <a-card :bordered="false" class="table-list">
      <a-table :dataSource="dataSource" :columns="columns" />
    </a-card>
  </div>
</template>
<script setup>
  import { App } from "../../../bindings/go-site-clone";
  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '大小',
      dataIndex: 'size',
      key: 'size',
    },
    {
      title: '权限',
      dataIndex: 'mode',
      key: 'mode',
    },
    {
      title: '修改时间',
      key: 'modTime',
      dataIndex: 'modTime',
    },
    {
      title: 'Action',
      key: 'action',
    },
  ];

  const dataSource = ref([]);
  const siteList = ref([])
  const getList = async ()=> {
    App.GetDownloadList().then((res)=>{   
      dataSource.value = JSON.parse(JSON.stringify(res))
    })
  }
  onMounted(()=> {
    getList()
  })
</script>