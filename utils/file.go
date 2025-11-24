package utils

import (
	"bytes"
	"fmt"
	"go-site-clone/config"
	"go-site-clone/types"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type File struct{}

type FileDir struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Mode    string    `json:"mode"`
	ModTime time.Time `json:"modTime"`
}

// 下载文件到本地
func (*File) Download(uri string) string {
	return downloadWithRetry(uri, false, 3, nil)
}

// DownloadWithOptions 带配置选项的下载
func (*File) DownloadWithOptions(uri string, options *types.DownloadOptions) string {
	return downloadWithRetry(uri, false, 3, options)
}

// downloadWithRetry 带重试的下载函数
func downloadWithRetry(uri string, isHTML bool, maxRetries int, options *types.DownloadOptions) string {
	// 如果没有提供配置，使用默认配置
	if options == nil {
		defaultOpts := types.DefaultDownloadOptions()
		options = &defaultOpts
	}
	// 解析文件链接及路径
	u, err := url.Parse(uri)
	if err != nil {
		log.Printf("URL解析失败 %s: %v", uri, err)
		return ""
	}

	host := u.Hostname()
	if u.Port() != "" {
		host += ":" + u.Port()
	}

	// 获取路径部分
	filePath := u.Path
	if isHTML && (filePath == "" || filePath == "/") {
		filePath = "/index.html"
	} else if isHTML && path.Base(filePath) == path.Dir(filePath) {
		filePath = filePath + "/index.html"
	}

	appConfig, _ := config.LoadConfig()
	fp := filepath.Join(appConfig.SiteFileDir, host, filePath)

	// 检查文件是否已存在
	if _, err := os.Stat(fp); err == nil {
		log.Printf("文件已存在，跳过: %s", fp)
		return fp
	}

	// 重试逻辑
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("重试下载 (%d/%d): %s", attempt+1, maxRetries, uri)
			time.Sleep(time.Second * time.Duration(attempt)) // 递增延迟
		}

		// 创建带超时的HTTP客户端
		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		resp, err := client.Get(uri)
		if err != nil {
			lastErr = fmt.Errorf("HTTP请求失败: %w", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
			continue
		}

		// 创建本地文件
		outFile, err := CreateFileWithDirs(fp)
		if err != nil {
			resp.Body.Close()
			lastErr = fmt.Errorf("文件创建失败: %w", err)
			continue
		}

		// 如果是HTML文件，需要替换资源路径
		if isHTML {
			// 读取HTML内容
			htmlContent, err := io.ReadAll(resp.Body)
			resp.Body.Close()

			if err != nil {
				outFile.Close()
				os.Remove(fp)
				lastErr = fmt.Errorf("读取HTML内容失败: %w", err)
				continue
			}

			// 替换资源路径
			modifiedHTML := replaceHTMLResourcePaths(string(htmlContent), uri, options)

			// 写入修改后的HTML
			_, err = outFile.WriteString(modifiedHTML)
			outFile.Close()

			if err != nil {
				os.Remove(fp)
				lastErr = fmt.Errorf("写入HTML文件失败: %w", err)
				continue
			}
		} else {
			// 非 HTML 文件，直接写入
			_, err = io.Copy(outFile, resp.Body)
			resp.Body.Close()
			outFile.Close()

			if err != nil {
				os.Remove(fp) // 删除不完整的文件
				lastErr = fmt.Errorf("写入文件失败: %w", err)
				continue
			}
		}

		log.Printf("下载完成: %s", fp)
		return fp
	}

	log.Printf("下载失败（已重试%d次）%s: %v", maxRetries, uri, lastErr)
	return ""
}

// 创建目录
func CreateFileWithDirs(filePath string) (*os.File, error) {
	// 取出目录部分
	dir := filepath.Dir(filePath)

	// 确保目录存在（递归创建）
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建文件（如果已存在会清空）
	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}

	return f, nil
}

// 下载html文件到本地
func (*File) HTMLDownload(uri string) string {
	return downloadWithRetry(uri, true, 3, nil)
}

// HTMLDownloadWithOptions 带配置选项的HTML下载
func (*File) HTMLDownloadWithOptions(uri string, options *types.DownloadOptions) string {
	return downloadWithRetry(uri, true, 3, options)
}

