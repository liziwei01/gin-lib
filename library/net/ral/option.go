/*
 * @Author: liziwei01
 * @Date: 2023-11-02 02:52:48
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-02 02:55:10
 * @Description: file content
 */
package ral

import (
	"github.com/liziwei01/gin-lib/library/extension/option"
	"github.com/liziwei01/gin-lib/library/logit"
	"github.com/liziwei01/gin-lib/library/net/connector"
)

// sessionConfig 在一次Ral调用过程中的配置信息
type sessionConfig struct {
	logFields []logit.Field
	opt       *option.Dynamic
}

// 将自定义的option 和 base option 进行merge
// 优先读取到自定义的
func (sc *sessionConfig) withServicerOpt(base option.Option) option.Option {
	cf := make(map[interface{}]interface{})
	sc.opt.Range(func(k, v interface{}) bool {
		cf[k] = v
		return true
	})
	return option.NewFixed(base, cf)
}

// ROption ral方法的配置选项
type ROption interface {
	apply(cf *sessionConfig)
}

const (
	// ROptKeyConnector ROptConnector会用到的特殊key，值为Connector
	ROptKeyConnector = "Ral-Connector"
)

// ROptConnector 给当前会话设置Connector
//
// 使用场景：
//
//	如给请求指定的ip发送请求，可以传入一个特殊的Connector,
//	&connector.Single{Addr:xxx}
func ROptConnector(cc connector.Connector) ROption {
	return newFuncROption(func(cf *sessionConfig) {
		cf.opt.Set(ROptKeyConnector, cc)
	})
}

// ROptConnTimeOut 给当前会话设置连接超时，单位毫秒
func ROptConnTimeOut(timeoutMs int) ROption {
	return newFuncROption(func(cf *sessionConfig) {
		cf.opt.Set(option.KeyConnTimeOut, timeoutMs)
	})
}

// ROptReadTimeOut 给当前会话设置读超时，单位毫秒
func ROptReadTimeOut(timeoutMs int) ROption {
	return newFuncROption(func(cf *sessionConfig) {
		cf.opt.Set(option.KeyReadTimeOut, timeoutMs)
	})
}

// ROptWriteTimeOut 给当前会话设置写超时，单位毫秒
func ROptWriteTimeOut(timeoutMs int) ROption {
	return newFuncROption(func(cf *sessionConfig) {
		cf.opt.Set(option.KeyWriteTimeOut, timeoutMs)
	})
}

// ROptRetry 给当前会话设置重试次数
func ROptRetry(retry int) ROption {
	return newFuncROption(func(cf *sessionConfig) {
		cf.opt.Set(option.KeyRetry, retry)
	})
}

// ROptLogFields 额外的日志字段，会记录到ral-worker.log里去
func ROptLogFields(fs ...logit.Field) ROption {
	return newFuncROption(func(cf *sessionConfig) {
		cf.logFields = append(cf.logFields, fs...)
	})
}

type funcROption struct {
	f func(cf *sessionConfig)
}

func (fdo *funcROption) apply(cf *sessionConfig) {
	fdo.f(cf)
}

func newFuncROption(f func(cf *sessionConfig)) *funcROption {
	return &funcROption{
		f: f,
	}
}

func ralOpts(base option.Option, opts ...ROption) *sessionConfig {
	config := &sessionConfig{
		opt: option.NewDynamic(base),
	}
	for _, opt := range opts {
		opt.apply(config)
	}
	return config
}
