/*
 * @Author: liziwei01
 * @Date: 2023-10-31 21:56:25
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-03 22:59:02
 * @Description: 异步写入，而不是立刻写入，可以节约cpu资源，提高性能
 */
package writer

import (
	"io"
	"sync"
	"time"
)

// NewAsync 创建一个异步的writer
//
//	bufSize 异步队列大小
//	timeout 写超时时间，可以为0，若为0将不超时，阻塞写；若设置为>0的值，当writeTo消费比实际写入多，buf满了将丢弃当前数据
//	writeTo 实际用于写入的同步writer
func NewAsync(bufSize int, timeout time.Duration, writeTo io.WriteCloser) io.WriteCloser {
	w := &asyncWriter{
		msgs:    make(chan []byte, bufSize),
		timeout: timeout,
		raw:     writeTo,
		done:    make(chan struct{}),
	}
	go w.consumer()
	return w
}

// 实现io.WriteCloser接口，用于异步写入
// 包装了一个用于同步写入的writer，和两个chan，用于异步写入
type asyncWriter struct {
	// msgs 用于存放要写入的数据，consumer消费者会从msgs中读取数据，写入raw中
	msgs chan []byte
	// closed 用于标记a是否已经关闭
	closed bool
	// timeout 写超时。0则不超时
	timeout time.Duration

	// 同步写入的writer
	raw io.WriteCloser
	// done 用于通知consumer消费者，a已经关闭
	done chan struct{}
	// mu 用于保护closed字段
	mu sync.Mutex
}

// consumer 消费者 消费写入msgs chan中的数据，直到a关闭。一般以协程的方式运行
func (a *asyncWriter) consumer() {
	for p := range a.msgs {
		_, _ = a.raw.Write(p)
	}
	a.done <- struct{}{}
}

func (a *asyncWriter) Write(p []byte) (n int, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return 0, io.ErrClosedPipe
	}

	// 如果timeout为0，那么就是阻塞写，不会丢弃数据
	// 将数据写入msgs chan中，排队等待consumer消费
	if a.timeout == 0 {
		a.msgs <- p
		return len(p), nil
	}

	// 如果timeout不为0，那么就是非阻塞写，如果p在timeout时间内没有被msgs传递到消费者，那么就会被丢弃
	select {
	case a.msgs <- p:
		return len(p), nil
	case <-time.After(a.timeout):
		return 0, ErrWriteTimeout
	}
}

// Close 关闭异步写入
// 实现io.WriteCloser Close接口
// 1给done chan发送一个信号，通知consumer消费者，a已经关闭，2再正常
func (a *asyncWriter) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.closed {
		return nil
	}

	close(a.msgs)
	<-a.done

	a.closed = true
	return a.raw.Close()
}

var _ io.WriteCloser = (*asyncWriter)(nil)
