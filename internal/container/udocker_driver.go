// internal/container/udocker_driver.go
// +build android // 关键：仅安卓环境编译，或者开发阶段注释掉这行以便测试

package container

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type UdockerDriver struct {
	// 可以包含安卓执行路径等配置
}

func NewUdockerDriver() *UdockerDriver {
	return &UdockerDriver{}
}

func (u *UdockerDriver) CreateAndStartContainer(ctx context.Context, opts *CreateOptions) (string, error) {
	// 1. 构建 UDocker 命令参数
	args := []string{"run", "-d", "--name", opts.ContainerName}

	// 2. 处理端口映射
	// 安卓安全沙箱：强制绑定 127.0.0.1
	for containerPort, hostPort := range opts.PortMap {
		args = append(args, "-p", fmt.Sprintf("127.0.0.1:%s:%s", hostPort, containerPort))
	}

	// 3. 处理挂载卷
	for hostPath, containerPath := range opts.Volumes {
		// UDocker 挂载格式：-v host:container
		args = append(args, "-v", fmt.Sprintf("%s:%s", hostPath, containerPath))
	}

	// 4. 处理环境变量
	for _, env := range opts.Env {
		args = append(args, "-e", env)
	}

	args = append(args, opts.Image)

	// 5. 执行命令 (通过 exec 调用 UDocker 二进制)
	// 注意：安卓下路径需使用绝对路径
	cmd := exec.CommandContext(ctx, "udocker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("udocker run failed: %s", string(output))
	}

	// UDocker run 直接返回容器 ID 或名称
	return strings.TrimSpace(string(output)), nil
}

func (u *UdockerDriver) StopContainer(ctx context.Context, containerID string) error {
	cmd := exec.CommandContext(ctx, "udocker", "stop", containerID)
	return cmd.Run()
}

func (u *UdockerDriver) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	args := []string{"rm", containerID}
	if force {
		args = []string{"rm", "-f", containerID}
	}
	cmd := exec.CommandContext(ctx, "udocker", args...)
	return cmd.Run()
}

func (u *UdockerDriver) ListContainers(ctx context.Context) ([]ContainerInfo, error) {
	// 解析 udocker ps 输出
	cmd := exec.CommandContext(ctx, "udocker", "ps", "-a")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	// 这里需要简单的字符串解析逻辑，将 UDocker 输出转为 ContainerInfo
	// 暂时返回空，后续补全解析逻辑
	_ = string(out) 
	return []ContainerInfo{}, nil
}
