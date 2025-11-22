package services

import (
	"bytes"
	"fmt"
	"go-site-clone/config"
	"go-site-clone/types"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"
)

type NginxService struct {
	Process *exec.Cmd
	Running bool
}

// NginxSiteConfig 站点配置结构
type NginxSiteConfig struct {
	ID      string   // 站点唯一ID
	Name    string   // 站点名称
	Domains []string // 域名列表
	Port    int      // 监听端口
	Path    string   // 网站根目录路径
	Index   string   // 默认首页文件
	Enabled bool     // 是否启用
}

// 启动 Nginx
func (n *NginxService) StartNginx() error {
	if n.Running {
		return fmt.Errorf("nginx 已经在运行中")
	}

	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 获取 nginx 可执行文件路径
	nginxPath := appConfig.Nginx.NginxPath
	if nginxPath == "" {
		nginxPath = "plugin/nginx/nginx.exe"
	}

	// 检查 nginx 文件是否存在
	if _, err := os.Stat(nginxPath); os.IsNotExist(err) {
		return fmt.Errorf("nginx 可执行文件不存在: %s", nginxPath)
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(nginxPath)
	if err != nil {
		return fmt.Errorf("获取 Nginx 路径失败: %v", err)
	}

	// 启动 nginx
	cmd := exec.Command(absPath)
	cmd.Dir = filepath.Dir(absPath)

	// 设置输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("启动 Nginx 失败: %v, stderr: %s", err, stderr.String())
	}

	n.Process = cmd
	n.Running = true

	// 保存 PID 到文件，用于后续检测
	pidFile := filepath.Join(filepath.Dir(absPath), "logs", "nginx.pid")
	pid := cmd.Process.Pid
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		log.Printf("警告: 保存 PID 文件失败: %v", err)
	}

	log.Printf("Nginx 启动成功, PID: %d", pid)
	return nil
}

// 停止 Nginx
func (n *NginxService) StopNginx() error {
	if !n.Running {
		return fmt.Errorf("nginx 未运行")
	}

	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	nginxPath := appConfig.Nginx.NginxPath
	if nginxPath == "" {
		nginxPath = "plugin/nginx/nginx.exe"
	}

	// 获取绝对路径
	absPath, _ := filepath.Abs(nginxPath)
	nginxDir := filepath.Dir(absPath)

	// 首先尝试优雅停止
	stopCmd := exec.Command(nginxPath, "-s", "quit")
	stopCmd.Dir = filepath.Dir(nginxPath)
	err = stopCmd.Run()

	// 等待一段时间让进程退出
	time.Sleep(1 * time.Second)

	// 检查进程是否还在运行
	stillRunning := false
	if runtime.GOOS == "windows" {
		// Windows 下检查所有 nginx.exe 进程
		stillRunning = n.checkAnyNginxProcess(nginxDir)
	}

	// 如果还有进程在运行，强制杀死
	if stillRunning || err != nil {
		log.Println("优雅停止失败或进程仍在运行，执行强制停止")

		if runtime.GOOS == "windows" {
			// Windows 下使用 taskkill 强制杀死所有 nginx.exe 进程
			n.killAllNginxProcesses(nginxDir)
		} else {
			// Linux/Mac 下使用 nginx -s stop
			stopCmd := exec.Command(nginxPath, "-s", "stop")
			stopCmd.Dir = filepath.Dir(nginxPath)
			stopCmd.Run()
		}

		// 再等待一下确保进程退出
		time.Sleep(500 * time.Millisecond)
	}

	n.Running = false
	n.Process = nil

	// 删除 PID 文件
	pidFile := filepath.Join(nginxDir, "logs", "nginx.pid")
	os.Remove(pidFile) // 忽略错误

	log.Println("Nginx 已停止")
	return nil
}

