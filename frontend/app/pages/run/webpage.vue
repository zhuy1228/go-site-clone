<template>
  <div class="webpage-page">
    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="gradient-blob blob-1"></div>
      <div class="gradient-blob blob-2"></div>
      <div class="gradient-blob blob-3"></div>
    </div>

    <!-- 主内容区 -->
    <div class="main-container">
      <!-- 页面标题 -->
      <div class="page-header">
        <div class="header-content">
          <div class="title-section">
            <global-outlined class="page-icon" />
            <div>
              <h1 class="page-title">网站托管</h1>
              <p class="page-subtitle">使用 Nginx 托管和预览已下载的网站</p>
            </div>
          </div>
          <a-space>
            <a-button size="large" @click="handleRefresh" :loading="loading">
              <template #icon><reload-outlined /></template>
              刷新
            </a-button>
            <a-button 
              type="primary" 
              size="large" 
              :danger="nginxRunning"
              @click="toggleNginx"
              :loading="nginxLoading"
            >
              <template #icon>
                <poweroff-outlined v-if="nginxRunning" />
                <play-circle-outlined v-else />
              </template>
              {{ nginxRunning ? '停止 Nginx' : '启动 Nginx' }}
            </a-button>
          </a-space>
        </div>
      </div>

      <!-- Nginx 状态卡片 -->
      <a-card class="status-card" :bordered="false">
        <div class="status-header">
          <div class="status-title">
            <api-outlined class="title-icon" />
            <span>Nginx 服务状态</span>
          </div>
          <a-badge 
            :status="nginxRunning ? 'processing' : 'default'" 
            :text="nginxRunning ? '运行中' : '未启动'" 
          />
        </div>
        
        <a-row :gutter="24" class="status-info">
          <a-col :xs="24" :sm="12" :lg="6">
            <div class="info-item">
              <div class="info-label">
                <thunderbolt-outlined />
                服务状态
              </div>
              <div class="info-value" :style="{ color: nginxRunning ? '#52c41a' : '#999' }">
                {{ nginxRunning ? '正常运行' : '已停止' }}
              </div>
            </div>
          </a-col>
          <a-col :xs="24" :sm="12" :lg="6">
            <div class="info-item">
              <div class="info-label">
                <cluster-outlined />
                托管站点
              </div>
              <div class="info-value">{{ hostingSites.length }} 个</div>
            </div>
          </a-col>
          <a-col :xs="24" :sm="12" :lg="6">
            <div class="info-item">
              <div class="info-label">
                <clock-circle-outlined />
                运行时长
              </div>
              <div class="info-value">{{ uptime }}</div>
            </div>
          </a-col>
          <a-col :xs="24" :sm="12" :lg="6">
            <div class="info-item">
              <div class="info-label">
                <fund-outlined />
                总访问量
              </div>
              <div class="info-value">{{ totalVisits }}</div>
            </div>
          </a-col>
        </a-row>
      </a-card>

      <!-- 网站托管列表 -->
      <a-card class="sites-card" :bordered="false">
        <template #title>
          <div class="card-title-wrapper">
            <folder-open-outlined class="title-icon" />
            <span>网站托管管理</span>
            <a-badge :count="hostingSites.length" :number-style="{ backgroundColor: '#1890ff' }" />
          </div>
        </template>
        <template #extra>
          <a-button type="primary" @click="showAddModal">
            <template #icon><plus-outlined /></template>
            添加站点
          </a-button>
        </template>

        <div v-if="loading" class="loading-state">
          <a-spin size="large" tip="加载中..." />
        </div>

        <div v-else-if="hostingSites.length === 0" class="empty-state">
          <inbox-outlined class="empty-icon" />
          <p class="empty-text">暂无托管的网站</p>
          <a-button type="primary" @click="showAddModal">
            <template #icon><plus-outlined /></template>
            添加第一个站点
          </a-button>
        </div>

        <a-table 
          v-else
          :columns="columns" 
          :data-source="hostingSites"
          :pagination="{ pageSize: 10, showSizeChanger: true, showTotal: (total: number) => `共 ${total} 个站点` }"
          :scroll="{ x: 1050 }"
          class="hosting-table"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'site'">
              <div class="site-cell">
                <a-avatar :style="{ backgroundColor: getAvatarColor(record.name) }" size="small">
                  {{ record.name.charAt(0).toUpperCase() }}
                </a-avatar>
                <div class="site-details">
                  <div class="site-name">{{ record.name }}</div>
                  <div class="site-path">{{ record.path }}</div>
                </div>
              </div>
            </template>

            <template v-else-if="column.key === 'domain'">
              <div class="domain-cell">
                <a-tag color="blue" v-for="domain in record.domains" :key="domain">
                  {{ domain }}
                </a-tag>
              </div>
            </template>

            <template v-else-if="column.key === 'port'">
              <a-tag color="cyan">{{ record.port }}</a-tag>
            </template>

            <template v-else-if="column.key === 'status'">
              <a-badge 
                :status="record.enabled ? 'success' : 'default'"
                :text="record.enabled ? '已启用' : '已禁用'"
              />
            </template>

            <template v-else-if="column.key === 'action'">
              <a-space size="small">
                <a-button type="link" size="small" @click="openSite(record)" :disabled="!nginxRunning || !record.enabled">
                  <template #icon><eye-outlined /></template>
                  访问
                </a-button>
                <a-button type="link" size="small" @click="openFolder(record)">
                  <template #icon><folder-outlined /></template>
                  文件夹
                </a-button>
                <a-button type="link" size="small" @click="editSite(record)">
                  <template #icon><edit-outlined /></template>
                  编辑
                </a-button>
                <a-switch 
                  v-model:checked="record.enabled" 
                  size="small"
                  @change="toggleSiteStatus(record)"
                />
                <a-popconfirm
                  title="确定要删除此站点吗？"
                  ok-text="确定"
                  cancel-text="取消"
                  @confirm="deleteSite(record)"
                >
                  <a-button type="link" size="small" danger>
                    <template #icon><delete-outlined /></template>
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-card>
    </div>

    <!-- 添加/编辑站点对话框 -->
    <a-modal
      v-model:open="modalVisible"
      :title="editingRecord ? '编辑站点' : '添加站点'"
      width="700px"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
    >
      <a-form :model="formData" layout="vertical" ref="formRef">
        <a-form-item label="选择网站" name="name" :rules="[{ required: true, message: '请选择网站' }]">
          <a-select v-model:value="formData.name" placeholder="请选择已下载的网站" @change="onSiteSelect">
            <a-select-option v-for="site in availableSites" :key="site.name" :value="site.name">
              {{ site.name }}
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="域名配置" name="domains" :rules="[{ required: true, message: '请至少添加一个域名' }]">
          <a-select
            v-model:value="formData.domains"
            mode="tags"
            placeholder="输入域名后按回车添加"
            :token-separators="[',', ' ']"
          >
          </a-select>
          <div class="form-tip">
            <info-circle-outlined />
            <span>可以添加多个域名,如 example.com 或 127.0.0.1</span>
          </div>
        </a-form-item>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="监听端口" name="port" :rules="[{ required: true, message: '请输入端口' }]">
              <a-input-number
                v-model:value="formData.port"
                :min="80"
                :max="65535"
                placeholder="80"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="默认首页" name="index">
              <a-input v-model:value="formData.index" placeholder="index.html" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="网站路径" name="path" :rules="[{ required: true, message: '请输入网站路径' }]">
          <a-input v-model:value="formData.path" placeholder="相对路径，如: www/example.com" />
          <div class="form-tip">
            <info-circle-outlined />
            <span>相对于项目根目录的路径</span>
          </div>
        </a-form-item>

        <a-form-item label="启用状态">
          <a-switch v-model:checked="formData.enabled" />
          <span style="margin-left: 12px; color: #666;">{{ formData.enabled ? '启用后立即生效' : '暂不启用' }}</span>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { useRouter } from 'vue-router';
