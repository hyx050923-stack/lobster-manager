package linux

import (
	"os"
	"os/exec"
)

// CommandExecutor 负责执行命令，处理sudo权限
type CommandExecutor struct {
	UseSudo bool // 是否使用sudo执行
}

// NewExecutor 根据当前用户是否为root决定是否使用sudo
func NewExecutor() *CommandExecutor {
	// 如果当前用户是root，则不需要sudo
	if os.Geteuid() == 0 {
		return &CommandExecutor{UseSudo: false}
	}
	// 否则尝试使用sudo
	return &CommandExecutor{UseSudo: true}
}

// Run 执行命令，返回输出和错误
func (e *CommandExecutor) Run(command string, args ...string) (string, error) {
	var cmd *exec.Cmd
	if e.UseSudo {
		// 构建 sudo 命令
		sudoArgs := append([]string{"-n"}, command) // -n 表示非交互式，如果要求密码则失败
		sudoArgs = append(sudoArgs, args...)
		cmd = exec.Command("sudo", sudoArgs...)
	} else {
		cmd = exec.Command(command, args...)
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// CanSudoWithoutPassword 检查是否可以在不输入密码的情况下使用sudo
func CanSudoWithoutPassword() bool {
	cmd := exec.Command("sudo", "-n", "true")
	err := cmd.Run()
	return err == nil
}