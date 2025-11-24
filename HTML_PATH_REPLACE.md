# HTML 资源路径替换功能说明

## 📋 功能概述

在下载网站时，自动将 HTML 文件中的所有资源绝对路径转换为相对路径，使得离线浏览时资源能够正确加载。

---

## ✨ 支持的资源类型

### 1. **CSS 样式表**
```html
<!-- 替换前 -->
<link href="https://example.com/css/style.css" rel="stylesheet">

<!-- 替换后 -->
<link href="./css/style.css" rel="stylesheet">
```

### 2. **JavaScript 脚本**
```html
<!-- 替换前 -->
<script src="https://example.com/js/app.js"></script>

<!-- 替换后 -->
<script src="./js/app.js"></script>
```

### 3. **图片资源**
```html
<!-- 替换前 -->
<img src="https://example.com/images/logo.png" alt="Logo">

<!-- 替换后 -->
<img src="./images/logo.png" alt="Logo">
```

### 4. **视频资源**
```html
<!-- 替换前 -->
<video src="https://example.com/videos/demo.mp4"></video>
<source src="https://example.com/videos/demo.webm">

<!-- 替换后 -->
<video src="./videos/demo.mp4"></video>
<source src="./videos/demo.webm">
```

### 5. **音频资源**
```html
<!-- 替换前 -->
<audio src="https://example.com/audio/music.mp3"></audio>

<!-- 替换后 -->
<audio src="./audio/music.mp3"></audio>
```

### 6. **CSS 中的 URL**
```css
/* 替换前 */
background-image: url('https://example.com/images/bg.jpg');
background: url("https://example.com/images/pattern.png");
background: url(https://example.com/images/texture.jpg);

/* 替换后 */
background-image: url('./images/bg.jpg');
background: url("./images/pattern.png");
background: url(./images/texture.jpg);
```

---

## 🔧 技术实现

### 核心函数

#### 1. `replaceHTMLResourcePaths()`
主函数，负责替换 HTML 中的所有资源路径。

```go
func replaceHTMLResourcePaths(htmlContent string, baseURL string) string {
    // 解析基础 URL
    // 使用正则表达式匹配各种资源标签
    // 调用 convertToRelativePath 转换每个 URL
    // 返回修改后的 HTML
}
```

#### 2. `convertToRelativePath()`
将绝对 URL 转换为相对路径。

**判断逻辑**：
- ✅ 只替换同域名的资源
- ❌ 跳过外部域名资源
- ❌ 跳过 data:、javascript:、mailto: 等协议
- ❌ 跳过已经是相对路径的资源

```go
func convertToRelativePath(resourceURL, baseURL, baseHost, baseScheme, basePath string) string {
    // 1. 验证 URL 类型
    // 2. 检查是否同域名
    // 3. 计算相对路径
    // 4. 返回转换结果
}
```

#### 3. `calculateRelativePath()`
计算两个路径之间的相对路径。

**示例**：
```go
basePath := "/pages/about/"
targetPath := "/css/style.css"

// 结果：../../css/style.css
```

---

## 📊 路径转换示例

### 示例 1：同级目录
```
当前页面：https://example.com/index.html
资源路径：https://example.com/style.css

转换结果：./style.css
```

### 示例 2：子目录
```
当前页面：https://example.com/index.html
资源路径：https://example.com/css/main.css

转换结果：./css/main.css
```

### 示例 3：上级目录
```
当前页面：https://example.com/pages/about.html
资源路径：https://example.com/css/style.css

转换结果：../css/style.css
```

### 示例 4：深层嵌套
```
当前页面：https://example.com/blog/2024/11/post.html
资源路径：https://example.com/images/logo.png

转换结果：../../../images/logo.png
```

### 示例 5：外部资源（不替换）
```
当前页面：https://example.com/index.html
资源路径：https://cdn.example.com/jquery.js

转换结果：https://cdn.example.com/jquery.js（保持不变）
```

---

## 🚫 不会被替换的资源

### 1. **Data URL**
```html
<img src="data:image/png;base64,iVBORw0KG...">
<!-- 保持不变 -->
```

### 2. **JavaScript 伪协议**
```html
<a href="javascript:void(0)">Click</a>
<!-- 保持不变 -->
```

### 3. **锚点链接**
```html
<a href="#section1">Go to Section 1</a>
<!-- 保持不变 -->
```

### 4. **外部域名资源**
```html
<script src="https://cdn.jsdelivr.net/npm/vue@3"></script>
<!-- 保持不变（不同域名）-->
```

### 5. **已经是相对路径**
```html
<img src="./images/photo.jpg">
<link href="../css/style.css">
<!-- 保持不变 -->
```

---

## 🎯 使用场景

### 1. **离线浏览**
下载整站后，可以在没有网络的情况下正常浏览，所有资源都能正确加载。

### 2. **网站备份**
保存网站的完整副本，包括所有页面和资源，路径关系正确。

