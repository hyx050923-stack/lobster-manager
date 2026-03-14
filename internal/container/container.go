//go:build !android
package container

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"time"
)

type CreateOptions struct {
	Image          string
	ContainerName  string
	PortMap        map[string]string // 容器端口->宿主机端口
	Volumes        map[string]string // 宿主机路径->容器路径
	Env            []string
	SecurityConfig *container.HostConfig
}

// CreateAndStartContainer 创建并启动容器
func (m *Manager) CreateAndStartContainer(opts *CreateOptions) (string, error) {
	// 配置端口映射
	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}
	for containerPort, hostPort := range opts.PortMap {
		port, err := nat.NewPort("tcp", containerPort)
		if err != nil {
			return "", err
		}
		exposedPorts[port] = struct{}{}
		portBindings[port] = []nat.PortBinding{
			{
				HostIP:   "127.0.0.1", // 强制本地回环，可配置化
				HostPort: hostPort,
			},
		}
	}

	// 挂载卷
	var mounts []mount.Mount
	for hostPath, containerPath := range opts.Volumes {
		mounts = append(mounts, mount.Mount{
			Type:     mount.TypeBind,
			Source:   hostPath,
			Target:   containerPath,
			ReadOnly: true, // 默认只读，可由权限配置决定
		})
	}

	// 准备容器配置
	containerConfig := &container.Config{
		Image:        opts.Image,
		ExposedPorts: exposedPorts,
		Env:          opts.Env,
		// 默认命令，由镜像决定
	}

	// 合并安全配置
	hostConfig := opts.SecurityConfig
	if hostConfig == nil {
		hostConfig = DefaultSecurityConfig()
	}
	hostConfig.PortBindings = portBindings
	hostConfig.Mounts = mounts

	// 创建容器
	// 修正：移除了多余的 nil 参数 (platform)
	resp, err := m.cli.ContainerCreate(m.ctx, containerConfig, hostConfig, nil, opts.ContainerName)
	if err != nil {
		return "", err
	}

	// 启动容器
	err = m.cli.ContainerStart(m.ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

// StopContainer 停止容器
func (m *Manager) StopContainer(containerID string) error {
	timeout := 15 * time.Second
	return m.cli.ContainerStop(m.ctx, containerID, &timeout)
}

// RemoveContainer 删除容器
func (m *Manager) RemoveContainer(containerID string, force bool) error {
	return m.cli.ContainerRemove(m.ctx, containerID, types.ContainerRemoveOptions{Force: force})
}

// ListContainers 列出所有龙虾容器（根据镜像名过滤）
func (m *Manager) ListContainers() ([]types.Container, error) {
	filter := filters.NewArgs()
	filter.Add("ancestor", OpenClawImage)
	return m.cli.ContainerList(m.ctx, types.ContainerListOptions{
		All:     true,
		Filters: filter,
	})
}
