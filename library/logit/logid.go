/*
 * @Author: liziwei01
 * @Date: 2023-10-30 12:12:18
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 11:40:21
 * @Description: 使用统一的requestID跟踪请求、处理日志、调试问题、分析性能
 */
package logit

import (
	"strconv"
	"sync/atomic"
	"time"
)

// NewRequestID 获取一个新的requestID，默认是 uint32 的字符串
var NewRequestID = func() string {
	return strconv.FormatUint(uint64(NewRequestIDUint32()), 10)
}

var idx uint32

func rid() uint32 {
	id := atomic.AddUint32(&idx, 1)
	if id < 65535 {
		return id
	}
	atomic.StoreUint32(&idx, 0)
	return rid()
}

// NewRequestIDUint32 获取一个新的requestID
func NewRequestIDUint32() uint32 {
	usec := now().UnixNano()
	reqid := usec&0x7FFFFFFF | 0x80000000
	// 通过rid()让同一时间生成的requestID不重复
	return uint32(reqid) + rid()
}

var now = func() time.Time {
	return time.Now()
}
