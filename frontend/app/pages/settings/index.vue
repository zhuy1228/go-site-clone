<template>
  <div class="settings-page">
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
            <setting-outlined class="page-icon" />
            <div>
              <h1 class="page-title">系统设置</h1>
              <p class="page-subtitle">配置应用参数和系统行为</p>
            </div>
          </div>
          <a-space>
            <a-button size="large" @click="handleReset">
              <template #icon><undo-outlined /></template>
              重置默认
            </a-button>
            <a-button type="primary" size="large" @click="handleSave" :loading="saving">
              <template #icon><save-outlined /></template>
              保存配置
            </a-button>
          </a-space>
        </div>
      </div>

      <!-- 配置面板 -->
      <a-row :gutter="24">
        <a-col :xs="24" :lg="16">
          <!-- 基础配置 -->
          <a-card class="config-card" :bordered="false">
            <template #title>
              <div class="card-title-wrapper">
                <control-outlined class="title-icon" />
                <span>基础配置</span>
              </div>
            </template>

            <a-form :model="config" layout="vertical">
              <a-form-item label="应用名称">
                <a-input 
                  v-model:value="config.appName" 
                  placeholder="Go Site Clone"
                  size="large"
                >
                  <template #prefix>
                    <appstore-outlined />
                  </template>
                </a-input>
              </a-form-item>

              <a-form-item label="监听端口">
                <a-input-number 
                  v-model:value="config.port" 
                  :min="1024" 
                  :max="65535"
                  placeholder="6997"
                  size="large"
                  style="width: 100%"
                >
                  <template #prefix>
                    <api-outlined />
                  </template>
                </a-input-number>
                <div class="form-tip">
                  <info-circle-outlined />
                  <span>应用启动时监听的端口号,修改后需要重启应用</span>
                </div>
              </a-form-item>
            </a-form>
          </a-card>

          <!-- 路径配置 -->
          <a-card class="config-card" :bordered="false">
            <template #title>
              <div class="card-title-wrapper">
                <folder-outlined class="title-icon" />
                <span>路径配置</span>
              </div>
            </template>

            <a-form :model="config" layout="vertical">
              <a-form-item label="网站文件保存目录">
                <a-input 
                  v-model:value="config.siteFileDir" 
                  placeholder="www"
                  size="large"
                >
                  <template #prefix>
                    <folder-open-outlined />
                  </template>
                  <template #suffix>
                    <a-button type="link" size="small" @click="selectDirectory('siteFileDir')">
                      浏览
                    </a-button>
                  </template>
                </a-input>
                <div class="form-tip">
                  <info-circle-outlined />
                  <span>下载的网站文件存储位置,相对于应用根目录</span>
                </div>
              </a-form-item>

              <a-form-item label="打包软件保存目录">
                <a-input 
                  v-model:value="config.packSiteFileDir" 
                  placeholder="site-dist"
                  size="large"
                >
                  <template #prefix>
                    <inbox-outlined />
                  </template>
                  <template #suffix>
                    <a-button type="link" size="small" @click="selectDirectory('packSiteFileDir')">
                      浏览
                    </a-button>
                  </template>
                </a-input>
                <div class="form-tip">
                  <info-circle-outlined />
                  <span>打包后的应用程序存储位置</span>
                </div>
              </a-form-item>
            </a-form>
          </a-card>

          <!-- 下载配置 -->
          <a-card class="config-card" :bordered="false">
            <template #title>
              <div class="card-title-wrapper">
                <cloud-download-outlined class="title-icon" />
                <span>下载配置</span>
              </div>
            </template>

            <a-form :model="downloadConfig" layout="vertical">
              <a-row :gutter="16">
                <a-col :xs="24" :sm="12">
                  <a-form-item label="并发下载数">
                    <a-input-number 
                      v-model:value="downloadConfig.concurrency" 
                      :min="1" 
                      :max="20"
                      size="large"
                      style="width: 100%"
                    />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :sm="12">
                  <a-form-item label="超时时间(秒)">
                    <a-input-number 
                      v-model:value="downloadConfig.timeout" 
                      :min="10" 
                      :max="300"
                      size="large"
                      style="width: 100%"
                    />
                  </a-form-item>
                </a-col>
              </a-row>

              <a-row :gutter="16">
                <a-col :xs="24" :sm="12">
                  <a-form-item label="重试次数">
                    <a-input-number 
                      v-model:value="downloadConfig.retryTimes" 
                      :min="0" 
                      :max="10"
                      size="large"
                      style="width: 100%"
                    />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :sm="12">
                  <a-form-item label="User-Agent">
                    <a-select 
                      v-model:value="downloadConfig.userAgent" 
                      size="large"
                    >
                      <a-select-option value="chrome">Chrome</a-select-option>
                      <a-select-option value="firefox">Firefox</a-select-option>
                      <a-select-option value="safari">Safari</a-select-option>
                      <a-select-option value="edge">Edge</a-select-option>
                    </a-select>
                  </a-form-item>
                </a-col>
              </a-row>

              <a-form-item label="下载选项">
                <a-space direction="vertical" style="width: 100%">
                  <a-checkbox v-model:checked="downloadConfig.downloadImages">
                    下载图片资源
                  </a-checkbox>
                  <a-checkbox v-model:checked="downloadConfig.downloadCSS">
                    下载 CSS 样式
                  </a-checkbox>
                  <a-checkbox v-model:checked="downloadConfig.downloadJS">
                    下载 JavaScript 脚本
                  </a-checkbox>
                  <a-checkbox v-model:checked="downloadConfig.downloadFonts">
                    下载字体文件
                  </a-checkbox>
                </a-space>
              </a-form-item>
            </a-form>
          </a-card>

          <!-- Nginx 配置 -->
          <a-card class="config-card" :bordered="false">
            <template #title>
              <div class="card-title-wrapper">
                <global-outlined class="title-icon" />
                <span>Nginx 配置</span>
              </div>
            </template>

            <a-form :model="nginxConfig" layout="vertical">
              <a-row :gutter="16">
                <a-col :xs="24" :sm="12">
                  <a-form-item label="默认端口">
                    <a-input-number 
                      v-model:value="nginxConfig.defaultPort" 
                      :min="80" 
                      :max="65535"
                      size="large"
                      style="width: 100%"
                    />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :sm="12">
                  <a-form-item label="工作进程数">
                    <a-input-number 
                      v-model:value="nginxConfig.workerProcesses" 
                      :min="1" 
                      :max="16"
                      size="large"
                      style="width: 100%"
                    />
                  </a-form-item>
                </a-col>
              </a-row>

              <a-form-item label="Nginx 选项">
                <a-space direction="vertical" style="width: 100%">
                  <a-checkbox v-model:checked="nginxConfig.autoStart">
                    应用启动时自动启动 Nginx
                  </a-checkbox>
                  <a-checkbox v-model:checked="nginxConfig.gzip">
                    启用 Gzip 压缩
                  </a-checkbox>
                  <a-checkbox v-model:checked="nginxConfig.accessLog">
                    记录访问日志
                  </a-checkbox>
                </a-space>
              </a-form-item>
            </a-form>
          </a-card>

          <!-- 界面配置 -->
          <a-card class="config-card" :bordered="false">
            <template #title>
              <div class="card-title-wrapper">
                <bg-colors-outlined class="title-icon" />
                <span>界面配置</span>
              </div>
            </template>

            <a-form :model="uiConfig" layout="vertical">
              <a-form-item label="主题模式">
                <a-radio-group v-model:value="uiConfig.theme" size="large" button-style="solid">
                  <a-radio-button value="light">
                    <sun-outlined /> 浅色
                  </a-radio-button>
                  <a-radio-button value="dark">
                    <moon-outlined /> 深色
                  </a-radio-button>
                  <a-radio-button value="auto">
                    <sync-outlined /> 自动
                  </a-radio-button>
                </a-radio-group>
              </a-form-item>

              <a-form-item label="语言">
                <a-select v-model:value="uiConfig.language" size="large">
                  <a-select-option value="zh-CN">简体中文</a-select-option>
                  <a-select-option value="zh-TW">繁體中文</a-select-option>
                  <a-select-option value="en-US">English</a-select-option>
                  <a-select-option value="ja-JP">日本語</a-select-option>
                </a-select>
              </a-form-item>

              <a-form-item label="界面选项">
                <a-space direction="vertical" style="width: 100%">
                  <a-checkbox v-model:checked="uiConfig.showNotification">
                    显示系统通知
                  </a-checkbox>
                  <a-checkbox v-model:checked="uiConfig.autoRefresh">
                    自动刷新数据
                  </a-checkbox>
                  <a-checkbox v-model:checked="uiConfig.compactMode">
                    紧凑模式
                  </a-checkbox>
                </a-space>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <!-- 右侧信息栏 -->
        <a-col :xs="24" :lg="8">
          <!-- 系统信息 -->
          <a-card class="info-card" :bordered="false">
            <template #title>
              <div class="card-title-wrapper">
                <info-circle-outlined class="title-icon" />
                <span>系统信息</span>
              </div>
            </template>

            <div class="info-list">
              <div class="info-item-row">
                <span class="info-label">应用版本</span>
                <span class="info-value">v1.0.0</span>
              </div>
              <a-divider style="margin: 12px 0" />
              <div class="info-item-row">
                <span class="info-label">运行环境</span>
                <span class="info-value">Windows</span>
              </div>
              <a-divider style="margin: 12px 0" />
              <div class="info-item-row">
                <span class="info-label">配置文件</span>
                <span class="info-value">config.yaml</span>
              </div>
              <a-divider style="margin: 12px 0" />
              <div class="info-item-row">
                <span class="info-label">数据目录</span>
                <a-button type="link" size="small" @click="openDataDir">
                  打开
                </a-button>
              </div>
            </div>
          </a-card>

          <!-- 快捷操作 -->
          <a-card class="info-card" :bordered="false">
            <template #title>
              <div class="card-title-wrapper">
                <thunderbolt-outlined class="title-icon" />
                <span>快捷操作</span>
              </div>
            </template>

            <a-space direction="vertical" style="width: 100%" :size="12">
              <a-button block size="large" @click="clearCache">
                <template #icon><clear-outlined /></template>
                清理缓存
              </a-button>
              <a-button block size="large" @click="exportConfig">
                <template #icon><export-outlined /></template>
                导出配置
              </a-button>
              <a-button block size="large" @click="importConfig">
                <template #icon><import-outlined /></template>
                导入配置
              </a-button>
              <a-button block size="large" @click="checkUpdate">
                <template #icon><sync-outlined /></template>
                检查更新
              </a-button>
            </a-space>
          </a-card>

          <!-- 帮助提示 -->
          <a-card class="info-card help-card" :bordered="false">
            <template #title>
              <div class="card-title-wrapper">
                <question-circle-outlined class="title-icon" />
                <span>配置说明</span>
              </div>
            </template>

            <div class="help-content">
              <p><strong>注意事项：</strong></p>
              <ul>
                <li>修改端口号后需要重启应用才能生效</li>
                <li>路径配置使用相对路径或绝对路径</li>
                <li>并发下载数过高可能影响性能</li>
                <li>配置会自动保存到 config.yaml</li>
              </ul>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { message, Modal } from 'ant-design-vue';
