# BBolt 存储快速开始

## 🚀 已完成的工作

✅ 安装 BBolt 依赖 (`go.etcd.io/bbolt@latest`)
✅ 创建存储服务基础架构 (`storage/store.go`)
✅ 实现站点数据 CRUD (`storage/site_store.go`)
✅ 实现下载记录管理 (`storage/download_store.go`)
✅ 集成到应用主服务 (`app.go`)
✅ 更新类型定义 (`types/site_config.go`)
✅ 修复前端类型绑定 (`webpage.vue`)

## 📁 新增文件

```
go-site-clone/
├── storage/                      # 存储层（新增）
│   ├── store.go                  # BBolt 核心封装
│   ├── site_store.go             # 站点配置存储
│   └── download_store.go         # 下载记录存储
├── types/
│   └── site_config.go            # 站点配置类型定义（新增）
├── data/                         # 数据目录（运行时自动创建）
│   └── site-clone.db             # BBolt 数据库文件
└── STORAGE.md                    # 存储方案文档（新增）
```

## 🎯 核心功能

### 1. 自动初始化
应用启动时自动创建数据库和必要的 buckets：
- `sites` - 站点配置
- `download_records` - 下载记录
- `access_logs` - 访问日志（预留）
- `settings` - 系统设置（预留）

### 2. 双重存储
- **数据库**: BBolt 存储（快速查询）
- **配置文件**: Nginx conf 文件（运行时使用）
- 自动同步两者状态

### 3. 数据持久化
- 站点配置自动保存
- 下载记录自动记录
- 支持数据查询和统计

## 🔧 使用方法

### 前端调用示例

```typescript
import { App } from "../../../bindings/go-site-clone";

// 1. 添加站点
await App.AddNginxSite({
  ID: "",
  Name: "example.com",
  Domains: ["example.com", "127.0.0.1"],
  Port: 8080,
  Path: "www/example.com",
  Index: "index.html",
  Enabled: true
});

// 2. 获取所有站点
const sites = await App.GetAllNginxSites();

// 3. 删除站点
await App.DeleteNginxSite("example.com");

// 4. 添加下载记录
await App.AddDownloadRecord({
  id: "",
  url: "https://example.com",
  site_name: "example.com",
  status: "processing",
  total_files: 100,
  downloaded: 0
});

// 5. 获取下载统计
const stats = await App.GetDownloadStats();
```

## 📊 可用 API

### 站点管理
- `AddNginxSite(site)` - 添加站点
- `GetAllNginxSites()` - 获取所有站点
- `UpdateNginxSite(site)` - 更新站点
- `DeleteNginxSite(name)` - 删除站点
- `EnableNginxSite(name)` - 启用站点
- `DisableNginxSite(name)` - 禁用站点

### 下载记录
- `AddDownloadRecord(record)` - 添加记录
- `GetAllDownloadRecords()` - 获取所有记录
- `GetRecentDownloadRecords(limit)` - 获取最近记录
- `GetDownloadStats()` - 获取统计信息
- `DeleteDownloadRecord(id)` - 删除记录
- `ClearOldDownloadRecords(days)` - 清理旧记录

### 数据备份
- `BackupDatabase(path)` - 备份数据库

## 🎨 数据结构

### 站点配置
```go
type NginxSiteConfig struct {
    ID      string   // 唯一标识
    Name    string   // 站点名称
    Domains []string // 域名列表
    Port    int      // 端口号
    Path    string   // 文件路径
    Index   string   // 默认首页
    Enabled bool     // 是否启用
}
```

### 下载记录
```go
type DownloadRecord struct {
    ID          string    // 记录ID
    URL         string    // 下载URL
    SiteName    string    // 站点名称
    Status      string    // 状态: success/failed/processing
    TotalFiles  int       // 总文件数
    Downloaded  int       // 已下载数
    CSSCount    int       // CSS 数量
    ScriptCount int       // JS 数量
    ImageCount  int       // 图片数量
    VideoCount  int       // 视频数量
    ErrorMsg    string    // 错误信息
    StartTime   time.Time // 开始时间
    EndTime     time.Time // 结束时间
    Duration    int64     // 耗时（秒）
}
```

## 💡 最佳实践

### 1. 错误处理
```typescript
try {
  await App.AddNginxSite(siteConfig);
  message.success('添加成功');
} catch (error) {
  message.error('添加失败: ' + error.message);
}
```

### 2. 数据同步
添加、更新、删除站点时，系统会自动：
- 更新数据库
- 生成/更新/删除 nginx 配置文件
- 如果 nginx 运行中，自动重载配置

### 3. 定期备份
```typescript
// 建议定期备份数据库
const backupPath = `backups/backup-${Date.now()}.db`;
await App.BackupDatabase(backupPath);
```

### 4. 清理旧数据
```typescript
// 清理 30 天前的下载记录
const deleted = await App.ClearOldDownloadRecords(30);
console.log(`清理了 ${deleted} 条记录`);
```

## 📝 注意事项

1. **数据库位置**: `data/site-clone.db`（自动创建）
2. **备份方式**: 直接复制 db 文件或使用 API
3. **并发安全**: BBolt 内部处理，无需担心
4. **性能**: 适合 < 10,000 站点，< 100,000 记录
5. **迁移**: 首次运行自动从 nginx 配置迁移

## 🔍 调试

### 查看数据库内容
可以使用 bbolt 命令行工具：
```bash
go install go.etcd.io/bbolt/cmd/bbolt@latest
bbolt dump data/site-clone.db
```

### 数据库统计
```bash
bbolt stats data/site-clone.db
```

## 📚 更多信息

详细文档请查看 [STORAGE.md](./STORAGE.md)

## ✅ 验证安装

运行应用后检查：
1. `data/` 目录是否创建
2. `data/site-clone.db` 文件是否存在
3. 添加站点后，数据库中是否有记录
4. 控制台是否显示 "数据库初始化成功"

完成！🎉
