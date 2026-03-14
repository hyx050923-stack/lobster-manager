// internal/envcheck/android/checker.go
// +build android

package android

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"claw-butler/internal/envcheck/common" // 假设你的模块名是 claw-butler
)

const (
	MinAndroidVersion = 7.0
	MinStorageGB      = 2.0
)

type AndroidChecker struct {} 

func NewChecker() *AndroidChecker {
	return &AndroidChecker{}
}

// CheckAll 执行所有环境检测
func (c *AndroidChecker) CheckAll() []common.CheckResult {
	results := make([]common.CheckResult, 0)
	
	// 1. 检测架构
	results = append(results, c.checkArch())
	
	// 2. 检测系统版本
	results = append(results, c.checkOSVersion())
	
	// 3. 检测存储空间
	results = append(results, c.checkStorage())

	return results
}

// checkArch 检测 CPU 架构 (ARM/ARM64)
func (c *AndroidChecker) checkArch() common.CheckResult {
	arch := runtime.GOARCH
	res := common.CheckResult{
		Item: "CPU架构",
	}

	if arch == "arm" || arch == "arm64" {
		res.Status = true
		res.Message = "当前架构: " + arch + " (支持)"
	} else {
		res.Status = false
		res.Message = "不支持的架构: " + arch
		res.FixTip = "当前机型不适配，仅支持 ARM 架构设备"
	}
	return res
}

// checkOSVersion 检测安卓版本 (读取 build.prop)
func (c *AndroidChecker) checkOSVersion() common.CheckResult {
	res := common.CheckResult{
		Item: "系统版本",
	}

	// 安卓系统版本信息存储在 /system/build.prop 中
	// key: ro.build.version.release
	versionStr := c.getProp("ro.build.version.release")
	
	if versionStr == "" {
		res.Status = false
		res.Message = "无法读取系统版本"
		res.FixTip = "系统环境异常，请确认是否为标准安卓系统"
		return res
	}

	// 解析版本号 (如 "10", "11", "7.1.1")
	ver, err := strconv.ParseFloat(strings.Split(versionStr, ".")[0], 64)
	if err != nil {
		res.Status = false
		res.Message = "版本解析失败: " + versionStr
		return res
	}

	if ver >= MinAndroidVersion {
		res.Status = true
		res.Message = "安卓 " + versionStr + " (满足要求 ≥7.0)"
	} else {
		res.Status = false
		res.Message = "安卓 " + versionStr + " (版本过低)"
		res.FixTip = "系统版本过低，请升级至安卓 7.0 以上"
	}

	return res
}

// checkStorage 检测存储空间 (需 ≥2GB)
func (c *AndroidChecker) checkStorage() common.CheckResult {
	res := common.CheckResult{
		Item: "存储空间",
	}

	// 获取应用私有目录路径 (无需额外权限)
	// 这通常对应 /data/data/<pkg_name> 或类似路径
	execPath, err := os.Getwd()
	if err != nil {
		execPath = "/data/local/tmp" // 兜底路径
	}

	var stat syscall.Statfs_t
	if err := syscall.Statfs(execPath, &stat); err != nil {
		res.Status = false
		res.Message = "无法读取存储信息"
		res.FixTip = "请授予应用存储权限"
		return res
	}

	// 计算可用空间 (GB)
	// Available = BlocksAvailable * FragmentSize
	availableGB := float64(stat.Bavail*uint64(stat.Bsize)) / 1024 / 1024 / 1024

	if availableGB >= MinStorageGB {
		res.Status = true
		res.Message = "可用空间: " + strconv.FormatFloat(availableGB, 'f', 2, 64) + " GB (充足)"
	} else {
		res.Status = false
		res.Message = "可用空间: " + strconv.FormatFloat(availableGB, 'f', 2, 64) + " GB (不足)"
		res.FixTip = "存储空间不足 2GB，请清理手机垃圾后重试"
	}
	return res
}

// getProp 模拟安卓 getprop 命令，直接读取文件
func (c *AndroidChecker) getProp(key string) string {
	file, err := os.Open("/system/build.prop")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, key+"=") {
			return strings.TrimPrefix(line, key+"=")
		}
	}
	return ""
}