import {
  SettingOutlined,
  UndoOutlined,
  SaveOutlined,
  ControlOutlined,
  FolderOutlined,
  FolderOpenOutlined,
  CloudDownloadOutlined,
  GlobalOutlined,
  BgColorsOutlined,
  InfoCircleOutlined,
  ThunderboltOutlined,
  ClearOutlined,
  ExportOutlined,
  ImportOutlined,
  SyncOutlined,
  QuestionCircleOutlined,
  AppstoreOutlined,
  ApiOutlined,
  InboxOutlined
} from '@ant-design/icons-vue';

// 保存状态
const saving = ref(false);

// 基础配置
const config = reactive({
  appName: 'Go Site Clone',
  port: 6997,
  siteFileDir: 'www',
  packSiteFileDir: 'site-dist'
});

// 下载配置
const downloadConfig = reactive({
  concurrency: 5,
  timeout: 30,
  retryTimes: 3,
  userAgent: 'chrome',
  downloadImages: true,
  downloadCSS: true,
  downloadJS: true,
  downloadFonts: true
});

// Nginx 配置
const nginxConfig = reactive({
  defaultPort: 80,
  workerProcesses: 4,
  autoStart: false,
  gzip: true,
  accessLog: true
});

// 界面配置
const uiConfig = reactive({
  theme: 'light',
  language: 'zh-CN',
  showNotification: true,
  autoRefresh: true,
  compactMode: false
});