// 重启 Nginx
func (n *NginxService) RestartNginx() error {
	if n.Running {
		if err := n.StopNginx(); err != nil {
			return fmt.Errorf("停止 Nginx 失败: %v", err)
		}
		// 等待一段时间确保完全停止
		time.Sleep(2 * time.Second)
	}

	return n.StartNginx()
}

// 重载配置
func (n *NginxService) ReloadNginx() error {
	if !n.Running {
		return fmt.Errorf("nginx 未运行")
	}

	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	nginxPath := appConfig.Nginx.NginxPath
	if nginxPath == "" {
		nginxPath = "plugin/nginx/nginx.exe"
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(nginxPath)
	if err != nil {
		return fmt.Errorf("获取 nginx 路径失败: %v", err)
	}

	// 执行重载命令
	reloadCmd := exec.Command(absPath, "-s", "reload")
	reloadCmd.Dir = filepath.Dir(absPath)

	if err := reloadCmd.Run(); err != nil {
		return fmt.Errorf("重载配置失败: %v", err)
	}

	log.Println("Nginx 配置已重载")
	return nil
}

// 检查 Nginx 状态
func (n *NginxService) CheckStatus() (bool, error) {
	appConfig, err := config.LoadConfig()
	if err != nil {
		return false, fmt.Errorf("加载配置失败: %v", err)
	}

	nginxPath := appConfig.Nginx.NginxPath
	if nginxPath == "" {
		nginxPath = "plugin/nginx/nginx.exe"
	}

	// 获取 nginx 目录的绝对路径
	nginxDir := filepath.Dir(nginxPath)
	absNginxDir, err := filepath.Abs(nginxDir)
	if err != nil {
		return false, err
	}

	// 读取 PID 文件
	pidFile := filepath.Join(absNginxDir, "logs", "nginx.pid")
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		// PID 文件不存在，说明没有运行或异常退出
		return false, nil
	}

	// 解析 PID
	var pid int
	if _, err := fmt.Sscanf(string(pidData), "%d", &pid); err != nil {
		return false, fmt.Errorf("解析 PID 文件失败: %v", err)
	}

	if runtime.GOOS == "windows" {
		// Windows 下检查进程是否存在并验证路径
		return n.checkWindowsProcess(pid, absNginxDir)
	} else {
		// Linux/Mac 下检查进程
		return n.checkUnixProcess(pid)
	}
}

// 测试配置文件
func (n *NginxService) TestConfig() error {
	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	nginxPath := appConfig.Nginx.NginxPath
	if nginxPath == "" {
		nginxPath = "plugin/nginx/nginx.exe"
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(nginxPath)
	if err != nil {
		return fmt.Errorf("获取 nginx 路径失败: %v", err)
	}

	// 执行配置测试
	testCmd := exec.Command(absPath, "-t")
	testCmd.Dir = filepath.Dir(absPath)

	var stderr bytes.Buffer
	testCmd.Stderr = &stderr

	if err := testCmd.Run(); err != nil {
		return fmt.Errorf("配置文件测试失败: %s", stderr.String())
	}

	log.Println("Nginx 配置文件测试通过")
	return nil
}

// 添加站点配置
func (n *NginxService) AddSite(site types.NginxSiteConfig) error {
	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 生成配置文件内容
	confContent, err := n.generateSiteConfig(site)
	if err != nil {
		return fmt.Errorf("生成配置失败: %v", err)
	}

	// 根据启用状态选择目录
	hostsDir := filepath.Join(filepath.Dir(appConfig.Nginx.NginxConfPath), "hosts")
	disabledDir := filepath.Join(filepath.Dir(appConfig.Nginx.NginxConfPath), "hosts.disabled")

	var targetDir string
	if site.Enabled {
		targetDir = hostsDir
	} else {
		targetDir = disabledDir
	}

	// 确保目录存在
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 写入配置文件
	confPath := filepath.Join(targetDir, fmt.Sprintf("%s.conf", site.Name))
	if err := os.WriteFile(confPath, []byte(confContent), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	log.Printf("站点配置已添加: %s (启用: %v)", confPath, site.Enabled)

	// 如果启用且 Nginx 正在运行,重载配置
	if site.Enabled && n.Running {
		if err := n.ReloadNginx(); err != nil {
			// 重载失败只记录警告，不影响添加操作
			log.Printf("警告: 重载 nginx 配置失败: %v", err)
		}
	}

	return nil
}

// 删除站点配置
func (n *NginxService) DeleteSite(siteName string) error {
	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	hostsDir := filepath.Join(filepath.Dir(appConfig.Nginx.NginxConfPath), "hosts")
	disabledDir := filepath.Join(filepath.Dir(appConfig.Nginx.NginxConfPath), "hosts.disabled")

	// 先尝试从启用目录删除
	confPath := filepath.Join(hostsDir, fmt.Sprintf("%s.conf", siteName))
	err = os.Remove(confPath)
	if err != nil && os.IsNotExist(err) {
		// 如果启用目录没有，尝试从禁用目录删除
		confPath = filepath.Join(disabledDir, fmt.Sprintf("%s.conf", siteName))
		err = os.Remove(confPath)
	}

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("站点配置不存在: %s", siteName)
		}
		return fmt.Errorf("删除配置文件失败: %v", err)
	}

	log.Printf("站点配置已删除: %s", siteName)

	// 如果 Nginx 正在运行,重载配置
	if n.Running {
		if err := n.ReloadNginx(); err != nil {
			// 重载失败只记录警告，不影响删除操作
			log.Printf("警告: 重载 nginx 配置失败: %v", err)
		}
	}

	return nil
}

