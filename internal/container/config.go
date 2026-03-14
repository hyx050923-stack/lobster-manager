//go:build !android
package container

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

// 默认安全配置
func DefaultSecurityConfig() *container.HostConfig {
	return &container.HostConfig{
		// 只读根文件系统（但通常需要一些可写目录，如/tmp，可以使用tmpfs）
		ReadonlyRootfs: true,
		// 添加必要的临时目录为可写tmpfs
		Tmpfs: map[string]string{
			"/tmp": "rw,noexec,nosuid,size=128m",
			"/var": "rw,noexec,nosuid,size=128m",
		},
		// 禁用所有Capabilities，只添加必要的
		CapAdd:  strslice.StrSlice{},
		CapDrop: strslice.StrSlice{"ALL"},
		// 资源限制（默认给2核2G）
		Resources: container.Resources{
			NanoCPUs: 2 * 1e9,           // 2 CPU cores
			Memory:   2 * 1024 * 1024 * 1024, // 2GB
		},
		// 网络模式：默认使用bridge，但后续通过端口映射只暴露本地
		NetworkMode: "bridge",
		// 安全选项
		SecurityOpt: []string{
			"no-new-privileges:true", // 禁止提权
			"seccomp=unconfined",     // 根据实际需求，可能需要更严格的seccomp
		},
	}
}

// 网络配置：只映射到本地回环
func GetPortBindings() map[nat.Port][]nat.PortBinding {
	portBindings := make(map[nat.Port][]nat.PortBinding)
	// 绑定8080端口
	port := nat.Port("8080/tcp")
	portBindings[port] = []nat.PortBinding{
		{
			HostIP:   "0.0.0.0",
			HostPort: "8080",
		},
	}
	return portBindings
}
