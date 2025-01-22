/*
 * @Author: liziwei01
 * @Date: 2023-11-03 13:49:58
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 04:06:28
 * @Description: file content
 */
package logit

import (
	"time"
)

// TimeCost 计时器
type TimeCost struct {
	key string
	t   time.Time
}

// NewTimeCost 计时 field
func NewTimeCost(key string) func() Field {
	tc := TimeCost{
		key: key,
		t:   time.Now(),
	}
	return tc.Stop
}

// Stop 停止计时，并返回一个Field
func (tc *TimeCost) Stop() Field {
	d := time.Since(tc.t)
	return Duration(tc.key, d)
}
