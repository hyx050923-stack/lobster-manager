package runtime

import (
	"os/exec"
)

type UDocker struct {
	Proot  string
	Rootfs string
}

func (u *UDocker) RunUDocker(args ...string) ([]byte, error) {

	base := []string{
		"-S",
		u.Rootfs,
		"udocker",
	}

	base = append(base, args...)

	cmd := exec.Command(u.Proot, base...)

	return cmd.CombinedOutput()
}
