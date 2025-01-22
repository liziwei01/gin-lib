/*
 * @Author: liziwei01
 * @Date: 2023-11-03 13:46:36
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 14:41:29
 * @Description: 消费者，用于更新配置
 */
package option

import (
	"github.com/liziwei01/gin-lib/library/extension/messager"
)

var (
	_ messager.Consumer = &Updater{}

	// Message 表示数据是配置更新信息，Data 是 Option
	Message = messager.AcquireMessageType()
)

// Updater 数据可更新
// 实现了 messager.Consumer 接口
type Updater struct {
	Setter Setter
}

// NewMessager 新消息
func NewMessager(opt Option) messager.Messager {
	return messager.New(Message, opt)
}

// FromMessager 从 Messager 中提取 Option，如果类型不对会得到nil
func FromMessager(m messager.Messager) Option {
	if m.Type() != Message {
		return nil
	}
	if ret, ok := m.Data().(Option); ok {
		return ret
	}
	return nil
}

// Consume 消费消息
func (ou *Updater) Consume(m messager.Messager) {
	opt := FromMessager(m)
	if opt != nil {
		opt.Range(func(k, v interface{}) bool {
			ou.Setter.Set(k, v)
			return true
		})
	}
}
