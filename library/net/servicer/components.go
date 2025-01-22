/*
 * @Author: liziwei01
 * @Date: 2023-11-03 13:39:45
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-03 13:39:46
 * @Description: file content
 */
package servicer

import (
	"fmt"

	"github.com/liziwei01/gin-lib/library/extension/option"
	"github.com/liziwei01/gin-lib/library/net/connector"

	"github.com/liziwei01/gin-lib/library/net/discoverer"
)

// Components 定义了一个 Servicer 的部件集合，用于创建新 Servicer 时使用
type Components struct {
	Connector  connector.Connector
	Discoverer discoverer.Discoverer
	Option     option.Option
}

// TrySetProxy 尝试给Connector设置代理
func (c *Components) TrySetProxy() error {
	setProxy, ok := c.Connector.(connector.SetProxy)
	if !ok {
		return nil
	}
	val, has := option.Value(c.Option, ConfKeyProxy)
	if !has {
		return nil
	}

	config, ok := val.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid %q section, expect is map[string]interface{}", ConfKeyProxy)
	}

	return setProxy.SetProxy(config)
}
