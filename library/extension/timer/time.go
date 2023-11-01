/*
 * @Author: liziwei01
 * @Date: 2023-10-31 20:05:52
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 11:34:36
 * @Description: 打包时间函数
 */
package timer

import (
	"time"
)

// 将当前时间函数替换为可配置的函数，方便测试
var nowFunc = time.Now