// replaceHTMLResourcePaths 替换HTML中的资源路径为相对路径
func replaceHTMLResourcePaths(htmlContent string, baseURL string, options *types.DownloadOptions) string {
	if options == nil {
		defaultOpts := types.DefaultDownloadOptions()
		options = &defaultOpts
	}
	parsedBase, err := url.Parse(baseURL)
	if err != nil {
		return htmlContent
	}

	baseHost := parsedBase.Hostname()
	if parsedBase.Port() != "" {
		baseHost += ":" + parsedBase.Port()
	}
	baseScheme := parsedBase.Scheme
	basePath := parsedBase.Path

	// 如果基础路径不是以/结尾，去掉文件名部分
	if basePath != "" && basePath != "/" {
		basePath = path.Dir(basePath)
	}
	if basePath == "." {
		basePath = "/"
	}

	modified := htmlContent

	// 1. 替换 CSS 链接: <link href="..." rel="stylesheet">
	cssRegex := regexp.MustCompile(`(<link[^>]*?href=["'])([^"']+)(["'][^>]*?>)`)
	modified = cssRegex.ReplaceAllStringFunc(modified, func(match string) string {
		parts := cssRegex.FindStringSubmatch(match)
		if len(parts) == 4 {
			originalURL := parts[2]
			if newPath := convertToRelativePath(originalURL, baseURL, baseHost, baseScheme, basePath, "css", options); newPath != "" {
				return parts[1] + newPath + parts[3]
			}
		}
		return match
	})

	// 2. 替换 JavaScript: <script src="...">
	scriptRegex := regexp.MustCompile(`(<script[^>]*?src=["'])([^"']+)(["'][^>]*?>)`)
	modified = scriptRegex.ReplaceAllStringFunc(modified, func(match string) string {
		parts := scriptRegex.FindStringSubmatch(match)
		if len(parts) == 4 {
			originalURL := parts[2]
			if newPath := convertToRelativePath(originalURL, baseURL, baseHost, baseScheme, basePath, "script", options); newPath != "" {
				return parts[1] + newPath + parts[3]
			}
		}
		return match
	})

	// 3. 替换图片: <img src="...">
	imgRegex := regexp.MustCompile(`(<img[^>]*?src=["'])([^"']+)(["'][^>]*?>)`)
	modified = imgRegex.ReplaceAllStringFunc(modified, func(match string) string {
		parts := imgRegex.FindStringSubmatch(match)
		if len(parts) == 4 {
			originalURL := parts[2]
			if newPath := convertToRelativePath(originalURL, baseURL, baseHost, baseScheme, basePath, "image", options); newPath != "" {
				return parts[1] + newPath + parts[3]
			}
		}
		return match
	})

	// 4. 替换视频: <video src="..."> 和 <source src="...">
	videoRegex := regexp.MustCompile(`(<(?:video|source)[^>]*?src=["'])([^"']+)(["'][^>]*?>)`)
	modified = videoRegex.ReplaceAllStringFunc(modified, func(match string) string {
		parts := videoRegex.FindStringSubmatch(match)
		if len(parts) == 4 {
			originalURL := parts[2]
			if newPath := convertToRelativePath(originalURL, baseURL, baseHost, baseScheme, basePath, "video", options); newPath != "" {
				return parts[1] + newPath + parts[3]
			}
		}
		return match
	})

	// 5. 替换音频: <audio src="...">
	audioRegex := regexp.MustCompile(`(<audio[^>]*?src=["'])([^"']+)(["'][^>]*?>)`)
	modified = audioRegex.ReplaceAllStringFunc(modified, func(match string) string {
		parts := audioRegex.FindStringSubmatch(match)
		if len(parts) == 4 {
			originalURL := parts[2]
			if newPath := convertToRelativePath(originalURL, baseURL, baseHost, baseScheme, basePath, "video", options); newPath != "" {
				return parts[1] + newPath + parts[3]
			}
		}
		return match
	})

	// 6. 替换 CSS 中的 url(): url("...") 或 url('...') 或 url(...)
	cssURLRegex := regexp.MustCompile(`(url\(["']?)([^"')]+)(["']?\))`)
	modified = cssURLRegex.ReplaceAllStringFunc(modified, func(match string) string {
		parts := cssURLRegex.FindStringSubmatch(match)
		if len(parts) == 4 {
			originalURL := parts[2]
			// 跳过 data: 和 # 开头的URL
			if !strings.HasPrefix(originalURL, "data:") && !strings.HasPrefix(originalURL, "#") {
				if newPath := convertToRelativePath(originalURL, baseURL, baseHost, baseScheme, basePath, "image", options); newPath != "" {
					return parts[1] + newPath + parts[3]
				}
			}
		}
		return match
	})

	return modified
}

