/*
 * @Author: liziwei01
 * @Date: 2022-03-04 15:42:58
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-03-20 19:43:15
 * @Description: file content
 */
package oss

// Config 配置
type Config struct {
	// Service的名字, 必选
	Name string

	OSS struct {
		Endpoint        string
		AccessKeyID     string
		AccessKeySecret string
	}
}
