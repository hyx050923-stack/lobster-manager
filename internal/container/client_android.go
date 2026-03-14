//go:build android

package container

import "errors"

type Manager struct{}

func NewManager() (*Manager, error) {
	return &Manager{}, nil
}

func (m *Manager) Close() error {
	return nil
}

var ErrNotSupported = errors.New("docker client not supported on android")