/*
 * @Author: liziwei01
 * @Date: 2021-04-19 15:00:00
 * @LastEditTime: 2023-05-13 05:01:37
 * @LastEditors: liziwei01
 * @Description: main
 * @FilePath: /github.com/liziwei01/gin-lib/main.go
 */
package main

import (
	"log"

	"github.com/liziwei01/gin-lib/bootstrap"
	"github.com/liziwei01/gin-lib/httpapi"
)

/**
 * @description: main
 * @param {*}
 * @return {*}
 */
func main() {
	appServer, err := bootstrap.Setup()
	if err != nil {
		log.Fatalln(err)
	}
	// 注册接口路由
	httpapi.InitRouters(appServer.Handler)

	appServer.Start()
}
