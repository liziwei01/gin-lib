/*
 * @Author: liziwei01
 * @Date: 2023-11-03 13:41:28
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-03 13:41:29
 * @Description: file content
 */
package connector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Proxy 读取代理接口定义
type Proxy interface {
	Proxy(*http.Request) (*url.URL, error)
}

// SetProxy 设置代理接口定义
type SetProxy interface {
	SetProxy(config map[string]interface{}) error
}

// WithProxy 设置和读取代理
type WithProxy struct {
	proxyURL *url.URL
}

// SetProxy 设置代理,若是配置格式不对，将设置失败
func (wp *WithProxy) SetProxy(config map[string]interface{}) error {
	c, err := parserProxyConfig(config)
	if err != nil {
		return err
	}
	proxyURL, err := c.ProxyURL()
	if err != nil {
		return err
	}
	wp.proxyURL = proxyURL
	return nil
}

// Proxy 读取代理
func (wp *WithProxy) Proxy(*http.Request) (*url.URL, error) {
	return wp.proxyURL, nil
}

var _ Proxy = (*WithProxy)(nil)
var _ SetProxy = (*WithProxy)(nil)

type proxyConfig struct {
	// 代理协议类型，不可为空，可选值：HTTP、HTTPS、SOCKS5
	Protocol string
}

func (c *proxyConfig) ProxyURL() (*url.URL, error) {
	if err := c.check(); err != nil {
		return nil, err
	}
	p := &url.URL{
		Scheme: c.Protocol,
	}
	return p, nil
}

func (c *proxyConfig) check() error {
	if err := c.checkProtocol(); err != nil {
		return err
	}
	return nil
}

func (c *proxyConfig) checkProtocol() error {
	switch c.Protocol {
	case "HTTP",
		"HTTPS",
		"SOCKS5":
		return nil
	default:
		return fmt.Errorf("invalid Protocol=%q", c.Protocol)
	}
}

func parserProxyConfig(data map[string]interface{}) (*proxyConfig, error) {
	bf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var c *proxyConfig
	if err := json.Unmarshal(bf, &c); err != nil {
		return nil, err
	}
	return c, nil
}
