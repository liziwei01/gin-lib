/*
 * @Author: liziwei01
 * @Date: 2022-03-04 22:06:10
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 09:43:15
 * @Description: file content
 */
package bootstrap

import (
	"context"

	"github.com/liziwei01/gin-lib/library/logit"
	"github.com/liziwei01/gin-lib/library/request"
	"github.com/liziwei01/gin-lib/middleware"

	"github.com/gin-gonic/gin"
)

func InitMust(ctx context.Context) {
	InitLog(ctx)
	InitMiddleware(ctx)
}

func InitLog(ctx context.Context) {
	logit.SetServiceLogger(ctx)
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
	idGenerator := request.RequestIDMiddleware()
	libLogger := logit.LogitMiddleware()
	handler.Use(ginRecovery, idGenerator, libLogger)
	return handler
}