import {
  GlobalOutlined,
  PlayCircleOutlined,
  PoweroffOutlined,
  ApiOutlined,
  FolderOpenOutlined,
  ReloadOutlined,
  InboxOutlined,
  PlusOutlined,
  EyeOutlined,
  FolderOutlined,
  EditOutlined,
  DeleteOutlined,
  ThunderboltOutlined,
  ClusterOutlined,
  ClockCircleOutlined,
  FundOutlined,
  InfoCircleOutlined
} from '@ant-design/icons-vue';
import { App } from "../../../bindings/go-site-clone";

const router = useRouter();

// 状态管理
const loading = ref(false);
const nginxRunning = ref(false);
const nginxLoading = ref(false);
const uptime = ref('--');
const totalVisits = ref(0);

// 托管站点列表
const hostingSites = ref<any[]>([]);
const availableSites = ref<any[]>([]);

// 对话框
const modalVisible = ref(false);
const editingRecord = ref<any>(null);
const formRef = ref();
const formData = ref({
  name: '',
  domains: ['127.0.0.1'],
  port: 8080,
  path: '',
  index: 'index.html',
  enabled: true
});

// 表格列配置
const columns = [
  {
    title: '网站',
    key: 'site',
    width: 280,
    ellipsis: true,
  },
  {
    title: '域名',
    key: 'domain',
    width: 250,
    ellipsis: true,
  },
  {
    title: '端口',
    key: 'port',
    width: 100,
    align: 'center',
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    align: 'center',
  },
  {
    title: '操作',
    key: 'action',
    width: 320,
    align: 'center',
    fixed: 'right',
  },
];

