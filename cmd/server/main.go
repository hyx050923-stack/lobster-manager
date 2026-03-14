package main

import (
	"context" // 【新增】导入 context 包
	"log"
	"net/http"

	"github.com/hyx050923-stack/lobster-manager/internal/container"
	"github.com/hyx050923-stack/lobster-manager/internal/deploy"
	"github.com/hyx050923-stack/lobster-manager/internal/envcheck/android"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化依赖组件
	ctx := context.Background() // 【新增】创建根 Context

	// 初始化容器管理器 (传入 ctx)
	containerMgr := container.NewManager(ctx)
	
	// 初始化安卓环境检测器
	envChecker := android.NewChecker()
	
	// 初始化部署服务 (注入依赖)
	deploySvc := deploy.NewService(envChecker, containerMgr)

	// 2. 启动 HTTP 服务
	r := gin.Default()

	// 接口：一键部署
	r.POST("/api/deploy", func(c *gin.Context) {
		result := deploySvc.Deploy()
		
		if result.Success {
			c.JSON(http.StatusOK, result)
		} else {
			c.JSON(http.StatusOK, result)
		}
	})

	// 接口：环境检测
	r.GET("/api/envcheck", func(c *gin.Context) {
		results := envChecker.CheckAll()
		c.JSON(http.StatusOK, results)
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	log.Println("Server starting on :38080")
	r.Run(":38080")
}
