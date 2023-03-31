/*
 * @Author: liziwei01
 * @Date: 2022-03-05 15:45:31
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-06-30 05:22:13
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
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if isRealease() != true {
			// 线下无限制.
			c.Next()
			return
		} else if !signConf.Enable {
			// 签名校验未开启.
			c.Next()
			return
		} else if checkNoSignPath(path) == true {
			// 不需要sign校验的接口.
			c.Next()
			return
		} else if !request.CheckSignValid(c.Request, signConf.Sign) {
			// sign校验失败.
			response.StdSignCheckFailed(c)
			c.Abort()
			return
		} else {
			// sign校验成功.
			c.Next()
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