// 更新站点配置
func (n *NginxService) UpdateSite(site types.NginxSiteConfig) error {
	// 先删除旧配置
	if err := n.DeleteSite(site.Name); err != nil {
		// 如果配置不存在,继续添加新配置
		if !strings.Contains(err.Error(), "不存在") {
			return err
		}
	}

	// 添加新配置
	return n.AddSite(site)
}

// 获取所有站点配置
func (n *NginxService) GetAllSites() ([]types.NginxSiteConfig, error) {
	appConfig, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %v", err)
	}

	hostsDir := filepath.Join(filepath.Dir(appConfig.Nginx.NginxConfPath), "hosts")
	disabledDir := filepath.Join(filepath.Dir(appConfig.Nginx.NginxConfPath), "hosts.disabled")

	// 检查目录是否存在
	if _, err := os.Stat(hostsDir); os.IsNotExist(err) {
		os.MkdirAll(hostsDir, 0755)
	}
	if _, err := os.Stat(disabledDir); os.IsNotExist(err) {
		os.MkdirAll(disabledDir, 0755)
	}

	var sites []types.NginxSiteConfig

	// 读取启用的站点
	enabledFiles, err := os.ReadDir(hostsDir)
	if err == nil {
		for _, file := range enabledFiles {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".conf") {
				siteName := strings.TrimSuffix(file.Name(), ".conf")
				confPath := filepath.Join(hostsDir, file.Name())
				site, err := n.parseSiteConfig(confPath, siteName, true)
				if err != nil {
					log.Printf("警告: 解析配置文件失败 %s: %v", file.Name(), err)
					continue
				}
				sites = append(sites, site)
			}
		}
	}

	// 读取禁用的站点
	disabledFiles, err := os.ReadDir(disabledDir)
	if err == nil {
		for _, file := range disabledFiles {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".conf") {
				siteName := strings.TrimSuffix(file.Name(), ".conf")
				confPath := filepath.Join(disabledDir, file.Name())
				site, err := n.parseSiteConfig(confPath, siteName, false)
				if err != nil {
					log.Printf("警告: 解析配置文件失败 %s: %v", file.Name(), err)
					continue
				}
				sites = append(sites, site)
			}
		}
	}

	return sites, nil
}

