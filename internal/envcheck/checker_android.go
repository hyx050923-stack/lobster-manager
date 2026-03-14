package envcheck

import (
	"os/exec"
)

type AndroidChecker struct{}

func (c *AndroidChecker) CheckProot() error {

	cmd := exec.Command("proot", "--version")
	return cmd.Run()
}

func (c *AndroidChecker) CheckUDocker() error {

	cmd := exec.Command("udocker", "--version")
	return cmd.Run()
}

func (c *AndroidChecker) CheckRootfs() error {

	cmd := exec.Command("test", "-d", "/data/data/app/files/rootfs")
	return cmd.Run()
}