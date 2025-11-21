<template>
  <div class="pack-page">
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
            <appstore-outlined class="page-icon" />
            <div>
              <h1 class="page-title">应用打包</h1>
              <p class="page-subtitle">将下载的网站打包成独立应用程序</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 打包向导 -->
      <a-card class="wizard-card" :bordered="false">
        <a-steps :current="currentStep" class="pack-steps">
          <a-step title="选择网站" />
          <a-step title="配置应用" />
          <a-step title="打包设置" />
          <a-step title="开始打包" />
        </a-steps>

        <!-- 步骤内容 -->
        <div class="step-content">
          <!-- 步骤1: 选择网站 -->
          <div v-if="currentStep === 0" class="step-panel">
            <div class="step-header">
              <h3 class="step-title">选择要打包的网站</h3>
              <a-button @click="getDownloadList" :loading="loading" size="small">
                <template #icon><reload-outlined /></template>
                刷新列表
              </a-button>
            </div>

            <div v-if="loading" class="loading-state">
              <a-spin size="large" tip="加载中..." />
            </div>

            <div v-else-if="availableSites.length === 0" class="empty-state">
              <inbox-outlined class="empty-icon" />
              <p class="empty-text">暂无已下载的网站</p>
              <a-button type="primary" @click="goToDownload">
                <template #icon><cloud-download-outlined /></template>
                去下载网站
              </a-button>
            </div>

            <div v-else class="sites-list">
              <a-radio-group v-model:value="selectedSite" class="site-radio-group">
                <a-card 
                  v-for="site in availableSites" 
                  :key="site.id"
                  class="site-card"
                  :class="{ 'selected': selectedSite === site.id }"
                  :bordered="false"
                  @click="selectedSite = site.id"
                >
                  <a-radio :value="site.id" class="site-radio">
                    <div class="site-content">
                      <a-avatar :style="{ backgroundColor: getAvatarColor(site.name) }">
                        {{ site.name.charAt(0).toUpperCase() }}
                      </a-avatar>
                      <div class="site-details">
                        <div class="site-name">{{ site.name }}</div>
                        <div class="site-info">
                          <span>{{ site.size }}</span>
                          <span>•</span>
                          <span>{{ formatTime(site.modTime) }}</span>
                        </div>
                      </div>
                    </div>
                  </a-radio>
                </a-card>
              </a-radio-group>
            </div>
          </div>

          <!-- 步骤2: 配置应用 -->
          <div v-if="currentStep === 1" class="step-panel">
            <h3 class="step-title">配置应用信息</h3>
            <a-form layout="vertical" class="config-form">
              <a-row :gutter="24">
                <a-col :span="24">
                  <a-form-item label="应用名称">
                    <a-input v-model:value="appConfig.name" placeholder="请输入应用名称" size="large" />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :sm="12">
                  <a-form-item label="应用版本">
                    <a-input v-model:value="appConfig.version" placeholder="1.0.0" size="large" />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :sm="12">
                  <a-form-item label="作者">
                    <a-input v-model:value="appConfig.author" placeholder="请输入作者名称" size="large" />
                  </a-form-item>
                </a-col>
                <a-col :span="24">
                  <a-form-item label="应用描述">
                    <a-textarea 
                      v-model:value="appConfig.description" 
                      placeholder="请输入应用描述"
                      :rows="4"
                    />
                  </a-form-item>
                </a-col>
              </a-row>
            </a-form>
          </div>

          <!-- 步骤3: 打包设置 -->
          <div v-if="currentStep === 2" class="step-panel">
            <h3 class="step-title">打包设置</h3>
            <a-form layout="vertical" class="config-form">
              <a-form-item label="目标平台">
                <a-checkbox-group v-model:value="packConfig.platforms" class="platform-group">
                  <a-card class="platform-card" :bordered="false">
                    <a-checkbox value="windows">
                      <windows-outlined class="platform-icon" />
                      <span>Windows</span>
                    </a-checkbox>
                  </a-card>
                  <a-card class="platform-card" :bordered="false">
                    <a-checkbox value="macos">
                      <apple-outlined class="platform-icon" />
                      <span>macOS</span>
                    </a-checkbox>
                  </a-card>
                  <a-card class="platform-card" :bordered="false">
                    <a-checkbox value="linux">
                      <linux-outlined class="platform-icon" />
                      <span>Linux</span>
                    </a-checkbox>
                  </a-card>
                </a-checkbox-group>
              </a-form-item>

              <a-row :gutter="24">
                <a-col :xs="24" :sm="12">
                  <a-form-item label="窗口宽度">
                    <a-input-number 
                      v-model:value="packConfig.width" 
                      :min="800"
                      class="full-width"
                      size="large"
                    />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :sm="12">
                  <a-form-item label="窗口高度">
                    <a-input-number 
                      v-model:value="packConfig.height" 
                      :min="600"
                      class="full-width"
                      size="large"
                    />
                  </a-form-item>
                </a-col>
              </a-row>

              <a-form-item label="输出目录">
                <a-input 
                  v-model:value="packConfig.outputDir" 
                  placeholder="选择输出目录"
                  size="large"
                >
                  <template #suffix>
                    <folder-outlined class="input-icon" />
                  </template>
                </a-input>
              </a-form-item>
            </a-form>
          </div>

          <!-- 步骤4: 开始打包 -->
          <div v-if="currentStep === 3" class="step-panel">
            <h3 class="step-title">准备打包</h3>
            <div class="pack-summary">
              <a-descriptions :column="1" bordered>
                <a-descriptions-item label="网站名称">{{ selectedSiteName }}</a-descriptions-item>
                <a-descriptions-item label="应用名称">{{ appConfig.name }}</a-descriptions-item>
                <a-descriptions-item label="版本号">{{ appConfig.version }}</a-descriptions-item>
                <a-descriptions-item label="目标平台">
                  <a-tag v-for="platform in packConfig.platforms" :key="platform" color="blue">
                    {{ platform }}
                  </a-tag>
                </a-descriptions-item>
                <a-descriptions-item label="窗口尺寸">
                  {{ packConfig.width }} x {{ packConfig.height }}
                </a-descriptions-item>
              </a-descriptions>

              <div class="pack-action">
                <a-button 
                  type="primary" 
                  size="large" 
                  :loading="packing"
                  @click="startPacking"
                  block
                >
                  <template #icon><rocket-outlined /></template>
                  {{ packing ? '打包中...' : '开始打包' }}
                </a-button>
              </div>

              <!-- 打包进度 -->
              <div v-if="packing" class="pack-progress">
                <a-progress :percent="packProgress" status="active" />
                <p class="progress-text">{{ packProgressText }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 步骤操作 -->
        <div class="step-actions">
          <a-button v-if="currentStep > 0" @click="prevStep" size="large">
            <template #icon><left-outlined /></template>
            上一步
          </a-button>
          <a-button 
            v-if="currentStep < 3" 
            type="primary" 
            @click="nextStep"
            :disabled="!canProceed"
            size="large"
          >
            下一步
            <template #icon><right-outlined /></template>
          </a-button>
        </div>
      </a-card>

      <!-- 帮助提示 -->
      <a-card class="help-card" :bordered="false">
        <template #title>
          <div class="card-title-wrapper">
            <question-circle-outlined class="title-icon" />
            <span>帮助说明</span>
          </div>
        </template>
        <div class="help-content">
          <p><strong>打包流程：</strong></p>
          <ol>
            <li>选择要打包的已下载网站</li>
            <li>填写应用的基本信息（名称、版本、作者等）</li>
            <li>配置打包参数（目标平台、窗口大小等）</li>
            <li>确认信息无误后开始打包</li>
          </ol>
          <p><strong>注意事项：</strong></p>
          <ul>
            <li>打包过程可能需要几分钟，请耐心等待</li>
            <li>确保有足够的磁盘空间用于存储打包文件</li>
            <li>不同平台的打包文件格式不同（Windows: .exe, macOS: .app, Linux: .AppImage）</li>
          </ul>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { useRouter } from 'vue-router';
import {
  AppstoreOutlined,
  WindowsOutlined,
  AppleOutlined,
  FolderOutlined,
  RocketOutlined,
  LeftOutlined,
  RightOutlined,
  QuestionCircleOutlined,
  ReloadOutlined,
  InboxOutlined,
  CloudDownloadOutlined
} from '@ant-design/icons-vue';
import { App } from "../../../bindings/go-site-clone";

const router = useRouter();

// 步骤控制
const currentStep = ref(0);
const selectedSite = ref<string>('');
const packing = ref(false);
const packProgress = ref(0);
const packProgressText = ref('');
const loading = ref(false);

// 从后端获取已下载的网站列表
const availableSites = ref<any[]>([]);

// 获取网站列表
const getDownloadList = async () => {
  loading.value = true;
  try {
    const res = await App.GetDownloadList();
    availableSites.value = res.map((site: any, index: number) => ({
      id: index.toString(),
      name: site.name,
      size: formatSize(site.size),
      files: '未知', // 如果后端提供文件数量可以使用
      modTime: site.modTime,
      rawData: site
    }));
  } catch (error) {
    console.error('获取列表失败:', error);
    message.error('获取站点列表失败');
  } finally {
    loading.value = false;
  }
};

// 格式化文件大小
const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i];
};

