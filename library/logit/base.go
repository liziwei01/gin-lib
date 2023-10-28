/*
 * @Author: liziwei01
 * @Date: 2022-03-04 15:40:52
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-10-28 09:39:15
 * @Description: 使用百度的log库
 */
package logit

import (
	"context"

	lib "github.com/baidu/go-lib/log"
	"github.com/baidu/go-lib/log/log4go"
	"github.com/liziwei01/gin-lib/library/env"
)

var (
	Logger      log4go.Logger  // 日志对象
	levelStr    = "INFO"       // levelStr以上级别的日志都会打印到stdout
	logDir      = env.LogDir() // 日志文件存放目录
	hasStdOut   = true         // 是否打印到stdout
	when        = "H"          // 每小时生成一个日志文件
	backupCount = 5            // 保留5个日志文件
)

/**
 * @description: all the log are recorded under ./log
 * @param {string} programName
 * @return {*}
 */
func Init(ctx context.Context, programName string) error {
	var err error
	Logger, err = lib.Create(programName, levelStr, logDir, hasStdOut, when, backupCount)
	return err
}
