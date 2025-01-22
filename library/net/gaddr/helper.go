/*
 * @Author: liziwei01
 * @Date: 2023-11-20 18:38:39
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-20 18:38:40
 * @Description: file content
 */
package gaddr

import (
	"net"
)

type optKey uint8

const (
	// optKeyRemoteIDC ip所在idc
	optKeyRemoteIDC optKey = iota
)

// MustSetOption 给Addr设置上扩展属性，若不可设置，将panic
func MustSetOption(addr net.Addr, key interface{}, val interface{}) {
	addr.(OptionSetter).OptionSet(key, val)
}

// OptionGet 从Addr读取扩展属性，若不存在，将返回nil
func OptionGet(addr net.Addr, key interface{}) interface{} {
	if addr == nil {
		return nil
	}
	opt, ok := addr.(HasOption)
	if !ok {
		return nil
	}
	return opt.Option().Value(key)
}

// OptionGetString 读取一个string的选项值
func OptionGetString(addr net.Addr, key interface{}) string {
	val := OptionGet(addr, key)
	if val == nil {
		return ""
	}
	if v, ok := val.(string); ok {
		return v
	}
	return ""
}

// MustSetRemoteIDC 给addr 设置 remote idc 信息
func MustSetRemoteIDC(addr net.Addr, idc string) {
	MustSetOption(addr, optKeyRemoteIDC, idc)
}

// RemoteIDC 获取remote idc
func RemoteIDC(addr net.Addr) string {
	return OptionGetString(addr, optKeyRemoteIDC)
}

// SetRemoteIDC 给addr 设置 remote idc 信息
// Deprecated 请使用MustSetRemoteIDC
func SetRemoteIDC(addr net.Addr, idc string) (ret bool) {
	defer func() {
		if re := recover(); re != nil {
			ret = false
		}
	}()
	MustSetRemoteIDC(addr, idc)
	return true
}
