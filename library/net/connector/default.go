/*
 * @Author: liziwei01
 * @Date: 2023-11-04 13:24:07
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 15:29:07
 * @Description: file content
 */
package connector

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/liziwei01/gin-lib/library/extension/messager"
	"github.com/liziwei01/gin-lib/library/extension/option"
	"github.com/liziwei01/gin-lib/library/logit"
	"github.com/liziwei01/gin-lib/library/net/addresspicker"
	_ "github.com/liziwei01/gin-lib/library/net/addresspicker/roundrobin"
	"github.com/liziwei01/gin-lib/library/net/connpool"
	"github.com/liziwei01/gin-lib/library/net/gaddr"
	"github.com/liziwei01/gin-lib/library/net/gconn"
)

const (
	defaultKeepAlivePeriod = 1 * time.Second
	defaultTCPLinger       = 1
	defaultTCPNoDelay      = true
	defaultTCPKeepAlive    = true
)

var (
	// ErrNoAddressPicker 因为Connector依赖AddressPicker，如果没有AddressPicker返回此错误。
	ErrNoAddressPicker = errors.New("connector can't find available AddressPicker")

	// Default is a messager consumer
	_ messager.Consumer = &Default{}
)

// Default 默认的连接器
type Default struct {
	logit.WithLogger

	WithProxy // 允许设置代理

	dialer *net.Dialer

	ap addresspicker.AddressPicker
	cp connpool.ConnPool

	// connTimeout 来自 option.Option
	connTimeout int64

	// TCP 各种参数
	tcpLinger          int
	tcpNoDelay         bool
	tcpKeepAlive       bool
	tcpKeepAlivePeriod time.Duration

	// bns option
	optUpdater *option.Updater
	opt        option.Option
}

// DefaultConnectorOpt 以下是 DefaultConnector 的 Opt 设置
type DefaultConnectorOpt interface {
	SetConnectorOpt(*Default)
}

// DCOptDialer 使用自定义的 Dialer
func DCOptDialer(dialer *net.Dialer) DefaultConnectorOpt {
	return &dcOptDialer{dialer: dialer}
}

type dcOptDialer struct {
	dialer *net.Dialer
}

func (o *dcOptDialer) SetConnectorOpt(d *Default) {
	d.dialer = o.dialer
}

// DCOptConnPool 设置 DefaultConnector 的 ConnPool
func DCOptConnPool(cp connpool.ConnPool) DefaultConnectorOpt {
	return &dcOptConnPool{ConnPool: cp}
}

type dcOptConnPool struct {
	connpool.ConnPool
}

func (o *dcOptConnPool) SetConnectorOpt(d *Default) {
	d.cp = o.ConnPool
}

// DCOptTCPKeepAlive 设置 Connector 是否开启 KeepAlive
func DCOptTCPKeepAlive(keepAlive bool) DefaultConnectorOpt {
	return &dcOptTCPKeepAlive{keepAlive: keepAlive}
}

type dcOptTCPKeepAlive struct {
	keepAlive bool
}

func (o *dcOptTCPKeepAlive) SetConnectorOpt(d *Default) {
	d.tcpKeepAlive = o.keepAlive
}

// DCOptTCPKeepAlivePeriod 设置 Connector 是否开启 KeepAlive Period
func DCOptTCPKeepAlivePeriod(period time.Duration) DefaultConnectorOpt {
	return &dcOptTCPKeepAlivePeriod{period: period}
}

type dcOptTCPKeepAlivePeriod struct {
	period time.Duration
}

func (o *dcOptTCPKeepAlivePeriod) SetConnectorOpt(d *Default) {
	d.tcpKeepAlivePeriod = o.period
}

// DCOptTCPLinger 设置 Connector 是否开启 KeepAlive
func DCOptTCPLinger(linger int) DefaultConnectorOpt {
	return &dcOptTCPLinger{linger: linger}
}

type dcOptTCPLinger struct {
	linger int
}

func (o *dcOptTCPLinger) SetConnectorOpt(d *Default) {
	d.tcpLinger = o.linger
}

// DCOptTCPNoDelay 设置 Connector 是否开启 KeepAlive
func DCOptTCPNoDelay(noDelay bool) DefaultConnectorOpt {
	return &dcOptTCPNoDelay{noDelay: noDelay}
}

type dcOptTCPNoDelay struct {
	noDelay bool
}

func (o *dcOptTCPNoDelay) SetConnectorOpt(d *Default) {
	d.tcpNoDelay = o.noDelay
}

// DefaultConnector 的 Opt 设置结束，以下是 DefaultConnector 的实现。

// NewDefault 创建新默认对象
func NewDefault(ap addresspicker.AddressPicker, opts ...DefaultConnectorOpt) *Default {
	d := &Default{
		ap:                 ap,
		connTimeout:        int64(option.ConnTimeout(nil)),
		tcpLinger:          defaultTCPLinger,
		tcpNoDelay:         defaultTCPNoDelay,
		tcpKeepAlive:       defaultTCPKeepAlive,
		tcpKeepAlivePeriod: defaultKeepAlivePeriod,
		opt:                option.NewDynamic(nil),
	}
	if optset, ok := d.opt.(option.Setter); ok {
		d.optUpdater = &option.Updater{Setter: optset}
	}

	for _, opt := range opts {
		opt.SetConnectorOpt(d)
	}
	if d.dialer == nil {
		d.dialer = &net.Dialer{
			KeepAlive: defaultKeepAlivePeriod,
		}
	}

	return d
}