// convertToRelativePath 将绝对URL转换为相对路径
func convertToRelativePath(resourceURL, baseURL, baseHost, baseScheme, basePath, resourceType string, options *types.DownloadOptions) string {
	// 跳过空值、data:, javascript:, mailto:, # 等
	if resourceURL == "" ||
		strings.HasPrefix(resourceURL, "data:") ||
		strings.HasPrefix(resourceURL, "javascript:") ||
		strings.HasPrefix(resourceURL, "mailto:") ||
		strings.HasPrefix(resourceURL, "#") {
		return ""
	}

	// 如果已经是相对路径，不处理
	if !strings.HasPrefix(resourceURL, "http://") &&
		!strings.HasPrefix(resourceURL, "https://") &&
		!strings.HasPrefix(resourceURL, "//") {
		return ""
	}

	// 处理 // 开头的URL
	if strings.HasPrefix(resourceURL, "//") {
		resourceURL = baseScheme + ":" + resourceURL
	}

	// 解析资源URL
	parsedResource, err := url.Parse(resourceURL)
	if err != nil {
		return ""
	}

	resourceHost := parsedResource.Hostname()
	if parsedResource.Port() != "" {
		resourceHost += ":" + parsedResource.Port()
	}

	// 使用配置判断是否应该下载该资源
	if !options.ShouldDownloadExternal(resourceHost, baseHost, resourceType) {
		// 不下载外部资源，保持原URL
		return ""
	} // 获取资源路径
	resourcePath := parsedResource.Path
	if resourcePath == "" {
		resourcePath = "/"
	}

	// 计算相对路径
	relativePath := calculateRelativePath(basePath, resourcePath)

	log.Printf("路径替换: %s -> %s", resourceURL, relativePath)

	return relativePath
}

// calculateRelativePath 计算从 basePath 到 targetPath 的相对路径
func calculateRelativePath(basePath, targetPath string) string {
	if basePath == "" {
		basePath = "/"
	}
	if targetPath == "" {
		targetPath = "/"
	}

	// 分割路径
	baseParts := strings.Split(strings.Trim(basePath, "/"), "/")
	targetParts := strings.Split(strings.Trim(targetPath, "/"), "/")

	// 找到公共前缀
	commonLen := 0
	for i := 0; i < len(baseParts) && i < len(targetParts); i++ {
		if baseParts[i] == targetParts[i] {
			commonLen++
		} else {
			break
		}
	}

	// 构建相对路径
	var result bytes.Buffer

	// 添加 ../ 返回上级目录
	upLevels := len(baseParts) - commonLen
	for i := 0; i < upLevels; i++ {
		if i > 0 {
			result.WriteString("/")
		}
		result.WriteString("..")
	}

	// 添加目标路径的剩余部分
	for i := commonLen; i < len(targetParts); i++ {
		if result.Len() > 0 {
			result.WriteString("/")
		}
		result.WriteString(targetParts[i])
	}

	// 如果结果为空，返回当前目录
	if result.Len() == 0 {
		return "./"
	}

	// 如果不是以 . 或 / 开头，添加 ./
	resultStr := result.String()
	if !strings.HasPrefix(resultStr, ".") && !strings.HasPrefix(resultStr, "/") {
		return "./" + resultStr
	}

	return resultStr
}

// 获取文件夹列表
func (*File) GetFileDirList(filePath string) []FileDir {
	var fileList []FileDir
	entries, err := os.ReadDir(filePath)
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				fmt.Println("无法获取属性:", err)
				continue
			}
			fileList = append(fileList, FileDir{
				Name:    entry.Name(),
				Size:    info.Size(),
				Mode:    info.Mode().String(),
				ModTime: info.ModTime(),
			})
		}
	}
	return fileList
}
