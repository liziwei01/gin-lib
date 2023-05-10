/*
 * @Author: liziwei01
 * @Date: 2021-07-23 15:23:10
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-05-09 23:58:15
 * @Description: tinycache
 * @FilePath: /gdp-config-platform/library/tinycaches/tinycache.go
 */
package tinycache

import (
	"time"
)

type object struct {
	val        string
	updateTime time.Time
	expireTime time.Duration
}

// val, ok, expired
func (c *client) Get(key string) string {
	if key == "" {
		return ""
	}
	if obj, ok := c.db.Load(key); ok {
		if !c.expired(obj.(object)) {
			return obj.(object).val
		}
	}
	return ""
}

func (c *client) Set(key string, value string, expireTime ...time.Duration) {
	if key == "" || value == "" {
		return
	}
	obj := object{
		val:        value,
		updateTime: time.Now(),
		expireTime: c.getExpiryDefault(expireTime, time.Duration(c.conf.ExpireTime)),
	}
	c.db.Store(key, obj)
}

func (c *client) expired(obj object) bool {
	if obj.expireTime == 0 {
		return false
	}
	return time.Since(obj.updateTime) > obj.expireTime
}

func (c *client) getExpiryDefault(expireTime []time.Duration, defau ...time.Duration) time.Duration {
	if len(expireTime) != 0 {
		return expireTime[0]
	}
	if len(defau) != 0 {
		return defau[0]
	}
	return 0
}
