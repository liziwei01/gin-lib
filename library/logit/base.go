/*
 * @Author: liziwei01
 * @Date: 2022-03-04 15:40:52
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-10-31 14:43:00
 * @Description: 提供默认Service Logger
 */
package logit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/liziwei01/gin-lib/library/env"
)

const (
	// logPath log 配置文件路径
	logPath = "logit/"
	suffix  = ".toml"
	svrName = "service"
)

var (
	// conf file root path
	configPath = env.Default.ConfDir()
	loggers    map[string]Logger
	// SvrLogger 默认service log
	SvrLogger Logger
	initMux sync.Mutex
)

func init() {
	loggers = make(map[string]Logger)
}

// SetServiceLogger 设置service logger 需要在程序启动初始化调用，最先调用
func SetServiceLogger(ctx context.Context) error {
	var err error
	SvrLogger, err = GetLogger(ctx, svrName)
	if err != nil {
		return err
	}
	return nil
}

// GetLogger 获取 logger
func GetLogger(ctx context.Context, logName string) (Logger, error) {
	// 先尝试获取
	if client, hasSet := loggers[logName]; hasSet {
		if client != nil {
			return client, nil
		}
	}
	// 没有则重新设置
	client, err := SetLogger(ctx, logName)
	if client != nil {
		return client, nil
	}
	return nil, err
}

// SetLogger 设置 logger
func SetLogger(ctx context.Context, logName string) (Logger, error) {
	// 互斥锁
	initMux.Lock()
	defer initMux.Unlock()
	// 初始化
	logger, err := initLogger(ctx, logName)
	if err == nil {
		// 添加
		loggers[logName] = logger
		return logger, nil
	}
	// 抛异常
	// log error
	return nil, err
}

// initLogger 初始化日志
func initLogger(ctx context.Context, logName string) (Logger, error) {
	fileAbs, err := filepath.Abs(filepath.Join(configPath, logPath, logName+suffix))
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(fileAbs); !os.IsNotExist(err) {
		Logger, err := NewLogger(ctx, OptConfigFile(filepath.Join(logPath, logName+suffix)))
		if err != nil {
			return nil, err
		}
		return Logger, nil
	}
	return nil, fmt.Errorf("log conf not exist")
}
