/*
 * @Author: liziwei01
 * @Date: 2022-06-30 05:14:06
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-07-02 22:05:53
 * @Description: file content
 */
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/liziwei01/gin-lib/library/cookie"
)

var (
	GIN_FILE_DOWNLOAD_COOKIE_NAME = "GFD_email"
)

func CheckLoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if ck, err := c.Cookie(GIN_FILE_DOWNLOAD_COOKIE_NAME); err == nil {
			ck, err = cookie.Decode(GIN_FILE_DOWNLOAD_COOKIE_NAME, ck)
			if err != nil {
				c.AbortWithStatus(401)
				return
			}
			c.Set("email", ck)
			c.Next()
			return
		} else {
			c.Redirect(302, "/gin-lib/user/login")
			// c.AbortWithStatus(401)
			return
		}
	}
}
