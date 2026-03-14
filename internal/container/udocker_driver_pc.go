// internal/container/udocker_driver_pc.go
//go:build !android
// +build !android

package container

import (
	"context"
	"fmt"
	"time"
)

// UdockerDriver PC端模拟实现
type UdockerDriver struct{}

func NewUdockerDriver() *UdockerDriver {
	return &UdockerDriver{}
}

func (u *UdockerDriver) CreateAndStartContainer(ctx context.Context, opts *CreateOptions) (string, error) {
	// PC端模拟：打印日志，返回成功
	fmt.Printf("[PC模拟] 创建容器: %s (镜像: %s)\n", opts.ContainerName, opts.Image)
	time.Sleep(1 * time.Second)
	return "pc-simulated-id-123", nil
}

func (u *UdockerDriver) StopContainer(ctx context.Context, containerID string) error {
	fmt.Printf("[PC模拟] 停止容器: %s\n", containerID)
	return nil
}

func (u *UdockerDriver) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	fmt.Printf("[PC模拟] 删除容器: %s\n", containerID)
	return nil
}

func (u *UdockerDriver) ListContainers(ctx context.Context) ([]ContainerInfo, error) {
	return []ContainerInfo{}, nil
}
