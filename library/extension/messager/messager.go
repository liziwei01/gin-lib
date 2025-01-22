/*
 * @Author: liziwei01
 * @Date: 2023-11-03 13:44:40
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 14:38:13
 * @Description: message queue 消息队列
 */
package messager

import (
	"sync/atomic"
)

type (
	// MessageType 消息类型
	MessageType uint32

	// Messager Servicer 中信息从Discovery传递到其它组件的接口定义
	Messager interface {
		// MessageType 返回消息的类型
		Type() MessageType

		// MessageData 返回消息的数据
		Data() interface{}
	}
)

var typeid uint32

// 自增
func AcquireMessageType() MessageType {
	return MessageType(atomic.AddUint32(&typeid, 1))
}

type message struct {
	_type MessageType
	_data interface{}
}

// New 创建新Messager
func New(t MessageType, d interface{}) Messager {
	return &message{
		_type: t,
		_data: d,
	}
}

func (m *message) Type() MessageType {
	return m._type
}

func (m *message) Data() interface{} {
	return m._data
}
