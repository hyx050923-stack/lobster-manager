package main

import (
	"fmt"
	"log"

	"github.com/hyx050923-stack/lobster-manager/internal/envcheck"
)

func main() {
	fmt.Println("龙虾管家 - 环境检测工具")
	rep, err := envcheck.Detect()
	if err != nil {
		log.Fatalf("检测失败: %v", err)
	}
	fmt.Printf("平台: %s\n", rep.Platform)
	fmt.Printf("Docker就绪: %v\n", rep.DockerReady)
	if len(rep.Issues) > 0 {
		fmt.Println("发现问题:")
		for _, issue := range rep.Issues {
			fmt.Printf(" - %s\n", issue)
		}
		fmt.Println("修复建议:")
		for _, advice := range rep.FixAdvice {
			fmt.Printf(" - %s\n", advice)
		}
	} else {
		fmt.Println("环境一切正常！")
	}
}