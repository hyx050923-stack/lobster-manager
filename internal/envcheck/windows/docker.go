//go:build windows

package windows

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// FindDockerPath 尝试查找docker.exe的路径
func FindDockerPath() (string, error) {
	// 常见安装路径
	paths := []string{
		`C:\Program Files\Docker\Docker\resources\bin\docker.exe`,
		`C:\Program Files\Docker\Docker\resources\bin\docker.exe`,
	}
	// 也可以通过where命令查找
	cmd := exec.Command("where", "docker")
	out, err := cmd.Output()
	if err == nil {
		// 返回第一个找到的路径
		first := strings.Split(string(out), "\n")[0]
		return strings.TrimSpace(first), nil
	}
	// 遍历常见路径
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", exec.ErrNotFound
}

// CheckDockerInstalled 检测Docker是否安装
func CheckDockerInstalled() bool {
	path, err := FindDockerPath()
	if err != nil {
		return false
	}
	// 尝试执行 docker version
	cmd := exec.Command(path, "version")
	err = cmd.Run()
	return err == nil
}

// CheckDockerWSL2Backend 检查Docker是否使用WSL2后端（通过查看Docker Desktop设置或执行docker info）
func CheckDockerWSL2Backend() bool {
	// 执行 docker info 检查是否包含 WSL2
	cmd := exec.Command("docker", "info", "--format", "{{.OSType}}")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	// 如果输出包含windows，可能是Windows容器；如果包含linux，则可能使用WSL2
	// 更准确的：检查是否存在 "WSL2" 字样
	cmd2 := exec.Command("docker", "info")
	out2, _ := cmd2.CombinedOutput()
	return strings.Contains(string(out2), "WSL2")
}