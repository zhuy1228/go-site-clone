# 下载配置功能实现文档

## 📋 功能概述

实现了灵活的下载配置系统，允许用户选择下载本站资源、所有资源或自定义资源。

---

## 🎯 配置选项

### 1. **下载模式 (Mode)**

```typescript
type DownloadMode = 'same-domain' | 'all-resources' | 'custom'
```

#### same-domain（默认）
- ✅ 只下载同域名资源
- ✅ 外部CDN保持原链接
- ✅ 文件小，下载快
- ⚠️ 需要网络才能完整浏览

#### all-resources
- ✅ 下载所有资源（包括外部CDN）
- ✅ 完全离线可用
- ⚠️ 文件大，下载慢
- ⚠️ 可能违反CDN使用条款

#### custom
- ✅ 自定义下载规则
- ✅ 指定特定域名下载
- ✅ 平衡体积和可用性

---

## 🔧 API 使用

### Go 后端

```go
// 获取默认配置
options := types.DefaultDownloadOptions()

// 自定义配置
options := types.DownloadOptions{
    Mode: types.DownloadModeCustom,
    CustomDomains: []string{
        "cdn.example.com",
        "static.example.com",
    },
    SkipLargeFiles: true,
    MaxFileSize: 10, // MB
    DownloadExternalCSS: true,
    DownloadExternalJS: true,
    DownloadExternalImages: false,
    DownloadExternalVideos: false,
}

// 使用配置下载
App.DownloadSiteWithOptions(uri, resources, options)
```

### TypeScript 前端

```typescript
interface DownloadOptions {
  mode: 'same-domain' | 'all-resources' | 'custom'
  customDomains: string[]
  skipLargeFiles: boolean
  maxFileSize: number
  downloadExternalCSS: boolean
  downloadExternalJS: boolean
  downloadExternalImages: boolean
  downloadExternalVideos: boolean
}

// 调用后端API
const options: DownloadOptions = {
  mode: 'custom',
  customDomains: ['cdn.jsdelivr.net'],
  skipLargeFiles: true,
  maxFileSize: 10,
  downloadExternalCSS: true,
  downloadExternalJS: true,
  downloadExternalImages: false,
  downloadExternalVideos: false,
}

await App.DownloadSiteWithOptions(url, resources, options)
```

---

## 🎨 前端UI示例

### Vue组件代码

```vue
<template>
  <div class="download-config">
    <!-- 下载模式选择 -->
    <a-card title="下载设置" class="config-card">
      <a-form-item label="下载模式">
        <a-radio-group v-model:value="downloadConfig.mode">
          <a-radio value="same-domain">
            <div class="radio-content">
              <div class="radio-title">仅本站资源</div>
              <div class="radio-desc">只下载同域名资源，体积小速度快</div>
            </div>
          </a-radio>
          <a-radio value="all-resources">
            <div class="radio-content">
              <div class="radio-title">全部资源</div>
              <div class="radio-desc">下载所有资源包括CDN，完全离线可用</div>
            </div>
          </a-radio>
          <a-radio value="custom">
            <div class="radio-content">
              <div class="radio-title">自定义</div>
              <div class="radio-desc">指定特定域名进行下载</div>
            </div>
          </a-radio>
        </a-radio-group>
      </a-form-item>

      <!-- 自定义域名（当mode=custom时显示） -->
      <a-form-item 
        v-if="downloadConfig.mode === 'custom'" 
        label="自定义域名"
      >
        <a-select
          v-model:value="downloadConfig.customDomains"
          mode="tags"
          placeholder="输入域名后按回车，如：cdn.example.com"
          :options="[]"
        />
        <div class="hint-text">
          <info-circle-outlined />
          只下载指定域名的资源，一行一个
        </div>
      </a-form-item>

      <!-- 资源类型过滤 -->
      <a-form-item 
        v-if="downloadConfig.mode !== 'same-domain'" 
        label="外部资源类型"
      >
        <a-checkbox-group v-model:value="selectedResourceTypes">
          <a-checkbox value="css">
            <file-text-outlined /> CSS样式表
          </a-checkbox>
          <a-checkbox value="js">
            <code-outlined /> JavaScript脚本
          </a-checkbox>
          <a-checkbox value="images">
            <picture-outlined /> 图片资源
          </a-checkbox>
          <a-checkbox value="videos">
            <video-camera-outlined /> 视频资源
          </a-checkbox>
        </a-checkbox-group>
      </a-form-item>

      <!-- 文件大小限制 -->
      <a-form-item label="文件大小限制">
        <a-switch 
          v-model:checked="downloadConfig.skipLargeFiles" 
          checked-children="启用" 
          un-checked-children="禁用"
        />
        <a-input-number
          v-if="downloadConfig.skipLargeFiles"
          v-model:value="downloadConfig.maxFileSize"
          :min="1"
          :max="100"
          addon-after="MB"
          class="size-input"
        />
        <div class="hint-text">
          <info-circle-outlined />
          跳过超过指定大小的文件，节省磁盘空间
        </div>
      </a-form-item>

      <!-- 预估信息 -->
      <a-alert 
        v-if="estimatedSize" 
        type="info" 
        class="estimate-alert"
      >
        <template #message>
          <div class="estimate-info">
            <database-outlined />
            预估下载大小: {{ estimatedSize }}
          </div>
          <div class="estimate-info">
            <clock-circle-outlined />
            预估下载时间: {{ estimatedTime }}
          </div>
        </template>
      </a-alert>
    </a-card>

    <!-- 操作按钮 -->
    <div class="action-buttons">
      <a-button @click="resetConfig">重置配置</a-button>
      <a-button type="primary" @click="startDownload">
        <download-outlined />
        开始下载
      </a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import { App } from '@/bindings/go-site-clone'

interface DownloadConfig {
  mode: 'same-domain' | 'all-resources' | 'custom'
  customDomains: string[]
  skipLargeFiles: boolean
  maxFileSize: number
  downloadExternalCSS: boolean
  downloadExternalJS: boolean
  downloadExternalImages: boolean
  downloadExternalVideos: boolean
}

const downloadConfig = ref<DownloadConfig>({
  mode: 'same-domain',
  customDomains: [],
  skipLargeFiles: true,
  maxFileSize: 10,
  downloadExternalCSS: false,
  downloadExternalJS: false,
  downloadExternalImages: false,
  downloadExternalVideos: false,
})

const selectedResourceTypes = ref<string[]>([])

// 监听资源类型选择变化
watch(selectedResourceTypes, (types) => {
  downloadConfig.value.downloadExternalCSS = types.includes('css')
  downloadConfig.value.downloadExternalJS = types.includes('js')
  downloadConfig.value.downloadExternalImages = types.includes('images')
  downloadConfig.value.downloadExternalVideos = types.includes('videos')
})

// 预估下载大小
const estimatedSize = computed(() => {
  // 根据资源数量和配置估算
  return '约 50 MB'
})

// 预估下载时间
const estimatedTime = computed(() => {
  // 根据网速和文件数量估算
  return '约 2-5 分钟'
})

// 重置配置
const resetConfig = () => {
  downloadConfig.value = {
    mode: 'same-domain',
    customDomains: [],
    skipLargeFiles: true,
    maxFileSize: 10,
    downloadExternalCSS: false,
    downloadExternalJS: false,
    downloadExternalImages: false,
    downloadExternalVideos: false,
  }
  selectedResourceTypes.value = []
  message.success('配置已重置')
}

// 开始下载
const startDownload = async () => {
  try {
    message.loading('准备下载...', 0)
    await App.DownloadSiteWithOptions(
      currentUrl.value,
      resources.value,
      downloadConfig.value
    )
    message.success('下载完成！')
  } catch (error) {
    message.error('下载失败：' + error)
  }
}
</script>

<style scoped>
.download-config {
  max-width: 800px;
  margin: 0 auto;
}

.config-card {
  margin-bottom: 20px;
}

.radio-content {
  padding: 8px 0;
}

.radio-title {
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 4px;
}

.radio-desc {
  font-size: 12px;
  color: #888;
}

.hint-text {
  margin-top: 8px;
  font-size: 12px;
  color: #666;
  display: flex;
  align-items: center;
  gap: 4px;
}

.size-input {
  margin-left: 12px;
  width: 120px;
}

.estimate-alert {
  margin-top: 16px;
}

.estimate-info {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 4px 0;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
```

