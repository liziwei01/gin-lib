/*
 * @Author: liziwei01
 * @Date: 2023-10-28 20:20:34
 * @LastEditors: liziwei01
 * @LastEditTime: 2024-12-24 14:40:55
 * @Description: file content
 */
package main

import (
	"reflect"
)

func main() {
	fun := ReturnAll("hello")
	a := fun.(func() string)()
	println(a)
}

// ReturnAll 返回一个原样返回传入参数的函数
func ReturnAll(args ...interface{}) interface{} {
	vs := make([]reflect.Value, 0, len(args))
	vts := make([]reflect.Type, 0, len(args))
	for _, v := range args {
		tmp := reflect.ValueOf(v)
		vs = append(vs, tmp)
		vts = append(vts, tmp.Type())
	}

	funTyp := reflect.FuncOf(nil, vts, false)
	funVal := func(_ []reflect.Value) []reflect.Value {
		return vs
	}
	return reflect.MakeFunc(funTyp, funVal).Interface()
}
