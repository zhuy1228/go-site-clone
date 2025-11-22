package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EnvStatus 环境状态结构
type EnvStatus struct {
	HasGo        bool   `json:"hasGo"`
	GoVersion    string `json:"goVersion"`
	HasWails     bool   `json:"hasWails"`
	WailsVersion string `json:"wailsVersion"`
	GoPath       string `json:"goPath"`
	WailsPath    string `json:"wailsPath"`
}

// 检测当前电脑是否有GO环境
func CheckGoEnv() bool {
	cmd := exec.Command("go", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// 检测当前电脑是否有Wails3 环境
func CheckWailsEnv() bool {
	cmd := exec.Command("wails3", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// GetEnvStatus 获取详细的环境状态
func GetEnvStatus() (*EnvStatus, error) {
	status := &EnvStatus{}

	// 检查系统Go环境
	cmd := exec.Command("go", "version")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		status.HasGo = true
		status.GoVersion = strings.TrimSpace(out.String())
		// 获取go路径
		goPath, _ := exec.LookPath("go")
		status.GoPath = goPath
	} else {
		// 检查本地安装的Go
		localGoPath := getLocalGoPath()
		if localGoPath != "" {
			cmd := exec.Command(localGoPath, "version")
			var localOut bytes.Buffer
			cmd.Stdout = &localOut
			if err := cmd.Run(); err == nil {
				status.HasGo = true
				status.GoVersion = strings.TrimSpace(localOut.String())
				status.GoPath = localGoPath
			}
		}
	}

	// 检查系统Wails环境
	cmd = exec.Command("wails3", "version")
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		status.HasWails = true
		status.WailsVersion = strings.TrimSpace(out.String())
		// 获取wails路径
		wailsPath, _ := exec.LookPath("wails3")
		status.WailsPath = wailsPath
	} else {
		// 检查本地安装的Wails
		localWailsPath := getLocalWailsPath()
		if localWailsPath != "" {
			cmd := exec.Command(localWailsPath, "version")
			var localOut bytes.Buffer
			cmd.Stdout = &localOut
			if err := cmd.Run(); err == nil {
				status.HasWails = true
				status.WailsVersion = strings.TrimSpace(localOut.String())
				status.WailsPath = localWailsPath
			}
		}
	}

	return status, nil
}

// getLocalGoPath 获取本地安装的Go路径
func getLocalGoPath() string {
	possiblePaths := []string{
		"plugin/go/1.25.3/bin/go.exe",
		"plugin/go/bin/go.exe",
	}

	for _, p := range possiblePaths {
		if _, err := os.Stat(p); err == nil {
			absPath, _ := filepath.Abs(p)
			return absPath
		}
	}
	return ""
}

// getLocalWailsPath 获取本地安装的Wails路径
func getLocalWailsPath() string {
	possiblePaths := []string{
		"plugin/wails3/wails.exe",
		"plugin/wails3/wails3.exe",
	}

	for _, p := range possiblePaths {
		if _, err := os.Stat(p); err == nil {
			absPath, _ := filepath.Abs(p)
			return absPath
		}
	}
	return ""
}

// InstallProgress 安装进度信息
type InstallProgress struct {
	Step    string `json:"step"`    // 当前步骤
	Percent int    `json:"percent"` // 进度百分比
	Message string `json:"message"` // 消息
	Error   string `json:"error"`   // 错误信息
}

// GetInstallScriptPath 获取安装脚本路径
func GetInstallScriptPath(tool string) (string, error) {
	var scriptName string
	if tool == "go" {
		scriptName = "go.bat"
	} else if tool == "wails" {
		scriptName = "wails.bat"
	} else {
		return "", fmt.Errorf("未知的工具: %s", tool)
	}

	scriptPath := filepath.Join("install", scriptName)
	absPath, err := filepath.Abs(scriptPath)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("安装脚本不存在: %s", absPath)
	}

	return absPath, nil
}
