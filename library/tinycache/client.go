/*
 * @Author: liziwei01
 * @Date: 2023-05-09 23:32:35
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-05-09 23:58:24
 * @Description: file content
 */
package tinycache

import (
	"sync"
	"time"
)

var (
	// 初始化互斥锁
	mu sync.Mutex
)

type Client interface {
	// GetStr 获取value
	Get(key string) string
	// Set 将字符串值 value 关联到 key
	// 如果 key 已经存在， SET 将覆盖旧值 无视类型
	// 过期时间为 nanoseconds 纳秒
	Set(key string, value string, expireTime ...time.Duration)

	expired(obj object) bool

	getExpiryDefault(expireTime []time.Duration, defau ...time.Duration) time.Duration

	name() string
	expireTime() int64
}

type client struct {
	conf *Config
	db   sync.Map
}

func New(config *Config) Client {
	c := &client{
		conf: config,
		db:   sync.Map{},
	}
	return c
}

func (c *client) name() string {
	return c.conf.Name
}

func (c *client) expireTime() int64 {
	return c.conf.ExpireTime
}
