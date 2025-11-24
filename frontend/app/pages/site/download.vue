<template>
  <div class="download-page">
    <!-- 面包屑导航 -->
    <a-breadcrumb class="page-breadcrumb">
      <a-breadcrumb-item>
        <home-outlined />
        <span>仿站</span>
      </a-breadcrumb-item>
      <a-breadcrumb-item>整站下载</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="gradient-blob blob-1"></div>
      <div class="gradient-blob blob-2"></div>
      <div class="gradient-blob blob-3"></div>
    </div>

    <!-- 主要内容区域 -->
    <main class="main-container">
      <!-- 搜索区域 -->
      <div :class="['search-area', { 'minimized': hasSearched }]">
        <!-- Logo 和标题 -->
        <div class="header-section" v-if="!hasSearched">
          <div class="brand-logo">
            <cloud-download-outlined class="logo-icon" />
            <h1 class="brand-title">整站下载工具</h1>
          </div>
          <p class="brand-subtitle">快速抓取网站资源，一键下载整站内容</p>
        </div>

        <!-- 搜索框 -->
        <div class="search-box-wrapper">
          <div class="search-box" :class="{ 'focused': isInputFocused, 'compact': hasSearched }">
            <global-outlined class="search-prefix-icon" />
            <a-input
              ref="searchInputRef"
              v-model:value="searchKeyword"
              class="search-input"
              placeholder="请输入要下载的网站地址，例如：https://example.com"
              size="large"
              @focus="onInputFocus"
              @blur="onInputBlur"
              @press-enter="handleSearch"
            />
            <a-button 
              type="primary" 
              class="search-action-btn"
              :loading="searchLoading"
              @click="handleSearch"
            >
              <template #icon><search-outlined /></template>
              {{ hasSearched ? '重新分析' : '开始分析' }}
            </a-button>
          </div>
        </div>

        <!-- 快捷提示 -->
        <div class="quick-hints" v-if="!hasSearched && !searchKeyword">
          <div class="hints-header">
            <history-outlined />
            <span>最近使用</span>
          </div>
          <div class="hints-list">
            <a-tag 
              v-for="tip in searchTips" 
              :key="tip"
              class="hint-tag"
              @click="applySearchTip(tip)"
            >
              {{ tip }}
            </a-tag>
          </div>
        </div>
      </div>

      <!-- 资源统计结果 -->
      <transition name="fade-slide">
        <div v-if="hasSearched" class="results-area">
          <!-- 结果头部 -->
          <div class="results-header">
            <div class="result-info">
              <check-circle-outlined class="success-icon" />
              <span class="info-text">分析完成，用时 {{ searchTime }} 秒</span>
            </div>
            <a-button 
              type="primary" 
              size="large"
              :loading="isDownload"
              @click="downloadResource"
              class="download-main-btn"
            >
              <template #icon><download-outlined /></template>
              {{ isDownload ? '下载中...' : '开始下载' }}
            </a-button>
          </div>

          <!-- 下载选项配置 -->
          <div class="download-options-card">
            <a-card :bordered="false" size="small">
              <template #title>
                <span style="font-size: 14px;">
                  <setting-outlined style="margin-right: 8px;" />
                  下载选项
                </span>
              </template>
              <a-space direction="vertical" style="width: 100%;" :size="16">
                <!-- 下载模式选择 -->
                <div class="option-row">
                  <div class="option-label">
                    <span class="label-text">下载模式</span>
                    <a-tooltip placement="top">
                      <template #title>
                        <div style="max-width: 300px;">
                          <p><b>仅本站资源：</b>只下载同域名的资源，外部CDN资源保持原链接。文件小，下载快，需要网络查看。</p>
                          <p style="margin-top: 8px;"><b>包含外部资源：</b>下载所有资源包括CDN。文件大，完全离线可用。</p>
                        </div>
                      </template>
                      <question-circle-outlined style="margin-left: 4px; color: #8c8c8c; cursor: help;" />
                    </a-tooltip>
                  </div>
                  <a-radio-group v-model:value="downloadOptions.mode" button-style="solid">
                    <a-radio-button value="same-domain">
                      <cloud-outlined />
                      仅本站资源
                    </a-radio-button>
                    <a-radio-button value="all-resources">
                      <global-outlined />
                      包含外部资源
                    </a-radio-button>
                  </a-radio-group>
                </div>

                <!-- 外部资源类型选择 -->
                <div v-if="downloadOptions.mode === 'all-resources'" class="option-row">
                  <div class="option-label">
                    <span class="label-text">外部资源类型</span>
                  </div>
                  <a-space wrap>
                    <a-checkbox v-model:checked="downloadOptions.downloadExternalCSS">
                      <file-text-outlined /> CSS样式
                    </a-checkbox>
                    <a-checkbox v-model:checked="downloadOptions.downloadExternalJS">
                      <code-outlined /> JavaScript
                    </a-checkbox>
                    <a-checkbox v-model:checked="downloadOptions.downloadExternalImages">
                      <file-image-outlined /> 图片
                    </a-checkbox>
                    <a-checkbox v-model:checked="downloadOptions.downloadExternalVideos">
                      <video-camera-outlined /> 视频
                    </a-checkbox>
                  </a-space>
                </div>

                <!-- 文件大小限制 -->
                <div class="option-row">
                  <div class="option-label">
                    <span class="label-text">跳过超大文件</span>
                    <a-tooltip title="跳过超过指定大小的文件，避免下载时间过长">
                      <question-circle-outlined style="margin-left: 4px; color: #8c8c8c; cursor: help;" />
                    </a-tooltip>
                  </div>
                  <a-space>
                    <a-switch v-model:checked="downloadOptions.skipLargeFiles" />
                    <span v-if="downloadOptions.skipLargeFiles">
                      最大 
                      <a-input-number 
                        v-model:value="downloadOptions.maxFileSize" 
                        :min="1" 
                        :max="100" 
                        size="small"
                        style="width: 80px;"
                      /> MB
                    </span>
                  </a-space>
                </div>

                <!-- 提示信息 -->
                <a-alert 
                  v-if="downloadOptions.mode === 'all-resources'"
                  message="注意事项"
                  type="warning"
                  show-icon
                >
                  <template #description>
                    <ul style="margin: 4px 0; padding-left: 20px; font-size: 12px;">
                      <li>下载外部资源会增加文件体积和下载时间</li>
                      <li>某些CDN资源可能有防盗链或访问限制</li>
                      <li>建议先使用"仅本站资源"模式测试</li>
                    </ul>
                  </template>
                </a-alert>
              </a-space>
            </a-card>
          </div>

          <!-- 资源筛选 -->
          <div class="resource-filter">
            <a-radio-group v-model:value="filterType" size="large" button-style="solid">
              <a-radio-button value="all">
                <appstore-outlined />
                全部资源
              </a-radio-button>
              <a-radio-button value="images">
                <file-image-outlined />
                图片
              </a-radio-button>
              <a-radio-button value="styles">
                <file-text-outlined />
                样式
              </a-radio-button>
              <a-radio-button value="scripts">
                <code-outlined />
                脚本
              </a-radio-button>
              <a-radio-button value="media">
                <video-camera-outlined />
                媒体
              </a-radio-button>
            </a-radio-group>
          </div>

          <!-- 资源卡片列表 -->
          <div class="resources-grid">
            <!-- 页面文件卡片 -->
            <a-card class="resource-card" :bordered="false" hoverable>
              <div class="card-header">
                <div class="header-left">
                  <file-outlined class="card-icon pages-icon" />
                  <div class="header-info">
                    <h3 class="card-title">页面文件</h3>
                    <p class="card-desc">HTML 页面及文档</p>
                  </div>
                </div>
                <a-badge :count="searchResults?.dom.length || 0" :number-style="{ backgroundColor: '#52c41a' }" />
              </div>
              <a-divider style="margin: 12px 0" />
              <div class="card-content">
                <div class="stat-row">
                  <span class="stat-label">可下载数量</span>
                  <span class="stat-value">{{ searchResults?.dom.length || 0 }}</span>
                </div>
                <a-progress 
                  v-if="isDownload && searchResults?.dom.length > 0" 
                  :percent="domDownloadProgress" 
                  :stroke-color="{ '0%': '#108ee9', '100%': '#87d068' }"
                  :show-info="true"
                />
              </div>
            </a-card>

            <!-- 图片文件卡片 -->
            <a-card class="resource-card" :bordered="false" hoverable>
              <div class="card-header">
                <div class="header-left">
                  <file-image-outlined class="card-icon images-icon" />
                  <div class="header-info">
                    <h3 class="card-title">图片文件</h3>
                    <p class="card-desc">JPG、PNG、SVG 等</p>
                  </div>
                </div>
                <a-badge :count="searchResults?.image.length || 0" :number-style="{ backgroundColor: '#1890ff' }" />
              </div>
              <a-divider style="margin: 12px 0" />
              <div class="card-content">
                <div class="stat-row">
                  <span class="stat-label">资源总数</span>
                  <span class="stat-value">{{ searchResults?.image.length || 0 }}</span>
                </div>
                <div class="stat-row primary">
                  <span class="stat-label">本站资源</span>
                  <span class="stat-value primary">{{ resourceData?.imageD.length || 0 }}</span>
                </div>
                <div class="stat-row secondary">
                  <span class="stat-label">外站资源</span>
                  <span class="stat-value">{{ resourceData?.image.length || 0 }}</span>
                </div>
                <a-progress 
                  v-if="isDownload && resourceData?.imageD.length > 0" 
                  :percent="imageDownloadProgress"
                  :stroke-color="{ '0%': '#108ee9', '100%': '#87d068' }"
                />
              </div>
            </a-card>

            <!-- 样式文件卡片 -->
            <a-card class="resource-card" :bordered="false" hoverable>
              <div class="card-header">
                <div class="header-left">
                  <file-text-outlined class="card-icon styles-icon" />
                  <div class="header-info">
                    <h3 class="card-title">样式文件</h3>
                    <p class="card-desc">CSS 样式表</p>
                  </div>
                </div>
                <a-badge :count="searchResults?.css.length || 0" :number-style="{ backgroundColor: '#722ed1' }" />
              </div>
              <a-divider style="margin: 12px 0" />
              <div class="card-content">
                <div class="stat-row">
                  <span class="stat-label">资源总数</span>
                  <span class="stat-value">{{ searchResults?.css.length || 0 }}</span>
                </div>
                <div class="stat-row primary">
                  <span class="stat-label">本站资源</span>
                  <span class="stat-value primary">{{ resourceData?.cssD.length || 0 }}</span>
                </div>
                <div class="stat-row secondary">
                  <span class="stat-label">外站资源</span>
                  <span class="stat-value">{{ resourceData?.css.length || 0 }}</span>
                </div>
                <a-progress 
                  v-if="isDownload && resourceData?.cssD.length > 0" 
                  :percent="cssDownloadProgress"
                  :stroke-color="{ '0%': '#722ed1', '100%': '#b37feb' }"
                />
              </div>
            </a-card>

            <!-- 脚本文件卡片 -->
            <a-card class="resource-card" :bordered="false" hoverable>
              <div class="card-header">
                <div class="header-left">
                  <code-outlined class="card-icon scripts-icon" />
                  <div class="header-info">
                    <h3 class="card-title">脚本文件</h3>
                    <p class="card-desc">JavaScript 脚本</p>
                  </div>
                </div>
                <a-badge :count="searchResults?.script.length || 0" :number-style="{ backgroundColor: '#fa8c16' }" />
              </div>
              <a-divider style="margin: 12px 0" />
              <div class="card-content">
                <div class="stat-row">
                  <span class="stat-label">资源总数</span>
                  <span class="stat-value">{{ searchResults?.script.length || 0 }}</span>
                </div>
                <div class="stat-row primary">
                  <span class="stat-label">本站资源</span>
                  <span class="stat-value primary">{{ resourceData?.scriptD.length || 0 }}</span>
                </div>
                <div class="stat-row secondary">
                  <span class="stat-label">外站资源</span>
                  <span class="stat-value">{{ resourceData?.script.length || 0 }}</span>
                </div>
                <a-progress 
                  v-if="isDownload && resourceData?.scriptD.length > 0" 
                  :percent="scriptDownloadProgress"
                  :stroke-color="{ '0%': '#fa8c16', '100%': '#ffc53d' }"
                />
              </div>
            </a-card>

            <!-- 媒体文件卡片 -->
            <a-card class="resource-card" :bordered="false" hoverable>
              <div class="card-header">
                <div class="header-left">
                  <video-camera-outlined class="card-icon media-icon" />
                  <div class="header-info">
                    <h3 class="card-title">媒体文件</h3>
                    <p class="card-desc">视频、音频等</p>
                  </div>
                </div>
                <a-badge :count="searchResults?.video.length || 0" :number-style="{ backgroundColor: '#eb2f96' }" />
              </div>
              <a-divider style="margin: 12px 0" />
              <div class="card-content">
                <div class="stat-row">
                  <span class="stat-label">资源总数</span>
                  <span class="stat-value">{{ searchResults?.video.length || 0 }}</span>
                </div>
                <div class="stat-row primary">
                  <span class="stat-label">本站资源</span>
                  <span class="stat-value primary">{{ resourceData?.videoD.length || 0 }}</span>
                </div>
                <div class="stat-row secondary">
                  <span class="stat-label">外站资源</span>
                  <span class="stat-value">{{ resourceData?.video.length || 0 }}</span>
                </div>
                <a-progress 
                  v-if="isDownload && resourceData?.videoD.length > 0" 
                  :percent="videoDownloadProgress"
                  :stroke-color="{ '0%': '#eb2f96', '100%': '#f759ab' }"
                />
              </div>
            </a-card>
          </div>

          <!-- 使用提示 -->
          <a-alert
            message="下载说明"
            description="默认只下载本站资源，外站资源需要额外配置。下载的文件将保存在配置的目录中。"
            type="info"
            show-icon
            closable
            class="download-tips"
          />
        </div>
      </transition>
    </main>
  </div>
