/*
 * @Author: liziwei01
 * @Date: 2023-11-04 15:22:47
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 15:22:48
 * @Description: file content
 */
package connpool

import (
	"net"

	"github.com/liziwei01/gin-lib/library/logit"
)

// ConnPool 连接池接口
type ConnPool interface {
	logit.Binder

	// 获取一个独享的连接
	Conn(net.Addr) net.Conn

	// 放进一个可用的连接
	PutConn(net.Conn)

	// 清除不在列表中的连接
	CleanExcept(addresses []net.Addr) error
}
