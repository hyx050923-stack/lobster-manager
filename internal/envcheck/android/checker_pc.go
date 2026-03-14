//go:build !android
// +build !android

package android

import (
	"github.com/hyx050923-stack/lobster-manager/internal/envcheck/common"
)

// AndroidChecker 在非安卓平台下的模拟实现
type AndroidChecker struct{}

func NewChecker() *AndroidChecker {
	return &AndroidChecker{}
}

// CheckAll 模拟返回成功的检测结果，方便在 PC 上调试
func (c *AndroidChecker) CheckAll() []common.CheckResult {
	return []common.CheckResult{
		{
			Item:    "CPU架构",
			Status:  true,
			Message: "PC 模拟环境: amd64 (通过)",
			FixTip:  "",
		},
		{
			Item:    "系统版本",
			Status:  true,
			Message: "PC 模拟环境: Linux Kernel (通过)",
			FixTip:  "",
		},
		{
			Item:    "存储空间",
			Status:  true,
			Message: "PC 模拟环境: 100GB 可用 (通过)",
			FixTip:  "",
		},
	}
}
