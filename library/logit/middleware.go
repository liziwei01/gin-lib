/*
 * @Author: liziwei01
 * @Date: 2023-09-13 16:55:50
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 10:10:05
 * @Description: 基于gin日志中间件重写，打印每次接口访问的请求信息
 */
package logit

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 创建一个全局的 sync.Pool
var fieldsPool = sync.Pool{
	New: func() interface{} {
		// 这个函数会在获取一个新的对象时调用，如果 Pool 中没有可用的对象
		return make([]Field, 0, 10) // 预分配10个元素的空间
	},
}

// LogitMiddleware instance a Logger middleware
func LogitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Start timer
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		requestID := ctx.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = NewRequestID()
		}
		ctx.Writer.Header().Set("X-Request-ID", requestID)
		SetRequestID(ctx, requestID)

		// Process request
		ctx.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		// 从 Pool 中获取一个 fields 对象
		fields := fieldsPool.Get().([]Field)
		defer fieldsPool.Put(fields) // 确保在函数结束时将 fields 对象放回 Pool

		fields[0] = String("requestID", requestID)
		fields[1] = Int("statusCode", ctx.Writer.Status())
		fields[2] = Duration("latency", time.Now().Sub(start))
		fields[3] = String("ip", ctx.ClientIP())
		fields[4] = String("method", ctx.Request.Method)
		fields[5] = String("path", path)

		err := ctx.Errors.ByType(gin.ErrorTypePrivate).String()
		fields[6] = Error("err", fmt.Errorf(err))

		if err == "" {
			SvrLogger.Notice(ctx, err, fields...)
		} else {
			SvrLogger.Warning(ctx, err, fields...)
		}
	}
}
