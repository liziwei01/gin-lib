/*
 * @Author: liziwei01
 * @Date: 2022-03-04 22:06:10
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 20:52:06
 * @Description: file content
 */
package bootstrap

import (
	"context"
	"path/filepath"

	"github.com/liziwei01/gin-lib/library/env"
	"github.com/liziwei01/gin-lib/library/logit"
	"github.com/liziwei01/gin-lib/library/net/servicer"
	"github.com/liziwei01/gin-lib/middleware"

	"github.com/gin-gonic/gin"
)

func InitMust(ctx context.Context) {
	InitServiceLogger(ctx)
	InitMiddleware(ctx)
}

// 默认业务日志单独初始化
func InitServiceLogger(ctx context.Context) {
	logit.SetServiceLogger(ctx)
}

// 加载mysql，redis等服务配置
func InitServicer(ctx context.Context) {
	pattern := filepath.Join(env.ConfDir(), "servicer", "*.toml")
	servicer.MustLoad(ctx, servicer.LoadOptFilesGlob(pattern, false))
}

func InitMiddleware(ctx context.Context) {
	middleware.Init(ctx)
}

// InitHandler 用*gin.Engine作http handler
func InitHandler(app *AppServer) *gin.Engine {
	gin.SetMode(app.Config.RunMode)
	handler := gin.New()
	handler.ContextWithFallback = true
	// 注册log recover中间件
	ginRecovery := gin.Recovery()
	libLogger := middleware.GinLoggerMiddleware()
	// ginLogger := gin.Logger()
	handler.Use(ginRecovery)
	handler.Use(libLogger)
	// handler.Use(ginLogger)
	return handler
}
