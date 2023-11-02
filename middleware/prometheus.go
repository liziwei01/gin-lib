/*
 * @Author: liziwei01
 * @Date: 2023-11-01 20:15:58
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 20:15:59
 * @Description: file content
 */
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/liziwei01/gin-lib/library/metrics"
)

// prometheusMiddleware 更新 metrics 和相应的计数器
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics.TotalRequests.WithLabelValues(c.Request.URL.Path).Inc()
		c.Next()
	}
}
