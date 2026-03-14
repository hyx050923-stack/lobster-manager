package deploy

import (
	"fmt"
	"time"

	"github.com/hyx050923-stack/lobster-manager/internal/container"
	"github.com/hyx050923-stack/lobster-manager/internal/envcheck"
)

type DeployResult struct {
	Success     bool     `json:"success"`
	ContainerID string   `json:"container_id,omitempty"`
	WebURL      string   `json:"web_url,omitempty"` // 例如 http://localhost:8080
	Logs        []string `json:"logs"`
	Error       string   `json:"error,omitempty"`
}

// Deploy 执行一键部署
func Deploy() *DeployResult {
	logs := []string{}
	result := &DeployResult{Logs: logs}

	// 1. 环境检测
	logs = append(logs, "开始环境检测...")
	rep, err := envcheck.Detect()
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("环境检测失败: %v", err)
		return result
	}
	logs = append(logs, fmt.Sprintf("平台: %s", rep.Platform))
	if rep.DockerReady {
		logs = append(logs, "环境已就绪")
	} else {
		logs = append(logs, "环境不满足要求，尝试自动修复...")
		// 2. 自动修复
		fixResult, err := envcheck.Fix()
		if err != nil || !fixResult.Success {
			result.Success = false
			result.Error = "环境修复失败，请手动检查"
			logs = append(logs, fixResult.Logs...)
			if err != nil {
				logs = append(logs, err.Error())
			}
			result.Logs = logs
			return result
		}
		logs = append(logs, fixResult.Logs...)
		logs = append(logs, "环境修复完成")
	}

	// 3. 创建容器管理器
	mgr, err := container.NewManager()
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("无法连接Docker: %v", err)
		result.Logs = logs
		return result
	}
	defer mgr.Close()

	// 4. 拉取镜像
	logs = append(logs, "拉取OpenClaw镜像...")
	err = mgr.PullImage(container.OpenClawImage, func(current, total int64, status string) {
		// 这里可以推送进度到前端（通过WebSocket），简单起见先记录最后状态
		logs = append(logs, fmt.Sprintf("拉取进度: %d/%d %s", current, total, status))
	})
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("镜像拉取失败: %v", err)
		result.Logs = logs
		return result
	}
	logs = append(logs, "镜像拉取完成")

	// 5. 创建并启动容器
	logs = append(logs, "创建安全容器...")
	opts := &container.CreateOptions{
		Image:        container.OpenClawImage,
		ContainerName: fmt.Sprintf("openclaw-%d", time.Now().Unix()),
		PortMap: map[string]string{
			"8080": "8080", // 假设OpenClaw默认端口8080，映射到宿主机8080
		},
		Volumes:        map[string]string{}, // 默认不挂载任何宿主机目录，后续权限配置可添加
		Env:            []string{},
		SecurityConfig: container.DefaultSecurityConfig(),
	}
	// 注意：端口映射只绑定了127.0.0.1，实现本地访问
	containerID, err := mgr.CreateAndStartContainer(opts)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("容器启动失败: %v", err)
		result.Logs = logs
		return result
	}

	logs = append(logs, fmt.Sprintf("容器启动成功，ID: %s", containerID))
	result.Success = true
	result.ContainerID = containerID
	result.WebURL = "http://localhost:8080" // 固定，实际可配置
	result.Logs = logs
	return result
}