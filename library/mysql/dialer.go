/*
 * @Author: liziwei01
 * @Date: 2023-10-29 10:44:05
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 18:54:20
 * @Description: file content
 */
package mysql

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/liziwei01/gin-lib/library/logit"
	"github.com/liziwei01/gin-lib/library/net/gaddr"
	"github.com/liziwei01/gin-lib/library/net/servicer"
)

const dialNetType = "gin_mysql"
const dialNameHostSeparator = "|"

var mapper = &sync.Map{}

func init() {
	mysql.RegisterDialContext(dialNetType, dialFunc)
}

func dialFunc(ctx context.Context, nameAndHost string) (net.Conn, error) {
	ctx = logit.WithContext(ctx)

	arr := strings.SplitN(nameAndHost, dialNameHostSeparator, 2)
	if len(arr) != 2 {
		return nil, fmt.Errorf("not support addr: %q, expect format: name|host:port", nameAndHost)
	}
	name := arr[0]
	hostPort := arr[1]

	// 由于拨号过程不能中不能传递context以及其他参数
	// 所以这里使用了一个全局的mapper
	val, has := mapper.Load(name)
	if !has {
		// 不应该运行到这里
		return nil, fmt.Errorf("mysql dial with unkonwn servicer %q", name)
	}
	ser := val.(servicer.Servicer)

	netAddr := gaddr.New("tcp", hostPort)

	return ser.Connector().Connect(ctx, netAddr)
}
