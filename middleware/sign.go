/*
 * @Author: liziwei01
 * @Date: 2022-03-05 15:45:31
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 09:51:03
 * @Description: file content
 */
package middleware

import (
	"strings"

	"github.com/liziwei01/gin-lib/library/request"
	"github.com/liziwei01/gin-lib/library/response"

	"github.com/gin-gonic/gin"
)

// 走接口签名校验防止接口被刷.
func CheckSignMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if isRealease() != true {
			// 线下无限制.
			ctx.Next()
			return
		} else if !signConf.Enable {
			// 签名校验未开启.
			ctx.Next()
			return
		} else if checkNoSignPath(path) == true {
			// 不需要sign校验的接口.
			ctx.Next()
			return
		} else if !request.CheckSignValid(ctx.Request, signConf.Sign) {
			// sign校验失败.
			response.StdSignCheckFailed(ctx)
			ctx.Abort()
			return
		} else {
			// sign校验成功.
			ctx.Next()
			return
		}
	}
}

// 判断是否是不需要经过md5校验的接口.
func checkNoSignPath(path string) bool {
	for _, preSetPath := range signConf.NoSignPath {
		if strings.Contains(path, preSetPath) {
			return true
		}
	}
	return false
}
