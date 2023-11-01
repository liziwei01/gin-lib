/*
 * @Author: liziwei01
 * @Date: 2023-10-31 21:57:01
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 11:29:02
 * @Description: 获取调用栈的路径
 */
package logit

import (
	"bytes"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	callerKey = "caller"
	stackKey  = "stack"
)

var (
	pcsPool = sync.Pool{
		New: func() interface{} {
			return &stackPtr{
				pcs: make([]uintptr, 64),
			}
		},
	}
)

type stackPtr struct {
	pcs []uintptr
}

// Stack retrieve call stack
func Stack() Field {
	return StackWithSkip(3)
}

// StackWithSkip 返回调用栈的Field
func StackWithSkip(skip int) Field {
	buf := &bytes.Buffer{}
	stack := pcsPool.Get().(*stackPtr)
	defer pcsPool.Put(stack)
	callStackSize := runtime.Callers(skip, stack.pcs)
	frames := runtime.CallersFrames(stack.pcs[:callStackSize])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		buf.WriteString(frame.File)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(frame.Line))
		buf.WriteByte(';')
	}
	return String(stackKey, buf.String())
}

// CallerField 默认的获取调用栈的Field
func CallerField() Field {
	return CallerFieldWithSkip(2)
}

// CallerFieldWithSkip 获取调用栈
func CallerFieldWithSkip(skip int) Field {
	return String(callerKey, callerWithSkip(skip+1))
}

// callerWithSkip 获取调用栈的路径
// 如  xxx/xxx/xxx.go:80
func callerWithSkip(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	return strings.Join([]string{
		CallerPathClean(file),
		strconv.Itoa(line),
	}, ":")
}

const githubPath = "github.com/"

// CallerPathClean 对caller的文件路径进行精简
// 原始的是完整的路径，比较长，该方法可以将路径变短
var CallerPathClean = callerPathClean

func callerPathClean(file string) string {
	idx := strings.Index(file, githubPath)
	if idx < 0 {
		return file
	}
	return file[idx+len(githubPath):]
}
