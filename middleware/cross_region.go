/*
 * @Author: liziwei01
 * @Date: 2022-04-12 21:54:39
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-04-13 23:16:09
 * @Description: file content
 */
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CrossRegionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		ctx.Header("Access-Control-Max-Age", "1728000")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusOK)
			return
		}
		ctx.Next()
	}
}
