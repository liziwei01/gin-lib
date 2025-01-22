/*
 * @Author: liziwei01
 * @Date: 2023-11-02 01:49:00
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 14:45:03
 * @Description: 发到消息队列的消息体配置，链表结构的map
 */
package option

import (
	"fmt"
	"sync"
)

var (
	_ Option = &Fixed{}
	_ Option = &Dynamic{}
)

type (
	// Option 接口定义
	Option interface {
		// String 将 Option 以有意义的string进行描述，用于日志输出等
		String() string

		// Value 根据 key 获取配置中的 value
		Value(k interface{}) interface{}

		// Range 遍历不包括 Base 的所有的key-value对，遇到false中止
		Range(func(k, v interface{}) bool)

		// Base 返回当前 Option 的下一层 Option
		Base() Option
	}
	// Setter can set option
	Setter interface {
		Set(key interface{}, val interface{})
	}
)

// Fixed 固定的 Option 配置
// 使用 map[interface{}]interface{} 存储配置，线程不安全，所以禁止外部修改
// 需要修改配置时，使用 Dynamic
type Fixed struct {
	base Option
	kv   map[interface{}]interface{}
}

// NewFixed 创建Fixed Option
func NewFixed(base Option, m map[interface{}]interface{}) Option {
	return &Fixed{
		base: base,
		kv:   m,
	}
}

// Value 从map中读取，如果没有则从base中读取
func (mo *Fixed) Value(k interface{}) interface{} {
	if v, ok := mo.kv[k]; ok {
		return v
	}
	if mo.base != nil {
		return mo.base.Value(k)
	}
	return nil
}

// Range 遍历读取
func (mo *Fixed) Range(f func(k, v interface{}) bool) {
	for k, v := range mo.kv {
		if !f(k, v) {
			break
		}
	}
}

// String 序列化，调试用
func (mo *Fixed) String() string {
	ret := fmt.Sprintf("FixedOption(%v)", mo.kv)
	if mo.Base() != nil {
		ret += " " + mo.Base().String()
	}
	return ret
}

// Base 基类
func (mo *Fixed) Base() Option {
	return mo.base
}

// Dynamic 可变更的 Option 配置
type Dynamic struct {
	base Option
	kv   sync.Map
}

// NewDynamic 创建可变更的 Option 配置
func NewDynamic(base Option) *Dynamic {
	return &Dynamic{
		base: base,
	}
}

// Value 从map中读取，如果没有则从base中读取
func (do *Dynamic) Value(k interface{}) interface{} {
	if v, ok := do.kv.Load(k); ok {
		return v
	}
	if do.base == nil {
		return nil
	}
	return do.base.Value(k)
}

// Range 遍历读取
func (do *Dynamic) Range(f func(k, v interface{}) bool) {
	do.kv.Range(f)
}

// String 序列化，调试用
func (do *Dynamic) String() string {
	ret := "DynamicOption("
	do.kv.Range(func(k, v interface{}) bool {
		ret += fmt.Sprintf("%v=%v,", k, v)
		return true
	})
	ret += ")"
	if do.base != nil {
		ret += " " + do.base.String()
	}
	return ret
}

// Base 基类
func (do *Dynamic) Base() Option {
	return do.base
}

// Set 设置
func (do *Dynamic) Set(k, v interface{}) {
	do.kv.Store(k, v)
}