</template>

<script setup>
import { 
  SearchOutlined, 
  CloudDownloadOutlined, 
  GlobalOutlined,
  HistoryOutlined,
  CheckCircleOutlined,
  DownloadOutlined,
  AppstoreOutlined,
  FileImageOutlined,
  FileTextOutlined,
  CodeOutlined,
  VideoCameraOutlined,
  FileOutlined,
  HomeOutlined
} from '@ant-design/icons-vue'
import { App } from "../../../bindings/go-site-clone"
const { $message } = useNuxtApp()
const [messageApi, contextHolder] = $message.useMessage()
import { Events } from '@wailsio/runtime'
// 搜索状态
const searchKeyword = ref('')
const hasSearched = ref(false)
const searchLoading = ref(false)
const searchTime = ref(0)
const searchInputRef = ref()
const isInputFocused = ref(false)
const filterType = ref('all')
// 下载进度
const domDownloadProgress = ref(0);
const cssDownloadProgress = ref(0);
const scriptDownloadProgress = ref(0);
const imageDownloadProgress = ref(0);
const videoDownloadProgress = ref(0);
const isDownload = ref(false)

// 下载选项配置
const downloadOptions = ref({
  mode: 'same-domain', // 'same-domain' | 'all-resources' | 'custom'
  customDomains: [],
  skipLargeFiles: true,
  maxFileSize: 10, // MB
  downloadExternalCSS: true,
  downloadExternalJS: true,
  downloadExternalImages: true,
  downloadExternalVideos: false, // 视频通常较大，默认不下载
})

