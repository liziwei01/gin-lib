/*
 * @Author: liziwei01
 * @Date: 2022-03-24 23:28:35
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-10-28 13:43:48
 * @Description: 日志中间件打印每次接口访问的请求信息，重写gin的日志格式供中间件使用
 */
package logit

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LogitMiddleware instance a Logger middleware with baidu/go-lib/log/log4go.
func LogitMiddleware() gin.HandlerFunc {
	formatter := func(param gin.LogFormatterParams) string {
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}
		return fmt.Sprintf("[GIN] [requestID]=%d, [code]=%3d, [latency]=%v, [ip]=%s, [method]=%s, [path]=%#v, [err]=%s",
			param.Keys["requestID"],
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	}

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

		if param.ErrorMessage == "" {
			Logger.Info(formatter(param))
		} else if param.ErrorMessage != "" {
			Logger.Error(formatter(param))
		}
	}
}
