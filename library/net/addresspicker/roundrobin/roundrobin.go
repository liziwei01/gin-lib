/*
 * @Author: liziwei01
 * @Date: 2023-11-04 14:46:21
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 14:46:22
 * @Description: roundrobin 依次遍历轮询策略的负载均衡器
 */
package roundrobin

import (
	"context"
	"net"
	"sync"

	"github.com/liziwei01/gin-lib/library/net/addresspicker"
)

const (
	// APStrategyRoundRobin rr
	APStrategyRoundRobin = "RoundRobin"
)

// RoundRobin implement AddressPicker
type RoundRobin struct {
	addrs []net.Addr
	idx   int
	mtx   sync.Mutex
}

func init() {
	RegisterRoundRobin(APStrategyRoundRobin)
}

// RegisterRoundRobin 注册
func RegisterRoundRobin(strategy string) {
	addresspicker.Register(strategy, func(f func(interface{}) error) (
		addresspicker.AddressPicker, error) {
		return New(), nil
	})
}

// New address picker
func New() *RoundRobin {
	return &RoundRobin{
		idx: -1,
	}
}

// Name 名字
func (rr *RoundRobin) Name() string {
	return APStrategyRoundRobin
}

// SetAddresses change addresses in picker
func (rr *RoundRobin) SetAddresses(addrs []net.Addr) error {
	rr.mtx.Lock()
	defer rr.mtx.Unlock()

	rr.addrs = addrs
	return nil
}

// Pick return a net address
func (rr *RoundRobin) Pick(context.Context, ...interface{}) (net.Addr, error) {
	rr.mtx.Lock()
	defer rr.mtx.Unlock()

	if len(rr.addrs) == 0 {
		return nil, addresspicker.ErrNoAddress
	}

	rr.idx = (rr.idx + 1) % len(rr.addrs)
	return rr.addrs[rr.idx], nil
}