// 搜索提示
const searchTips = [
  'https://www.example.com',
  'https://www.github.com',
  'https://www.npmjs.com'
]

// 模拟搜索结果
const searchResults = ref({
  script: [],
  css: [],
  image: [],
  dom: [],
  video: [],
})

const resourceData = ref({
  script: [],
  scriptD: [],
  css: [],
  cssD: [],
  image: [],
  imageD: [],
  video: [],
  videoD: []
})

const isValidURL = (str) => { try { new URL(str); return true; } catch { return false; }}

// 下载当前页面资源
const downloadResource = async () => {
  if(!isDownload.value) {
    isDownload.value = true
    // 重置所有进度
    domDownloadProgress.value = 0
    cssDownloadProgress.value = 0
    scriptDownloadProgress.value = 0
    imageDownloadProgress.value = 0
    videoDownloadProgress.value = 0
    
    try {
      const mode = downloadOptions.value.mode
      const optionsInfo = mode === 'same-domain' 
        ? '仅下载本站资源' 
        : '下载所有资源（包含外部CDN）'
      
      messageApi.info(`开始下载网站资源...（${optionsInfo}）`)
      
      await App.DownloadSiteWithOptions(searchKeyword.value.trim(), searchResults.value, downloadOptions.value)
      
      messageApi.success('网站资源下载完成！')
      isDownload.value = false
    } catch (error) {
      console.error('下载失败:', error)
      messageApi.error('下载失败，请重试')
      isDownload.value = false
    }

  } else {
    messageApi.warning('正在下载中，请稍候...')
  }
}