// 解析站点配置文件
func (n *NginxService) parseSiteConfig(confPath string, siteName string, enabled bool) (types.NginxSiteConfig, error) {
	content, err := os.ReadFile(confPath)
	if err != nil {
		return types.NginxSiteConfig{}, err
	}

	config := types.NginxSiteConfig{
		Name:    siteName,
		Enabled: enabled,
		Domains: []string{},
	}

	// 简单解析配置文件内容
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 解析 listen 端口
		if strings.HasPrefix(line, "listen") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				portStr := strings.TrimSuffix(parts[1], ";")
				var port int
				fmt.Sscanf(portStr, "%d", &port)
				config.Port = port
			}
		}

		// 解析 server_name 域名
		if strings.HasPrefix(line, "server_name") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				for i := 1; i < len(parts); i++ {
					domain := strings.TrimSuffix(parts[i], ";")
					if domain != "" {
						config.Domains = append(config.Domains, domain)
					}
				}
			}
		}

		// 解析 root 路径
		if strings.HasPrefix(line, "root") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				path := strings.Trim(parts[1], "\";")
				config.Path = path
			}
		}

		// 解析 index
		if strings.HasPrefix(line, "index") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				index := strings.TrimSuffix(parts[1], ";")
				config.Index = index
			}
		}
	}

	return config, nil
}

// 启用站点
func (n *NginxService) EnableSite(siteName string) error {
	// 对于当前的配置结构,所有在 hosts 目录下的配置都是启用的
	// 如果需要禁用,可以移动到其他目录
	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	hostsDir := filepath.Join(filepath.Dir(appConfig.Nginx.NginxConfPath), "hosts")
	disabledDir := filepath.Join(filepath.Dir(appConfig.Nginx.NginxConfPath), "hosts.disabled")

	// 从 disabled 移动到 hosts
	srcPath := filepath.Join(disabledDir, fmt.Sprintf("%s.conf", siteName))
	dstPath := filepath.Join(hostsDir, fmt.Sprintf("%s.conf", siteName))

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("站点配置不存在或已启用")
	}

	if err := os.MkdirAll(hostsDir, 0755); err != nil {
		return err
	}

	if err := os.Rename(srcPath, dstPath); err != nil {
		return fmt.Errorf("启用站点失败: %v", err)
	}

	// 重载配置
	if n.Running {
		if err := n.ReloadNginx(); err != nil {
			log.Printf("警告: 重载 nginx 配置失败: %v", err)
		}
	}

	return nil
}

// 禁用站点
func (n *NginxService) DisableSite(siteName string) error {
	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	hostsDir := filepath.Join(filepath.Dir(appConfig.Nginx.NginxConfPath), "hosts")
	disabledDir := filepath.Join(filepath.Dir(appConfig.Nginx.NginxConfPath), "hosts.disabled")

	// 从 hosts 移动到 disabled
	srcPath := filepath.Join(hostsDir, fmt.Sprintf("%s.conf", siteName))
	dstPath := filepath.Join(disabledDir, fmt.Sprintf("%s.conf", siteName))

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("站点配置不存在或已禁用")
	}

	if err := os.MkdirAll(disabledDir, 0755); err != nil {
		return err
	}

	if err := os.Rename(srcPath, dstPath); err != nil {
		return fmt.Errorf("禁用站点失败: %v", err)
	}

	// 重载配置
	if n.Running {
		if err := n.ReloadNginx(); err != nil {
			log.Printf("警告: 重载 nginx 配置失败: %v", err)
		}
	}

	return nil
}

