/*
 * @Author: liziwei01
 * @Date: 2022-03-10 20:56:10
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-03-10 20:56:11
 * @Description: file content
 */
package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// md5 string.
func (u *UMd5) Md5String(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