// 搜索处理函数
const handleSearch = async () => {
  if (searchLoading.value == true) {
    messageApi.info('当前已有网址在获取详情!');
    return
  }

  if (!searchKeyword.value.trim()) {
    messageApi.warning('请输入网址')
    return
  }
  if (!isValidURL(searchKeyword.value.trim())) {
    messageApi.error('请输入有效的网址，例如：https://example.com')
    return
  }
  searchLoading.value = true
  const startTime = Date.now()
  isDownload.value = false
  
  // 清空之前的资源数据
  resourceData.value = {
    script: [],
    scriptD: [],
    css: [],
    cssD: [],
    image: [],
    imageD: [],
    video: [],
    videoD: []
  }
  
  nextTick(() => {
    setTimeout(() => {
      searchInputRef.value?.focus()
      App.GetResources(searchKeyword.value.trim()).then((res) => {
        if(res) {
          searchResults.value = res
          const urlData = new URL(searchKeyword.value.trim())
          
          // 分类图片资源
          res.image?.forEach(url => {
            try {
              const item = new URL(url)
              if (urlData.hostname === item.hostname) {
                resourceData.value.imageD.push(url)
              } else {
                resourceData.value.image.push(url)
              }
            } catch (e) {
              console.warn('Invalid image URL:', url)
            }
          })
          
          // 分类CSS资源
          res.css?.forEach(url => {
            try {
              const item = new URL(url)
              if (urlData.hostname === item.hostname) {
                resourceData.value.cssD.push(url)
              } else {
                resourceData.value.css.push(url)
              }
            } catch (e) {
              console.warn('Invalid CSS URL:', url)
            }
          })
          
          // 分类脚本资源
          res.script?.forEach(url => {
            try {
              const item = new URL(url)
              if (urlData.hostname === item.hostname) {
                resourceData.value.scriptD.push(url)
              } else {
                resourceData.value.script.push(url)
              }
            } catch (e) {
              console.warn('Invalid script URL:', url)
            }
          })
          
          // 分类视频资源
          res.video?.forEach(url => {
            try {
              const item = new URL(url)
              if (urlData.hostname === item.hostname) {
                resourceData.value.videoD.push(url)
              } else {
                resourceData.value.video.push(url)
              }
            } catch (e) {
              console.warn('Invalid video URL:', url)
            }
          })
        }
        searchTime.value = ((Date.now() - startTime) / 1000).toFixed(2)
        hasSearched.value = true
        searchLoading.value = false
      }).catch((error) => {
        console.error('获取资源失败:', error)
        messageApi.error('获取网站资源失败，请检查网址是否正确')
        searchLoading.value = false
      })
    }, 400)
  })
} 