// 初始化
onMounted(() => {
  loadData();
});

// 加载数据
const loadData = async () => {
  await Promise.all([
    getAvailableSites(),
    getHostingSites(),
    checkNginxStatus()
  ]);
};

const checkNginxStatus = async () => {
  try {
    const status = await App.CheckNginxStatus();
    nginxRunning.value = status;
    console.log("Nginx running status:", status);
    
  } catch (error) {
    console.error('获取 Nginx 状态失败:', error);
  }
};

// 获取可用网站列表(已下载的)
const getAvailableSites = async () => {
  try {
    const res = await App.GetDownloadList();
    availableSites.value = res || [];
  } catch (error) {
    console.error('获取网站列表失败:', error);
  }
};

// 获取托管列表
const getHostingSites = async () => {
  loading.value = true;
  try {
    const sites = await App.GetAllNginxSites();
    hostingSites.value = sites || [];
  } catch (error) {
    console.error('获取托管列表失败:', error);
    message.error('获取托管列表失败');
  } finally {
    loading.value = false;
  }
};

// 刷新数据
const handleRefresh = () => {
  loadData();
  message.success('刷新成功');
};

// 启停 Nginx
const toggleNginx = async () => {
  nginxLoading.value = true;
  try {
    if (nginxRunning.value) {
      // TODO: 调用后端接口停止 Nginx
      await App.StopNginx();
      nginxRunning.value = false;
      uptime.value = '--';
      message.success('Nginx 已停止');
    } else {
      // TODO: 调用后端接口启动 Nginx
      await App.StartNginx();
      nginxRunning.value = true;
      uptime.value = '00:00:00';
      message.success('Nginx 已启动');
    }
  } catch (error) {
    message.error(nginxRunning.value ? '停止失败' : '启动失败');
  } finally {
    nginxLoading.value = false;
  }
};

// 显示添加对话框
const showAddModal = () => {
  editingRecord.value = null;
  formData.value = {
    name: '',
    domains: ['127.0.0.1'],
    port: 8080,
    path: '',
    index: 'index.html',
    enabled: true
  };
  modalVisible.value = true;
};

