/*
 * @Author: liziwei01
 * @Date: 2023-11-04 15:33:02
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 15:33:02
 * @Description: file content
 */
package gconn

import (
	"errors"
)

var (
	// ErrDialFailed connect error
	ErrDialFailed = errors.New("connect dial failed")
)
