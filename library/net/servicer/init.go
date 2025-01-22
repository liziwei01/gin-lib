/*
 * @Author: liziwei01
 * @Date: 2023-11-03 13:38:42
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 23:05:41
 * @Description: 载入日志组件。ral.InitDefault会调用InitDefault，业务只需要调用ral.InitDefault即可
 */
package servicer

import (
	"context"
	"fmt"

	"github.com/liziwei01/gin-lib/library/conf"
	"github.com/liziwei01/gin-lib/library/logit"
)

// DefaultLoggerConfNames 默认 ral 日志配置文件
var DefaultLoggerConfNames = []string{
	"logit/ral.toml",
}

// DefaultWorkerLoggerConfNames 默认 ral-worker 日志的配置文件
var DefaultWorkerLoggerConfNames = []string{
	"logit/ral-worker.toml",
}

// DefaultLogger 用于打印servicer组件自身运行信息的logger
// 如服务启动、停止、服务信息更新（服务发现）
var DefaultLogger logit.Logger

// DefaultWorkerLogger 用于打印使用servier过程中，发送请求执行详情的logger
//
// 如使用ral发送一个http请求后，记录执行耗时等信息
// 如使用redis、mysql client 发送请求后，记录执行耗时、状态等信息
var DefaultWorkerLogger logit.Logger

// InitDefault 默认的初始化
//
// 加载 DefaultLogger DefaultWorkerLogger DefaultMapper
// 之后还需要调用 servicer.LoadWithFileGlob 方法加载服务
func InitDefault(ctx context.Context) error {
	if DefaultLogger == nil {
		var err error
		for _, name := range DefaultLoggerConfNames {
			if conf.Exists(name) {
				DefaultLogger, err = logit.NewLogger(ctx, logit.OptConfigFile(name))
				break
			}
		}
		if err != nil {
			return fmt.Errorf("init ral logger failed: %w", err)
		}
	}

	if DefaultWorkerLogger == nil {
		var err error
		for _, name := range DefaultWorkerLoggerConfNames {
			if conf.Exists(name) {
				DefaultWorkerLogger, err = logit.NewLogger(ctx, logit.OptConfigFile(name))
				break
			}
		}
		if err != nil {
			return fmt.Errorf("init ral-worker logger failed: %w", err)
		}
	}

	if DefaultMapper == nil {
		DefaultMapper = NewMapper()
	}
	return nil
}
