package envcheck

import (
	"runtime"

	"github.com/hyx050923-stack/lobster-manager/internal/envcheck/linux"
)

type Report struct {
	Platform    string
	DockerReady bool
	Issues      []string
	FixAdvice   []string
}

func Detect() (*Report, error) {
	var rep Report
	rep.Platform = runtime.GOOS

	switch runtime.GOOS {
	case "linux":
		lrep, err := linux.Detect()
		if err != nil {
			return nil, err
		}
		rep.DockerReady = lrep.DockerInstalled && lrep.DockerServiceActive && lrep.UserInDockerGroup
		rep.FixAdvice = linux.GenerateFixAdvice(lrep)
		// 可进一步收集Issues
		if !lrep.DockerInstalled {
			rep.Issues = append(rep.Issues, "Docker未安装")
		} else if !lrep.DockerServiceActive {
			rep.Issues = append(rep.Issues, "Docker服务未运行")
		} else if !lrep.UserInDockerGroup {
			rep.Issues = append(rep.Issues, "用户不在docker组")
		}
	case "windows":
		// TODO: 实现Windows检测
		rep.DockerReady = false
		rep.Issues = []string{"Windows检测尚未实现"}
	default:
		rep.DockerReady = false
		rep.Issues = []string{"不支持的操作系统"}
	}

	return &rep, nil
}