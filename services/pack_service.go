package services

import (
	"fmt"
	"go-site-clone/utils"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// PackService 打包服务
type PackService struct{}

// PackConfig 打包配置
type PackConfig struct {
	SitePath    string   // 网站文件夹路径
	AppName     string   // 应用名称
	Version     string   // 版本号
	Author      string   // 作者
	Description string   // 描述
	Platforms   []string // 目标平台
	Width       int      // 窗口宽度
	Height      int      // 窗口高度
	OutputDir   string   // 输出目录
}

// Wails项目模板
const mainGoTemplate = `package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "{{.AppName}}",
		Description: "{{.Description}}",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "{{.AppName}}",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(255, 255, 255),
		URL:              "/",
		Width:            {{.Width}},
		Height:           {{.Height}},
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
`

const goModTemplate = `module {{.ModuleName}}

go 1.21

require github.com/wailsapp/wails/v3 v3.0.0-alpha.27
`

// CreateWailsProject 创建Wails项目结构
func (p *PackService) CreateWailsProject(config PackConfig, projectDir string) error {
	// 1. 创建项目目录结构
	dirs := []string{
		projectDir,
		filepath.Join(projectDir, "frontend"),
		filepath.Join(projectDir, "build"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败 %s: %v", dir, err)
		}
	}

	// 1.1 复制项目根目录的 Taskfile.yml
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}
	srcTaskfile := filepath.Join(currentDir, "Taskfile.yml")
	dstTaskfile := filepath.Join(projectDir, "Taskfile.yml")
	if err := p.copyFile(srcTaskfile, dstTaskfile); err != nil {
		return fmt.Errorf("复制Taskfile.yml失败: %v", err)
	}

	// 1.2 复制整个 build 目录
	srcBuildDir := filepath.Join(currentDir, "build")
	dstBuildDir := filepath.Join(projectDir, "build")
	if err := p.copyDirectory(srcBuildDir, dstBuildDir); err != nil {
		return fmt.Errorf("复制build目录失败: %v", err)
	}

	// 2. 复制网站文件到frontend目录
	frontendDir := filepath.Join(projectDir, "frontend")
	if err := p.copyDirectory(config.SitePath, frontendDir); err != nil {
		return fmt.Errorf("复制网站文件失败: %v", err)
	}

	// 2.1 创建一个空的 package.json，避免构建系统报错
	packageJSON := `{
  "name": "static-site",
  "version": "1.0.0",
  "private": true,
  "scripts": {
    "build": "echo 'Static site, no build needed'",
    "generate": "echo 'Static site, no build needed'"
  }
}`
	packageJSONPath := filepath.Join(frontendDir, "package.json")
	if err := os.WriteFile(packageJSONPath, []byte(packageJSON), 0644); err != nil {
		return fmt.Errorf("创建package.json失败: %v", err)
	}

	// 2.2 创建一个空的 dist 目录（如果构建系统需要）
	distDir := filepath.Join(frontendDir, "dist")
	if err := os.MkdirAll(distDir, 0755); err != nil {
		return fmt.Errorf("创建dist目录失败: %v", err)
	}

	// 3. 生成main.go
	mainGoPath := filepath.Join(projectDir, "main.go")
	tmpl, err := template.New("main").Parse(mainGoTemplate)
	if err != nil {
		return fmt.Errorf("解析main.go模板失败: %v", err)
	}

	mainFile, err := os.Create(mainGoPath)
	if err != nil {
		return fmt.Errorf("创建main.go失败: %v", err)
	}
	defer mainFile.Close()

	templateData := map[string]interface{}{
		"AppName":     config.AppName,
		"Description": config.Description,
		"Width":       config.Width,
		"Height":      config.Height,
	}

	if err := tmpl.Execute(mainFile, templateData); err != nil {
		return fmt.Errorf("生成main.go失败: %v", err)
	}

	// 4. 生成go.mod
	goModPath := filepath.Join(projectDir, "go.mod")
	modTmpl, err := template.New("gomod").Parse(goModTemplate)
	if err != nil {
		return fmt.Errorf("解析go.mod模板失败: %v", err)
	}

	modFile, err := os.Create(goModPath)
	if err != nil {
		return fmt.Errorf("创建go.mod失败: %v", err)
	}
	defer modFile.Close()

	modData := map[string]interface{}{
		"ModuleName": strings.ToLower(strings.ReplaceAll(config.AppName, " ", "-")),
	}

	if err := modTmpl.Execute(modFile, modData); err != nil {
		return fmt.Errorf("生成go.mod失败: %v", err)
	}

	return nil
}

// BuildWailsApp 构建Wails应用
func (p *PackService) BuildWailsApp(projectDir string, outputDir string, platforms []string) error {
	// 获取环境状态
	envStatus, err := utils.GetEnvStatus()
	if err != nil {
		return fmt.Errorf("检查环境失败: %v", err)
	}

	// 确定使用的Go命令
	goCmd := "go"
	if envStatus.GoPath != "" {
		goCmd = envStatus.GoPath
	}

	// 1. 执行 go mod tidy
	tidyCmd := exec.Command(goCmd, "mod", "tidy")
	tidyCmd.Dir = projectDir
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy 失败: %v\n%s", err, output)
	}

	// 2. 直接使用 go build 编译，跳过复杂的构建流程
	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 确定输出文件名
	appName := filepath.Base(projectDir)
	outputFile := filepath.Join(outputDir, appName+".exe")

	// 构建命令，使用生产环境标签
	buildCmd := exec.Command(goCmd, "build",
		"-tags", "production",
		"-ldflags", "-w -s -H windowsgui",
		"-o", outputFile,
		".")
	buildCmd.Dir = projectDir

	// 设置环境变量
	buildCmd.Env = append(os.Environ(),
		"CGO_ENABLED=1",
		"GOOS=windows",
	)

	output, err := buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go build 失败: %v\n%s", err, output)
	}

	return nil
}

// copyDirectory 递归复制目录
func (p *PackService) copyDirectory(src, dst string) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// 读取源目录
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归复制子目录
			if err := p.copyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 复制文件
			if err := p.copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile 复制单个文件
func (p *PackService) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// 复制文件权限
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

// copyBuildOutput 复制构建输出
func (p *PackService) copyBuildOutput(buildDir, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// 检查 buildDir 是否存在
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		return fmt.Errorf("构建目录不存在: %s", buildDir)
	}

	entries, err := os.ReadDir(buildDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(buildDir, entry.Name())
		dstPath := filepath.Join(outputDir, entry.Name())

		if entry.IsDir() {
			if err := p.copyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := p.copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
