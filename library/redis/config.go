/*
 * @Author: liziwei01
 * @Date: 2022-03-04 15:42:58
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-03-30 23:16:33
 * @Description: file content
 */
package redis

// Config 配置
type Config struct {
	// Service的名字, 必选
	Name string

	// 各种自定义的参数, 全部非必选
	// 写数据超时
	WriteTimeOut int
	// 读数据超时
	ReadTimeOut int
	// 请求失败后的重试次数: 总请求次数 = Retry + 1
	Retry int

	// 资源定位: 手动配置 - 使用IP、端口
	Resource struct {
		Manual struct {
			Host string
			Port int
		}
	}

	Redis struct {
		Password string
		DB       int
	}
}
