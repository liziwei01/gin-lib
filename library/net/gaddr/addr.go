/*
 * @Author: liziwei01
 * @Date: 2023-11-04 15:37:27
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 15:37:28
 * @Description: file content
 */
package gaddr

import (
	"net"

	"github.com/liziwei01/gin-lib/library/extension/option"
)

type (
	// WeightedAddr 带优先级和权重的地址
	WeightedAddr interface {
		net.Addr
		// Priority 声明地址的优先级，Priority越大，优先级越高，相同Weight的地址应该用优先级更高的这个
		Priority() int
		// Weight 声明地址的权重，Weight越大，权重越大，负载均衡时应该按权重分配流量
		Weight() int64
	}
)

// OptionSetter 允许设置Option
// 包含 OptionSet 方法
type OptionSetter interface {
	OptionSet(key interface{}, value interface{})
}

// HasOption 是否有option
type HasOption interface {
	Option() option.Option
}

// New create net.Addr
func New(network, address string) net.Addr {
	return &addr{
		network: network,
		address: address,
		opt:     option.NewDynamic(nil),
	}
}

type addr struct {
	network string
	address string
	opt     *option.Dynamic
}

func (a *addr) Network() string {
	return a.network
}

func (a *addr) String() string {
	return a.address
}

func (a *addr) OptionSet(key interface{}, value interface{}) {
	a.opt.Set(key, value)
}

func (a *addr) Option() option.Option {
	return a.opt
}

type weightedAddr struct {
	network  string
	address  string
	priority int
	weight   int64
}

// NewWeighted 生成新对象
func NewWeighted(network, address string, priority int, weight int64) WeightedAddr {
	return &weightedAddr{
		network:  network,
		address:  address,
		priority: priority,
		weight:   weight,
	}
}

func (wa *weightedAddr) Network() string {
	return wa.network
}

func (wa *weightedAddr) String() string {
	return wa.address
}

func (wa *weightedAddr) Priority() int {
	return wa.priority
}

func (wa *weightedAddr) Weight() int64 {
	return wa.weight
}
