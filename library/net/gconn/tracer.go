/*
 * @Author: liziwei01
 * @Date: 2023-11-04 15:30:08
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 15:30:09
 * @Description: file content
 */
package gconn

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// StatTracer trace info of GConn
type StatTracer interface {
	Tracer
	// Context
	Context() context.Context

	// Local address
	Local() net.Addr
	// Remote address
	Remote() net.Addr

	// BeginTime
	BeginTime() time.Time
	// CloseTime
	EndTime() time.Time
	// LastReadTime
	LastReadTime() time.Time
	// LastWriteTime
	LastWroteTime() time.Time
	// ReadDeadline
	ReadDeadline() time.Time
	// WriteDeadline
	WriteDeadline() time.Time
	// Deadline
	Deadline() time.Time

	// ReadSize
	ReadSize() int64
	// WroteSize
	WroteSize() int64

	// TracerError
	Error() error
	// SetError
	SetError(error)
}

type (
	// Tracer 是一个空接口，用来判断一个对象是不是可以加到tracer列表
	Tracer interface {
		TraceGConn()
	}
	// AfterCreateTracer trace after gconn created
	AfterCreateTracer interface {
		TraceAfterCreate(GConn)
	}
	// BeforeReadTracer trace before read
	BeforeReadTracer interface {
		TraceBeforeRead(GConn, []byte)
	}
	// AfterReadTracer trace after read
	AfterReadTracer interface {
		TraceAfterRead(GConn, []byte, int, error)
	}
	// BeforeWriteTracer trace before write
	BeforeWriteTracer interface {
		TraceBeforeWrite(GConn, []byte)
	}
	// AfterWriteTracer trace after write
	AfterWriteTracer interface {
		TraceAfterWrite(GConn, []byte, int, error)
	}
	// BeforeCloseTracer trace before close
	BeforeCloseTracer interface {
		TraceBeforeClose(GConn)
	}
	// AfterCloseTracer trace after close
	AfterCloseTracer interface {
		TraceAfterClose(GConn, error)
	}
	// BeforeSetReadDeadlineTracer trace before SetReadDeadline
	BeforeSetReadDeadlineTracer interface {
		TraceBeforeSetReadDeadline(GConn, time.Time)
	}
	// AfterSetReadDeadlineTracer trace after SetReadDeadline
	AfterSetReadDeadlineTracer interface {
		TraceAfterSetReadDeadline(GConn, error)
	}
	// BeforeSetWriteDeadlineTracer trace before SetWriteDeadline
	BeforeSetWriteDeadlineTracer interface {
		TraceBeforeSetWriteDeadline(GConn, time.Time)
	}
	// AfterSetWriteDeadlineTracer trace after SetWriteDeadline
	AfterSetWriteDeadlineTracer interface {
		TraceAfterSetWriteDeadline(GConn, error)
	}
	// BeforeSetDeadlineTracer trace before SetDeadline
	BeforeSetDeadlineTracer interface {
		TraceBeforeSetDeadline(GConn, time.Time)
	}
	// AfterSetDeadlineTracer trace after SetDeadline
	AfterSetDeadlineTracer interface {
		TraceAfterSetDeadline(GConn, error)
	}
)

// statTracer
type statTracer struct {
	ctx    context.Context
	local  net.Addr
	remote net.Addr

	// times 连接创建时间
	beginTime     time.Time
	endTime       time.Time
	lastReadTime  time.Time
	lastWroteTime time.Time
	readDeadline  time.Time
	writeDeadline time.Time
	deadline      time.Time

	// sizes
	readSize  int64
	wroteSize int64

	err error

	mu  sync.RWMutex
	rmu sync.Locker
}

// NewStatTracer create new StatTracer
func NewStatTracer(ctx context.Context, local, remote net.Addr, now time.Time) StatTracer {
	st := &statTracer{
		ctx:       ctx,
		local:     local,
		remote:    remote,
		beginTime: now,
	}
	st.rmu = st.mu.RLocker()
	return st
}

