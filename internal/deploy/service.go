// internal/deploy/service.go
package deploy

import (
	"fmt"
	"time"

	"github.com/hyx050923-stack/lobster-manager/internal/container"
	"github.com/hyx050923-stack/lobster-manager/internal/envcheck/common"
)

// Service 部署服务结构体
type Service struct {
	checker common.Checker     // 注入检测接口 (安卓/PC)
	manager *container.Manager // 注入容器管理器
}

// NewService 构造函数
func NewService(c common.Checker, m *container.Manager) *Service {
	return &Service{
		checker: c,
		manager: m,
	}
}

// DeployResult 保持不变
type DeployResult struct {
	Success     bool     `json:"success"`
	ContainerID string   `json:"container_id,omitempty"`
	WebURL      string   `json:"web_url,omitempty"`
	Logs        []string `json:"logs"`
	Error       string   `json:"error,omitempty"`
}

// Deploy 执行一键部署 (改为方法)
func (s *Service) Deploy() *DeployResult {
	logs := []string{}
	result := &DeployResult{Logs: logs}

	// 1. 环境检测 (调用接口)
	logs = append(logs, "开始环境检测...")
	// 注意：CheckAll 返回的是 []common.CheckResult，这里需要遍历判断
	checkResults := s.checker.CheckAll() 
	hasError := false
	for _, r := range checkResults {
		logs = append(logs, r.Message)
		if !r.Status {
			hasError = true
			result.Error = r.FixTip // 直接把修复建议放到错误里
		}
	}
	
	if hasError {
		result.Success = false
		result.Logs = logs
		return result
	}

	// 2. 准备容器配置 (安全沙箱)
	// 注意：这里不再调用 NewManager，而是使用 s.manager
	opts := &container.CreateOptions{
		Image:         container.OpenClawImage,
		ContainerName: fmt.Sprintf("openclaw-%d", time.Now().Unix()),
		PortMap: map[string]string{
			"8080": "18080", // 【重要修改】映射到宿主机 18080，避免与 Go 服务 8080 冲突
		},
		Volumes:        map[string]string{},
		Env:            []string{},
		SecurityConfig: container.DefaultSecurityConfig(), // 这里的实现需要适配 UDocker
	}

	// 3. 启动容器 (Manager 内部包含 Pull 逻辑)
	logs = append(logs, "正在部署 OpenClaw (首次可能较慢)...")
	containerID, err := s.manager.CreateAndStartContainer(opts)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("容器启动失败: %v", err)
		result.Logs = logs
		return result
	}

	logs = append(logs, "部署成功！")
	result.Success = true
	result.ContainerID = containerID
	// 【重要修改】前端访问地址改为 18080
	result.WebURL = "http://127.0.0.1:18080" 
	result.Logs = logs
	return result
}
