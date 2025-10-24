<template>
  <div class="google-search-page">
    <!-- 背景装饰元素 -->
    <div class="background-elements">
      <div class="floating-circle circle-1"></div>
      <div class="floating-circle circle-2"></div>
      <div class="floating-circle circle-3"></div>
    </div>

    <!-- 主要内容区域 -->
    <main class="main-content">
      <!-- 搜索容器 - 带动画效果 -->
      <div 
        :class="['search-section', { 'collapsed': hasSearched }]"
      >
        <!-- Google Logo -->
        <div class="logo-section" v-if="!hasSearched">
          <div class="google-logo">
            <span class="logo-letter g">G</span>
            <span class="logo-letter o1">o</span>
            <span class="logo-letter o2">o</span>
            <span class="logo-letter g2">g</span>
            <span class="logo-letter l">l</span>
            <span class="logo-letter e">e</span>
          </div>
          <div class="country-indicator">中国</div>
        </div>

        <!-- 搜索输入区域 -->
        <div class="search-input-container">
          <div class="search-input-wrapper" :class="{ focused: isInputFocused }">
            <SearchOutlined class="search-icon" />
            <a-input
              ref="searchInputRef"
              v-model:value="searchKeyword"
              class="enhanced-search-input"
              placeholder="在 Google 上搜索，或者输入一个网址"
              size="large"
              @focus="onInputFocus"
              @blur="onInputBlur"
              @press-enter="handleSearch"
            />
            <div class="input-right-icons">
              <span class="voice-icon" @click="handleVoiceSearch">🎤</span>
              <span class="lens-icon" @click="handleImageSearch">📷</span>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="action-section" v-if="!hasSearched">
          <a-button 
            type="primary" 
            class="search-btn primary-btn"
            @click="handleSearch"
            :loading="searchLoading"
          >
            <template #icon><SearchOutlined /></template>
            Google 搜索
          </a-button>
          <a-button class="search-btn secondary-btn" @click="handleFeelingLucky">
            手气不错
          </a-button>
        </div>

        <!-- 快速提示 -->
        <div class="quick-tips" v-if="!hasSearched && !searchKeyword">
          <div class="tips-title">搜索历史</div>
          <div class="tips-grid">
            <div 
              v-for="tip in searchTips" 
              :key="tip"
              class="tip-item"
              @click="applySearchTip(tip)"
            >
              {{ tip }}
            </div>
          </div>
        </div>
      </div>

      <!-- 搜索结果区域 -->
      <transition name="results-fade">
        <div v-if="hasSearched" class="results-section">
          <!-- 结果统计和筛选 -->
          <div class="results-meta">
            <div class="meta-info">
              用时 {{ searchTime }} 秒 <a-button @click="downloadResource">下载</a-button>
            </div>
            <div class="filter-options">
              <a-radio-group v-model:value="filterType" size="small">
                <a-radio-button value="all">全部</a-radio-button>
                <a-radio-button value="images">图片</a-radio-button>
                <a-radio-button value="news">样式</a-radio-button>
                <a-radio-button value="videos">脚本</a-radio-button>
                <a-radio-button value="videos">媒体</a-radio-button>
              </a-radio-group>
            </div>
          </div>

          <!-- 搜索结果列表 -->
          <div class="results-content">
            <a-card class="results-card" :bordered="false">
              <a-list
                item-layout="vertical"
                :loading="searchLoading"
                :pagination="paginationConfig"
              >
                <a-list-item class="result-item">
                  <a-list-item-meta description="页面文件过多可能是当前站点包含博客、教程页面。">
                    <template #title>
                      <a class="result-title">页面文件</a>
                    </template>
                  </a-list-item-meta>
                  <div class="result-content">
                    <div class="result-snippet">可下载：{{ searchResults?.dom.length }}</div>
                  </div>
                </a-list-item>

                <a-list-item class="result-item">
                  <a-list-item-meta description="统计包含站点所有图片资源，包含站外、站内，默认只下载站内资源。">
                    <template #title>
                      <a class="result-title">图片文件</a>
                    </template>
                  </a-list-item-meta>
                  <div class="result-content">
                    <div class="result-snippet">总数：{{ searchResults?.image.length }}</div>
                    <div class="result-snippet">本站资源：{{ resourceData?.imageD.length }}</div>
                    <div class="result-snippet">其他资源：{{ resourceData?.image.length }}</div>
                  </div>
                </a-list-item>
                <a-list-item class="result-item">
                  <a-list-item-meta description="统计包含站点所有CSS资源，包含站外、站内，默认只下载站内资源。">
                    <template #title>
                      <a class="result-title">样式文件</a>
                    </template>
                  </a-list-item-meta>
                  <div class="result-content">
                    <div class="result-snippet">总数：{{ searchResults?.css.length }}</div>
                    <div class="result-snippet">本站资源：{{ resourceData?.cssD.length }}</div>
                    <div class="result-snippet">其他资源：{{ resourceData?.css.length }}</div>
                  </div>
                </a-list-item>
                <a-list-item class="result-item">
                  <a-list-item-meta  description="统计包含站点所有脚本资源，包含站外、站内，默认只下载站内资源。">
                    <template #title>
                      <a class="result-title">脚本文件</a>
                    </template>
                  </a-list-item-meta>
                  <div class="result-content">
                    <div class="result-snippet">总数：{{ searchResults?.script.length }}</div>
                    <div class="result-snippet">本站资源：{{ resourceData?.scriptD.length }}</div>
                    <div class="result-snippet">其他资源：{{ resourceData?.script.length }}</div>
                  </div>
                </a-list-item>
                <a-list-item class="result-item">
                  <a-list-item-meta description="统计包含站点所有媒体资源，包含站外、站内，默认只下载站内资源。">
                    <template #title>
                      <a class="result-title">媒体文件</a>
                    </template>
                  </a-list-item-meta>
                  <div class="result-content">
                    <div class="result-snippet">总数：{{ searchResults?.video.length }}</div>
                    <div class="result-snippet">本站资源：{{ resourceData?.videoD.length }}</div>
                    <div class="result-snippet">其他资源：{{ resourceData?.video.length }}</div>
                  </div>
                </a-list-item>
                <template #loadMore>
                  <div v-if="searchLoading" class="loading-more">
                    <a-spin />
                  </div>
                </template>
              </a-list>
            </a-card>

            <!-- 相关搜索 -->
            <div class="related-searches">
              <h3 class="related-title">相关搜索</h3>
              <div class="related-tags">
                <a-tag 
                  v-for="related in relatedSearches" 
                  :key="related"
                  class="related-tag"
                  @click="applyRelatedSearch(related)"
                >
                  {{ related }}
                </a-tag>
              </div>
            </div>
          </div>
        </div>
      </transition>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { SearchOutlined, SettingOutlined, MenuOutlined } from '@ant-design/icons-vue'