// 编辑站点
const editSite = (record: any) => {
  editingRecord.value = record;
  formData.value = { ...record };
  modalVisible.value = true;
};

// 网站选择改变
const onSiteSelect = (value: string) => {
  const site = availableSites.value.find(s => s.name === value);
  if (site) {
    // 设置网站路径（从工作目录）
    formData.value.path = `www/${value}`;
    formData.value.domains = [value, '127.0.0.1'];
  }
};

// 确定对话框
const handleModalOk = async () => {
  try {
    await formRef.value?.validate();
    
    if (editingRecord.value) {
      // 更新站点配置
      await App.UpdateNginxSite({
        ID: "",
        Name: formData.value.name,
        Domains: formData.value.domains,
        Port: formData.value.port,
        Path: formData.value.path,
        Index: formData.value.index,
        Enabled: formData.value.enabled
      });
      message.success('更新成功');
    } else {
      // 添加站点配置
      await App.AddNginxSite({
        ID: "",
        Name: formData.value.name,
        Domains: formData.value.domains,
        Port: formData.value.port,
        Path: formData.value.path,
        Index: formData.value.index,
        Enabled: formData.value.enabled
      });
      message.success('站点配置创建成功');
    }
    
    // 重新加载站点列表
    await getHostingSites();
    
    // 如果 nginx 正在运行，重载配置
    if (nginxRunning.value) {
      await App.ReloadNginx();
      message.info('Nginx 配置已重载');
    }
    
    modalVisible.value = false;
  } catch (error: any) {
    console.error('操作失败:', error);
    message.error(error.message || '操作失败');
  }
};

// 取消对话框
const handleModalCancel = () => {
  modalVisible.value = false;
};

// 切换站点状态
const toggleSiteStatus = async (record: any) => {
  try {
    if (record.enabled) {
      await App.EnableNginxSite(record.name);
      message.success('站点已启用');
    } else {
      await App.DisableNginxSite(record.name);
      message.success('站点已禁用');
    }
    
    // 如果 nginx 正在运行，重载配置
    if (nginxRunning.value) {
      await App.ReloadNginx();
    }
  } catch (error: any) {
    console.error('切换状态失败:', error);
    message.error(error.message || '操作失败');
    // 回滚状态
    record.enabled = !record.enabled;
  }
};

// 删除站点
const deleteSite = async (record: any) => {
  try {
    await App.DeleteNginxSite(record.name);
    message.success('站点配置已删除');
    
    // 重新加载站点列表
    await getHostingSites();
    
    // 如果 nginx 正在运行，重载配置
    if (nginxRunning.value) {
      await App.ReloadNginx();
    }
  } catch (error: any) {
    console.error('删除失败:', error);
    message.error(error.message || '删除失败');
  }
};

// 访问站点
const openSite = (record: any) => {
  const url = `http://${record.domains[0]}:${record.port}`;
  window.open(url, '_blank');
};

// 打开文件夹
const openFolder = async (record: any) => {
  try {
    await App.OpenSiteFileDir(record.name);
    message.success('已打开文件夹');
  } catch (error) {
    console.error('打开文件夹失败:', error);
    message.error('打开文件夹失败');
  }
};

// 获取头像颜色
const getAvatarColor = (name: string) => {
  const colors = ['#1890ff', '#52c41a', '#722ed1', '#fa8c16', '#eb2f96'];
  const index = name.charCodeAt(0) % colors.length;
  return colors[index];
};
</script>

<style scoped>
/* 页面布局 */
.webpage-page {
  min-height: calc(100vh - 120px);
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  position: relative;
  overflow: hidden;
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
}

.blob-2 {
  width: 500px;
  height: 500px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  bottom: -15%;
  left: -15%;
  animation-delay: 12s;
}