// 生成站点配置文件内容
func (n *NginxService) generateSiteConfig(site types.NginxSiteConfig) (string, error) {
	// 如果用户提供了路径，使用用户的路径
	var absPath string
	var err error

	if site.Path != "" {
		// 用户提供的路径（相对路径或绝对路径）
		if filepath.IsAbs(site.Path) {
			absPath = site.Path
		} else {
			// 相对路径，转换为绝对路径
			absPath, err = filepath.Abs(site.Path)
			if err != nil {
				return "", fmt.Errorf("获取站点路径失败: %v", err)
			}
		}
	} else {
		// 如果没有提供路径，使用配置中的默认目录
		appConfig, err := config.LoadConfig()
		if err != nil {
			return "", err
		}

		sitePath := filepath.Join(appConfig.SiteFileDir, site.Name)
		absPath, err = filepath.Abs(sitePath)
		if err != nil {
			return "", fmt.Errorf("获取站点路径失败: %v", err)
		}
	}

	// Windows 路径处理（nginx 需要使用正斜杠）
	if runtime.GOOS == "windows" {
		absPath = strings.ReplaceAll(absPath, "\\", "/")
	}

	// 默认首页
	if site.Index == "" {
		site.Index = "index.html"
	}

	// 生成域名列表字符串
	domains := strings.Join(site.Domains, " ")

	// 使用模板生成配置
	tmpl, err := template.New("nginx").Parse(types.NGINX_CONF_TMP)
	if err != nil {
		return "", err
	}

	// 获取当前时间
	createTime := time.Now().Format("2006-01-02 15:04:05")

	data := map[string]interface{}{
		"SiteName":   site.Name,
		"Port":       site.Port,
		"Host":       domains,
		"FilePath":   absPath,
		"Index":      site.Index,
		"CreateTime": createTime,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// 获取 Nginx 日志
func (n *NginxService) GetAccessLog(lines int) ([]string, error) {
	appConfig, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %v", err)
	}

	logPath := filepath.Join(filepath.Dir(appConfig.Nginx.NginxPath), "logs", "access.log")
	return n.readLogLines(logPath, lines)
}

// 获取错误日志
func (n *NginxService) GetErrorLog(lines int) ([]string, error) {
	appConfig, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %v", err)
	}

	logPath := filepath.Join(filepath.Dir(appConfig.Nginx.NginxPath), "logs", "error.log")
	return n.readLogLines(logPath, lines)
}

// 读取日志文件指定行数
func (n *NginxService) readLogLines(filePath string, lines int) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	allLines := strings.Split(string(content), "\n")
	start := len(allLines) - lines
	if start < 0 {
		start = 0
	}

	return allLines[start:], nil
}

// 清空日志
func (n *NginxService) ClearLogs() error {
	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	logsDir := filepath.Join(filepath.Dir(appConfig.Nginx.NginxPath), "logs")

	// 清空 access.log
	accessLog := filepath.Join(logsDir, "access.log")
	if err := os.WriteFile(accessLog, []byte(""), 0644); err != nil {
		return fmt.Errorf("清空访问日志失败: %v", err)
	}

	// 清空 error.log
	errorLog := filepath.Join(logsDir, "error.log")
	if err := os.WriteFile(errorLog, []byte(""), 0644); err != nil {
		return fmt.Errorf("清空错误日志失败: %v", err)
	}

	log.Println("Nginx 日志已清空")
	return nil
}

// Windows 下检查进程
func (n *NginxService) checkWindowsProcess(pid int, nginxDir string) (bool, error) {
	// 使用 tasklist 检查进程是否存在
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	outputStr := string(output)
	// 如果输出包含 "INFO: No tasks"，说明进程不存在
	if strings.Contains(outputStr, "INFO:") || !strings.Contains(outputStr, "nginx.exe") {
		return false, nil
	}

	// 进程存在，进一步验证进程路径（可选，但更严格）
	// 使用 wmic 获取进程的可执行文件路径
	cmd = exec.Command("wmic", "process", "where", fmt.Sprintf("ProcessId=%d", pid), "get", "ExecutablePath", "/format:list")
	output, err = cmd.Output()
	if err != nil {
		// 如果获取路径失败，但进程存在，认为是运行中
		return true, nil
	}

	outputStr = string(output)
	// 检查路径是否包含我们的 nginx 目录
	nginxDirNorm := strings.ToLower(strings.ReplaceAll(nginxDir, "/", "\\"))
	outputStrNorm := strings.ToLower(outputStr)

	return strings.Contains(outputStrNorm, nginxDirNorm), nil
}

