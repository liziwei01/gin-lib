/*
 * @Author: liziwei01
 * @Date: 2023-05-09 23:35:17
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-05-09 23:36:49
 * @Description: file content
 */
package tinycache

// Config 配置
type Config struct {
	// Service的名字, 必选
	Name string

	// 各种自定义的参数, 全部非必选
	// 超时
	ExpireTime int64
}
