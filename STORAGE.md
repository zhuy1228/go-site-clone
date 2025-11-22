# BBolt 本地化存储方案实施文档

## 概述

本项目已成功集成 BBolt 嵌入式数据库，实现高效的本地化持久化存储。BBolt 是纯 Go 实现的键值数据库，单文件存储，无需安装额外依赖。

## 架构设计

### 数据库文件
- **位置**: `data/site-clone.db`
- **格式**: BBolt 数据库文件（二进制）
- **特点**: 单文件、可复制、可备份

### Bucket 结构

```
site-clone.db
├── sites/              # 站点配置数据
│   ├── example.com → {NginxSiteConfig + 时间戳}
│   └── test.com → {NginxSiteConfig + 时间戳}
├── download_records/   # 下载记录
│   ├── 1234567890 → {DownloadRecord}
│   └── 1234567891 → {DownloadRecord}
├── access_logs/        # 访问日志（预留）
└── settings/          # 系统设置（预留）
```

## 核心功能

### 1. 站点配置管理

#### 数据结构
```go
type NginxSiteConfig struct {
    ID      string   `json:"id"`
    Name    string   `json:"name"`
    Domains []string `json:"domains"`
    Port    int      `json:"port"`
    Path    string   `json:"path"`
    Index   string   `json:"index"`
    Enabled bool     `json:"enabled"`
}

type SiteRecord struct {
    NginxSiteConfig
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### 可用方法
- `AddSite(site)` - 添加站点配置
- `GetSite(name)` - 获取单个站点
- `GetAllSites()` - 获取所有站点
- `UpdateSite(site)` - 更新站点配置
- `DeleteSite(name)` - 删除站点配置
- `UpdateSiteStatus(name, enabled)` - 更新启用状态
- `GetSitesByStatus(enabled)` - 按状态筛选站点
- `SiteExists(name)` - 检查站点是否存在

### 2. 下载记录管理

#### 数据结构
```go
type DownloadRecord struct {
    ID          string    `json:"id"`
    URL         string    `json:"url"`
    SiteName    string    `json:"site_name"`
    Status      string    `json:"status"`      // success, failed, processing
    TotalFiles  int       `json:"total_files"`
    Downloaded  int       `json:"downloaded"`
    CSSCount    int       `json:"css_count"`
    ScriptCount int       `json:"script_count"`
    ImageCount  int       `json:"image_count"`
    VideoCount  int       `json:"video_count"`
    ErrorMsg    string    `json:"error_msg"`
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    Duration    int64     `json:"duration"`
}
```

#### 可用方法
- `AddDownloadRecord(record)` - 添加下载记录
- `GetDownloadRecord(id)` - 获取单条记录
- `GetAllDownloadRecords()` - 获取所有记录
- `UpdateDownloadRecord(record)` - 更新记录
- `DeleteDownloadRecord(id)` - 删除记录
- `GetDownloadRecordsBySite(name)` - 按站点筛选
- `GetDownloadRecordsByStatus(status)` - 按状态筛选
- `GetRecentDownloadRecords(limit)` - 获取最近 N 条
- `GetDownloadStats()` - 获取统计信息
- `ClearOldDownloadRecords(days)` - 清理旧记录

## 使用示例

### 前端调用示例

#### 1. 添加站点配置
```typescript
import { App } from "../../../bindings/go-site-clone";

