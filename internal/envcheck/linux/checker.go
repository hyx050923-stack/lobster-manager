package linux

type EnvReport struct {
	DockerInstalled     bool
	DockerVersion       string
	DockerServiceActive bool
	UserInDockerGroup   bool
}

func Detect() (*EnvReport, error) {
	rep := &EnvReport{}
	installed, ver := CheckDockerInstalled()
	rep.DockerInstalled = installed
	rep.DockerVersion = ver

	if installed {
		rep.DockerServiceActive = CheckDockerServiceActive()
		inGroup, _ := CheckUserInDockerGroup()
		rep.UserInDockerGroup = inGroup
	}
	return rep, nil
}

// GenerateFixAdvice 根据报告生成可读的修复建议
func GenerateFixAdvice(report *EnvReport) []string {
	var advice []string
	if !report.DockerInstalled {
		advice = append(advice,
			"Docker未安装。请执行以下命令安装：",
			"curl -fsSL https://get.docker.com | sudo sh",
			"安装后请重新运行检测。")
	} else if !report.DockerServiceActive {
		advice = append(advice,
			"Docker服务未运行。请启动服务：",
			"sudo systemctl start docker",
			"并设置为开机自启：sudo systemctl enable docker")
	}
	if report.DockerInstalled && !report.UserInDockerGroup {
		advice = append(advice,
			"当前用户不在docker组，为避免每次使用sudo，请将用户加入docker组：",
			"sudo usermod -aG docker $USER",
			"然后重新登录或执行 'newgrp docker' 使权限生效。")
	}
	if len(advice) == 0 {
		advice = append(advice, "环境已就绪，无需修复。")
	}
	return advice
}