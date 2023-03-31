/*
 * @Author: liziwei01
 * @Date: 2022-03-04 15:42:58
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-03-21 17:41:26
 * @Description: file content
 */
package email

// Config 配置
type Config struct {
	// Service的名字, 必选
	Name string

	// 资源定位: 手动配置 - 使用IP、端口
	Resource struct {
		Manual struct {
			Host string
			Port int
		}
	}

	Email struct {
		Address  string
		Password string
	}
}
