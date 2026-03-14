//go:build linux

package envcheck

import (
    "github.com/hyx050923-stack/lobster-manager/internal/envcheck/linux"
)

type Report struct {
    Platform    string   `json:"platform"`
    DockerReady bool     `json:"docker_ready"`
    Issues      []string `json:"issues"`
    FixAdvice   []string `json:"fix_advice"`
}

func Detect() (*Report, error) {
    rep := &Report{Platform: "linux"}
    lrep, err := linux.Detect()
    if err != nil {
        return nil, err
    }
    rep.DockerReady = lrep.DockerInstalled && lrep.DockerServiceActive && lrep.UserInDockerGroup
    rep.FixAdvice = linux.GenerateFixAdvice(lrep)
    if !lrep.DockerInstalled {
        rep.Issues = append(rep.Issues, "Docker未安装")
    } else if !lrep.DockerServiceActive {
        rep.Issues = append(rep.Issues, "Docker服务未运行")
    } else if !lrep.UserInDockerGroup {
        rep.Issues = append(rep.Issues, "用户不在docker组")
    }
    return rep, nil
}