// 输入框焦点事件
const onInputFocus = () => {
  isInputFocused.value = true
}

const onInputBlur = () => {
  isInputFocused.value = false
}

// 应用搜索提示
const applySearchTip = (tip) => {
  searchKeyword.value = tip
  handleSearch()
}

// 页面加载后自动聚焦搜索框
onMounted(() => {
  nextTick(() => {
    searchInputRef.value?.focus()
  })
  Events.On("download:css", (event)=> {
    if(resourceData.value.cssD.length > 0 && event.data && event.data[0] !== undefined) {
      const progress = Math.min(Math.ceil(((event.data[0] + 1) / resourceData.value.cssD.length) * 100), 100)
      cssDownloadProgress.value = progress
    }
  })
  
  Events.On("download:script", (event)=> {
    if(resourceData.value.scriptD.length > 0 && event.data && event.data[0] !== undefined) {
      const progress = Math.min(Math.ceil(((event.data[0] + 1) / resourceData.value.scriptD.length) * 100), 100)
      scriptDownloadProgress.value = progress
    }
  })
  
  Events.On("download:image", (event)=> {
    if(resourceData.value.imageD.length > 0 && event.data && event.data[0] !== undefined) {
      const progress = Math.min(Math.ceil(((event.data[0] + 1) / resourceData.value.imageD.length) * 100), 100)
      imageDownloadProgress.value = progress
    }
  })
  
  Events.On("download:video", (event)=> {
    if(resourceData.value.videoD.length > 0 && event.data && event.data[0] !== undefined) {
      const progress = Math.min(Math.ceil(((event.data[0] + 1) / resourceData.value.videoD.length) * 100), 100)
      videoDownloadProgress.value = progress
    }
  })
  
  Events.On("download:dom", (event)=> {
    if(searchResults.value.dom.length > 0 && event.data && event.data[0] !== undefined) {
      const progress = Math.min(Math.ceil(((event.data[0] + 1) / searchResults.value.dom.length) * 100), 100)
      domDownloadProgress.value = progress
    }
  })
})
</script>

