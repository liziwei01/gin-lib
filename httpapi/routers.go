/*
 * @Author: liziwei01
 * @Date: 2022-03-03 16:04:46
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 09:11:02
 * @Description: 路由分发
 */

package httpapi

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/liziwei01/gin-lib/library/logit"

	// libRouters "github.com/liziwei01/gin-lib/modules/mod1/routers"

	"github.com/gin-gonic/gin"
)

/**
 * @description: start http server and start listening
 * @param {*}
 * @return {*}
 */
func InitRouters(handler *gin.Engine) {
	// 暂时解决跨域问题
	handler.Use(cors.Default())
	// handler.Use(middleware.CrossRegionMiddleware())
	// init routers
	router := handler.Group("/")
	// libRouters.Init(router)

	// safe router
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello! THis is gin-lib. Welcome!")

		logit.SvrLogger.Debug(ctx, "safe router debug test", logit.String("fieldKey", "fieldValue"))
		logit.SvrLogger.Trace(ctx, "safe router trace test", logit.String("fieldKey", "fieldValue"))
		logit.SvrLogger.Notice(ctx, "safe router notice test", logit.String("fieldKey", "fieldValue"))
		logit.SvrLogger.Warning(ctx, "safe router warning test", logit.String("fieldKey", "fieldValue"))
		logit.SvrLogger.Error(ctx, "safe router error test", logit.String("fieldKey", "fieldValue"))
		logit.SvrLogger.Fatal(ctx, "safe router fatal test", logit.String("fieldKey", "fieldValue"))
		panic("safe router panic test")
	})
}
