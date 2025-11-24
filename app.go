package main

import (
	"bytes"
	"fmt"
	"go-site-clone/config"
	"go-site-clone/services"
	"go-site-clone/storage"
	"go-site-clone/types"
	"go-site-clone/utils"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-rod/rod/lib/proto"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	app          *application.App
	nginxService *services.NginxService
	store        *storage.Store
}

var siteService services.SiteService

// 服务启动时初始化 nginx 服务和数据库
func (a *App) OnStartup() {
	// 初始化数据库
	dbPath := "data/site-clone.db"
	store, err := storage.NewStore(dbPath)
	if err != nil {
		log.Printf("初始化数据库失败: %v", err)
	} else {
		a.store = store
		log.Printf("数据库初始化成功")
	}

	// 初始化 nginx 服务
	a.nginxService = &services.NginxService{}

	// 检测 nginx 是否已经在运行
	running, err := a.nginxService.CheckStatus()
	if err == nil && running {
		a.nginxService.Running = true
		fmt.Println("nginx 服务已检测到正在运行")
	} else {
		a.nginxService.Running = false
		fmt.Println("nginx 服务未运行")
	}
}

// 服务关闭时不关闭 nginx，但关闭数据库
func (a *App) OnShutdown() {
	// 关闭数据库连接
	if a.store != nil {
		if err := a.store.Close(); err != nil {
			log.Printf("关闭数据库失败: %v", err)
		} else {
			log.Println("数据库已关闭")
		}
	}
	// 不做其他事，让 nginx 继续运行
	fmt.Println("应用关闭，nginx 服务保持运行状态")
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
// func (a *App) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {

// 	return nil
// }

// func (a *App) ServiceShutdown(ctx context.Context) error {

//		return nil
//	}
//
// ========== 网站克隆相关方法 ==========
func (a *App) GetResources(rawURL string) *services.ResourcesList {
	RequestParamsAll, RouterAll := siteService.GetAllResources(rawURL)
	// 再次去重 根据url链接去除重复项
	dataList := siteService.DeduplicationRequestByUrl(RequestParamsAll)
	resources := &services.ResourcesList{
		Dom: RouterAll,
	}
	// 将页面及资源一起返回
	for _, v := range dataList {
		if v.Type == proto.NetworkResourceTypeImage {
			resources.Image = append(resources.Image, v.URL)
		}
		if v.Type == proto.NetworkResourceTypeStylesheet {
			resources.CSS = append(resources.CSS, v.URL)
		}
		if v.Type == proto.NetworkResourceTypeScript {
			resources.Script = append(resources.Script, v.URL)
		}
		if v.Type == proto.NetworkResourceTypeMedia {
			resources.Video = append(resources.Video, v.URL)
		}
	}
	return resources
}

// 下载网站资源
func (a *App) DownloadSite(uri string, obj services.ResourcesList) bool {
	// 将页面及资源一起返回
	parsed, _ := url.Parse(uri)
	var File utils.File
	if len(obj.CSS) > 0 {
		for k, v := range obj.CSS {
			u, _ := url.Parse(v)
			if parsed.Hostname() == u.Hostname() {
				File.Download(v)
			}
			a.app.Event.Emit("download:css", k)
		}
	}
	if len(obj.Script) > 0 {
		for k, v := range obj.Script {
			u, _ := url.Parse(v)
			if parsed.Hostname() == u.Hostname() {
				File.Download(v)
			}
			a.app.Event.Emit("download:script", k)
		}
	}
	if len(obj.Image) > 0 {
		for k, v := range obj.Image {
			u, _ := url.Parse(v)
			if parsed.Hostname() == u.Hostname() {
				File.Download(v)
			}
			a.app.Event.Emit("download:image", k)
		}
	}
	if len(obj.Video) > 0 {
		for k, v := range obj.Video {
			u, _ := url.Parse(v)
			if parsed.Hostname() == u.Hostname() {
				File.Download(v)
			}
			a.app.Event.Emit("download:video", k)
		}
	}
	if len(obj.Dom) > 0 {
		for k, v := range obj.Dom {
			u, _ := url.Parse(v)
			if parsed.Hostname() == u.Hostname() {
				File.HTMLDownload(v)
			}
			a.app.Event.Emit("download:dom", k)
		}
	}

	return true
}

// 获取本地已下载网站列表
func (a *App) GetDownloadList() []utils.FileDir {
	return siteService.GetLocalSiteList()
}

// 打开网站文件夹
func (a *App) OpenSiteFileDir(pathDir string) bool {
	appConfig, _ := config.LoadConfig()
	newPath := filepath.Join(appConfig.SiteFileDir, pathDir)
	cmd := exec.Command("explorer", newPath)
	err := cmd.Start()
	if err != nil {
		fmt.Println("打开文件夹失败:", err)
		return false
	}
	return true
}

// 删除网站文件夹
func (a *App) DeleteSiteFileDir(pathDir string) bool {
	appConfig, _ := config.LoadConfig()
	newPath := filepath.Join(appConfig.SiteFileDir, pathDir)
	err := os.RemoveAll(newPath)
	if err != nil {
		fmt.Println("删除文件夹失败:", err)
		return false
	}
	return true
}

// ========== Nginx 相关方法 ==========

// 启动 Nginx
func (a *App) StartNginx() error {
	return a.nginxService.StartNginx()
}

// 停止 Nginx
func (a *App) StopNginx() error {
	return a.nginxService.StopNginx()
}

// 重启 Nginx
func (a *App) RestartNginx() error {
	return a.nginxService.RestartNginx()
}

// 重载 Nginx 配置
func (a *App) ReloadNginx() error {
	return a.nginxService.ReloadNginx()
}

// 检查 Nginx 状态
func (a *App) CheckNginxStatus() (bool, error) {
	running, err := a.nginxService.CheckStatus()
	if err == nil {
		a.nginxService.Running = running
	}
	return running, err
}

// 测试 Nginx 配置
func (a *App) TestNginxConfig() error {
	return a.nginxService.TestConfig()
}

// 添加站点配置
func (a *App) AddNginxSite(site types.NginxSiteConfig) error {
	// 先保存到数据库
	if a.store != nil {
		if err := a.store.AddSite(site); err != nil {
			return fmt.Errorf("保存站点到数据库失败: %v", err)
		}
	}
	// 再创建 nginx 配置文件
	return a.nginxService.AddSite(site)
}

// 删除站点配置
func (a *App) DeleteNginxSite(siteName string) error {
	// 从数据库删除
	if a.store != nil {
		if err := a.store.DeleteSite(siteName); err != nil {
			log.Printf("警告: 从数据库删除站点失败: %v", err)
		}
	}
	// 删除 nginx 配置文件
	return a.nginxService.DeleteSite(siteName)
}

// 更新站点配置
func (a *App) UpdateNginxSite(site types.NginxSiteConfig) error {
	// 更新数据库
	if a.store != nil {
		if err := a.store.UpdateSite(site); err != nil {
			return fmt.Errorf("更新数据库失败: %v", err)
		}
	}
	// 更新 nginx 配置
	return a.nginxService.UpdateSite(site)
}

// 获取所有站点配置
func (a *App) GetAllNginxSites() ([]types.NginxSiteConfig, error) {
	// 优先从数据库读取
	if a.store != nil {
		sites, err := a.store.GetAllSites()
		if err == nil && len(sites) > 0 {
			return sites, nil
		}
	}
	// 如果数据库为空，从 nginx 配置文件读取并同步到数据库
	sites, err := a.nginxService.GetAllSites()
	if err != nil {
		return nil, err
	}
	// 同步到数据库
	if a.store != nil {
		for _, site := range sites {
			a.store.AddSite(site)
		}
	}
	return sites, nil
}

// 启用站点
func (a *App) EnableNginxSite(siteName string) error {
	// 更新数据库
	if a.store != nil {
		if err := a.store.UpdateSiteStatus(siteName, true); err != nil {
			log.Printf("警告: 更新数据库失败: %v", err)
		}
	}
	return a.nginxService.EnableSite(siteName)
}

// 禁用站点
func (a *App) DisableNginxSite(siteName string) error {
	// 更新数据库
	if a.store != nil {
		if err := a.store.UpdateSiteStatus(siteName, false); err != nil {
			log.Printf("警告: 更新数据库失败: %v", err)
		}
	}
	return a.nginxService.DisableSite(siteName)
}

// 获取 Nginx 访问日志
func (a *App) GetNginxAccessLog(lines int) ([]string, error) {
	return a.nginxService.GetAccessLog(lines)
}

// 获取 Nginx 错误日志
func (a *App) GetNginxErrorLog(lines int) ([]string, error) {
	return a.nginxService.GetErrorLog(lines)
}

// 清空 Nginx 日志
func (a *App) ClearNginxLogs() error {
	return a.nginxService.ClearLogs()
}

// 获取可用端口号（从 startPort 开始查找）
func (a *App) getAvailablePort(startPort int) int {
	// 获取所有已配置的站点
	sites, err := a.nginxService.GetAllSites()
	if err != nil {
		return startPort
	}

	// 收集已使用的端口
	usedPorts := make(map[int]bool)
	for _, site := range sites {
		if site.Port > 0 {
			usedPorts[site.Port] = true
		}
	}

	// 查找可用端口
	port := startPort
	for {
		if !usedPorts[port] {
			return port
		}
		port++
		if port > 65535 {
			return startPort // 如果超出范围，返回默认端口
		}
	}
}

// ========== 下载记录相关方法 ==========

// AddDownloadRecord 添加下载记录
func (a *App) AddDownloadRecord(record storage.DownloadRecord) error {
	if a.store == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.store.AddDownloadRecord(record)
}

// GetAllDownloadRecords 获取所有下载记录
func (a *App) GetAllDownloadRecords() ([]storage.DownloadRecord, error) {
	if a.store == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return a.store.GetAllDownloadRecords()
}

// GetRecentDownloadRecords 获取最近的下载记录
func (a *App) GetRecentDownloadRecords(limit int) ([]storage.DownloadRecord, error) {
	if a.store == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return a.store.GetRecentDownloadRecords(limit)
}

// GetDownloadStats 获取下载统计
func (a *App) GetDownloadStats() (map[string]interface{}, error) {
	if a.store == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return a.store.GetDownloadStats()
}

// DeleteDownloadRecord 删除下载记录
func (a *App) DeleteDownloadRecord(id string) error {
	if a.store == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.store.DeleteDownloadRecord(id)
}

// ClearOldDownloadRecords 清理旧的下载记录
func (a *App) ClearOldDownloadRecords(days int) (int, error) {
	if a.store == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}
	return a.store.ClearOldDownloadRecords(days)
}

