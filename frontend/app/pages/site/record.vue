<template>
  <div class="record-page">
    <!-- 面包屑导航 -->
    <a-breadcrumb class="page-breadcrumb">
      <a-breadcrumb-item>
        <home-outlined />
        <span>仿站</span>
      </a-breadcrumb-item>
      <a-breadcrumb-item>
        <history-outlined />
        <span>仿站记录</span>
      </a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="gradient-blob blob-1"></div>
      <div class="gradient-blob blob-2"></div>
    </div>

    <!-- 主内容区 -->
    <div class="main-container">
      <!-- 页面头部 -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-left">
            <folder-open-outlined class="header-icon" />
            <div class="header-info">
              <h1 class="page-title">仿站记录</h1>
              <p class="page-subtitle">查看和管理已下载的网站资源</p>
            </div>
          </div>
          <div class="header-right">
            <a-space size="middle">
              <a-button @click="handleRefresh" :loading="loading">
                <template #icon><reload-outlined /></template>
                刷新列表
              </a-button>
              <a-badge :count="dataSource.length" :number-style="{ backgroundColor: '#52c41a' }">
                <a-button type="primary">
                  <template #icon><database-outlined /></template>
                  总计站点
                </a-button>
              </a-badge>
            </a-space>
          </div>
        </div>
      </div>

      <!-- 数据表格卡片 -->
      <a-card :bordered="false" class="table-card">
        <template #title>
          <div class="card-title-wrapper">
            <appstore-outlined class="title-icon" />
            <span>站点列表</span>
          </div>
        </template>
        
        <a-table 
          :dataSource="dataSource" 
          :columns="columns"
          :loading="loading"
          :pagination="paginationConfig"
          :scroll="{ x: 950 }"
          class="custom-table"
        >
          <template #bodyCell="{ column, record, index }">
            <!-- 名称列 -->
            <template v-if="column.key === 'name'">
              <div class="name-cell">
                <a-avatar 
                  :style="{ backgroundColor: getAvatarColor(index) }" 
                  :size="40"
                  shape="square"
                >
                  <template #icon><global-outlined /></template>
                </a-avatar>
                <div class="name-info">
                  <div class="site-name">{{ record.name }}</div>
                  <div class="site-path">{{ record.name }}</div>
                </div>
              </div>
            </template>
            
            <!-- 大小列 -->
            <template v-if="column.key === 'size'">
              <a-tag color="blue" class="size-tag">
                <file-outlined />
                {{ formatSize(record.size) }}
              </a-tag>
            </template>
            
            <!-- 权限列 -->
            <template v-if="column.key === 'mode'">
              <a-tag color="purple">
                <lock-outlined />
                {{ record.mode }}
              </a-tag>
            </template>
            
            <!-- 修改时间列 -->
            <template v-if="column.key === 'modTime'">
              <div class="time-cell">
                <clock-circle-outlined class="time-icon" />
                <span>{{ formatTime(record.modTime) }}</span>
              </div>
            </template>
            
            <!-- 操作列 -->
            <template v-if="column.key === 'action'">
              <a-space>
                <a-button 
                  type="link" 
                  @click="openFileDir(record.name)"
                  class="action-btn"
                >
                  <template #icon><folder-open-outlined /></template>
                  打开文件夹
                </a-button>
                <a-button 
                  type="link" 
                  danger
                  @click="handleDelete(record)"
                  class="action-btn"
                >
                  <template #icon><delete-outlined /></template>
                  删除
                </a-button>
              </a-space>
            </template>
          </template>
          
          <!-- 空状态 -->
          <template #emptyText>
            <a-empty description="暂无仿站记录">
              <template #image>
                <inbox-outlined style="font-size: 64px; color: #d9d9d9;" />
              </template>
            </a-empty>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>
</template>
<script setup>
import { 
  HomeOutlined,
  HistoryOutlined,
  FolderOpenOutlined,
  ReloadOutlined,
  DatabaseOutlined,
  AppstoreOutlined,
  GlobalOutlined,
  FileOutlined,
  LockOutlined,
  ClockCircleOutlined,
  DeleteOutlined,
  InboxOutlined
} from '@ant-design/icons-vue'
import { App } from "../../../bindings/go-site-clone"
import { Modal, message } from 'ant-design-vue'

// 表格列配置
const columns = [
  {
    title: '站点名称',
    dataIndex: 'name',
    key: 'name',
    width: 300,
    ellipsis: true,
  },
  {
    title: '文件大小',
    dataIndex: 'size',
    key: 'size',
    width: 130,
    sorter: (a, b) => a.size - b.size,
  },
  {
    title: '权限',
    dataIndex: 'mode',
    key: 'mode',
    width: 120,
  },
  {
    title: '最后修改时间',
    key: 'modTime',
    dataIndex: 'modTime',
    width: 180,
    sorter: (a, b) => new Date(a.modTime) - new Date(b.modTime),
  },
  {
    title: '操作',
    key: 'action',
    width: 220,
    fixed: 'right',
  },
]

// 分页配置
const paginationConfig = {
  pageSize: 10,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条记录`,
  pageSizeOptions: ['10', '20', '50', '100'],
}

const dataSource = ref([])
const loading = ref(false)

// 获取网站列表
const getList = async () => {
  loading.value = true
  try {
    const res = await App.GetDownloadList()
    dataSource.value = JSON.parse(JSON.stringify(res))
  } catch (error) {
    console.error('获取列表失败:', error)
    message.error('获取站点列表失败')
  } finally {
    loading.value = false
  }
}

// 刷新列表
const handleRefresh = () => {
  getList()
  message.success('刷新成功')
}

// 打开文件夹
const openFileDir = async (name) => {
  try {
    await App.OpenSiteFileDir(name)
    message.success('已打开文件夹')
  } catch (error) {
    console.error('打开文件夹失败:', error)
    message.error('打开文件夹失败')
  }
}

// 删除记录
const handleDelete = (record) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除站点 "${record.name}" 吗？此操作不可恢复。`,
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk: async () => {
      try {
        await App.DeleteSiteFileDir(record.name)
        message.success('删除成功')
        getList()
      } catch (error) {
        console.error('删除失败:', error)
        message.error('删除失败')
      }
    },
  })
}

