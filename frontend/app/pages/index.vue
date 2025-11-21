<template>
  <div class="dashboard-page">
    <!-- 面包屑导航 -->
    <a-breadcrumb class="page-breadcrumb">
      <a-breadcrumb-item>
        <home-outlined />
        <span>控制台</span>
      </a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="gradient-blob blob-1"></div>
      <div class="gradient-blob blob-2"></div>
      <div class="gradient-blob blob-3"></div>
    </div>

    <!-- 主内容区 -->
    <div class="main-container">
      <!-- 欢迎区域 -->
      <div class="welcome-section">
        <div class="welcome-content">
          <div class="welcome-text">
            <h1 class="welcome-title">
              <rocket-outlined class="welcome-icon" />
              欢迎使用整站下载工具
            </h1>
            <p class="welcome-subtitle">快速克隆任何网站，一键下载所有资源文件</p>
          </div>
          <div class="welcome-actions">
            <a-button type="primary" size="large" @click="goToDownload">
              <template #icon><cloud-download-outlined /></template>
              开始下载
            </a-button>
            <a-button size="large" @click="goToRecords">
              <template #icon><history-outlined /></template>
              查看记录
            </a-button>
          </div>
        </div>
        <div class="welcome-illustration">
          <cloud-server-outlined class="illustration-icon" />
        </div>
      </div>

      <!-- 统计卡片 -->
      <div class="stats-grid">
        <a-card class="stat-card" :bordered="false" hoverable>
          <div class="stat-content">
            <div class="stat-icon sites-icon">
              <global-outlined />
            </div>
            <div class="stat-info">
              <div class="stat-label">已下载站点</div>
              <div class="stat-value">{{ stats.totalSites }}</div>
              <div class="stat-trend positive">
                <arrow-up-outlined />
                <span>较上周 +{{ stats.sitesTrend }}%</span>
              </div>
            </div>
          </div>
        </a-card>

        <a-card class="stat-card" :bordered="false" hoverable>
          <div class="stat-content">
            <div class="stat-icon files-icon">
              <file-outlined />
            </div>
            <div class="stat-info">
              <div class="stat-label">总文件数</div>
              <div class="stat-value">{{ formatNumber(stats.totalFiles) }}</div>
              <div class="stat-trend positive">
                <arrow-up-outlined />
                <span>较上周 +{{ stats.filesTrend }}%</span>
              </div>
            </div>
          </div>
        </a-card>

        <a-card class="stat-card" :bordered="false" hoverable>
          <div class="stat-content">
            <div class="stat-icon storage-icon">
              <database-outlined />
            </div>
            <div class="stat-info">
              <div class="stat-label">占用空间</div>
              <div class="stat-value">{{ stats.totalStorage }}</div>
              <div class="stat-trend positive">
                <arrow-up-outlined />
                <span>较上周 +{{ stats.storageTrend }}%</span>
              </div>
            </div>
          </div>
        </a-card>

        <a-card class="stat-card" :bordered="false" hoverable>
          <div class="stat-content">
            <div class="stat-icon success-icon">
              <check-circle-outlined />
            </div>
            <div class="stat-info">
              <div class="stat-label">成功率</div>
              <div class="stat-value">{{ stats.successRate }}%</div>
              <div class="stat-trend positive">
                <arrow-up-outlined />
                <span>较上周 +{{ stats.rateTrend }}%</span>
              </div>
            </div>
          </div>
        </a-card>
      </div>

      <!-- 功能卡片 -->
      <a-row :gutter="24" class="feature-cards">
        <a-col :xs="24" :sm="12" :lg="8">
          <a-card class="feature-card" :bordered="false" hoverable @click="goToDownload">
            <div class="feature-icon-wrapper download-feature">
              <cloud-download-outlined class="feature-icon" />
            </div>
            <h3 class="feature-title">整站下载</h3>
            <p class="feature-desc">输入网址，自动分析并下载网站的所有资源文件，包括HTML、CSS、JS、图片等</p>
            <a-button type="link" class="feature-link">
              开始使用
              <right-outlined />
            </a-button>
          </a-card>
        </a-col>

        <a-col :xs="24" :sm="12" :lg="8">
          <a-card class="feature-card" :bordered="false" hoverable @click="goToRecords">
            <div class="feature-icon-wrapper records-feature">
              <history-outlined class="feature-icon" />
            </div>
            <h3 class="feature-title">下载记录</h3>
            <p class="feature-desc">查看所有已下载的网站记录，管理文件，支持快速打开本地文件夹查看内容</p>
            <a-button type="link" class="feature-link">
              查看记录
              <right-outlined />
            </a-button>
          </a-card>
        </a-col>

        <a-col :xs="24" :sm="12" :lg="8">
          <a-card class="feature-card" :bordered="false" hoverable>
            <div class="feature-icon-wrapper settings-feature">
              <setting-outlined class="feature-icon" />
            </div>
            <h3 class="feature-title">高级设置</h3>
            <p class="feature-desc">自定义下载配置，设置保存路径、并发数、超时时间等参数，优化下载体验</p>
            <a-button type="link" class="feature-link">
              配置选项
              <right-outlined />
            </a-button>
          </a-card>
        </a-col>
      </a-row>

      <!-- 使用指南 -->
      <a-card class="guide-card" :bordered="false">
        <template #title>
          <div class="card-title-wrapper">
            <book-outlined class="title-icon" />
            <span>快速开始</span>
          </div>
        </template>
        
        <a-row :gutter="[24, 24]">
          <a-col :xs="24" :md="8">
            <div class="guide-step">
              <div class="step-number">1</div>
              <h4 class="step-title">输入网址</h4>
              <p class="step-desc">在整站下载页面输入你想要下载的网站完整地址</p>
            </div>
          </a-col>
          <a-col :xs="24" :md="8">
            <div class="guide-step">
              <div class="step-number">2</div>
              <h4 class="step-title">分析资源</h4>
              <p class="step-desc">系统自动分析网站结构，统计所有可下载的资源文件</p>
            </div>
          </a-col>
          <a-col :xs="24" :md="8">
            <div class="guide-step">
              <div class="step-number">3</div>
              <h4 class="step-title">开始下载</h4>
              <p class="step-desc">点击下载按钮，等待下载完成，即可在本地查看文件</p>
            </div>
          </a-col>
        </a-row>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { 
  HomeOutlined,
  RocketOutlined,
  CloudDownloadOutlined,
  CloudServerOutlined,
  HistoryOutlined,
  GlobalOutlined,
  FileOutlined,
  DatabaseOutlined,
  CheckCircleOutlined,
  ArrowUpOutlined,
  RightOutlined,
  SettingOutlined,
  BookOutlined
} from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 统计数据
const stats = reactive({
  totalSites: 24,
  sitesTrend: 15,
  totalFiles: 15680,
  filesTrend: 23,
  totalStorage: '2.8 GB',
  storageTrend: 18,
  successRate: 98.5,
  rateTrend: 2.5
})