---

## 📊 配置场景示例

### 场景1：快速备份（默认）
```json
{
  "mode": "same-domain",
  "skipLargeFiles": true,
  "maxFileSize": 10
}
```
**适用**：快速备份网站，保持文件结构  
**优点**：快速、文件小  
**缺点**：需要网络才能查看外部资源

---

### 场景2：完全离线
```json
{
  "mode": "all-resources",
  "downloadExternalCSS": true,
  "downloadExternalJS": true,
  "downloadExternalImages": true,
  "downloadExternalVideos": true,
  "skipLargeFiles": true,
  "maxFileSize": 50
}
```
**适用**：需要完全离线浏览  
**优点**：不依赖网络  
**缺点**：文件大、下载慢

---

### 场景3：只下载常用CDN
```json
{
  "mode": "custom",
  "customDomains": [
    "cdn.jsdelivr.net",
    "unpkg.com",
    "cdnjs.cloudflare.com"
  ],
  "downloadExternalCSS": true,
  "downloadExternalJS": true,
  "downloadExternalImages": false,
  "downloadExternalVideos": false,
  "skipLargeFiles": true,
  "maxFileSize": 10
}
```
**适用**：平衡体积和可用性  
**优点**：关键资源离线可用  
**缺点**：部分资源仍需网络

---

## ⚠️ 注意事项

### 1. **法律合规**
- 下载CDN资源可能违反服务条款
- 仅用于个人学习和备份
- 不得用于商业用途

### 2. **性能影响**
- 下载所有资源会显著增加时间
- 建议设置合理的文件大小限制
- 大型网站可能需要数小时

### 3. **存储空间**
- 全资源模式可能需要数GB空间
- 定期清理不需要的下载
- 建议保留至少50%可用空间

### 4. **网络礼仪**
- 避免过于频繁的请求
- 尊重网站的robots.txt
- 不要对同一网站重复下载

---

## 🚀 后续优化

1. **智能预估**
   - 根据资源列表预估下载大小和时间
   - 提供取消和暂停功能

2. **下载队列**
   - 支持批量下载多个网站
   - 队列管理和优先级

3. **增量下载**
   - 检测已下载的文件
   - 只下载新增或修改的资源

4. **压缩存储**
   - 下载后自动压缩
   - 节省磁盘空间

---

## 📝 总结

下载配置功能提供了灵活的资源管理：

✅ **三种模式** - 本站/全部/自定义  
✅ **类型过滤** - 选择性下载资源类型  
✅ **大小限制** - 跳过超大文件  
✅ **域名白名单** - 精确控制下载范围  
✅ **用户友好** - 直观的配置界面  

现在用户可以根据不同需求选择最合适的下载策略！
