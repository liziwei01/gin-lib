/*
 * @Author: liziwei01
 * @Date: 2023-11-02 02:53:46
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-02 02:53:47
 * @Description: file content
 */
package option

import (
	"time"
)

const (
	// DefaultConnTimeout 5 seconds。
	DefaultConnTimeout = 5 * time.Second

	// DefaultReadTimeout 5 seconds。
	DefaultReadTimeout = 5 * time.Second

	// DefaultWriteTimeout 5 seconds。
	DefaultWriteTimeout = 5 * time.Second
)

const (
	// KeyConnTimeOut 连接超时的key:ConnTimeOut，其值为 int 类型，单位为毫秒
	KeyConnTimeOut = "ConnTimeOut"

	// KeyReadTimeOut 读超时的key:ReadTimeOut，其值为 int 类型，单位为毫秒
	KeyReadTimeOut = "ReadTimeOut"

	// KeyWriteTimeOut 写超时的key:WriteTimeOut，其值为 int 类型，单位为毫秒
	KeyWriteTimeOut = "WriteTimeOut"

	// KeyRetry 重试次数key:Retry，其值为 int 类型
	KeyRetry = "Retry"

	// KeyMetrics int 是否增加 metrics，1 为true
	KeyMetrics = "Metrics"
)

// Value 从opt中查找最外层的key的值，找到返回值和true，opt为nil或未找到返回def和false。
func Value(opt Option, key interface{}) (interface{}, bool) {
	if opt == nil {
		return nil, false
	}
	val := opt.Value(key)
	return val, val != nil
}

// Int 从opt中查找最外层的key的值，找到返回值和true，未找到返回def和false。
func Int(opt Option, key interface{}, def int) (int, bool) {
	if val, ok := Value(opt, key); ok {
		switch v := val.(type) {
		case int:
			return v, true
		case int64:
			return int(v), true
		case int32:
			return int(v), true
		case uint:
			return int(v), true
		case uint64:
			return int(v), true
		case uint32:
			return int(v), true
		}
	}
	return def, false
}

// Int64 从opt中查找最外层的key的值，找到返回值和true，未找到返回def和false。
func Int64(opt Option, key interface{}, def int64) (int64, bool) {
	if val, ok := Value(opt, key); ok {
		switch v := val.(type) {
		case int64:
			return v, true
		case int:
			return int64(v), true
		case int32:
			return int64(v), true
		case uint:
			return int64(v), true
		case uint64:
			return int64(v), true
		case uint32:
			return int64(v), true
		}
	}
	return def, false
}

// String 从opt中查找最外层的key的值，找到返回值和true，未找到返回def和false。
func String(opt Option, key interface{}, def string) (string, bool) {
	if val, ok := Value(opt, key); ok {
		if v, tok := val.(string); tok {
			return v, true
		}
	}
	return def, false
}

// Bool 从opt中查找最外层的key的值，找到返回值和true，未找到返回def和false
func Bool(opt Option, key interface{}, def bool) (value bool, exists bool) {
	if val, ok := Value(opt, key); ok {
		if v, tok := val.(bool); tok {
			return v, true
		}
	}
	return def, false
}

// TotalTimeout 根据opt获取总耗时
func TotalTimeout(opt Option) time.Duration {
	return ConnTimeout(opt) + ReadTimeout(opt) + WriteTimeout(opt)
}

// ConnTimeout 获取连接超时，若无返回默认值5s
func ConnTimeout(opt Option) time.Duration {

	// 关于查询keys=[]string{KeyConnTimeOut, "service_ctimeout"} 的说明：
	// Key=ConnTimeOut 是在直接配置文件里配置的,在配置文件里是可选的
	// Key=service_ctimeout 是配置在 bns-group 的配置里的，一般要求都是需要设置有效值的
	// 下面的ReadTimeout、WriteTimeout、Retry 和 ConnTimeout类似

	for _, key := range []string{KeyConnTimeOut, "service_ctimeout"} {
		if t, ok := Int(opt, key, 0); ok {
			return time.Millisecond * time.Duration(t)
		}
	}
	return DefaultConnTimeout
}

// ReadTimeout 获取读超时，若无返回默认值5s
func ReadTimeout(opt Option) time.Duration {
	for _, key := range []string{KeyReadTimeOut, "service_rtimeout"} {
		if t, ok := Int(opt, key, 0); ok {
			return time.Millisecond * time.Duration(t)
		}
	}
	return DefaultReadTimeout
}

// WriteTimeout 获取写超时，若无返回默认值5s
func WriteTimeout(opt Option) time.Duration {
	for _, optkey := range []string{KeyWriteTimeOut, "service_wtimeout"} {
		if t, ok := Int(opt, optkey, 0); ok {
			return time.Millisecond * time.Duration(t)
		}
	}
	return DefaultWriteTimeout
}

// Retry 获取重试次数，若无返回默认值0
func Retry(opt Option) int {
	for _, key := range []string{KeyRetry, "service_retry"} {
		if r, ok := Int(opt, key, 0); ok {
			return r
		}
	}
	return 0
}

// HasMetrics servicer 是否使用采集 metrics
func HasMetrics(opt Option) bool {
	if r, ok := Int(opt, KeyMetrics, 0); ok {
		return r == 1
	}
	return false
}
