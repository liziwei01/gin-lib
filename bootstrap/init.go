/*
 * @Author: liziwei01
 * @Date: 2022-03-04 22:06:10
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-04-06 17:05:32
 * @Description: file content
 */
package bootstrap

import (
	"context"

	"github.com/liziwei01/gin-lib/library/logit"
	"github.com/liziwei01/gin-lib/middleware"

	"github.com/gin-gonic/gin"
)

func InitMust(ctx context.Context) {
	InitLog(ctx)
	InitMiddleware(ctx)
}

func InitLog(ctx context.Context) {
	logit.Init(ctx, "gin-lib")
}

func InitMiddleware(ctx context.Context) {
	middleware.Init(ctx)
}

// InitHandler 获取http handler
func InitHandler(app *AppServer) *gin.Engine {
	gin.SetMode(app.Config.RunMode)
	handler := gin.New()
	// 注册log recover中间件
	ginRecovery := gin.Recovery()
	baiduLogger := logit.LogitMiddleware()
	handler.Use(ginRecovery, baiduLogger)
	return handler
}
