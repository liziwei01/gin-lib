/*
 * @Author: liziwei01
 * @Date: 2022-03-04 15:43:21
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 11:18:41
 * @Description: file content
 */
package tinycache

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/liziwei01/gin-lib/library/conf"
	"github.com/liziwei01/gin-lib/library/env"
)

const (
	// tinycache conf file path
	tinycachePath = "/servicer/"
	suffix        = ".toml"
)

var (
	// conf file root path
	configPath = env.Default.ConfDir()
	// tinycache client map, client use single instance mode
	clients map[string]Client
	// 初始化互斥锁
	initMux sync.Mutex
)

/**
 * @description:
 * @param {context.Context} ctx
 * @param {string} serviceName
 * @return {*}
 */
func GetClient(ctx context.Context, serviceName string) (Client, error) {
	// try to get from single instance map
	if client, hasSet := clients[serviceName]; hasSet {
		if client != nil {
			return client, nil
		}
	}
	// set a new instance
	client, err := setClient(serviceName)
	if client != nil {
		return client, nil
	}
	return nil, err
}

/**
 * @description: init tinycache client，considering concurrent set, lock
 * @param {string} serviceName
 * @return {*}
 */
func setClient(serviceName string) (Client, error) {
	// 互斥锁
	initMux.Lock()
	defer initMux.Unlock()
	// 初始化
	client, err := initClient(serviceName)
	if err == nil {
		if clients == nil {
			clients = make(map[string]Client)
		}
		// 添加
		clients[serviceName] = client
		return client, nil
	}
	return nil, err
}

/**
 * @description: according to conf service, read conf from conf file to init tinycache client
 * @param {string} serviceName
 * @return {*}
 */
func initClient(serviceName string) (Client, error) {
	var config *Config
	fileAbs, err := filepath.Abs(filepath.Join(configPath, tinycachePath, serviceName+suffix))
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(fileAbs); !os.IsNotExist(err) {
		conf.Default.Parse(fileAbs, &config)
		client := New(config)
		return client, nil
	}
	return nil, fmt.Errorf("conf file not exist")
}
