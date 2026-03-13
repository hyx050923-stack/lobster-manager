package linux

import (
	"os/exec"
	"strings"
)

// CheckDockerInstalled 检测Docker是否安装，返回是否安装及版本号
func CheckDockerInstalled() (installed bool, version string) {
	cmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	out, err := cmd.Output()
	if err != nil {
		return false, ""
	}
	ver := strings.TrimSpace(string(out))
	return ver != "", ver
}

// CheckDockerServiceActive 检测Docker服务是否运行（使用systemctl）
func CheckDockerServiceActive() bool {
	cmd := exec.Command("systemctl", "is-active", "docker")
	err := cmd.Run()
	return err == nil
}

// CheckUserInDockerGroup 检测当前用户是否在docker组中
func CheckUserInDockerGroup() (bool, error) {
	// 读取 /etc/group
	out, err := exec.Command("getent", "group", "docker").Output()
	if err != nil {
		// docker组不存在
		return false, nil
	}
	// 输出格式：docker:x:999:user1,user2
	parts := strings.Split(strings.TrimSpace(string(out)), ":")
	if len(parts) < 4 {
		return false, nil
	}
	members := strings.Split(parts[3], ",")
	// 获取当前用户名
	currentUserOut, err := exec.Command("whoami").Output()
	if err != nil {
		return false, err
	}
	currentUser := strings.TrimSpace(string(currentUserOut))
	for _, m := range members {
		if m == currentUser {
			return true, nil
		}
	}
	return false, nil
}