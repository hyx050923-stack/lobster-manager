package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyx050923-stack/lobster-manager/internal/deploy"
	"github.com/hyx050923-stack/lobster-manager/internal/envcheck"
)

func main() {
	r := gin.Default()

	// 环境检测接口
	r.GET("/api/envcheck", func(c *gin.Context) {
		rep, err := envcheck.Detect()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rep)
	})

	// 自动修复接口
	r.POST("/api/fix", func(c *gin.Context) {
		result, err := envcheck.Fix()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	// 一键部署接口
	r.POST("/api/deploy", func(c *gin.Context) {
		result := deploy.Deploy()
		if result.Success {
			c.JSON(http.StatusOK, result)
		} else {
			c.JSON(http.StatusInternalServerError, result)
		}
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	log.Println("Server starting on :8080")
	r.Run(":8080")
}