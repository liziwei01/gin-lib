/*
 * @Author: liziwei01
 * @Date: 2023-11-02 02:08:22
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 04:13:50
 * @Description: file content
 */
package logit

// Binder 可以与 Logger 绑定的组件：比如 RAL
type Binder interface {
	// SetLogger 设置打日志的Logger
	SetLogger(Logger)

	// Logger 获取Logger
	Logger() Logger
}

// WithLogger 默认的Binder实现
// 任何struct只要嵌入了WithLogger，就可以自动继承内部的logger对象和logger和WithLogger实现的两个方法
type WithLogger struct {
	logger Logger
}

// SetLogger 设置logger
func (b *WithLogger) SetLogger(logger Logger) {
	b.logger = logger
}

// Logger 返回logger
func (b *WithLogger) Logger() Logger {
	return b.logger
}

// AutoLogger 自动获取logger，若未设置，会返回DefaultLogger
func (b *WithLogger) AutoLogger() Logger {
	if b.logger != nil {
		return b.logger
	}
	return DefaultLogger
}
