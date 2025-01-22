/*
 * @Author: liziwei01
 * @Date: 2023-09-13 16:55:50
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 04:03:20
 * @Description: 基于gin日志中间件重写，打印每次接口访问的请求信息
 */
package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liziwei01/gin-lib/library/logit"
)

// 创建一个全局的 sync.Pool
var fieldsPool = sync.Pool{
	New: func() interface{} {
		// 这个函数会在获取一个新的对象时调用，如果 Pool 中没有可用的对象
		return make([]logit.Field, 7, 7) // 分配7个元素的空间
	},
}

// GinLoggerMiddleware instance a Logger middleware
func GinLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(logit.WithContext(ctx.Request.Context()))
		// Start timer
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		requestID := ctx.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = logit.NewRequestID()
		}
		ctx.Writer.Header().Set("X-Request-ID", requestID)
		logit.SetRequestID(ctx, requestID)

		// Process request
		ctx.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		// 从 Pool 中获取一个 fields 对象
		fields := fieldsPool.Get().([]logit.Field)
		defer fieldsPool.Put(fields) // 确保在函数结束时将 fields 对象放回 Pool

		fields[0] = logit.String("requestID", requestID)
		fields[1] = logit.Int("statusCode", ctx.Writer.Status())
		fields[2] = logit.Duration("latency", time.Now().Sub(start))
		fields[3] = logit.String("ip", ctx.ClientIP())
		fields[4] = logit.String("method", ctx.Request.Method)
		fields[5] = logit.String("path", path)

		err := ctx.Errors.ByType(gin.ErrorTypePrivate).String()
		fields[6] = logit.Error("err", fmt.Errorf(err))

		if err == "" {
			logit.SrvLogger.Notice(ctx, err, fields...)
		} else {
			logit.SrvLogger.Warning(ctx, err, fields...)
		}
	}
}
