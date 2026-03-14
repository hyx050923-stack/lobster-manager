package linux

import (
	"fmt"
	"strings"
)

// FixDockerInstall 安装Docker
func FixDockerInstall(executor *CommandExecutor) (string, error) {
	// 先确保有curl或wget
	// 使用官方脚本安装
	out, err := executor.Run("sh", "-c", "curl -fsSL https://get.docker.com | sh")
	if err != nil {
		return out, fmt.Errorf("安装Docker失败: %v", err)
	}
	return out, nil
}

// FixDockerStart 启动Docker服务
func FixDockerStart(executor *CommandExecutor) (string, error) {
	out, err := executor.Run("systemctl", "start", "docker")
	if err != nil {
		return out, fmt.Errorf("启动Docker服务失败: %v", err)
	}
	return out, nil
}

// FixDockerEnable 设置Docker开机自启
func FixDockerEnable(executor *CommandExecutor) (string, error) {
	out, err := executor.Run("systemctl", "enable", "docker")
	if err != nil {
		return out, fmt.Errorf("设置Docker开机自启失败: %v", err)
	}
	return out, nil
}

// FixAddUserToDockerGroup 将当前用户添加到docker组
func FixAddUserToDockerGroup(executor *CommandExecutor) (string, error) {
	// 获取当前用户名
	whoami, err := executor.Run("whoami")
	if err != nil {
		return "", fmt.Errorf("获取用户名失败: %v", err)
	}
	user := strings.TrimSpace(whoami)
	out, err := executor.Run("usermod", "-aG", "docker", user)
	if err != nil {
		return out, fmt.Errorf("添加用户到docker组失败: %v", err)
	}
	return out, nil
}

// FixAll 根据报告执行所有必要的修复
func FixAll(report *EnvReport) ([]string, error) {
	executor := NewExecutor()
	var logs []string
	var lastErr error

	if !report.DockerInstalled {
		logs = append(logs, "开始安装Docker...")
		out, err := FixDockerInstall(executor)
		logs = append(logs, out)
		if err != nil {
			lastErr = err
			logs = append(logs, fmt.Sprintf("安装Docker失败: %v", err))
		} else {
			logs = append(logs, "Docker安装成功")
		}
	}
	// 如果刚刚安装了Docker，可能服务已自动启动，但为了保险还是检查并启动
	if report.DockerInstalled || !report.DockerServiceActive {
		logs = append(logs, "启动Docker服务...")
		out, err := FixDockerStart(executor)
		logs = append(logs, out)
		if err != nil {
			lastErr = err
			logs = append(logs, fmt.Sprintf("启动Docker服务失败: %v", err))
		} else {
			logs = append(logs, "Docker服务启动成功")
		}
	}
	// 设置开机自启（可选）
	logs = append(logs, "设置Docker开机自启...")
	out, err := FixDockerEnable(executor)
	logs = append(logs, out)
	if err != nil {
		lastErr = err
		logs = append(logs, fmt.Sprintf("设置开机自启失败: %v", err))
	} else {
		logs = append(logs, "开机自启设置成功")
	}

	if !report.UserInDockerGroup {
		logs = append(logs, "将当前用户添加到docker组...")
		out, err := FixAddUserToDockerGroup(executor)
		logs = append(logs, out)
		if err != nil {
			lastErr = err
			logs = append(logs, fmt.Sprintf("添加到docker组失败: %v", err))
		} else {
			logs = append(logs, "用户已添加到docker组，请重新登录或执行 'newgrp docker' 使权限生效")
		}
	}
	return logs, lastErr
}