/*
 * @Author: liziwei01
 * @Date: 2022-03-03 16:04:46
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-05 16:41:06
 * @Description: 路由分发
 */

package httpapi

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/liziwei01/gin-lib/library/logit"
	"github.com/liziwei01/gin-lib/library/metrics"
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
	handler.Use(cors.Default())
	handler.Use(middleware.PrometheusMiddleware())
	// handler.Use(middleware.CrossRegionMiddleware())
	// init routers
	router := handler.Group("/")
	// libRouters.Init(router)
	router.GET("metrics", metrics.PrometheusHandler())
	router.POST("postChat", func(ctx *gin.Context) {
		logit.SrvLogger.Notice(ctx, "postChat Notice")
		// read the request body and send it back
		body, err := ctx.GetRawData()
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		logit.SrvLogger.Notice(ctx, "postChat", logit.String("body", string(body)))
		ctx.Data(http.StatusOK, "application/json", body)
	})

	// safe router
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello! THis is gin-lib. Welcome!")

		logit.SrvLogger.Debug(ctx, "safe router debug test", logit.String("fieldKey", "fieldValue"))
		logit.SrvLogger.Trace(ctx, "safe router trace test", logit.String("fieldKey", "fieldValue"))
		logit.SrvLogger.Notice(ctx, "safe router notice test", logit.String("fieldKey", "fieldValue"))
		logit.SrvLogger.Warning(ctx, "safe router warning test", logit.String("fieldKey", "fieldValue"))
		logit.SrvLogger.Error(ctx, "safe router error test", logit.String("fieldKey", "fieldValue"))
		logit.SrvLogger.Fatal(ctx, "safe router fatal test", logit.String("fieldKey", "fieldValue"))
		// panic("safe router panic test")
	})
}
