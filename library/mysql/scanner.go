/*
 * @Author: liziwei01
 * @Date: 2023-11-04 15:42:31
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 15:42:32
 * @Description: file content
 */

package mysql

import (
	"github.com/didi/gendry/scanner"
)

// ReadRowsClose 从查询结果(Rows)中读出结果数据 并自动调用rs.Close
//
// 注意：
// 若结果为空，不会返回error
// 一个Rows只能调用该方法一次，多次调用会返回错误
func ReadRowsClose(rs scanner.Rows, target interface{}) error {
	err := scanner.ScanClose(rs, target)
	if err == scanner.ErrEmptyResult {
		return nil
	}
	return err
}

var tagName = "ddb"

// SetTagName 设置scanner的tag，只允许设置一次
func SetTagName(name string) {
	tagName = name
	scanner.SetTagName(name)
}