// Pick 获取一个addr
func (c *Default) Pick(ctx context.Context, args ...interface{}) (net.Addr, error) {
	if c.ap == nil {
		return nil, ErrNoAddressPicker
	}
	return c.ap.Pick(ctx, args...)
}

// Connect 是一个同时支持从ConnPool中获取连接和
func (c *Default) Connect(ctx context.Context, addr net.Addr) (net.Conn, error) {
	// 检查 context 是否已经出错了
	if ctx.Err() != nil {
		return nil, fmt.Errorf("cancel dial with context error %w", ctx.Err())
	}

	// 连接池取不到，获取新连接
	if conn := c.pooledConn(ctx, addr); conn != nil {
		return conn, nil
	}

	return c.connect(ctx, addr)
}

// pooledConn 尝试从pool中获取一个连接
func (c *Default) pooledConn(ctx context.Context, addr net.Addr) net.Conn {
	if c.cp == nil {
		return nil
	}

	conn := c.cp.Conn(addr)
	if conn == nil {
		return nil
	}

	statTracer := gconn.NewStatTracer(ctx, conn.LocalAddr(), conn.RemoteAddr(), time.Now())
	return gconn.NewGConn(conn, statTracer, c.moreTracers()...)
}

// connect 创建一个新的连接
func (c *Default) connect(ctx context.Context, addr net.Addr) (net.Conn, error) {
	// 创建 net.Dialer，设置了一个默认的KeepAlive时间。
	// context 设置超时时间
	now := time.Now()
	ddl, ok := ctx.Deadline()
	if !ok || now.Add(c.getConnTimeout()).Before(ddl) {
		ddl = now.Add(c.getConnTimeout())
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, ddl)
		defer cancel()
	}

	logit.AddError(ctx, logit.Duration("conn_timeout", ddl.Sub(now)))

	conn, err := c.dialer.DialContext(ctx, addr.Network(), addr.String())

	if err != nil {
		c.triggerBadConn(ctx, nil, addr, now)
		return nil, fmt.Errorf("connector dial (network=%s, addr=%s) error with %w",
			addr.Network(), addr.String(), err)
	}

	if err = c.setSockOpt(conn); err != nil {
		return nil, err
	}

	statTracer := gconn.NewStatTracer(ctx, conn.LocalAddr(), addr, now)
	return gconn.NewGConn(conn, statTracer, c.moreTracers()...), nil
}

// setSockOpt improve socket performance
func (c *Default) setSockOpt(conn net.Conn) error {
	switch t := conn.(type) {
	case *net.TCPConn:
		if err := t.SetKeepAlive(c.tcpKeepAlive); err != nil {
			return err
		}
		if err := t.SetKeepAlivePeriod(c.tcpKeepAlivePeriod); err != nil {
			return err
		}
		if err := t.SetLinger(c.tcpLinger); err != nil {
			return err
		}
		if err := t.SetNoDelay(c.tcpNoDelay); err != nil {
			return err
		}
	}
	return nil
}

func (c *Default) triggerBadConn(ctx context.Context, localAddr net.Addr, remoteAddr net.Addr, now time.Time) {
	statTracer := gconn.NewStatTracer(ctx, localAddr, remoteAddr, now)
	statTracer.SetError(gconn.ErrDialFailed)
	gConn := gconn.NewGConn(nil, statTracer, c.moreTracers()...)
	// triger BeforeCloseTracer and AfterCloseTracer
	gConn.Close()
}

// Strategy 获取名称
func (c *Default) Strategy() string {
	if c.ap != nil {
		return c.ap.Name()
	}
	return ""
}

// moreTracers
func (c *Default) moreTracers() []gconn.Tracer {
	var tracers []gconn.Tracer
	for _, comp := range []interface{}{c.ap, c.cp} {
		if t, ok := comp.(gconn.Tracer); ok {
			tracers = append(tracers, t)
		}
	}
	return tracers
}

func (c *Default) getConnTimeout() time.Duration {
	return time.Duration(atomic.LoadInt64(&c.connTimeout))
}

func (c *Default) setConnTimeout(d time.Duration) {
	atomic.StoreInt64(&c.connTimeout, int64(d))
}

// Consume consume message as a MessageSubscriber
func (c *Default) Consume(msg messager.Messager) {
	if opt := option.FromMessager(msg); opt != nil {
		c.optUpdater.Consume(msg)
		c.setConnTimeout(option.ConnTimeout(c.opt))
	}
	if addrs := gaddr.FromMessager(msg); addrs != nil {
		if c.ap != nil {
			err := c.ap.SetAddresses(addrs)
			if err != nil {
				c.AutoLogger().Warning(context.Background(), "AddressPicker.SetAddresses error", logit.Error("error", err))
			}
		}
		if c.cp != nil {
			err := c.cp.CleanExcept(addrs)
			if err != nil {
				c.AutoLogger().Warning(context.Background(), "ConnPool.CleanExcept error", logit.Error("error", err))
			}
		}
	}
	if sub, ok := c.ap.(messager.Consumer); ok {
		sub.Consume(msg)
	}
	if sub, ok := c.cp.(messager.Consumer); ok {
		sub.Consume(msg)
	}
}

var _ Connector = (*Default)(nil)
