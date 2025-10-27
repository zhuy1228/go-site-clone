<template>
  <a-breadcrumb style="margin: 10px 0">
    <a-breadcrumb-item>仿站</a-breadcrumb-item>
    <a-breadcrumb-item>仿站记录</a-breadcrumb-item>
  </a-breadcrumb>
  <div>
    <a-card :bordered="false" class="table-list">
      <a-table :dataSource="dataSource" :columns="columns" >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <span>
              <a @click="openFileDir(record.name)">查看文件</a>
            </span>
          </template>
        </template>
      </a-table>
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
  // 获取网站列表
  const getList = async ()=> {
    App.GetDownloadList().then((res)=>{   
      dataSource.value = JSON.parse(JSON.stringify(res))
    })
  }
  // 打开文件夹
  const openFileDir = async (name) => {
    App.OpenSiteFileDir(name)
  }
  onMounted(()=> {
    getList()
  })
</script>