// TraceGConn fit Tracer
func (st *statTracer) TraceGConn() {}

// Context
func (st *statTracer) Context() context.Context {
	return st.ctx
}

// Local
func (st *statTracer) Local() net.Addr {
	return st.local
}

// Remote
func (st *statTracer) Remote() net.Addr {
	return st.remote
}

// BeginTime
func (st *statTracer) BeginTime() time.Time {
	st.rmu.Lock()
	defer st.rmu.Unlock()
	return st.beginTime
}

// EndTime
func (st *statTracer) EndTime() time.Time {
	st.rmu.Lock()
	defer st.rmu.Unlock()
	return st.endTime
}

// LastReadTime
func (st *statTracer) LastReadTime() time.Time {
	st.rmu.Lock()
	defer st.rmu.Unlock()
	return st.lastReadTime
}

// LastWroteTime
func (st *statTracer) LastWroteTime() time.Time {
	st.rmu.Lock()
	defer st.rmu.Unlock()
	return st.lastWroteTime
}

// ReadDeadline
func (st *statTracer) ReadDeadline() time.Time {
	st.rmu.Lock()
	defer st.rmu.Unlock()
	return st.readDeadline
}

// WriteDeadline
func (st *statTracer) WriteDeadline() time.Time {
	st.rmu.Lock()
	defer st.rmu.Unlock()
	return st.writeDeadline
}

// Deadline
func (st *statTracer) Deadline() time.Time {
	st.rmu.Lock()
	defer st.rmu.Unlock()
	return st.deadline
}

// ReadSize
func (st *statTracer) ReadSize() int64 {
	return atomic.LoadInt64(&st.readSize)
}

// WroteSize
func (st *statTracer) WroteSize() int64 {
	return atomic.LoadInt64(&st.wroteSize)
}

// TracerError
func (st *statTracer) Error() error {
	st.rmu.Lock()
	defer st.rmu.Unlock()
	return st.err
}

// SetTracerError
func (st *statTracer) SetError(err error) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.err = err
}

// TraceAfterRead
func (st *statTracer) TraceAfterRead(_ GConn, n int, err error) {
	atomic.AddInt64(&st.readSize, int64(n))
	if err != nil {
		st.SetError(err)
	}
}

// TraceAfterWrite
func (st *statTracer) TraceAfterWrite(_ GConn, n int, err error) {
	atomic.AddInt64(&st.wroteSize, int64(n))
	if err != nil {
		st.SetError(err)
	}
}

// TraceClose
func (st *statTracer) TraceAfterClose(_ GConn, err error) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.endTime = time.Now()
}

// BeforeSetReadDeadlineTracer trace before SetReadDeadline
func (st *statTracer) TraceBeforeSetReadDeadline(_ GConn, t time.Time) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.readDeadline = t
}

// AfterSetReadDeadlineTracer trace after SetReadDeadline
func (st *statTracer) TraceAfterSetReadDeadline(_ GConn, err error) {
	if err != nil {
		st.SetError(err)
	}
}

// BeforeSetWriteDeadlineTracer trace before SetWriteDeadline

func (st *statTracer) TraceBeforeSetWriteDeadline(_ GConn, t time.Time) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.writeDeadline = t
}

// AfterSetWriteDeadlineTracer trace after SetWriteDeadline
func (st *statTracer) TraceAfterSetWriteDeadline(_ GConn, err error) {
	if err != nil {
		st.SetError(err)
	}
}

// BeforeSetDeadlineTracer trace before SetDeadline
func (st *statTracer) TraceBeforeSetDeadline(_ GConn, t time.Time) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.deadline = t
}

// AfterSetDeadlineTracer trace after SetDeadline
func (st *statTracer) TraceAfterSetDeadline(_ GConn, err error) {
	if err != nil {
		st.SetError(err)
	}
}
