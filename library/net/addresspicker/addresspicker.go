/*
 * @Author: liziwei01
 * @Date: 2023-11-04 14:45:49
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 14:45:49
 * @Description: file content
 */
package addresspicker

import (
	"context"
	"errors"
	"fmt"
	"net"
)

var (
	// Factory 生成 AddressPicker的名字到工厂方法的映射
	Factory = make(map[string]NewFunc)

	// ErrUnknownType 不支持的类型
	ErrUnknownType = errors.New("unknown address picker type")

	// ErrNoAddress 未找到地址
	ErrNoAddress = errors.New("no address found")
)

type (
	// AddressPicker pick an address by it's own strategy
	AddressPicker interface {
		// Name 组件名字
		Name() string

		// SetAddress 修改AddressPicker 的地址库
		SetAddresses([]net.Addr) error

		// 获取一个地址
		Pick(context.Context, ...interface{}) (net.Addr, error)
	}
)

// NewFunc 创建新picker的函数类型
type NewFunc func(configure func(interface{}) error) (AddressPicker, error)

// Register 注册新的picker
func Register(key string, f func(configure func(interface{}) error) (
	AddressPicker, error)) {
	Factory[key] = NewFunc(f)
}

// New 创建新的picker对象
func New(key string, configure func(interface{}) error) (AddressPicker, error) {
	f, ok := Factory[key]
	if !ok {
		return nil, fmt.Errorf("%w %s", ErrUnknownType, key)
	}
	return f(configure)
}