<style scoped>
/* 页面布局 */
.download-page {
  min-height: calc(100vh - 120px);
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  position: relative;
  overflow: hidden;
}

.page-breadcrumb {
  padding: 16px 24px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
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
  filter: blur(80px);
  opacity: 0.3;
  animation: float 20s ease-in-out infinite;
}

.blob-1 {
  width: 500px;
  height: 500px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  top: -10%;
  right: -10%;
  animation-delay: 0s;
}

.blob-2 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  bottom: -10%;
  left: -10%;
  animation-delay: 7s;
}

.blob-3 {
  width: 350px;
  height: 350px;
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  top: 40%;
  left: 50%;
  animation-delay: 14s;
}

/* 主容器 */
.main-container {
  position: relative;
  z-index: 1;
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 24px;
}

/* 搜索区域 */
.search-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin-bottom: 40px;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.search-area.minimized {
  margin-bottom: 30px;
}

/* 品牌区域 */
.header-section {
  text-align: center;
  margin-bottom: 40px;
  animation: fadeInDown 0.8s ease-out;
}

.brand-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-bottom: 16px;
}

.logo-icon {
  font-size: 56px;
  color: #1890ff;
  filter: drop-shadow(0 4px 12px rgba(24, 144, 255, 0.3));
}

.brand-title {
  font-size: 42px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
}

.brand-subtitle {
  font-size: 16px;
  color: #666;
  margin: 0;
}

