/*
 * @Author: liziwei01
 * @Date: 2022-03-04 23:49:32
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-03-04 23:53:11
 * @Description: file content
 */
package utils

import "net/http"

func (u *URequest) Header(req *http.Request, name string) (value string, has bool) {
	vs := req.Header.Values(name)
	if len(vs) == 0 {
		return "", false
	}
	return vs[0], true
}

func (u *URequest) HeaderDefault(req *http.Request, name string, defaultValue string) string {
	if v, _ := u.Header(req, name); v != "" {
		return v
	}
	return defaultValue
}
