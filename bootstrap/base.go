/*
 * @Author: liziwei01
 * @Date: 2022-03-03 16:04:06
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-04-06 17:10:23
 * @Description: 读取配置文件, 初始化路由
 */
package bootstrap

import (
	"context"
	"log"

	"github.com/liziwei01/gin-lib/library/env"

	"github.com/gin-gonic/gin"
)

const (
	appConfPath = "./conf/app.toml"
)

// AppServer struct.
type AppServer struct {
	Handler *gin.Engine
	Ctx     context.Context
	Config  *Config
	Cancel  context.CancelFunc
}

// Setup 准备.
func Setup() (*AppServer, error) {
	appServer := &AppServer{}
	var (
		err error
	)
	appServer.Config, err = ParserAppConfig(appConfPath)
	if err != nil {
		return nil, err
	}
	env.Default = appServer.Config.Env
	appServer.Ctx, appServer.Cancel = context.WithCancel(context.Background())
	InitMust(context.WithValue(appServer.Ctx, "APPName", appServer.Config.APPName))
	appServer.Handler = InitHandler(appServer)

	return appServer, nil
}

// SetupScript 准备.
func SetupScript(conf ...string) (*AppServer, error) {
	cPath := appConfPath
	if len(conf) > 0 {
		cPath = conf[0]
	}
	appServer := &AppServer{}
	var (
		err error
	)
	appServer.Config, err = ParserAppConfig(cPath)
	if err != nil {
		return nil, err
	}
	env.Default = appServer.Config.Env
	appServer.Ctx, appServer.Cancel = context.WithCancel(context.Background())
	InitMust(appServer.Ctx)
	return appServer, nil
}

// Start 启动http服务器.
func (appServer *AppServer) Start() {
	defer appServer.Cancel()
	app := NewApp(appServer.Ctx, appServer.Config, appServer.Handler)
	log.Fatalln("server exit:", app.Start())
}
