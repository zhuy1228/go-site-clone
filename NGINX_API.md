# Nginx 服务 API 说明

## 功能说明

1. **应用启动时自动检测 nginx 状态**
   - 启动应用时会自动调用 `OnStartup()` 检测 nginx 是否在运行
   - 如果检测到 nginx 进程，会自动设置 `Running = true`
   - 不会自动启动或停止 nginx

2. **应用关闭时保持 nginx 运行**
   - 关闭应用时调用 `OnShutdown()`，但不会停止 nginx
   - nginx 进程继续在后台运行
   - 下次启动应用时会自动检测到运行状态

## 前端可调用方法

### 进程管理

#### 启动 Nginx
```typescript
import { App } from '../bindings/go-site-clone'

try {
  await App.StartNginx()
  console.log('Nginx 启动成功')
} catch (error) {
  console.error('Nginx 启动失败:', error)
}
```

#### 停止 Nginx
```typescript
try {
  await App.StopNginx()
  console.log('Nginx 停止成功')
} catch (error) {
  console.error('Nginx 停止失败:', error)
}
```

#### 重启 Nginx
```typescript
try {
  await App.RestartNginx()
  console.log('Nginx 重启成功')
} catch (error) {
  console.error('Nginx 重启失败:', error)
}
```

#### 重载配置
```typescript
try {
  await App.ReloadNginx()
  console.log('Nginx 配置重载成功')
} catch (error) {
  console.error('Nginx 配置重载失败:', error)
}
```

#### 检查状态
```typescript
try {
  const running = await App.CheckNginxStatus()
  console.log('Nginx 运行状态:', running)
} catch (error) {
  console.error('检查状态失败:', error)
}
```

#### 测试配置
```typescript
try {
  await App.TestNginxConfig()
  console.log('配置文件有效')
} catch (error) {
  console.error('配置文件测试失败:', error)
}
```

### 站点管理

#### 添加站点
```typescript
const siteConfig = {
  name: 'example.com',
  domains: ['example.com', '127.0.0.1'],
  port: 80,
  path: '/www/example.com',
  index: 'index.html',
  enabled: true
}

try {
  await App.AddNginxSite(siteConfig)
  console.log('站点添加成功')
} catch (error) {
  console.error('站点添加失败:', error)
}
```

#### 删除站点
```typescript
try {
  await App.DeleteNginxSite('example.com')
  console.log('站点删除成功')
} catch (error) {
  console.error('站点删除失败:', error)
}
```

#### 更新站点
```typescript
const updatedConfig = {
  name: 'example.com',
  domains: ['example.com', 'www.example.com'],
  port: 8080,
  path: '/www/example.com',
  index: 'index.html',
  enabled: true
}

try {
  await App.UpdateNginxSite(updatedConfig)
  console.log('站点更新成功')
} catch (error) {
  console.error('站点更新失败:', error)
}
```

#### 获取所有站点
```typescript
try {
  const sites = await App.GetAllNginxSites()
  console.log('站点列表:', sites)
} catch (error) {
  console.error('获取站点列表失败:', error)
}
```

#### 启用站点
```typescript
try {
  await App.EnableNginxSite('example.com')
  console.log('站点启用成功')
} catch (error) {
  console.error('站点启用失败:', error)
}
```

#### 禁用站点
```typescript
try {
  await App.DisableNginxSite('example.com')
  console.log('站点禁用成功')
} catch (error) {
  console.error('站点禁用失败:', error)
}
```

### 日志管理

#### 获取访问日志（最后100行）
```typescript
try {
  const logs = await App.GetNginxAccessLog(100)
  console.log('访问日志:', logs)
} catch (error) {
  console.error('获取访问日志失败:', error)
}
```

#### 获取错误日志（最后100行）
```typescript
try {
  const logs = await App.GetNginxErrorLog(100)
  console.log('错误日志:', logs)
} catch (error) {
  console.error('获取错误日志失败:', error)
}
```

#### 清空日志
```typescript
try {
  await App.ClearNginxLogs()
  console.log('日志清空成功')
} catch (error) {
  console.error('日志清空失败:', error)
}
```

## 数据结构

### NginxSiteConfig
```typescript
interface NginxSiteConfig {
  ID: string       // 站点唯一ID（可选）
  name: string     // 站点名称
  domains: string[] // 域名列表
  port: number     // 监听端口
  path: string     // 网站根目录路径
  index: string    // 默认首页文件
  enabled: boolean // 是否启用
}
```

## 工作流程示例

### 页面初始化时检查 nginx 状态
```typescript
import { ref, onMounted } from 'vue'
import { App } from '../bindings/go-site-clone'

const nginxRunning = ref(false)

onMounted(async () => {
  try {
    nginxRunning.value = await App.CheckNginxStatus()
  } catch (error) {
    console.error('检查 nginx 状态失败:', error)
  }
})
```

### 完整的站点托管流程
```typescript
// 1. 检查 nginx 状态
const running = await App.CheckNginxStatus()

// 2. 如果未运行，启动 nginx
if (!running) {
  await App.StartNginx()
}

// 3. 添加站点配置
await App.AddNginxSite({
  name: 'mysite.com',
  domains: ['mysite.com', '127.0.0.1'],
  port: 80,
  path: '/www/mysite.com',
  index: 'index.html',
  enabled: true
})

// 4. 重载配置使其生效
await App.ReloadNginx()

// 5. 验证站点已添加
const sites = await App.GetAllNginxSites()
console.log('当前托管站点:', sites)
```

## 注意事项

1. **Windows 路径处理**：路径会自动转换为 nginx 兼容格式（斜杠）
2. **配置文件位置**：站点配置保存在 `plugin/nginx/conf/hosts/` 目录
3. **禁用站点**：禁用的站点配置会移动到 `plugin/nginx/conf/hosts.disabled/` 目录
4. **自动重载**：添加、删除、更新站点时，如果 nginx 正在运行会自动重载配置
5. **持久化运行**：应用关闭后 nginx 继续运行，重新打开应用会自动检测状态
