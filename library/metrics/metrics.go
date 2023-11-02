/*
 * @Author: liziwei01
 * @Date: 2023-11-01 19:26:36
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 20:16:24
 * @Description: file content
 */
package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 创建新的 Prometheus metrics，例如计数器
var (
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"path"},
	)
)

func init() {
	// 注册 metrics
	prometheus.MustRegister(TotalRequests)
}

// prometheusHandler 返回一个处理程序，该处理程序调用 promhttp 包中的 HandlerFor
func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
