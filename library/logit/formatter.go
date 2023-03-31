/*
 * @Author: liziwei01
 * @Date: 2022-03-24 23:27:44
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-03-24 23:42:32
 * @Description: file content
 */
package logit

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[GIN] | %s %3d %s| %13v | %15s |%s %-7s %s %#v",
		// param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		// param.ErrorMessage,
	)
}
