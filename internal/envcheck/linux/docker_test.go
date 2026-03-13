package linux

import (
	"testing"
)

func TestCheckDockerInstalled(t *testing.T) {
	installed, ver := CheckDockerInstalled()
	t.Logf("Docker installed: %v, version: %s", installed, ver)
	// 不强制要求结果，仅测试无panic
}

func TestCheckDockerServiceActive(t *testing.T) {
	active := CheckDockerServiceActive()
	t.Logf("Docker service active: %v", active)
}

func TestCheckUserInDockerGroup(t *testing.T) {
	inGroup, err := CheckUserInDockerGroup()
	if err != nil {
		t.Errorf("CheckUserInDockerGroup error: %v", err)
	}
	t.Logf("User in docker group: %v", inGroup)
}