// 页面挂载时加载数据
onMounted(() => {
  getDownloadList();
});

// 应用配置
const appConfig = ref({
  name: '',
  version: '1.0.0',
  author: '',
  description: ''
});

// 打包配置
const packConfig = ref({
  platforms: ['windows'],
  width: 1280,
  height: 800,
  outputDir: ''
});

// 计算属性
const selectedSiteName = computed(() => {
  const site = availableSites.value.find(s => s.id === selectedSite.value);
  return site?.name || '';
});

const canProceed = computed(() => {
  if (currentStep.value === 0) return selectedSite.value !== '';
  if (currentStep.value === 1) return appConfig.value.name !== '';
  if (currentStep.value === 2) return packConfig.value.platforms.length > 0;
  return true;
});

// 步骤操作
const nextStep = () => {
  if (currentStep.value < 3) {
    currentStep.value++;
  }
};

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--;
  }
};

// 开始打包
const startPacking = () => {
  packing.value = true;
  packProgress.value = 0;
  
  // 模拟打包进度
  const interval = setInterval(() => {
    packProgress.value += 10;
    
    if (packProgress.value <= 30) {
      packProgressText.value = '正在准备打包环境...';
    } else if (packProgress.value <= 60) {
      packProgressText.value = '正在编译应用...';
    } else if (packProgress.value <= 90) {
      packProgressText.value = '正在生成安装包...';
    } else {
      packProgressText.value = '打包完成！';
    }
    
    if (packProgress.value >= 100) {
      clearInterval(interval);
      packing.value = false;
    }
  }, 500);
};

