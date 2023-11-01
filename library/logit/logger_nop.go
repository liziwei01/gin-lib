/*
 * @Author: liziwei01
 * @Date: 2023-10-30 12:15:13
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 11:38:16
 * @Description: 空logger
 */
package logit

import "context"

var (
	// NopLogger 黑洞，调用这个logger，什么也不做
	NopLogger Logger = &nopLogger{}
)

// nopLogger do nothing, just implement the interface
type nopLogger struct{}

func (nl *nopLogger) Debug(ctx context.Context, message string, fields ...Field)   {}
func (nl *nopLogger) Trace(ctx context.Context, message string, fields ...Field)   {}
func (nl *nopLogger) Notice(ctx context.Context, message string, fields ...Field)  {}
func (nl *nopLogger) Warning(ctx context.Context, message string, fields ...Field) {}
func (nl *nopLogger) Error(ctx context.Context, message string, fields ...Field)   {}
func (nl *nopLogger) Fatal(ctx context.Context, message string, fields ...Field)   {}
func (nl *nopLogger) Output(ctx context.Context, level Level, callDepth int, message string, fields ...Field) {}

var _ Logger = (*nopLogger)(nil)
