/*
 * @Author: liziwei01
 * @Date: 2023-11-04 15:32:26
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 15:32:26
 * @Description: provide GConn which wrap net.Conn and support extensive abilities
 */
package gconn

import (
	"errors"
	"net"
	"net/http"
	"time"
)

var (
	// ErrHijacked occurs while hijack more than once
	ErrHijacked = errors.New("connection is already hijacked")
)

// GConn 扩展了 net.Conn ，增加了callback的能力，提供了trace信息。
type GConn interface {
	// Conn
	net.Conn
	// StatTracer
	StatTracer() StatTracer
	// Hijack
	Hijack() (net.Conn, error)
}

// gConn
type gConn struct {
	// Conn
	_conn net.Conn
	// st
	st StatTracer
	// tracers
	tracers []Tracer
	// hijacked is whether this connection has been hijacked
	// by a HandlerFunc with the Hijacker interface.
	hijacked bool
}

// NewGConn create new GConn
func NewGConn(conn net.Conn, statTracer StatTracer, moreTracers ...Tracer) GConn {
	gconn := &gConn{
		_conn:   conn,
		st:      statTracer,
		tracers: append([]Tracer{statTracer}, moreTracers...),
	}

	for _, t := range moreTracers {
		if c, ok := t.(AfterCreateTracer); ok {
			c.TraceAfterCreate(gconn)
		}
	}

	return gconn
}

// Read
func (gc *gConn) Read(b []byte) (int, error) {
	for _, t := range gc.tracers {
		if c, ok := t.(BeforeReadTracer); ok {
			c.TraceBeforeRead(gc, b)
		}
	}

	n, err := gc._conn.Read(b)

	for _, t := range gc.tracers {
		if c, ok := t.(AfterReadTracer); ok {
			c.TraceAfterRead(gc, b, n, err)
		}
	}

	return n, err
}

// Write
func (gc *gConn) Write(b []byte) (int, error) {
	for _, t := range gc.tracers {
		if c, ok := t.(BeforeWriteTracer); ok {
			c.TraceBeforeWrite(gc, b)
		}
	}

	n, err := gc._conn.Write(b)

	for _, t := range gc.tracers {
		if c, ok := t.(AfterWriteTracer); ok {
			c.TraceAfterWrite(gc, b, n, err)
		}
	}

	return n, err
}

// Close
func (gc *gConn) Close() error {
	for _, t := range gc.tracers {
		if c, ok := t.(BeforeCloseTracer); ok {
			c.TraceBeforeClose(gc)
		}
	}

	var err error
	if !gc.hijacked && gc._conn != nil {
		err = gc._conn.Close()
	}

	for _, t := range gc.tracers {
		if c, ok := t.(AfterCloseTracer); ok {
			c.TraceAfterClose(gc, err)
		}
	}

	return err
}

func (gc *gConn) SetReadDeadline(ddl time.Time) error {
	for _, t := range gc.tracers {
		if c, ok := t.(BeforeSetDeadlineTracer); ok {
			c.TraceBeforeSetDeadline(gc, ddl)
		}
	}

	err := gc._conn.SetReadDeadline(ddl)

	for _, t := range gc.tracers {
		if c, ok := t.(AfterSetReadDeadlineTracer); ok {
			c.TraceAfterSetReadDeadline(gc, err)
		}
	}

	return err
}

func (gc *gConn) SetWriteDeadline(ddl time.Time) error {
	for _, t := range gc.tracers {
		if c, ok := t.(BeforeSetWriteDeadlineTracer); ok {
			c.TraceBeforeSetWriteDeadline(gc, ddl)
		}
	}

	err := gc._conn.SetWriteDeadline(ddl)

	for _, t := range gc.tracers {
		if c, ok := t.(AfterSetWriteDeadlineTracer); ok {
			c.TraceAfterSetWriteDeadline(gc, err)
		}
	}

	return err
}

func (gc *gConn) SetDeadline(ddl time.Time) error {
	for _, t := range gc.tracers {
		if c, ok := t.(BeforeSetDeadlineTracer); ok {
			c.TraceBeforeSetDeadline(gc, ddl)
		}
	}

	err := gc._conn.SetDeadline(ddl)

	for _, t := range gc.tracers {
		if c, ok := t.(AfterSetDeadlineTracer); ok {
			c.TraceAfterSetDeadline(gc, err)
		}
	}

	return err
}

func (gc *gConn) LocalAddr() net.Addr {
	return gc._conn.LocalAddr()
}

func (gc *gConn) RemoteAddr() net.Addr {
	return gc._conn.RemoteAddr()
}

// StatTracer
func (gc *gConn) StatTracer() StatTracer {
	return gc.st
}

// Hijack
func (gc *gConn) Hijack() (net.Conn, error) {
	if gc.hijacked {
		return nil, http.ErrHijacked
	}

	gc.hijacked = true
	return gc._conn, nil
}
