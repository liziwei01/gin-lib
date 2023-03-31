/*
 * @Author: liziwei01
 * @Date: 2022-03-24 23:28:35
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-06-27 18:50:48
 * @Description: file content
 */
package logit

import (
	"time"

	"github.com/gin-gonic/gin"
)

// LogitMiddleware instance a Logger middleware with baidu/gin-lib/log.
func LogitMiddleware() gin.HandlerFunc {
	formatter := defaultLogFormatter

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path

		Logger.Info(formatter(param))
		if param.ErrorMessage != "" {
			Logger.Error(param.ErrorMessage)
		}
	}
}
