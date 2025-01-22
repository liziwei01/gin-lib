/*
 * @Author: liziwei01
 * @Date: 2022-03-04 15:40:52
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 20:51:34
 * @Description: 提供默认业务日志Service Logger
 */
package logit

import (
	"context"
)

var (
	// SrvLogger 默认的Service Logger，用于打印业务日志
	SrvLogger Logger
)

// SetServiceLogger 设置service logger 需要在程序启动时初始化
func SetServiceLogger(ctx context.Context) error {
	var err error
	SrvLogger, err = NewLogger(ctx, OptConfigFile("logit/service.toml"))
	if err != nil {
		return err
	}
	return nil
}
