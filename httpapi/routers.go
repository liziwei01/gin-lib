/*
 * @Author: liziwei01
 * @Date: 2022-03-03 16:04:46
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-10-28 14:20:55
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

		logit.Logger.Fine("safe router fine test")
		logit.Logger.Debug("safe router debug test")
		logit.Logger.Trace("safe router trace test")
		logit.Logger.Info("safe router info test")
		logit.Logger.Warn("safe router warn test")
		logit.Logger.Warn("safe router error test")
		logit.Logger.Critical("safe router critical test")
		// panic("safe router panic test")
	})
}
