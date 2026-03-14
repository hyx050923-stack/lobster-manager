//go:build windows

package envcheck

import (
    "github.com/hyx050923-stack/lobster-manager/internal/envcheck/windows"
)

type Report struct {
    Platform    string   `json:"platform"`
    DockerReady bool     `json:"docker_ready"`
    Issues      []string `json:"issues"`
    FixAdvice   []string `json:"fix_advice"`
}

func Detect() (*Report, error) {
    rep := &Report{Platform: "windows"}
    wrep, err := windows.Detect()
    if err != nil {
        return nil, err
    }
    rep.DockerReady = wrep.WSL2Enabled && wrep.WSLDistributionExist && wrep.DockerInstalled && wrep.DockerWSL2Backend
    rep.FixAdvice = windows.GenerateFixAdvice(wrep)
    if !wrep.WSL2Enabled {
        rep.Issues = append(rep.Issues, "WSL2未启用")
    }
    if !wrep.WSLDistributionExist {
        rep.Issues = append(rep.Issues, "未安装WSL发行版")
    }
    if !wrep.DockerInstalled {
        rep.Issues = append(rep.Issues, "Docker未安装")
    } else if !wrep.DockerWSL2Backend {
        rep.Issues = append(rep.Issues, "Docker未使用WSL2后端")
    }
    return rep, nil
}