.blob-3 {
  width: 450px;
  height: 450px;
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  top: 50%;
  left: 50%;
  animation-delay: 8s;
}

/* 主容器 */
.main-container {
  position: relative;
  z-index: 1;
  max-width: 1400px;
  margin: 0 auto;
  padding: 32px 24px;
}

/* 页面标题 */
.page-header {
  background: white;
  border-radius: 20px;
  padding: 32px;
  margin-bottom: 24px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
  animation: fadeInDown 0.6s ease-out;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24px;
  flex-wrap: wrap;
}

.title-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.page-icon {
  font-size: 48px;
  color: #1890ff;
  padding: 16px;
  background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
  border-radius: 16px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0;
  color: #333;
}

.page-subtitle {
  font-size: 14px;
  color: #666;
  margin: 4px 0 0 0;
}

/* 卡片样式 */
.status-card,
.sites-card,
.config-card {
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.06);
  margin-bottom: 24px;
  animation: fadeInUp 0.6s ease-out;
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

/* 状态卡片 */
.status-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.status-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.status-info {
  margin-top: 24px;
}

.info-item {
  padding: 20px;
  background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
  border-radius: 12px;
  margin-bottom: 16px;
  border: 1px solid #f0f0f0;
  transition: all 0.3s ease;
}

.info-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  border-color: #1890ff;
}

.info-label {
  font-size: 13px;
  color: #666;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.info-value {
  font-size: 20px;
  font-weight: 700;
  color: #333;
}

/* 托管表格 */
.hosting-table {
  margin-top: 16px;
}

:deep(.hosting-table .ant-table) {
  background: transparent;
}

:deep(.hosting-table .ant-table-thead > tr > th) {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  font-weight: 600;
  color: #333;
  border-bottom: 2px solid #1890ff;
  padding: 16px;
}

:deep(.hosting-table .ant-table-tbody > tr) {
  transition: all 0.3s ease;
}

:deep(.hosting-table .ant-table-tbody > tr:hover) {
  background: #e6f7ff !important;
}

:deep(.hosting-table .ant-table-tbody > tr > td) {
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.site-cell {
  display: flex;
  align-items: center;
  gap: 12px;
  max-width: 100%;
}

.site-details {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.site-name {
  font-size: 14px;
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

.domain-cell {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  max-width: 100%;
  overflow: hidden;
}

.domain-cell :deep(.ant-tag) {
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 加载和空状态 */
.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 80px 20px;
}

.empty-icon {
  font-size: 100px;
  color: #d9d9d9;
  margin-bottom: 24px;
  opacity: 0.6;
}

.empty-text {
  font-size: 16px;
  color: #999;
  margin-bottom: 32px;
}

/* 表单提示 */
.form-tip {
  margin-top: 8px;
  padding: 8px 12px;
  background: #e6f7ff;
  border-radius: 6px;
  font-size: 13px;
  color: #1890ff;
  display: flex;
  align-items: center;
  gap: 6px;
}

/* 对话框样式优化 */
:deep(.ant-modal-header) {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-bottom: 2px solid #1890ff;
}

:deep(.ant-modal-title) {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

/* 按钮样式优化 */
:deep(.ant-btn-primary) {
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.2);
}

:deep(.ant-btn-primary:hover) {
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
}

/* 网站网格 */
.sites-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.site-item {
  border-radius: 12px;
  transition: all 0.3s ease;
}

.site-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.site-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.site-info {
  flex: 1;
  min-width: 0;
}

.site-url {
  font-size: 13px;
  color: #999;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 4px;
}

.site-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

/* 表单 */
.full-width {
  width: 100%;
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
@media (max-width: 768px) {
  .main-container {
    padding: 24px 16px;
  }

  .page-header {
    padding: 24px 16px;
  }

  .title-section {
    flex-direction: column;
    align-items: flex-start;
  }

  .sites-grid {
    grid-template-columns: 1fr;
  }
}
</style>