// 获取头像颜色
const getAvatarColor = (name: string) => {
  const colors = ['#1890ff', '#52c41a', '#722ed1', '#fa8c16', '#eb2f96'];
  const index = name.charCodeAt(0) % colors.length;
  return colors[index];
};

// 格式化时间
const formatTime = (dateStr: string) => {
  if (!dateStr) return '--';
  const date = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  
  const minute = 60 * 1000;
  const hour = 60 * minute;
  const day = 24 * hour;
  
  if (diff < hour) {
    return Math.floor(diff / minute) + ' 分钟前';
  } else if (diff < day) {
    return Math.floor(diff / hour) + ' 小时前';
  } else if (diff < 7 * day) {
    return Math.floor(diff / day) + ' 天前';
  } else {
    return date.toLocaleDateString('zh-CN');
  }
};

// 跳转到下载页面
const goToDownload = () => {
  router.push('/site/download');
};
</script>

<style scoped>
/* 页面布局 */
.pack-page {
  min-height: calc(100vh - 120px);
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  position: relative;
  overflow: hidden;
}

/* 背景装饰 - 复用之前的样式 */
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
  max-width: 1200px;
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
}

.title-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.page-icon {
  font-size: 48px;
  color: #722ed1;
  padding: 16px;
  background: linear-gradient(135deg, #f3e5f5 0%, #e1bee7 100%);
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

/* 向导卡片 */
.wizard-card {
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.06);
  margin-bottom: 24px;
  animation: fadeInUp 0.6s ease-out;
}

:deep(.wizard-card .ant-card-body) {
  padding: 40px;
}

.pack-steps {
  margin-bottom: 48px;
}

/* 步骤内容 */
.step-content {
  min-height: 400px;
  margin-bottom: 32px;
}

.step-panel {
  animation: fadeIn 0.4s ease-out;
}

.step-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.step-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin: 0;
  padding-bottom: 16px;
  border-bottom: 2px solid #f0f0f0;
  flex: 1;
}

