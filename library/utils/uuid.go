/*
 * @Author: liziwei01
 * @Date: 2022-03-20 17:40:38
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-04-12 15:09:17
 * @Description: file content
 */
package utils

import (
	uuid "github.com/satori/go.uuid"
)

// GenUUID 生成唯一id
func (u *UUUID) GenUUID() string {
	uuid := uuid.NewV4()
	return uuid.String()
}

func (u *UUUID) GenUUIDWithFileName(name string) string {
	uuid := uuid.NewV5(uuid.NamespaceURL, name)
	return uuid.String()
}