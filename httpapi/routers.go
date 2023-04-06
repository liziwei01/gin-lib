/*
 * @Author: liziwei01
 * @Date: 2022-03-03 16:04:46
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-04-06 17:15:22
 * @Description: 路由分发
 */

package httpapi

import (
	"net/http"

	"github.com/liziwei01/gin-lib/library/logit"
	"github.com/liziwei01/gin-lib/middleware"
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
	handler.Use(middleware.CrossRegionMiddleware())
	// router.Use(middleware.CheckTokenMiddleware(), middleware.GetFrequencyControlMiddleware(), middleware.PostFrequencyControlMiddleware(), middleware.MailFrequencyControlMiddleware())
	// init routers
	router := handler.Group("/")
	// libRouters.Init(router)

	// safe router
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello! THis is gin-lib. Welcome!")
		logit.Logger.Info("safe router test")
		logit.Logger.Error("safe router error test")
	})
}