.step-header .step-title {
  border-bottom: none;
  padding-bottom: 0;
}

/* 加载和空状态 */
.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-icon {
  font-size: 80px;
  color: #d9d9d9;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  color: #999;
  margin-bottom: 24px;
}

/* 网站列表 */
.sites-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.site-radio-group {
  width: 100%;
  display: contents;
}

.site-card {
  border: 2px solid #f0f0f0;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.site-card.selected {
  border-color: #1890ff;
  background: #e6f7ff;
}

.site-card:hover {
  border-color: #1890ff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.15);
}

.site-radio {
  width: 100%;
}

:deep(.site-radio .ant-radio) {
  display: none;
}

.site-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.site-details {
  flex: 1;
}

.site-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.site-info {
  font-size: 13px;
  color: #999;
  display: flex;
  gap: 8px;
}

/* 配置表单 */
.config-form {
  max-width: 800px;
}

/* 平台选择 */
.platform-group {
  width: 100%;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 16px;
}

.platform-card {
  border: 2px solid #f0f0f0;
  border-radius: 12px;
  text-align: center;
  padding: 8px;
  transition: all 0.3s ease;
}

:deep(.platform-card .ant-card-body) {
  padding: 20px;
}

.platform-card:has(:checked) {
  border-color: #1890ff;
  background: #e6f7ff;
}

.platform-icon {
  font-size: 32px;
  display: block;
  margin-bottom: 8px;
}

/* 打包摘要 */
.pack-summary {
  max-width: 600px;
  margin: 0 auto;
}

.pack-action {
  margin-top: 32px;
}

.pack-progress {
  margin-top: 24px;
  padding: 24px;
  background: #f8f9fa;
  border-radius: 12px;
}

.progress-text {
  text-align: center;
  margin: 12px 0 0 0;
  color: #666;
  font-size: 14px;
}

/* 步骤操作 */
.step-actions {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding-top: 24px;
  border-top: 1px solid #f0f0f0;
}

/* 帮助卡片 */
.help-card {
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.06);
  animation: fadeInUp 0.6s ease-out 0.1s both;
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

.help-content {
  color: #666;
  line-height: 1.8;
}

.help-content strong {
  color: #333;
  font-weight: 600;
}

.help-content ol,
.help-content ul {
  margin: 12px 0;
  padding-left: 24px;
}

.help-content li {
  margin: 8px 0;
}

/* 通用样式 */
.full-width {
  width: 100%;
}

.input-icon {
  color: #999;
  cursor: pointer;
}

.input-icon:hover {
  color: #1890ff;
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

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
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

  :deep(.wizard-card .ant-card-body) {
    padding: 24px 16px;
  }

  .sites-list {
    grid-template-columns: 1fr;
  }

  .platform-group {
    grid-template-columns: 1fr;
  }
}
</style>