/* 搜索框 */
.search-box-wrapper {
  width: 100%;
  max-width: 800px;
  animation: fadeInUp 0.8s ease-out 0.2s both;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 12px;
  background: white;
  border-radius: 50px;
  padding: 8px 8px 8px 24px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 2px solid transparent;
}

.search-box.focused {
  box-shadow: 0 8px 32px rgba(24, 144, 255, 0.2);
  border-color: #1890ff;
  transform: translateY(-2px);
}

.search-box.compact {
  border-radius: 24px;
  padding: 6px 6px 6px 20px;
}

.search-prefix-icon {
  font-size: 20px;
  color: #999;
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  border: none !important;
  box-shadow: none !important;
  background: transparent !important;
}

:deep(.search-input .ant-input) {
  font-size: 16px;
  background: transparent;
  border: none;
  box-shadow: none !important;
  padding: 0;
}

:deep(.search-input .ant-input:focus) {
  box-shadow: none !important;
}

.search-action-btn {
  border-radius: 40px;
  height: 44px;
  padding: 0 32px;
  font-weight: 500;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
  transition: all 0.3s ease;
}

.search-action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.4);
}

/* 快捷提示 */
.quick-hints {
  margin-top: 32px;
  animation: fadeIn 0.8s ease-out 0.4s both;
}

.hints-header {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #666;
  font-size: 14px;
  margin-bottom: 12px;
}

.hints-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
}

.hint-tag {
  cursor: pointer;
  padding: 8px 16px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(24, 144, 255, 0.2);
  transition: all 0.3s ease;
  font-size: 13px;
}

.hint-tag:hover {
  background: #1890ff;
  color: white;
  border-color: #1890ff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
}

/* 结果区域 */
.results-area {
  animation: fadeInUp 0.6s ease-out;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  padding: 24px;
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.06);
  margin-bottom: 24px;
}

.result-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.success-icon {
  font-size: 24px;
  color: #52c41a;
}

.info-text {
  font-size: 16px;
  color: #333;
  font-weight: 500;
}

.download-main-btn {
  height: 44px;
  padding: 0 32px;
  border-radius: 22px;
  font-size: 16px;
  font-weight: 500;
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.3);
  transition: all 0.3s ease;
}

.download-main-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(24, 144, 255, 0.4);
}

/* 资源筛选 */
.resource-filter {
  background: white;
  padding: 20px;
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.06);
  margin-bottom: 24px;
  display: flex;
  justify-content: center;
}

:deep(.resource-filter .ant-radio-group) {
  display: flex;
  gap: 8px;
}

:deep(.resource-filter .ant-radio-button-wrapper) {
  border-radius: 12px !important;
  border: none;
  padding: 0 20px;
  height: 40px;
  line-height: 38px;
  font-weight: 500;
  transition: all 0.3s ease;
}

:deep(.resource-filter .ant-radio-button-wrapper:not(:first-child)::before) {
  display: none;
}

:deep(.resource-filter .ant-radio-button-wrapper-checked) {
  background: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
}

/* 资源卡片网格 */
.resources-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.resource-card {
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.06);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.resource-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

:deep(.resource-card .ant-card-body) {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.card-icon {
  font-size: 32px;
  padding: 12px;
  border-radius: 12px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e9ecef 100%);
}

.pages-icon {
  color: #52c41a;
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
}

.images-icon {
  color: #1890ff;
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
}