// 格式化数字
const formatNumber = (num: number) => {
  return num.toLocaleString('zh-CN')
}

// 跳转到下载页面
const goToDownload = () => {
  router.push('/site/download')
}

// 跳转到记录页面
const goToRecords = () => {
  router.push('/site/record')
}
</script>

<style scoped>
/* 页面布局 */
.dashboard-page {
  min-height: calc(100vh - 64px);
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

/* 欢迎区域 */
.welcome-section {
  background: white;
  border-radius: 20px;
  padding: 48px;
  margin-bottom: 32px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 48px;
  overflow: hidden;
  position: relative;
  animation: fadeInDown 0.6s ease-out;
}

.welcome-content {
  flex: 1;
  z-index: 1;
}

.welcome-title {
  font-size: 36px;
  font-weight: 700;
  margin: 0 0 16px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  display: flex;
  align-items: center;
  gap: 16px;
}

.welcome-icon {
  font-size: 42px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.welcome-subtitle {
  font-size: 16px;
  color: #666;
  margin: 0 0 32px 0;
  line-height: 1.6;
}

.welcome-actions {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.welcome-actions .ant-btn {
  height: 48px;
  padding: 0 32px;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 500;
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.2);
  transition: all 0.3s ease;
}

.welcome-actions .ant-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(24, 144, 255, 0.3);
}

.welcome-illustration {
  flex-shrink: 0;
  width: 200px;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
  border-radius: 50%;
  animation: pulse 3s ease-in-out infinite;
}

.illustration-icon {
  font-size: 100px;
  color: #1890ff;
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
  animation: fadeInUp 0.6s ease-out 0.1s both;
}

.stat-card {
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

:deep(.stat-card .ant-card-body) {
  padding: 24px;
}

.stat-content {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  flex-shrink: 0;
}

.sites-icon {
  background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
  color: #1890ff;
}

.files-icon {
  background: linear-gradient(135deg, #f3e5f5 0%, #e1bee7 100%);
  color: #722ed1;
}

.storage-icon {
  background: linear-gradient(135deg, #fff3e0 0%, #ffe0b2 100%);
  color: #fa8c16;
}

.success-icon {
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
  color: #52c41a;
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #333;
  margin-bottom: 8px;
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 500;
}

.stat-trend.positive {
  color: #52c41a;
}

/* 功能卡片 */
.feature-cards {
  margin-bottom: 32px;
  animation: fadeInUp 0.6s ease-out 0.2s both;
}

.feature-card {
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  cursor: pointer;
  height: 100%;
  margin-bottom: 24px;
}

.feature-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.15);
}

:deep(.feature-card .ant-card-body) {
  padding: 32px;
  text-align: center;
}

.feature-icon-wrapper {
  width: 80px;
  height: 80px;
  margin: 0 auto 24px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.download-feature {
  background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
}

.records-feature {
  background: linear-gradient(135deg, #f3e5f5 0%, #e1bee7 100%);
}

.settings-feature {
  background: linear-gradient(135deg, #fff3e0 0%, #ffe0b2 100%);
}

.feature-card:hover .feature-icon-wrapper {
  transform: scale(1.1) rotate(5deg);
}

.feature-icon {
  font-size: 40px;
}

.download-feature .feature-icon {
  color: #1890ff;
}

.records-feature .feature-icon {
  color: #722ed1;
}

.settings-feature .feature-icon {
  color: #fa8c16;
}

.feature-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin: 0 0 12px 0;
}

.feature-desc {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
  margin: 0 0 20px 0;
  min-height: 64px;
}

.feature-link {
  padding: 0;
  font-weight: 500;
}

/* 使用指南 */
.guide-card {
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.06);
  animation: fadeInUp 0.6s ease-out 0.3s both;
}

:deep(.guide-card .ant-card-head) {
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

.guide-step {
  text-align: center;
  padding: 24px;
  background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.guide-step:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);
}

.step-number {
  width: 48px;
  height: 48px;
  margin: 0 auto 16px;
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 700;
}

.step-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
}

.step-desc {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
  margin: 0;
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

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .main-container {
    padding: 24px 16px;
  }
}

@media (max-width: 768px) {
  .welcome-section {
    flex-direction: column;
    padding: 32px 24px;
    text-align: center;
  }

  .welcome-title {
    font-size: 28px;
    justify-content: center;
  }

  .welcome-actions {
    justify-content: center;
  }

  .welcome-illustration {
    width: 150px;
    height: 150px;
  }

  .illustration-icon {
    font-size: 80px;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .stat-value {
    font-size: 24px;
  }
}

@media (max-width: 480px) {
  .page-breadcrumb {
    padding: 12px 16px;
  }

  .welcome-section {
    padding: 24px 16px;
  }

  .welcome-title {
    font-size: 24px;
    flex-direction: column;
    gap: 8px;
  }

  .welcome-icon {
    font-size: 32px;
  }

  .welcome-subtitle {
    font-size: 14px;
  }

  .welcome-actions .ant-btn {
    height: 40px;
    padding: 0 24px;
    font-size: 14px;
  }

  :deep(.feature-card .ant-card-body) {
    padding: 24px 16px;
  }

  .feature-desc {
    min-height: auto;
  }
}
</style>
