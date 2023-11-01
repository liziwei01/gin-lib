/*
 * @Author: liziwei01
 * @Date: 2023-10-28 12:26:10
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 09:48:50
 * @Description: file content
 */
package request

import (
	"github.com/gin-gonic/gin"
	"github.com/liziwei01/gin-lib/library/logit"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = logit.NewRequestID()
		}
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Set("requestID", requestID)
		c.Next()
	}
}
