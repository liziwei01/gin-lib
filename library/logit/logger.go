/*
 * @Author: liziwei01
 * @Date: 2023-10-30 12:10:40
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-10-31 21:38:39
 * @Description: 日志接口
 */

package logit

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

// Logger 接口定义
type Logger interface {
	// 正常使用 fields 打印日志的方式
	Debug(ctx context.Context, message string, fields ...Field)
	Trace(ctx context.Context, message string, fields ...Field)
	Notice(ctx context.Context, message string, fields ...Field)
	Warning(ctx context.Context, message string, fields ...Field)
	Error(ctx context.Context, message string, fields ...Field)
	Fatal(ctx context.Context, message string, fields ...Field)

	Output(ctx context.Context, level Level, callDepth int, message string, fields ...Field)
}

// NewLogger 创建一个新的logger
//
//	ctx 用于控制logger的writer的生命周期
//	如程序要退出，可将ctx 取消，如此writer将进行Close，可避免日志丢失
func NewLogger(ctx context.Context, opts ...Option) (l Logger, errResult error) {
	cfg, err := loggerCfg(opts...)
	if err != nil {
		return nil, err
	}
	nop := &nopLogger{}

	if cfg.nopLog() {
		return nop, nil
	}

	mapper := make(map[Level]Logger, 6) // 默认是6个日志等级
	closeFns := make([]func() error, 0, 6)

	var closeWritersFunc func() error

	if ctx != nil {
		closeWritersFunc = func() error {
			var builder strings.Builder
			for idx, fn := range closeFns {
				if e := fn(); e != nil {
					builder.WriteString(fmt.Sprintf("idx=%d error=%s;", idx, e))
				}
			}
			if builder.Len() == 0 {
				return nil
			}
			return fmt.Errorf("logger close with errors: %s", builder.String())
		}

		closeAll := func() {
			if e := closeWritersFunc(); e != nil {
				fmt.Fprintf(os.Stderr, "%s %v\n", time.Now(), e)
			}
		}

		defer func() {
			if errResult != nil {
				closeAll()
			}
		}()

		// 如果ctx 被取消，那么关闭所有的writer
		go func() {
			<-ctx.Done()
			closeAll()
		}()
	}

	// 每个日志分发规则对应一种文件后缀，每种文件后缀对应一个writer
	for idx, item := range cfg.Dispatch {
		if len(item.Levels) == 0 {
			continue
		}

		itemOpt := *cfg
		itemOpt.FileName += item.FileSuffix

		// 每个文件后缀新建一个writer
		awc, err := itemOpt.getWriter()

		if err != nil {
			return nil, fmt.Errorf("init logger (%d)(%q) failed: %w", idx, itemOpt.FileName, err)
		}
		lg := &SimpleLogger{
			PrefixFunc:       cfg.PrefixFunc,
			EncoderPool:      cfg.encoderPool,
			BeforeOutputFunc: cfg.BeforeOutputFunc,
			Writer:           awc,
		}
		closeFns = append(closeFns, awc.Close)

		// 某文件后缀所需要写入的所有日志等级，如：TRACE, NOTICE，按照日志等级归类为一个logger存入mapper
		for _, l := range item.Levels {
			if _, has := mapper[l]; !has {
				mapper[l] = MultiLogger(lg)
			} else {
				mapper[l] = MultiLogger(mapper[l], lg)
			}
		}
	}

	// 新的调度器，用于包含记录了所有logger的mapper，根据日志等级，返回对应的logger
	dl := newDispatcher(func(level Level) Logger {
		logger, has := mapper[level]
		if !has {
			return nop
		}
		return logger
	})

	dl.closeFunc = closeWritersFunc

	// 比如，调用Warning，某两个文件后缀都需要写入Warning级别的日志，那么会通过MultiLogger调用两次Output写入文件
	return dl, nil
}
