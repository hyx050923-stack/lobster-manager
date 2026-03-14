// internal/envcheck/common/types.go
package common

// CheckResult 单项检测结果
type CheckResult struct {
	Item    string // 检测项：OS/Arch/Storage/Dependency
	Status  bool   // true: 通过, false: 不通过
	Message string // 具体描述，如 "安卓版本: 10"
	FixTip  string // 如果失败，给小白的修复建议
}

// Checker 环境检测接口
type Checker interface {
	CheckAll() []CheckResult
}
