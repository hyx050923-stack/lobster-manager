package container

type ContainerInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

type Manager interface {

	// 创建并启动容器
	CreateAndStartContainer(opts *CreateOptions) (string, error)

	// 停止容器
	StopContainer(containerID string) error

	// 删除容器
	RemoveContainer(containerID string, force bool) error

	// 列出容器
	ListContainers() (interface{}, error)

	// 拉取镜像
	PullImage(image string) error

	Close() error
}
