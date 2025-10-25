package utils

import (
	"os/exec"
)

func CheckGoEnv() bool {
	cmd := exec.Command("go", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func CheckWailsEnv() bool {
	cmd := exec.Command("wails3", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
