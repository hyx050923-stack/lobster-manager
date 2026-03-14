package container

import (
	"os/exec"
)

type UDockerManager struct {
	ProotPath string
	Rootfs    string
}

func (m *UDockerManager) exec(args ...string) ([]byte, error) {

	cmdArgs := []string{
		"-S",
		m.Rootfs,
		"/usr/bin/udocker",
	}

	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command(m.ProotPath, cmdArgs...)

	return cmd.CombinedOutput()
}

func (m *UDockerManager) Run(image string) error {

	_, err := m.exec("run", image)

	return err
}

func (m *UDockerManager) List() ([]byte, error) {

	return m.exec("ps")
}

func (m *UDockerManager) Stop(id string) error {

	_, err := m.exec("stop", id)

	return err
}