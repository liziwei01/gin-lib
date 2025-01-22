/*
 * @Author: liziwei01
 * @Date: 2023-11-04 15:27:27
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 15:29:14
 * @Description: file content
 */
/*
 * Copyright(C) 2020 Baidu Inc. All Rights Reserved.
 * Author: Chen Xin (chenxin@baidu.com)
 * Date: 2020/04/19
 */

package gaddr

import (
	"net"

	"github.com/liziwei01/gin-lib/library/extension/messager"
)

var (
	// AddressMessage 表示数据是地址更新信息，Data 是 []net.Addr
	AddressMessage = messager.AcquireMessageType()
)

// NewMessager 创建一个携带一组地址的messager
func NewMessager(a []net.Addr) messager.Messager {
	return messager.New(AddressMessage, a)
}

// FromMessager 从 Messager 中提取 net.Addr，如果类型不对会得到nil
func FromMessager(m messager.Messager) []net.Addr {
	if m.Type() != AddressMessage {
		return nil
	}
	if ret, ok := m.Data().([]net.Addr); ok {
		return ret
	}
	return nil
}
