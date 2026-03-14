//go:build windows

package windows

type EnvReport struct {
	WSL2Enabled          bool
	WSLDistributionExist bool
	DockerInstalled      bool
	DockerWSL2Backend    bool // 是否使用WSL2后端（推荐）
}

func Detect() (*EnvReport, error) {
	rep := &EnvReport{}
	wsl2, _ := CheckWSL2Enabled()
	rep.WSL2Enabled = wsl2
	rep.WSLDistributionExist = CheckWSLDistributionExists()
	rep.DockerInstalled = CheckDockerInstalled()
	if rep.DockerInstalled {
		rep.DockerWSL2Backend = CheckDockerWSL2Backend()
	}
	return rep, nil
}

// GenerateFixAdvice 生成修复建议
func GenerateFixAdvice(report *EnvReport) []string {
	var advice []string
	if !report.WSL2Enabled {
		advice = append(advice,
			"WSL2功能未启用。请以管理员身份运行PowerShell，执行以下命令：",
			"dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart",
			"dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart",
			"然后重启电脑。")
	}
	if !report.WSLDistributionExist {
		advice = append(advice,
			"未安装WSL Linux发行版。请运行以下命令安装默认发行版（如Ubuntu）：",
			"wsl --install -d Ubuntu",
			"或从Microsoft Store安装。")
	}
	if !report.DockerInstalled {
		advice = append(advice,
			"Docker Desktop未安装。请从 https://www.docker.com/products/docker-desktop/ 下载安装。",
			"安装时请确保选择WSL2后端。")
	} else if !report.DockerWSL2Backend {
		advice = append(advice,
			"Docker当前未使用WSL2后端，建议切换到WSL2以获得更好性能和集成。",
			"可在Docker Desktop设置中启用：Settings -> General -> Use WSL 2 based engine。")
	}
	if len(advice) == 0 {
		advice = append(advice, "Windows环境已就绪。")
	}
	return advice
}