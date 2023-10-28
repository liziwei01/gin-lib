/*
 * @Author: liziwei01
 * @Date: 2023-10-28 12:26:10
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-10-28 13:20:13
 * @Description: file content
 */
package request

import (
	"github.com/gin-gonic/gin"
	"github.com/liziwei01/gin-lib/library/utils"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = utils.UUID.GenUUID()
		}
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Set("requestID", requestID)
		c.Next()
	}
}
