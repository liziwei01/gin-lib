/*
 * @Author: liziwei01
 * @Date: 2023-11-03 13:43:58
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 13:39:50
 * @Description: file content
 */
package discoverer

import (
	"context"
	"errors"
	"fmt"

	"github.com/liziwei01/gin-lib/library/env"
	"github.com/liziwei01/gin-lib/library/extension/messager"
	"github.com/liziwei01/gin-lib/library/logit"
)

var (
	// Factory 生成ResourceServicer的名字到工厂方法的映射
	Factory = make(map[string]NewDiscovererByConfigFunc)

	// ErrUnknownType 未知的discoverer类型
	ErrUnknownType = errors.New("unknown discoverer type")
)

type (
	// Discoverer start and run Servicer
	Discoverer interface {
		logit.Binder

		// String describe the service updater
		String() string

		// Discovery 完成一次同步服务Update并返回Messager
		Discovery(ctx context.Context) ([]messager.Messager, error)
	}
)

// NewDiscovererByConfigFunc 生成新对象的函数类型
type NewDiscovererByConfigFunc func(env env.AppEnv, configure func(interface{}) error) (Discoverer, error)

// RegisterDiscoverer 注册
func RegisterDiscoverer(key string, f NewDiscovererByConfigFunc) {
	Factory[key] = f
}

// New 新创建
func New(key string, e env.AppEnv, configure func(interface{}) error) (Discoverer, error) {
	f, ok := Factory[key]
	if !ok {
		return nil, fmt.Errorf("%w with %s", ErrUnknownType, key)
	}
	if e == nil {
		e = env.Default
	}
	return f(e, configure)
}
