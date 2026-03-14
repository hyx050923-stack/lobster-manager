//go:build windows

package windows

import (
	"os/exec"
	"strings"
)

// CheckWSL2Enabled 检测WSL2功能是否启用
func CheckWSL2Enabled() (bool, error) {
	// 执行 wsl --status 查看输出
	cmd := exec.Command("wsl", "--status")
	out, err := cmd.CombinedOutput()
	if err != nil {
		// wsl命令不存在，可能未安装WSL
		return false, nil
	}
	output := strings.ToLower(string(out))
	// 检查是否包含WSL2相关信息
	return strings.Contains(output, "wsl 2") || strings.Contains(output, "默认版本：2"), nil
}

// CheckWSLDistributionExists 检测是否存在已安装的WSL发行版（至少有一个才能运行容器）
func CheckWSLDistributionExists() bool {
	cmd := exec.Command("wsl", "-l", "-v")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	// 如果有至少一个发行版，通常第一行是标题，第二行开始是列表
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return len(lines) >= 2
}