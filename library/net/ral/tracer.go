/*
 * @Author: liziwei01
 * @Date: 2023-11-04 20:43:44
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 20:43:45
 * @Description: file content
 */
package ral

import (
	"context"
	"net"
	"time"

	"github.com/liziwei01/gin-lib/library/extension/option"
)

// RalTracer trace info of metrics
type RalTracer interface {
	DoTracer(*TracerMeta)
}

type ConnMeta struct {
	ConnStart time.Time
	ConnEnd   time.Time
	ConnErr   error

	PickStart time.Time
	PickEnd   time.Time
	PickErr   error

	RemoteAddr net.Addr
	Conn       net.Conn
}

// TracerMeta 统计使用参数
type TracerMeta struct {
	Ctx context.Context
	Opt option.Option

	API         string
	ServiceName string
	Err         error
	StatusCode  int

	ConnMetas []*ConnMeta

	DecodeStart time.Time
	DecodeEnd   time.Time

	DoStart time.Time
	DoEnd   time.Time

	PreloadStart time.Time
	PreloadEnd   time.Time
}

// NewTracerMeta _
func NewTracerMeta(ctx context.Context, api string, opt option.Option) *TracerMeta {

	if api == "" {
		api = "default"
	}

	serviceName, _ := option.String(opt, "Name", "unknow")

	return &TracerMeta{
		Ctx:         ctx,
		Opt:         opt,
		API:         api,
		ServiceName: serviceName,
		ConnMetas:   make([]*ConnMeta, 0, option.Retry(opt)+1),
	}
}