### 3. **网站迁移**
将网站从一个域名迁移到另一个域名，资源路径自动适配。

### 4. **本地开发**
在本地环境测试网站，无需配置虚拟主机。

---

## 🔍 处理流程

```
1. 下载 HTML 文件
   ↓
2. 读取 HTML 内容
   ↓
3. 解析基础 URL（当前页面地址）
   ↓
4. 使用正则表达式匹配资源标签
   ├─ <link href="...">
   ├─ <script src="...">
   ├─ <img src="...">
   ├─ <video src="...">
   ├─ <audio src="...">
   └─ url(...)
   ↓
5. 对每个匹配项：
   ├─ 提取 URL
   ├─ 判断是否需要替换
   ├─ 计算相对路径
   └─ 替换原 URL
   ↓
6. 保存修改后的 HTML 文件
```

---

## 📝 正则表达式说明

### CSS 链接
```go
regexp: `(<link[^>]*?href=[\"'])([^\"']+)([\"'][^>]*?>)`
匹配: <link ... href="URL" ...>
```

### JavaScript 脚本
```go
regexp: `(<script[^>]*?src=[\"'])([^\"']+)([\"'][^>]*?>)`
匹配: <script ... src="URL" ...>
```

### 图片
```go
regexp: `(<img[^>]*?src=[\"'])([^\"']+)([\"'][^>]*?>)`
匹配: <img ... src="URL" ...>
```

### 视频和 Source
```go
regexp: `(<(?:video|source)[^>]*?src=[\"'])([^\"']+)([\"'][^>]*?>)`
匹配: <video src="URL"> 或 <source src="URL">
```

### CSS URL
```go
regexp: `(url\\([\"']?)([^\"')]+)([\"']?\\))`
匹配: url("URL") 或 url('URL') 或 url(URL)
```

---

## ⚠️ 注意事项

### 1. **同域名限制**
只替换与当前页面同域名的资源，确保不破坏 CDN 等外部资源。

### 2. **编码问题**
假设 HTML 文件使用 UTF-8 编码，其他编码可能需要额外处理。

### 3. **动态加载**
通过 JavaScript 动态加载的资源可能无法自动替换，需要手动处理。

### 4. **绝对路径限制**
页面内如果使用了 `<base>` 标签，可能影响相对路径的解析。

### 5. **特殊字符**
URL 中包含特殊字符（如空格、中文）可能需要 URL 编码。

---

## 🔄 完整示例

### 原始 HTML
```html
<!DOCTYPE html>
<html>
<head>
    <link href="https://example.com/css/style.css" rel="stylesheet">
    <script src="https://example.com/js/app.js"></script>
    <style>
        body { background: url('https://example.com/images/bg.jpg'); }
    </style>
</head>
<body>
    <img src="https://example.com/images/logo.png" alt="Logo">
    <video src="https://example.com/videos/intro.mp4"></video>
    <script src="https://cdn.example.com/library.js"></script>
</body>
</html>
```

### 转换后 HTML
```html
<!DOCTYPE html>
<html>
<head>
    <link href="./css/style.css" rel="stylesheet">
    <script src="./js/app.js"></script>
    <style>
        body { background: url('./images/bg.jpg'); }
    </style>
</head>
<body>
    <img src="./images/logo.png" alt="Logo">
    <video src="./videos/intro.mp4"></video>
    <script src="https://cdn.example.com/library.js"></script>
    <!-- 注意：外部 CDN 资源保持不变 -->
</body>
</html>
```

---

## 🚀 性能优化

### 1. **正则表达式预编译**
可以将正则表达式编译为全局变量，避免重复编译：

```go
var (
    cssRegex    = regexp.MustCompile(`(<link[^>]*?href=[\"'])([^\"']+)([\"'][^>]*?>)`)
    scriptRegex = regexp.MustCompile(`(<script[^>]*?src=[\"'])([^\"']+)([\"'][^>]*?>)`)
    // ... 其他正则
)
```

### 2. **批量替换**
使用 `strings.Replacer` 进行批量替换可能更高效（但需要先收集所有需要替换的 URL）。

### 3. **HTML 解析器**
对于复杂的 HTML，可以考虑使用 `golang.org/x/net/html` 进行 DOM 解析，更准确但性能稍低。

---

## 🐛 调试日志

函数会输出详细的转换日志：

```
路径替换: https://example.com/css/style.css -> ./css/style.css
路径替换: https://example.com/images/logo.png -> ./images/logo.png
路径替换: https://example.com/js/app.js -> ./js/app.js
```

可以通过日志查看哪些资源被替换，哪些被跳过。

---

## 📌 总结

HTML 资源路径替换功能使得下载的网站能够：

✅ **离线完整浏览** - 所有资源路径正确  
✅ **保持原有结构** - 目录关系不变  
✅ **智能判断** - 只替换需要替换的资源  
✅ **保护外部资源** - CDN 等外部资源保持不变  
✅ **性能优化** - 使用正则批量处理  

这使得整站下载功能更加完善和实用！
