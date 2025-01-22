/*
 * @Author: liziwei01
 * @Date: 2023-11-02 02:54:53
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-02 02:54:54
 * @Description: file content
 */
package connector

import (
	"context"
	"net"
)

// Connector 是 Servicer 中用于发起连接的部件
// 但是为了便于Request/Response单独使用，Connector 并不要求实现 ServiceBinder 接口
type (
	Connector interface {
		// Pick 暴露 address picker 的 Pick 方法给 Request 使用
		Pick(context.Context, ...interface{}) (net.Addr, error)
		// Connect 创建一个新的网络连接
		Connect(ctx context.Context, addr net.Addr) (net.Conn, error)
	}

	// HasStrategy 带有负载均衡策略的Connector的策略名称
	HasStrategy interface {
		Strategy() string
	}
)