// 默认配置(用于重置)
const defaultConfig = {
  appName: 'Go Site Clone',
  port: 6997,
  siteFileDir: 'www',
  packSiteFileDir: 'site-dist'
};

// 初始化
onMounted(() => {
  loadConfig();
});

// 加载配置
const loadConfig = async () => {
  try {
    // TODO: 从后端加载配置
    message.success('配置加载成功');
  } catch (error) {
    console.error('加载配置失败:', error);
    message.error('加载配置失败');
  }
};

// 保存配置
const handleSave = async () => {
  saving.value = true;
  try {
    // TODO: 调用后端保存配置
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    const allConfig = {
      ...config,
      download: downloadConfig,
      nginx: nginxConfig,
      ui: uiConfig
    };
    
    console.log('保存配置:', allConfig);
    message.success('配置保存成功');
  } catch (error) {
    console.error('保存配置失败:', error);
    message.error('保存配置失败');
  } finally {
    saving.value = false;
  }
};

// 重置配置
const handleReset = () => {
  Modal.confirm({
    title: '确认重置',
    content: '确定要恢复所有默认配置吗？此操作不可撤销。',
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk: () => {
      Object.assign(config, defaultConfig);
      downloadConfig.concurrency = 5;
      downloadConfig.timeout = 30;
      downloadConfig.retryTimes = 3;
      downloadConfig.userAgent = 'chrome';
      downloadConfig.downloadImages = true;
      downloadConfig.downloadCSS = true;
      downloadConfig.downloadJS = true;
      downloadConfig.downloadFonts = true;
      
      nginxConfig.defaultPort = 80;
      nginxConfig.workerProcesses = 4;
      nginxConfig.autoStart = false;
      nginxConfig.gzip = true;
      nginxConfig.accessLog = true;
      
      uiConfig.theme = 'light';
      uiConfig.language = 'zh-CN';
      uiConfig.showNotification = true;
      uiConfig.autoRefresh = true;
      uiConfig.compactMode = false;
      
      message.success('已恢复默认配置');
    }
  });
};

