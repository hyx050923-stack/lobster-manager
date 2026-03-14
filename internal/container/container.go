// internal/container/container.go
package container

import (
	"context"
)

// ContainerInfo 统一的容器信息结构（不依赖 Docker SDK 的 types.Container）
type ContainerInfo struct {
	ID     string
	Name   string
	Status string
	Image  string
}

// Driver 定义容器驱动接口（核心解耦层）
type Driver interface {
	CreateAndStartContainer(ctx context.Context, opts *CreateOptions) (string, error)
	StopContainer(ctx context.Context, containerID string) error
	RemoveContainer(ctx context.Context, containerID string, force bool) error
	ListContainers(ctx context.Context) ([]ContainerInfo, error)
}

// Manager 容器管理器
type Manager struct {
	driver Driver
	ctx    context.Context
}

// NewManager 根据运行环境初始化不同的驱动
func NewManager(ctx context.Context) *Manager {
	var d Driver
	
	// 这里后续通过条件编译或配置判断环境
	// 目前为了跑通安卓 MVP，我们强制注入 UDocker 驱动
	// 实际生产环境可以通过 build tags 来区分
	d = NewUdockerDriver() 
	
	return &Manager{
		driver: d,
		ctx:    ctx,
	}
}

// CreateOptions 保持不变，作为通用的业务参数
type CreateOptions struct {
	Image          string
	ContainerName  string
	PortMap        map[string]string
	Volumes        map[string]string
	Env            []string
	SecurityConfig interface{} // 暂时改为 interface{}，由具体驱动处理
}

// --- 以下是暴露给外部的业务方法 ---

func (m *Manager) CreateAndStartContainer(opts *CreateOptions) (string, error) {
	return m.driver.CreateAndStartContainer(m.ctx, opts)
}

func (m *Manager) StopContainer(containerID string) error {
	return m.driver.StopContainer(m.ctx, containerID)
}

func (m *Manager) RemoveContainer(containerID string, force bool) error {
	return m.driver.RemoveContainer(m.ctx, containerID, force)
}

func (m *Manager) ListContainers() ([]ContainerInfo, error) {
	return m.driver.ListContainers(m.ctx)
}
