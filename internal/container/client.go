package container

import (
	"context"
	"github.com/docker/docker/client"
)

type Manager struct {
	cli *client.Client
	ctx context.Context
}

func NewManager() (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Manager{
		cli: cli,
		ctx: context.Background(),
	}, nil
}

func (m *Manager) Close() error {
	return m.cli.Close()
}