// BackupDatabase 备份数据库
func (a *App) BackupDatabase(backupPath string) error {
	if a.store == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.store.Backup(backupPath)
}

// SelectFolder 打开文件夹选择对话框
func (a *App) SelectFolder() (string, error) {
	result, err := a.app.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
		CanChooseDirectories: true,
		CanChooseFiles:       false,
		Title:                "选择网站文件夹",
		ButtonText:           "选择",
	}).PromptForSingleSelection()

	if err != nil {
		return "", err
	}
	return result, nil
}

// ========== 打包相关方法 ==========

// CheckEnvironment 检查Go和Wails环境
func (a *App) CheckEnvironment() (*utils.EnvStatus, error) {
	return utils.GetEnvStatus()
}

// InstallGo 安装Go环境
func (a *App) InstallGo(version string) error {
	scriptPath, err := utils.GetInstallScriptPath("go")
	if err != nil {
		return err
	}

	// 构建命令参数
	args := []string{"/c", scriptPath}
	if version != "" {
		args = append(args, version)
	}

	// 发送安装进度事件
	a.app.Event.Emit("install:progress", map[string]interface{}{
		"tool":    "go",
		"step":    "downloading",
		"percent": 10,
		"message": "正在下载 Go...",
	})

	// 在Windows上使用cmd.exe执行bat文件
	cmd := exec.Command("cmd.exe", args...)
	cmd.Dir = filepath.Dir(scriptPath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 启动命令但不等待完成(bat脚本会暂停)
	err = cmd.Start()
	if err != nil {
		a.app.Event.Emit("install:progress", map[string]interface{}{
			"tool":    "go",
			"step":    "error",
			"percent": 0,
			"error":   fmt.Sprintf("启动安装脚本失败: %v", err),
		})
		return fmt.Errorf("安装失败: %v", err)
	}

	// 在后台等待完成
	go func() {
		err := cmd.Wait()
		if err != nil {
			a.app.Event.Emit("install:progress", map[string]interface{}{
				"tool":    "go",
				"step":    "error",
				"percent": 0,
				"error":   stderr.String(),
			})
		} else {
			a.app.Event.Emit("install:progress", map[string]interface{}{
				"tool":    "go",
				"step":    "completed",
				"percent": 100,
				"message": "Go 安装完成",
			})
		}
	}()

	return nil
}

// InstallWails 安装Wails3环境
func (a *App) InstallWails() error {
	scriptPath, err := utils.GetInstallScriptPath("wails")
	if err != nil {
		return err
	}

	// 发送安装进度事件
	a.app.Event.Emit("install:progress", map[string]interface{}{
		"tool":    "wails",
		"step":    "installing",
		"percent": 10,
		"message": "正在安装 Wails3...",
	})

	// 在Windows上使用cmd.exe执行bat文件
	cmd := exec.Command("cmd.exe", "/c", scriptPath)
	cmd.Dir = filepath.Dir(scriptPath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 启动命令但不等待完成
	err = cmd.Start()
	if err != nil {
		a.app.Event.Emit("install:progress", map[string]interface{}{
			"tool":    "wails",
			"step":    "error",
			"percent": 0,
			"error":   fmt.Sprintf("启动安装脚本失败: %v", err),
		})
		return fmt.Errorf("安装失败: %v", err)
	}

	// 在后台等待完成
	go func() {
		err := cmd.Wait()
		if err != nil {
			a.app.Event.Emit("install:progress", map[string]interface{}{
				"tool":    "wails",
				"step":    "error",
				"percent": 0,
				"error":   stderr.String(),
			})
		} else {
			a.app.Event.Emit("install:progress", map[string]interface{}{
				"tool":    "wails",
				"step":    "completed",
				"percent": 100,
				"message": "Wails3 安装完成",
			})
		}
	}()

	return nil
}

// PackApp 打包应用
func (a *App) PackApp(packConfig map[string]interface{}) error {
	// 获取配置参数
	sitePath, _ := packConfig["sitePath"].(string)
	appName, _ := packConfig["appName"].(string)
	version, _ := packConfig["version"].(string)
	author, _ := packConfig["author"].(string)
	description, _ := packConfig["description"].(string)
	outputDir, _ := packConfig["outputDir"].(string)
	width, _ := packConfig["width"].(float64)
	height, _ := packConfig["height"].(float64)

	// 平台参数处理
	platformsInterface, _ := packConfig["platforms"].([]interface{})
	platforms := make([]string, 0)
	for _, p := range platformsInterface {
		if platform, ok := p.(string); ok {
			platforms = append(platforms, platform)
		}
	}

	if sitePath == "" || appName == "" {
		return fmt.Errorf("缺少必要参数: 网站路径和应用名称不能为空")
	}

	// 检查网站路径是否存在
	if _, err := os.Stat(sitePath); os.IsNotExist(err) {
		return fmt.Errorf("网站路径不存在: %s", sitePath)
	}

	// 发送打包进度
	a.app.Event.Emit("pack:progress", map[string]interface{}{
		"step":    "preparing",
		"percent": 5,
		"message": "正在准备打包环境...",
	})

	// 检查环境
	envStatus, err := utils.GetEnvStatus()
	if err != nil {
		return fmt.Errorf("检查环境失败: %v", err)
	}

	if !envStatus.HasGo || !envStatus.HasWails {
		return fmt.Errorf("缺少必要的环境: Go=%v, Wails=%v。请先安装开发环境", envStatus.HasGo, envStatus.HasWails)
	}

	// 设置默认值
	if version == "" {
		version = "1.0.0"
	}
	if description == "" {
		description = appName
	}
	if len(platforms) == 0 {
		platforms = []string{"windows"}
	}
	if width == 0 {
		width = 1280
	}
	if height == 0 {
		height = 800
	}

	// 设置输出目录
	if outputDir == "" {
		appConfig, _ := config.LoadConfig()
		outputDir = appConfig.PackSiteFileDir
		if outputDir == "" {
			outputDir = "site-dist"
		}
	}

	// 创建临时项目目录
	a.app.Event.Emit("pack:progress", map[string]interface{}{
		"step":    "creating",
		"percent": 15,
		"message": "正在创建Wails项目结构...",
	})

	tempDir := filepath.Join(os.TempDir(), "wails-pack-"+appName)
	os.RemoveAll(tempDir) // 清理旧的
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}

	// 使用打包服务创建项目
	packService := &services.PackService{}
	config := services.PackConfig{
		SitePath:    sitePath,
		AppName:     appName,
		Version:     version,
		Author:      author,
		Description: description,
		Platforms:   platforms,
		Width:       int(width),
		Height:      int(height),
		OutputDir:   outputDir,
	}

	log.Printf("开始打包网站: %s", sitePath)
	log.Printf("应用名称: %s", appName)
	log.Printf("临时目录: %s", tempDir)
	log.Printf("输出目录: %s", outputDir)

	a.app.Event.Emit("pack:progress", map[string]interface{}{
		"step":    "copying",
		"percent": 30,
		"message": "正在复制网站文件...",
	})

	// 创建Wails项目结构
	if err := packService.CreateWailsProject(config, tempDir); err != nil {
		a.app.Event.Emit("pack:progress", map[string]interface{}{
			"step":    "error",
			"percent": 0,
			"error":   fmt.Sprintf("创建项目失败: %v", err),
		})
		return fmt.Errorf("创建Wails项目失败: %v", err)
	}

	a.app.Event.Emit("pack:progress", map[string]interface{}{
		"step":    "building",
		"percent": 50,
		"message": "正在编译应用(这可能需要几分钟)...",
	})

	// 构建应用
	if err := packService.BuildWailsApp(tempDir, outputDir, platforms); err != nil {
		a.app.Event.Emit("pack:progress", map[string]interface{}{
			"step":    "error",
			"percent": 0,
			"error":   fmt.Sprintf("构建失败: %v", err),
		})
		return fmt.Errorf("构建应用失败: %v", err)
	}

	a.app.Event.Emit("pack:progress", map[string]interface{}{
		"step":    "completed",
		"percent": 100,
		"message": fmt.Sprintf("打包完成! 应用已保存到: %s", outputDir),
	})

	log.Printf("打包完成! 输出目录: %s", outputDir)
	return nil
}
