package utils

import (
	"os/exec"
)

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
