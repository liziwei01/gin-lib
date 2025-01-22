/*
 * @Author: liziwei01
 * @Date: 2023-11-05 00:25:38
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-05 00:26:41
 * @Description: file content
 */
package manual

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/liziwei01/gin-lib/library/env"
	"github.com/liziwei01/gin-lib/library/extension/messager"
	"github.com/liziwei01/gin-lib/library/logit"

	"github.com/liziwei01/gin-lib/library/net/discoverer"
	"github.com/liziwei01/gin-lib/library/net/gaddr"
)

// Manual Resource Discoverer
var (
	manualSubsetsOrder = []string{"all", "default", "tc"}
)

func init() {
	RegisterManualDiscoverer()
}

// RegisterManualDiscoverer 支持 Manual 格式的配置
func RegisterManualDiscoverer() {
	discoverer.RegisterDiscoverer("Manual", manualDiscoverer)
}

func manualDiscoverer(env env.AppEnv, configure func(interface{}) error) (discoverer.Discoverer, error) {
	conf := make(map[string][]*struct {
		Host string // 可以是ip，也可以是域名
		Port int    // 端口
	})
	if err := configure(&conf); err != nil {
		return nil, err
	}
	if len(conf) == 0 {
		return nil, errors.New("empty manual config")
	}
	idc := env.IDC()
	var addresses []net.Addr
	manualRsrc := conf[idc]
	// 如果没有对应的 IDC 信息，就按照manualSubsetsOrder的顺序查找
	for i := 0; len(manualRsrc) == 0 && i < len(manualSubsetsOrder); i++ {
		idc = manualSubsetsOrder[i]
		manualRsrc = conf[idc]
	}
	// 如果没有manualSubsetsOrder的配置项出现，就随机用一个
	if len(manualRsrc) == 0 {
		for k := range conf {
			idc = k
			manualRsrc = conf[k]
			break
		}
	}
	for _, confAddr := range manualRsrc {
		if len(confAddr.Host) > 0 {
			addr := gaddr.New("tcp", confAddr.Host+":"+strconv.Itoa(confAddr.Port))
			gaddr.MustSetRemoteIDC(addr, idc)
			addresses = append(addresses, addr)
		}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("manual config no hosts")
	}

	return NewStaticDiscoverer(addresses), nil
}

// StaticDiscoverer is an implementation of static service
type StaticDiscoverer struct {
	logit.WithLogger

	addresses []net.Addr
}

// NewStaticDiscoverer Discoverer of static
func NewStaticDiscoverer(addresses []net.Addr) *StaticDiscoverer {
	return &StaticDiscoverer{
		addresses: addresses,
	}
}

// String return describe of service updater
func (d *StaticDiscoverer) String() string {
	return fmt.Sprintf("StaticDiscoverer(addresses=%v)", d.addresses)
}

// Discovery 返回当前 IDC 的 address
func (d *StaticDiscoverer) Discovery(ctx context.Context) ([]messager.Messager, error) {
	return []messager.Messager{gaddr.NewMessager(d.addresses)}, nil
}