// 格式化文件大小
const formatSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i]
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '-'
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  
  // 小于1分钟
  if (diff < 60000) {
    return '刚刚'
  }
  // 小于1小时
  if (diff < 3600000) {
    return Math.floor(diff / 60000) + ' 分钟前'
  }
  // 小于24小时
  if (diff < 86400000) {
    return Math.floor(diff / 3600000) + ' 小时前'
  }
  // 小于7天
  if (diff < 604800000) {
    return Math.floor(diff / 86400000) + ' 天前'
  }
  
  // 格式化为完整日期
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取头像颜色
const getAvatarColor = (index) => {
  const colors = [
    '#1890ff', '#52c41a', '#722ed1', '#fa8c16', 
    '#eb2f96', '#13c2c2', '#faad14', '#f5222d'
  ]
  return colors[index % colors.length]
}

onMounted(() => {
  getList()
})
</script>

<style scoped>
/* 页面布局 */
.record-page {
  min-height: calc(100vh - 120px);
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  position: relative;
  overflow: hidden;
}

.page-breadcrumb {
  padding: 16px 24px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  z-index: 10;
  position: relative;
}

.page-breadcrumb span {
  margin-left: 6px;
}

/* 背景装饰 */
.background-decoration {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
}

.gradient-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.25;
  animation: float 25s ease-in-out infinite;
}

.blob-1 {
  width: 600px;
  height: 600px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  top: -15%;
  right: -15%;
  animation-delay: 0s;
}

.blob-2 {
  width: 500px;
  height: 500px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  bottom: -15%;
  left: -15%;
  animation-delay: 12s;
}

/* 主容器 */
.main-container {
  position: relative;
  z-index: 1;
  max-width: 1400px;
  margin: 0 auto;
  padding: 32px 24px;
}

/* 页面头部 */
.page-header {
  background: white;
  border-radius: 16px;
  padding: 28px 32px;
  margin-bottom: 24px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.06);
  animation: fadeInDown 0.6s ease-out;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
  flex: 1;
}

.header-icon {
  font-size: 48px;
  color: #1890ff;
  padding: 16px;
  background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-info {
  flex: 1;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 8px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.header-right {
  flex-shrink: 0;
}

/* 表格卡片 */
.table-card {
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.06);
  overflow: hidden;
  animation: fadeInUp 0.6s ease-out 0.2s both;
}

:deep(.table-card .ant-card-body) {
  padding: 0;
  overflow: hidden;
}

:deep(.table-card .ant-card-head) {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-bottom: 2px solid #f0f0f0;
  padding: 20px 24px;
}

.card-title-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.title-icon {
  font-size: 20px;
  color: #1890ff;
}

/* 自定义表格样式 */
.custom-table {
  background: transparent;
}

:deep(.custom-table .ant-table) {
  background: transparent;
}

:deep(.custom-table .ant-table-thead > tr > th) {
  background: #fafafa;
  font-weight: 600;
  color: #333;
  border-bottom: 2px solid #e8e8e8;
  padding: 16px;
}

:deep(.custom-table .ant-table-tbody > tr) {
  transition: all 0.3s ease;
}

:deep(.custom-table .ant-table-tbody > tr:hover) {
  background: rgba(24, 144, 255, 0.04);
  transform: translateX(4px);
  box-shadow: -4px 0 0 0 #1890ff;
}

:deep(.custom-table .ant-table-tbody > tr > td) {
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 名称单元格 */
.name-cell {
  display: flex;
  align-items: center;
  gap: 12px;
  max-width: 100%;
}

.name-info {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.site-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}

.site-path {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}

/* 大小标签 */
.size-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  font-weight: 500;
  border-radius: 6px;
}

/* 时间单元格 */
.time-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #666;
}

.time-icon {
  color: #999;
  font-size: 14px;
}

/* 操作按钮 */
.action-btn {
  padding: 4px 8px;
  transition: all 0.3s ease;
}

.action-btn:hover {
  transform: translateY(-2px);
}

/* 动画 */
@keyframes float {
  0%, 100% { 
    transform: translate(0, 0) rotate(0deg); 
  }
  33% { 
    transform: translate(40px, -40px) rotate(120deg); 
  }
  66% { 
    transform: translate(-30px, 30px) rotate(240deg); 
  }
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .main-container {
    padding: 24px 16px;
  }
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    align-items: stretch;
  }

  .header-left {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .header-right {
    width: 100%;
  }

  .header-right :deep(.ant-space) {
    width: 100%;
    justify-content: center;
  }

  .page-title {
    font-size: 24px;
  }

  .page-header {
    padding: 20px 16px;
  }

  :deep(.table-card .ant-card-body) {
    padding: 12px;
  }

  :deep(.custom-table .ant-table-thead > tr > th),
  :deep(.custom-table .ant-table-tbody > tr > td) {
    padding: 12px 8px;
  }

  .name-cell {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}

@media (max-width: 480px) {
  .page-breadcrumb {
    padding: 12px 16px;
    font-size: 12px;
  }

  .header-icon {
    font-size: 36px;
    padding: 12px;
  }

  .page-title {
    font-size: 20px;
  }

  .page-subtitle {
    font-size: 12px;
  }

  :deep(.custom-table .ant-table) {
    font-size: 12px;
  }
}
</style>