// 添加站点
await App.AddNginxSite({
  name: "example.com",
  domains: ["example.com", "127.0.0.1"],
  port: 8080,
  path: "www/example.com",
  index: "index.html",
  enabled: true
});
```

#### 2. 获取站点列表
```typescript
// 获取所有站点
const sites = await App.GetAllNginxSites();
console.log("站点列表:", sites);
```

#### 3. 添加下载记录
```typescript
// 开始下载时创建记录
const record = {
  id: "",  // 自动生成
  url: "https://example.com",
  site_name: "example.com",
  status: "processing",
  total_files: 100,
  downloaded: 0,
  start_time: new Date()
};
await App.AddDownloadRecord(record);
```

#### 4. 获取下载统计
```typescript
// 获取统计信息
const stats = await App.GetDownloadStats();
console.log("下载统计:", stats);
// 输出: { total: 10, success: 8, failed: 1, processing: 1, total_files: 1234 }
```

### 后端使用示例

#### 在其他服务中使用存储
```go
// 在 app.go 中
func (a *App) SomeMethod() {
    // 使用站点存储
    sites, err := a.store.GetAllSites()
    if err != nil {
        log.Printf("获取站点失败: %v", err)
        return
    }
    
    // 使用下载记录
    records, err := a.store.GetRecentDownloadRecords(10)
    if err != nil {
        log.Printf("获取记录失败: %v", err)
        return
    }
}
```

## 数据同步机制

### 双重存储
系统同时维护两种存储：
1. **BBolt 数据库**: 用于快速查询和持久化
2. **Nginx 配置文件**: 用于 nginx 实际运行

### 同步策略
- **添加站点**: 先保存数据库 → 生成配置文件
- **删除站点**: 先删除配置文件 → 删除数据库记录
- **更新站点**: 同步更新数据库和配置文件
- **查询站点**: 优先从数据库读取 → 如果为空则从配置文件迁移

## 数据备份

### 备份方法
```typescript
// 前端调用备份
await App.BackupDatabase("backups/site-clone-backup.db");
```

### 手动备份
直接复制 `data/site-clone.db` 文件即可

### 恢复数据
将备份文件复制回 `data/site-clone.db` 即可

## 性能特点

### 优势
- ✅ **读取速度快**: 索引查询，O(log n) 复杂度
- ✅ **写入安全**: ACID 事务支持
- ✅ **并发读取**: 支持多个 goroutine 并发读
- ✅ **内存占用小**: 仅加载必要数据
- ✅ **文件紧凑**: 自动压缩和优化

### 适用场景
- ✅ 站点数量: < 10,000
- ✅ 下载记录: < 100,000
- ✅ 并发读取: 高
- ✅ 并发写入: 中等

## 注意事项

### 1. 数据库初始化
应用启动时自动初始化，无需手动操作

### 2. 错误处理
所有存储操作都会返回 error，建议：
- 数据库操作失败时记录日志
- 关键操作失败时提示用户
- 不要因存储失败而中断主流程

### 3. 数据迁移
首次启动时，系统会自动从 nginx 配置文件迁移数据到数据库

### 4. 并发安全
BBolt 支持：
- ✅ 多个并发读取
- ⚠️ 同时只能有一个写入（已内部处理）

## 扩展功能

### 未来可添加的功能
1. **访问日志存储**: 记录每个站点的访问日志
2. **系统设置**: 存储用户配置和偏好
3. **定时任务**: 记录定时下载任务
4. **缓存管理**: 存储页面缓存信息
5. **用户数据**: 多用户支持（如需要）

### 添加新 Bucket 示例
```go
// 在 storage/store.go 添加
var BucketCustom = []byte("custom_data")

// 在 NewStore 中初始化
buckets := [][]byte{BucketSites, BucketDownload, BucketLogs, BucketSettings, BucketCustom}
```

## 维护建议

### 定期维护
1. **清理旧记录**: 定期清理超过 30 天的下载记录
   ```go
   deleted, _ := store.ClearOldDownloadRecords(30)
   ```

2. **备份数据库**: 建议每周自动备份
   ```go
   store.Backup("backups/weekly-backup.db")
   ```

3. **监控大小**: 数据库文件超过 100MB 时考虑归档

### 故障恢复
如果数据库损坏：
1. 删除 `data/site-clone.db`
2. 重启应用
3. 系统会从 nginx 配置文件重建数据库

## 总结

BBolt 存储方案为项目提供了：
- 🎯 **高性能**: 快速的键值存储
- 🔒 **高可靠**: ACID 事务保证
- 📦 **易维护**: 单文件、易备份
- 🚀 **易扩展**: 简单的 Bucket 模型
- 💪 **零依赖**: 纯 Go 实现，无需外部数据库

完美满足本地化、便携式存储需求！
