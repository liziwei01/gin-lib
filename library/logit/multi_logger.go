/*
 * @Author: liziwei01
 * @Date: 2023-10-31 19:47:57
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 11:40:46
 * @Description: logger多合一
 */
package logit

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

// MultiLogger 将多个logger转换为一个，实现日志多目标输出；类似io.MultiWriter
func MultiLogger(loggers ...Logger) Logger {
	allLoggers := make([]Logger, 0, len(loggers))
	for _, l := range loggers {
		if ml, ok := l.(*multiLogger); ok {
			allLoggers = append(allLoggers, ml.loggers...)
		} else {
			allLoggers = append(allLoggers, l)
		}
	}
	return &multiLogger{
		loggers: allLoggers,
	}
}

type multiLogger struct {
	loggers []Logger
}

func (m *multiLogger) Debug(ctx context.Context, message string, fields ...Field) {
	m.Output(ctx, DebugLevel, 1, message, fields...)
}

func (m *multiLogger) Trace(ctx context.Context, message string, fields ...Field) {
	m.Output(ctx, TraceLevel, 1, message, fields...)
}

func (m *multiLogger) Notice(ctx context.Context, message string, fields ...Field) {
	m.Output(ctx, NoticeLevel, 1, message, fields...)
}

func (m *multiLogger) Warning(ctx context.Context, message string, fields ...Field) {
	m.Output(ctx, WarningLevel, 1, message, fields...)
}

func (m *multiLogger) Error(ctx context.Context, message string, fields ...Field) {
	m.Output(ctx, ErrorLevel, 1, message, fields...)
}

func (m *multiLogger) Fatal(ctx context.Context, message string, fields ...Field) {
	m.Output(ctx, FatalLevel, 1, message, fields...)
}

// 一次性输出到多个logger
func (m *multiLogger) Output(ctx context.Context, level Level, callDepth int, message string, fields ...Field) {
	for _, l := range m.loggers {
		l.Output(ctx, level, callDepth+1, message, fields...)
	}
}

func (m *multiLogger) Close() error {
	if len(m.loggers) == 0 {
		return nil
	}
	var b strings.Builder

	for idx, lg := range m.loggers {
		if lc, ok := lg.(io.Closer); ok {
			if e := lc.Close(); e != nil {
				b.WriteString(fmt.Sprintf("logger[idx=%d].Close() with error: %s;", idx, e))
			}
		}
	}
	if b.Len() == 0 {
		return nil
	}
	return errors.New(b.String())
}

var _ Logger = (*multiLogger)(nil)
