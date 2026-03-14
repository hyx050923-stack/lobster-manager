package envcheck

import (
	"runtime"

	"github.com/hyx050923-stack/lobster-manager/internal/envcheck/linux"
)

type FixResult struct {
	Success bool     `json:"success"`
	Logs    []string `json:"logs"`
	Message string   `json:"message"`
}

// Fix 根据当前平台执行自动修复
func Fix() (*FixResult, error) {
	switch runtime.GOOS {
	case "linux":
		// 先检测当前环境
		rep, err := linux.Detect()
		if err != nil {
			return nil, err
		}
		// 执行修复
		logs, err := linux.FixAll(rep)
		result := &FixResult{
			Logs: logs,
		}
		if err != nil {
			result.Success = false
			result.Message = err.Error()
		} else {
			result.Success = true
			result.Message = "修复完成，部分操作可能需要重启或重新登录才能生效"
		}
		return result, nil
	case "windows":
		// TODO: 实现Windows修复
		return &FixResult{
			Success: false,
			Message: "Windows自动修复尚未实现",
		}, nil
	default:
		return &FixResult{
			Success: false,
			Message: "不支持的操作系统",
		}, nil
	}
}