// 选择目录
const selectDirectory = (field: string) => {
  // TODO: 调用系统文件选择器
  message.info('选择目录功能待实现');
};

// 打开数据目录
const openDataDir = () => {
  // TODO: 打开文件管理器
  message.info('打开数据目录功能待实现');
};

// 清理缓存
const clearCache = () => {
  Modal.confirm({
    title: '确认清理',
    content: '确定要清理所有缓存吗？',
    okText: '确定',
    cancelText: '取消',
    onOk: async () => {
      // TODO: 调用清理接口
      await new Promise(resolve => setTimeout(resolve, 500));
      message.success('缓存已清理');
    }
  });
};

// 导出配置
const exportConfig = () => {
  const allConfig = {
    ...config,
    download: downloadConfig,
    nginx: nginxConfig,
    ui: uiConfig
  };
  
  const blob = new Blob([JSON.stringify(allConfig, null, 2)], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'config-backup.json';
  a.click();
  URL.revokeObjectURL(url);
  
  message.success('配置已导出');
};

// 导入配置
const importConfig = () => {
  const input = document.createElement('input');
  input.type = 'file';
  input.accept = '.json';
  input.onchange = (e: any) => {
    const file = e.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (event: any) => {
        try {
          const imported = JSON.parse(event.target.result);
          Object.assign(config, imported);
          if (imported.download) Object.assign(downloadConfig, imported.download);
          if (imported.nginx) Object.assign(nginxConfig, imported.nginx);
          if (imported.ui) Object.assign(uiConfig, imported.ui);
          message.success('配置已导入');
        } catch (error) {
          message.error('配置文件格式错误');
        }
      };
      reader.readAsText(file);
    }
  };
  input.click();
};

// 检查更新
const checkUpdate = async () => {
  const hide = message.loading('检查更新中...', 0);
  try {
    await new Promise(resolve => setTimeout(resolve, 1500));
    hide();
    message.info('当前已是最新版本');
  } catch (error) {
    hide();
    message.error('检查更新失败');
  }
};
</script>

<style scoped>
/* 页面布局 */
.settings-page {
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

/* 卡片样式 */
.config-card,
.info-card {
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

/* 信息列表 */
.info-list {
  padding: 8px 0;
}

.info-item-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-size: 14px;
  color: #666;
}

.info-value {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

/* 帮助卡片 */
.help-card .help-content {
  color: #666;
  line-height: 1.8;
  font-size: 13px;
}

.help-content strong {
  color: #333;
  font-weight: 600;
}

.help-content ul {
  margin: 8px 0;
  padding-left: 20px;
}

.help-content li {
  margin: 6px 0;
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
}
</style>
