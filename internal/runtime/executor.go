package runtime

import "os/exec"

func Run(cmd string, args ...string) ([]byte, error) {

	c := exec.Command(cmd, args...)

	return c.CombinedOutput()
}