.styles-icon {
  color: #722ed1;
  background: linear-gradient(135deg, #f3e5f5 0%, #e1bee7 100%);
}

.scripts-icon {
  color: #fa8c16;
  background: linear-gradient(135deg, #fff3e0 0%, #ffe0b2 100%);
}

.media-icon {
  color: #eb2f96;
  background: linear-gradient(135deg, #fce4ec 0%, #f8bbd0 100%);
}

.header-info {
  flex: 1;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0 0 4px 0;
}

.card-desc {
  font-size: 13px;
  color: #999;
  margin: 0;
}

.card-content {
  margin-top: 8px;
}

.stat-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
}

.stat-row:last-of-type {
  border-bottom: none;
}

.stat-row.primary {
  background: linear-gradient(90deg, rgba(24, 144, 255, 0.05) 0%, transparent 100%);
  padding: 10px 12px;
  margin: 0 -12px;
  border-radius: 8px;
}

.stat-row.secondary {
  opacity: 0.7;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.stat-value.primary {
  color: #1890ff;
}

:deep(.card-content .ant-progress) {
  margin-top: 12px;
}

/* 下载提示 */
.download-tips {
  border-radius: 12px;
  border: none;
  background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
}

:deep(.download-tips .ant-alert-icon) {
  font-size: 18px;
}

/* 下载选项卡片 */
.download-options-card {
  margin-bottom: 24px;
}

.download-options-card :deep(.ant-card) {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  background: linear-gradient(135deg, #f0f9ff 0%, #ffffff 100%);
}

.download-options-card :deep(.ant-card-head) {
  border-bottom: 1px solid #e8f4ff;
  background: linear-gradient(90deg, #e6f7ff 0%, transparent 100%);
}

.download-options-card :deep(.ant-card-body) {
  padding: 20px;
}

.option-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 0;
}

.option-label {
  display: flex;
  align-items: center;
  min-width: 120px;
}

.label-text {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.option-row :deep(.ant-radio-group) {
  display: flex;
  gap: 8px;
}

.option-row :deep(.ant-radio-button-wrapper) {
  height: 36px;
  line-height: 34px;
  padding: 0 20px;
  border-radius: 6px !important;
  border: 1px solid #d9d9d9;
  display: flex;
  align-items: center;
  gap: 6px;
}

.option-row :deep(.ant-radio-button-wrapper-checked) {
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
  border-color: #1890ff;
  color: #fff;
}

.option-row :deep(.ant-radio-button-wrapper:not(.ant-radio-button-wrapper-checked):hover) {
  color: #1890ff;
  border-color: #1890ff;
}

.option-row :deep(.ant-checkbox-wrapper) {
  margin: 0;
  padding: 6px 12px;
  border-radius: 6px;
  transition: all 0.3s;
}

.option-row :deep(.ant-checkbox-wrapper:hover) {
  background: #f0f9ff;
}

.option-row :deep(.ant-checkbox-wrapper .anticon) {
  margin-right: 4px;
}

/* 动画 */
@keyframes float {
  0%, 100% { 
    transform: translate(0, 0) rotate(0deg); 
  }
  33% { 
    transform: translate(30px, -30px) rotate(120deg); 
  }
  66% { 
    transform: translate(-20px, 20px) rotate(240deg); 
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

.fade-slide-enter-active {
  animation: fadeInUp 0.6s ease-out;
}

.fade-slide-leave-active {
  animation: fadeInUp 0.4s ease-out reverse;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .main-container {
    padding: 20px 16px;
  }

  .brand-title {
    font-size: 32px;
  }

  .brand-subtitle {
    font-size: 14px;
  }

  .search-box {
    flex-direction: column;
    border-radius: 20px;
    padding: 16px;
    gap: 12px;
  }

  .search-input {
    width: 100%;
  }

  .search-action-btn {
    width: 100%;
  }

  .results-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .download-main-btn {
    width: 100%;
  }

  .resources-grid {
    grid-template-columns: 1fr;
  }

  :deep(.resource-filter .ant-radio-group) {
    flex-wrap: wrap;
  }
}

@media (max-width: 480px) {
  .page-breadcrumb {
    padding: 12px 16px;
    font-size: 12px;
  }

  .brand-logo {
    flex-direction: column;
    gap: 8px;
  }

  .logo-icon {
    font-size: 40px;
  }

  .brand-title {
    font-size: 24px;
  }

  .hints-list {
    flex-direction: column;
  }

  .hint-tag {
    text-align: center;
  }
}
</style>