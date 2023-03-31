/*
 * @Author: liziwei01
 * @Date: 2022-03-04 23:32:56
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-03-10 20:52:39
 * @Description: 接口token校验
 */

package middleware

import (
	"strings"

	"github.com/liziwei01/gin-lib/library/env"
	"github.com/liziwei01/gin-lib/library/response"
	"github.com/liziwei01/gin-lib/library/utils"

	"github.com/gin-gonic/gin"
)

// 走接口token校验防止后台get接口被刷.
func CheckTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		inputToken, ok := utils.Request.Header(c.Request, "token")
		if isRealease() != true {
			// 线下无限制.
			c.Next()
		} else if !tokenConf.Enable {
			// token校验未开启.
			c.Next()
		} else if checkNoTokenPath(path) == true {
			// 不需要token校验的接口.
			c.Next()
		} else if !ok || tokenConf.Token != inputToken {
			// token校验失败.
			response.StdTokenCheckFailed(c)
			c.Abort()
		} else {
			// token校验成功.
			c.Next()
		}
	}
}

// 判断是否是不需要经过token校验的接口.
func checkNoTokenPath(path string) bool {
	for _, preSetPath := range tokenConf.NoTokenPath {
		if strings.Contains(path, preSetPath) {
			return true
		}
	}
	return false
}

// 判断是否为线上环境.
func isRealease() bool {
	return env.RunMode() == env.RunModeRelease
}
