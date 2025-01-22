/*
 * @Author: liziwei01
 * @Date: 2023-11-03 13:47:15
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 13:36:41
 * @Description: message queue 消息队列
 */
package messager

import (
	"context"
	"fmt"
)

type (
	// Broker 支持多个消息生产者和消费者的消息撮合中介
	Broker struct {
		mch  chan Messager
		pros []Producer
		cons []Consumer

		handlerErr func(err error)
	}

	// Producer 可以生产消息的对象
	// 生产的消息会被发送到消息队列中
	Producer interface {
		Messager() <-chan Messager
	}

	// Consumer 可以消费消息的对象
	// 消息队列中的消息会多次消费，被每个消费者都消费一次
	Consumer interface {
		Consume(msg Messager)
	}
)

// NewMessageBroker 创建新的消息中介，bufsize 指定消息队列的buffer大小
func NewMessageBroker(bufSize int) *Broker {
	return &Broker{
		mch: make(chan Messager, bufSize),
	}
}

// AddProducer 添加消息生产者
func (q *Broker) AddProducer(pro Producer) {
	q.pros = append(q.pros, pro)
}

// AddConsumer 添加消息消费者
func (q *Broker) AddConsumer(con Consumer) {
	q.cons = append(q.cons, con)
}

// Publish 发送多条消息给消费者接收，此处需要保证订阅者成功收到消息并处理，需要同步模式。
// 所以不能直接写channel.
func (q *Broker) Publish(msgs []Messager) {
	for _, msg := range msgs {
		for _, con := range q.cons {
			con.Consume(msg)
		}
	}
}

// OnError 注册处理异常时候的回调方法
// 如可以用来打印日志
func (q *Broker) OnError(handler func(err error)) {
	q.handlerErr = handler
}

// Run 后台运行消息队列，传入的context控制消息队列是否结束
// 开始运行后，消息队列会将生产者的消息汇总到一个channel中，然后再将消息分发给消费者
func (q *Broker) Run(ctx context.Context) {
	// 先将pubs的数据汇总到mch channel这个message queue中
	for _, pro := range q.pros {
		go func(pub Producer) {
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-pub.Messager():
					if !ok { // channel was closed
						return
					}
					q.mch <- msg
				}
			}
		}(pro)
	}
	// 再启动pubsub协程将mq中的数据写到subscribers
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-q.pubsub(ctx):
			}
		}
	}()
}

// pubsub 将消息队列中的消息分发给消费者
func (q *Broker) pubsub(ctx context.Context) <-chan error {
	cherr := make(chan error)
	go func(cherr chan<- error) {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("messageBroker panic: %v", r)
				cherr <- err
				if q.handlerErr != nil {
					q.handlerErr(err)
				}
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-q.mch:
				for _, con := range q.cons {
					con.Consume(msg)
				}
			}
		}
	}(cherr)
	return cherr
}