import {App} from "../../../bindings/go-site-clone";
const { $message } = useNuxtApp();
const [messageApi, contextHolder] = $message.useMessage();
// 搜索状态
const searchKeyword = ref('')
const hasSearched = ref(false)
const searchLoading = ref(false)
const searchTime = ref(0)
const searchInputRef = ref()
const isInputFocused = ref(false)
const filterType = ref('all')

// 分页配置
const paginationConfig = {
  pageSize: 10,
  showSizeChanger: false,
  showQuickJumper: true,
  showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条结果`
}

// 搜索提示和建议
const searchTips = [
  'Vue.js 最新特性',
  'Nuxt 4 新功能',
  'Ant Design Vue 使用指南',
  '前端开发最佳实践'
]

const relatedSearches = [
  'Vue 3 Composition API',
  'Nuxt 服务端渲染',
  'Ant Design 主题定制',
  'TypeScript 教程'
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
const downloadResource = async ()=> {
  App.DownloadSite(searchResults.value)
}

// 搜索处理函数
const handleSearch = async () => {
  if (searchLoading.value == true) {
    messageApi.info('当前已有网址在获取详情!');
    return
  }

  if (!searchKeyword.value.trim()) {
    return
  }
  if (!isValidURL(searchKeyword.value.trim())) {
    return
  }
  searchLoading.value = true
  const startTime = Date.now()
  
  nextTick(() => {
    setTimeout(() => {
      searchInputRef.value?.focus()
      App.GetResources(searchKeyword.value.trim()).then((res) => {
        if(res) {
          searchResults.value = res
          const urlData = new URL(searchKeyword.value.trim())
          for (let i = 0; i < res.image.length; i++) {
            const item = new URL(res.image[i]);
            if (urlData.hostname == item.hostname) {
              resourceData.value.imageD.push(item)
            } else {
              resourceData.value.image.push(item)
            }
          }
          for (let i = 0; i < res.css.length; i++) {
            const item = new URL(res.css[i]);
            if (urlData.hostname == item.hostname) {
              resourceData.value.cssD.push(item)
            } else {
              resourceData.value.css.push(item)
            }
          }
          for (let i = 0; i < res.script.length; i++) {
            const item = new URL(res.script[i]);
            if (urlData.hostname == item.hostname) {
              resourceData.value.scriptD.push(item)
            } else {
              resourceData.value.script.push(item)
            }
          }
          for (let i = 0; i < res.video.length; i++) {
            const item = new URL(res.video[i]);
            if (urlData.hostname == item.hostname) {
              resourceData.value.videoD.push(item)
            } else {
              resourceData.value.video.push(item)
            }
          }
        }
        searchTime.value = (Date.now() - startTime) / 1000
        hasSearched.value = true
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

// 辅助功能
const handleFeelingLucky = () => {
  searchKeyword.value = 'Ant Design Vue 在 Nuxt 4 中的集成指南'
  handleSearch()
}

const handleVoiceSearch = () => {
  console.log('语音搜索功能')
}

const handleImageSearch = () => {
  console.log('图片搜索功能')
}

const applySearchTip = (tip) => {
  searchKeyword.value = tip
  handleSearch()
}

const applyRelatedSearch = (query) => {
  searchKeyword.value = query
  handleSearch()
}

// 页面加载后自动聚焦搜索框
onMounted(() => {
  nextTick(() => {
    searchInputRef.value?.focus()
  })
})
</script>

<style scoped>
.google-search-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow-x: hidden;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.background-elements {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
}

.floating-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;
}

.circle-1 {
  width: 200px;
  height: 200px;
  top: 10%;
  right: 10%;
  animation-delay: 0s;
}

.circle-2 {
  width: 150px;
  height: 150px;
  top: 60%;
  left: 5%;
  animation-delay: 2s;
}

.circle-3 {
  width: 100px;
  height: 100px;
  bottom: 20%;
  right: 20%;
  animation-delay: 4s;
}

.main-content {
  position: relative;
  z-index: 1;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* 搜索区域样式 */
.search-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: center;
}

.search-section.collapsed {
  justify-content: flex-start;
  align-items: flex-start;
  padding: 1rem 2rem;
  animation: slideUp 0.6s ease-out;
}

/* Logo 样式 */
.logo-section {
  text-align: center;
  margin-bottom: 2rem;
}

.google-logo {
  font-size: 5rem;
  font-weight: 400;
  margin-bottom: 0.5rem;
  display: flex;
  justify-content: center;
  gap: 0.5rem;
}

.logo-letter {
  display: inline-block;
  animation: bounce 0.5s ease;
}

.logo-letter.g { color: #4285f4; animation-delay: 0.1s; }
.logo-letter.o1 { color: #ea4335; animation-delay: 0.2s; }
.logo-letter.o2 { color: #fbbc05; animation-delay: 0.3s; }
.logo-letter.g2 { color: #4285f4; animation-delay: 0.4s; }
.logo-letter.l { color: #34a853; animation-delay: 0.5s; }
.logo-letter.e { color: #ea4335; animation-delay: 0.6s; }

.small-logo {
  font-size: 1.5rem;
  font-weight: 500;
  display: flex;
  gap: 2px;
}

.logo-small-g { color: #4285f4; }
.logo-small-o1 { color: #ea4335; }
.logo-small-o2 { color: #fbbc05; }
.logo-small-g2 { color: #4285f4; }
.logo-small-l { color: #34a853; }
.logo-small-e { color: #ea4335; }

.country-indicator {
  color: rgba(255, 255, 255, 0.8);
  font-size: 1.2rem;
}

/* 搜索输入框样式 */
.search-input-container {
  width: 100%;
  max-width: 584px;
  margin-bottom: 2rem;
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  background: white;
  border-radius: 24px;
  padding: 0.5rem 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.search-input-wrapper.focused {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
  transform: translateY(-2px);
  border-color: #4285f4;
}

.search-icon {
  color: #9aa0a6;
  margin-right: 0.5rem;
  font-size: 1.2rem;
}

.enhanced-search-input {
  border: none;
  box-shadow: none !important;
  background: transparent;
  flex: 1;
}

:deep(.enhanced-search-input .ant-input) {
  background: transparent;
  border: none;
  font-size: 1.1rem;
  padding: 8px 4px;
}

.input-right-icons {
  display: flex;
  gap: 0.5rem;
  margin-left: 0.5rem;
}

.voice-icon, .lens-icon {
  cursor: pointer;
  font-size: 1.2rem;
  transition: transform 0.2s ease;
  padding: 4px;
}

.voice-icon:hover, .lens-icon:hover {
  transform: scale(1.1);
}

/* 操作按钮样式 */
.action-section {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.search-btn {
  border-radius: 4px;
  padding: 0.5rem 1.5rem;
  font-weight: 500;
  transition: all 0.3s ease;
  border: none;
}

.primary-btn {
  background: #1a73e8;
  border-color: #1a73e8;
}

.primary-btn:hover {
  background: #1669d9;
  border-color: #1669d9;
  transform: translateY(-1px);
}

.secondary-btn {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.3);
  color: white;
}

.secondary-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.4);
  transform: translateY(-1px);
}

/* 快速提示样式 */
.quick-tips {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 1.5rem;
  max-width: 600px;
  width: 100%;
}

.tips-title {
  color: white;
  font-size: 1.1rem;
  margin-bottom: 1rem;
  font-weight: 500;
}

.tips-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 0.5rem;
}

.tip-item {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  padding: 0.75rem 1rem;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: center;
}

.tip-item:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-1px);
}

/* 搜索结果区域样式 */
.results-fade-enter-active {
  transition: opacity 0.5s ease, transform 0.5s ease;
}

.results-fade-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.results-section {
  background: white;
  border-radius: 12px 12px 0 0;
  margin-top: 1rem;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.1);
  flex: 1;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  border-bottom: 1px solid #e8eaed;
  background: white;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.search-box-mini {
  width: 400px;
}

:deep(.mini-search-input .ant-input) {
  border-radius: 4px;
  border: 1px solid #dfe1e5;
  height: 36px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.header-btn, .sign-in-btn {
  border-radius: 4px;
}

/* 结果统计 */
.results-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background: #f8f9fa;
  border-bottom: 1px solid #e8eaed;
}

.meta-info {
  color: #70757a;
  font-size: 0.9rem;
}

/* 结果项样式 */
.results-content {
  padding: 0 2rem 2rem;
}

.result-item {
  padding: 1.5rem 0;
  border-bottom: 1px solid #f0f0f0;
  position: relative;
}

.result-number {
  position: absolute;
  left: -2rem;
  top: 1.5rem;
  color: #70757a;
  font-size: 0.9rem;
  font-weight: 500;
}

.result-title {
  font-size: 1.3rem;
  color: #1a0dab;
  text-decoration: none;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: color 0.2s ease;
}

.result-title:hover {
  text-decoration: underline;
  color: #1a0dab;
}

.external-icon {
  font-size: 0.8rem;
  opacity: 0.7;
}

.result-avatar {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.result-url {
  color: #006621;
  font-size: 0.9rem;
}

.result-snippet {
  color: #4d5156;
  line-height: 1.6;
  margin: 0.5rem 0;
}

.result-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-top: 0.5rem;
}

.publish-date {
  color: #70757a;
  font-size: 0.9rem;
}

.feature-tag {
  font-size: 0.8rem;
}

/* 相关搜索 */
.related-searches {
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #f0f0f0;
}

.related-title {
  font-size: 1.1rem;
  color: #333;
  margin-bottom: 1rem;
}

.related-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.related-tag {
  cursor: pointer;
  transition: all 0.2s ease;
}

.related-tag:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

/* 页脚样式 */
.page-footer {
  background: #f2f2f2;
  border-top: 1px solid #e4e4e4;
  padding: 1rem 0;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-location {
  color: #70757a;
  font-size: 0.9rem;
}

.footer-links {
  display: flex;
  gap: 1.5rem;
}

.footer-link {
  color: #70757a;
  text-decoration: none;
  font-size: 0.9rem;
  transition: color 0.2s ease;
}

.footer-link:hover {
  color: #1a73e8;
}

/* 动画定义 */
@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(5deg); }
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

@keyframes slideUp {
  0% { 
    transform: translateY(0) scale(1);
    opacity: 1;
  }
  50% { 
    transform: translateY(-20px) scale(0.95);
    opacity: 0.8;
  }
  100% { 
    transform: translateY(0) scale(1);
    opacity: 1;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .google-logo {
    font-size: 3rem;
    flex-wrap: wrap;
  }
  
  .search-input-container {
    max-width: 100%;
  }
  
  .search-box-mini {
    width: 200px;
  }
  
  .tips-grid {
    grid-template-columns: 1fr;
  }
  
  .results-header {
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
  }
  
  .header-left, .header-right {
    width: 100%;
    justify-content: center;
  }
  
  .results-meta {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }
  
  .footer-content {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
}

/* 加载状态 */
.loading-more {
  display: flex;
  justify-content: center;
  padding: 2rem;
}
</style>