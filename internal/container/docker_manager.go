//go:build !android

package container

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerManager struct {
	cli *client.Client
	ctx context.Context
}

func NewDockerManager() (*DockerManager, error) {
	var _ Manager = (*UDockerManager)(nil)
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		return nil, err
	}

	return &DockerManager{
		cli: cli,
		ctx: context.Background(),
	}, nil
}

func (m *UDockerManager) PullImage(image string) error {

	_, err := m.exec("pull", image)

	return err
}