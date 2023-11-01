/*
 * @Author: liziwei01
 * @Date: 2022-03-04 21:44:14
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 09:50:33
 * @Description: 频控中间件
 */
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	rate "github.com/wallstreetcn/rate/redis"
)

func GetFrequencyControlMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !freqControlConf.Enable {
			// 不限制.
			ctx.Next()
		} else {
			// setup a 1 ops/s rate limiter.
			limiter := rate.NewLimiter(rate.Every(time.Second), 2, "a-sample-operation")
			if limiter.Allow() {
				// serve the user request
			} else {
				// reject the user request
			}
		}
	}
}

func PostFrequencyControlMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !freqControlConf.Enable {
			// 不限制.
			ctx.Next()
		}
	}
}

func MailFrequencyControlMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !freqControlConf.Enable {
			// 不限制.
			ctx.Next()
		}
	}
}
