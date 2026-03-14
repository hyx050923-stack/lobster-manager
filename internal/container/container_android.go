//go:build android

package container

import (
	"fmt"
)

type CreateOptions struct {
	Image         string
	ContainerName string
	PortMap       map[string]string
	Volumes       map[string]string
	Env           []string
}

func (m *Manager) CreateAndStartContainer(opts *CreateOptions) (string, error) {

	cmd := []string{"run", "-d"}

	for host, container := range opts.Volumes {
		cmd = append(cmd, "-v", fmt.Sprintf("%s:%s", host, container))
	}

	for cport, hport := range opts.PortMap {
		cmd = append(cmd, "-p", fmt.Sprintf("%s:%s", hport, cport))
	}

	cmd = append(cmd, opts.Image)

	_, err := m.runtime.RunUDocker(cmd...)

	if err != nil {
		return "", err
	}

	return opts.ContainerName, nil
}

func (m *Manager) StopContainer(id string) error {

	_, err := m.runtime.RunUDocker("stop", id)

	return err
}

func (m *Manager) RemoveContainer(id string, force bool) error {

	_, err := m.runtime.RunUDocker("rm", id)

	return err
}

func (m *Manager) ListContainers() ([]byte, error) {

	return m.runtime.RunUDocker("ps")
}
