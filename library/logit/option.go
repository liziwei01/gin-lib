/*
 * @Author: liziwei01
 * @Date: 2023-10-30 12:14:46
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-10-31 19:45:42
 * @Description: 用函数选项模式的方式，对logit进行配置
 */
package logit

import (
	"errors"
	"fmt"
	"io"
	"time"
)

// DefaultLogger 默认logger:内容输出到黑洞
var DefaultLogger Logger = &nopLogger{}

// Option FileLogger 的配置选项
// 使用接口隐藏func类型，避免外部直接调用
type Option interface {
	apply(*Config)
}

// 实现了Option接口
type funcOption struct {
	f func(*Config)
}

func (fdo *funcOption) apply(do *Config) {
	fdo.f(do)
}

// 将传入的func隐藏到funcOption中，用于内部调用
func newFuncOption(f func(*Config)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// 返回一个Option，用于设置错误
func errOption(err error) Option {
	return newFuncOption(func(config *Config) {
		if err != nil {
			config.err = err
		}
	})
}

// OptConfig 配置选项-整体配置
func OptConfig(c *Config) Option {
	return newFuncOption(func(config *Config) {
		*config = *c
	})
}

// OptConfigFile 配置选项-从文件读取配置
func OptConfigFile(confName string) Option {
	c, err := LoadConfig(confName)
	if err != nil {
		return errOption(err)
	}
	return OptConfig(c)
}

// OptSetConfigFn 配置选型-直接对Config进行修改
func OptSetConfigFn(fn func(c *Config)) Option {
	return newFuncOption(func(config *Config) {
		fn(config)
	})
}

// OptLogFileName 配置选项-设置日志文件名
//
// 如 OptLogFileName("ral/ral-worker.log")
func OptLogFileName(name string) Option {
	if name == "" {
		return errOption(errors.New("optLogFileName with empty name"))
	}
	return newFuncOption(func(config *Config) {
		config.FileName = name
	})
}

// OptRotateRule 配置选项-设置日志切分规则
//
//	如 1hour,1day,no,默认为1hour
func OptRotateRule(rule string) Option {
	if rule == "" {
		return errOption(errors.New("optRotateRule with empty rule"))
	}
	return newFuncOption(func(config *Config) {
		config.RotateRule = rule
	})
}

// OptMaxFileNum 配置选项-设置日志文件保留数
//
// 保留最多日志文件数，默认(当值为0时)为48，若<0,则不会清理
func OptMaxFileNum(num int) Option {
	return newFuncOption(func(config *Config) {
		config.MaxFileNum = num
	})
}

// OptBufferSize 配置选项-设置writer 的 BufferSize
func OptBufferSize(size int) Option {
	return newFuncOption(func(config *Config) {
		config.BufferSize = size
	})
}

// OptWriterTimeout 配置选项-设置writer 的 超时时间
func OptWriterTimeout(timeout time.Duration) Option {
	return newFuncOption(func(config *Config) {
		config.WriterTimeout = int(timeout / time.Millisecond)
	})
}

// OptWriter 配置选项-将writer替换掉
//
// 注意，如此替换writer 后，其他的writer配置参数将不生效
// 如BufferSize、WriterTimeout、RotateRule、MaxFileNum
func OptWriter(w io.WriteCloser) Option {
	return newFuncOption(func(config *Config) {
		config.writer = w
	})
}

// OptFlushDuration 配置选项-writer的刷新间隔(毫秒)
func OptFlushDuration(dur uint) Option {
	if dur == 0 {
		return errOption(fmt.Errorf("flushDuration=%d expect>0", dur))
	}
	return newFuncOption(func(config *Config) {
		config.FlushDuration = int(dur)
	})
}

// OptPrefixFunc 配置选项-设置日志前缀方法
func OptPrefixFunc(fn PrefixFunc) Option {
	return newFuncOption(func(config *Config) {
		config.PrefixFunc = fn
	})
}

// OptPrefixName 配置选项-设置日志前缀方法名称
func OptPrefixName(name string) Option {
	return newFuncOption(func(config *Config) {
		// 赋值为nil，这样避免之前已经设置过 PrefixFunc
		config.PrefixFunc = nil
		config.Prefix = name
	})
}

// loggerCfg 通过加载 Option 初始化生成一个 Config
func loggerCfg(opts ...Option) (*Config, error) {
	cfg := &Config{}
	for _, opt := range opts {
		opt.apply(cfg)
		// 还未解析，设置为false
		cfg.parsed = false
	}

	// Option 有错误，直接返回
	if cfg.err != nil {
		return nil, cfg.err
	}

	if err := cfg.parser(); err != nil {
		return nil, err
	}
	return cfg, nil
}