// Unix/Linux 下检查进程
func (n *NginxService) checkUnixProcess(pid int) (bool, error) {
	// 检查 /proc/pid 是否存在
	procPath := fmt.Sprintf("/proc/%d", pid)
	if _, err := os.Stat(procPath); os.IsNotExist(err) {
		return false, nil
	}

	// 读取进程的 cmdline 验证是否是 nginx
	cmdlinePath := filepath.Join(procPath, "cmdline")
	cmdline, err := os.ReadFile(cmdlinePath)
	if err != nil {
		return false, err
	}

	return strings.Contains(string(cmdline), "nginx"), nil
}

// 检查是否有任何 nginx 进程在运行（Windows）
func (n *NginxService) checkAnyNginxProcess(nginxDir string) bool {
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq nginx.exe", "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	outputStr := string(output)
	// 检查是否有 nginx.exe 进程
	if !strings.Contains(outputStr, "nginx.exe") {
		return false
	}

	// 进一步验证是否是本项目的 nginx（通过路径）
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if !strings.Contains(line, "nginx.exe") {
			continue
		}

		// 提取 PID
		fields := strings.Split(line, "\",\"")
		if len(fields) < 2 {
			continue
		}

		pidStr := strings.Trim(fields[1], "\"")
		var pid int
		if _, err := fmt.Sscanf(pidStr, "%d", &pid); err != nil {
			continue
		}

		// 检查这个进程的路径
		cmd = exec.Command("wmic", "process", "where", fmt.Sprintf("ProcessId=%d", pid), "get", "ExecutablePath", "/format:list")
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		outputStr := string(output)
		nginxDirNorm := strings.ToLower(strings.ReplaceAll(nginxDir, "/", "\\"))
		if strings.Contains(strings.ToLower(outputStr), nginxDirNorm) {
			return true
		}
	}

	return false
}

// 杀死所有本项目的 nginx 进程（Windows）
func (n *NginxService) killAllNginxProcesses(nginxDir string) error {
	// 获取所有 nginx.exe 进程
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq nginx.exe", "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "nginx.exe") {
		return nil // 没有 nginx 进程
	}

	// 解析每个进程并杀死属于本项目的
	lines := strings.Split(outputStr, "\n")
	nginxDirNorm := strings.ToLower(strings.ReplaceAll(nginxDir, "/", "\\"))

	for _, line := range lines {
		if !strings.Contains(line, "nginx.exe") {
			continue
		}

		// 提取 PID
		fields := strings.Split(line, "\",\"")
		if len(fields) < 2 {
			continue
		}

		pidStr := strings.Trim(fields[1], "\"")
		var pid int
		if _, err := fmt.Sscanf(pidStr, "%d", &pid); err != nil {
			continue
		}

		// 验证进程路径是否属于本项目
		cmd = exec.Command("wmic", "process", "where", fmt.Sprintf("ProcessId=%d", pid), "get", "ExecutablePath", "/format:list")
		output, err := cmd.Output()
		if err != nil {
			// 如果无法获取路径，为了安全起见，杀死这个进程
			// （因为我们已经通过其他方式确认这是 nginx.exe）
			killCmd := exec.Command("taskkill", "/F", "/PID", pidStr)
			killCmd.Run()
			log.Printf("已杀死 nginx 进程 PID: %s", pidStr)
			continue
		}

		// 检查路径是否匹配
		pathOutput := strings.ToLower(string(output))
		if strings.Contains(pathOutput, nginxDirNorm) {
			// 强制杀死进程
			killCmd := exec.Command("taskkill", "/F", "/PID", pidStr)
			if err := killCmd.Run(); err != nil {
				log.Printf("警告: 杀死进程 PID %s 失败: %v", pidStr, err)
			} else {
				log.Printf("已杀死 nginx 进程 PID: %s", pidStr)
			}
		}
	}

	return nil
}
