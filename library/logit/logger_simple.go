/*
 * @Author: liziwei01
 * @Date: 2023-10-31 20:11:34
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 11:38:36
 * @Description: 提供的基础logger，其他的logger都是基于此logger进行扩展
 */
package logit

import (
	"context"
	"io"

	"github.com/liziwei01/gin-lib/library/extension/pool"
)

// NewSimple 创建一个简单的logger
// 默认使用 text编码
func NewSimple(writer io.Writer) Logger {
	return &SimpleLogger{
		PrefixFunc:  DefaultPrefixFunc,
		Writer:      writer,
		EncoderPool: DefaultTextEncoderPool,

		// 此等级将打印所有的日志
		MinLevel: UnknownLevel,
	}
}

// SimpleLogger 默认的一个logger
type SimpleLogger struct {
	// 日志前缀
	PrefixFunc PrefixFunc

	// 在 Output 执行前执行
	BeforeOutputFunc BeforeOutputFunc

	// 由于Writer 是外部传入的，所以在SimpleLogger 内部，不能对其close
	Writer io.Writer

	// 编码器池
	EncoderPool EncoderPool

	// 最小日志等级，低于此等级的日志信息将不打印
	MinLevel Level
}

// Debug debug
func (sl *SimpleLogger) Debug(ctx context.Context, message string, fields ...Field) {
	sl.Output(ctx, DebugLevel, 1, message, fields...)
}

// Trace Trace
func (sl *SimpleLogger) Trace(ctx context.Context, message string, fields ...Field) {
	sl.Output(ctx, TraceLevel, 1, message, fields...)
}

// Notice Notice
func (sl *SimpleLogger) Notice(ctx context.Context, message string, fields ...Field) {
	sl.Output(ctx, NoticeLevel, 1, message, fields...)
}

// Warning Warning
func (sl *SimpleLogger) Warning(ctx context.Context, message string, fields ...Field) {
	sl.Output(ctx, WarningLevel, 1, message, fields...)
}

// Error Error
func (sl *SimpleLogger) Error(ctx context.Context, message string, fields ...Field) {
	sl.Output(ctx, ErrorLevel, 1, message, fields...)
}

// Fatal Fatal
func (sl *SimpleLogger) Fatal(ctx context.Context, message string, fields ...Field) {
	sl.Output(ctx, FatalLevel, 1, message, fields...)
}

// Output Output
func (sl *SimpleLogger) Output(ctx context.Context, level Level, callDepth int, message string, fields ...Field) {
	if level == UnknownLevel || level >= AllLevels {
		return
	}

	if sl.MinLevel > level || sl.MinLevel >= AllLevels {
		return
	}

	// 复用编码器，减少内存分配
	enc := sl.EncoderPool.Get()
	defer sl.EncoderPool.Put(enc)

	if sl.BeforeOutputFunc != nil {
		sl.BeforeOutputFunc(ctx, enc, level, callDepth+2)
	}

	// 字段的顺序
	// 1：ctx里存储的字段  2：全局meta fields  3：传入的临时补充字段

	// 10是meta fields数量的最大值
	fkv := make(map[string]Field, len(fields)+10)

	var metaFields []Field
	RangeMetaFields(ctx, func(f Field) error {
		metaFields = append(metaFields, f)
		fkv[f.Key()] = f
		return nil
	})

	// 如果有重复的key，后面的会覆盖前面的，如requestID
	for _, f := range fields {
		fkv[f.Key()] = f
	}

	// 开始编码
	// ctx 里存储的 fields 优先级 第一
	// 对每一个ctx里面存储的field，都应用一次传入的func，其目的是为了将field格式化写入到编码器中，形成类似key[value]的形式
	Range(ctx, func(f Field) error {
		// 只有在当前日志等级下，才会输出
		if f.Level().Is(level) {
			// 若字段之前在ctx，后面又在fields 里出现，则使用后面传入的
			if fn, has := fkv[f.Key()]; has {
				fn.AddTo(enc)
				// 不需要再打印第二次
				delete(fkv, f.Key())
			} else {
				f.AddTo(enc)
			}
		}
		return nil
	})

	if len(fkv) > 0 {
		// meta fields 优先级第二
		if len(metaFields) > 0 {
			for _, f := range metaFields {
				if lastField, has := fkv[f.Key()]; has {
					lastField.AddTo(enc)
					delete(fkv, f.Key())
				}
			}
		}

		// 补充字段 优先级第三
		for _, f := range fields {
			if lastField, has := fkv[f.Key()]; has {
				lastField.AddTo(enc)
			}
		}
	}

	// 最后添加message字段，形如 message[xxx]
	// 相当于 Field: String("message", message).AddTo(enc)
	enc.AddString("message", message)

	// 复用bytes.Buffer，减少内存分配
	buf := bpSL.Get()

	// 添加前缀，形如：NOTICE: 2023-11-01 07:38:36 /Users/liziwei01/Desktop/OpenSource/github.com/gin-lib/library/logit/logit_test.go:30
	if sl.PrefixFunc != nil {
		prefix := sl.PrefixFunc(ctx, level, callDepth+2)
		if len(prefix) > 0 {
			_, _ = buf.Write(prefix)
		}
	}
	// 添加日志内容，形如：keyT[valueT] keyN[valueN] keyW[valueW] keyE[valueE] message[test error]
	// 不应该出错，所以忽略返回值
	_, _ = enc.WriteTo(buf)

	
	logLine := make([]byte, buf.Len())
	copy(logLine, buf.Bytes())
	// 输出日志，一般是写入文件，sl.Writer是一个对应文件的io.Writer
	_, _ = sl.Writer.Write(logLine)

	bpSL.Put(buf)
}

var bpSL = pool.NewBytesPool()

var _ Logger = (*SimpleLogger)(nil)
