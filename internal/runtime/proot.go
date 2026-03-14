package runtime

import "os/exec"

type Proot struct {
	Path   string
	Rootfs string
}

func (p *Proot) Exec(cmd string, args ...string) ([]byte, error) {

	base := []string{
		"-S",
		p.Rootfs,
		cmd,
	}

	base = append(base, args...)

	return exec.Command(p.Path, base...).CombinedOutput()
}