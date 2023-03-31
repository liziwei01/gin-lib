/*
 * @Author: liziwei01
 * @Date: 2022-06-28 00:58:59
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-06-28 01:04:46
 * @Description: file content
 */
package cookie

import (
	"os"
	"path/filepath"

	"github.com/liziwei01/gin-lib/library/conf"
	"github.com/liziwei01/gin-lib/library/env"
)

const (
	// oss conf file path
	cookiePath  = "/middleware/"
	serviceName = "cookie"
	prefix      = ".toml"
)

var (
	// conf file root path
	configPath = env.Default.ConfDir()
	hashKey    = []byte("default-hashKey")
)

/**
 * @description:
 * @param {string} serviceName
 * @return {*}
 */
func init() {
	var config *Config
	fileAbs, _ := filepath.Abs(filepath.Join(configPath, cookiePath, serviceName+prefix))

	if _, err := os.Stat(fileAbs); !os.IsNotExist(err) {
		conf.Default.Parse(fileAbs, &config)
		hashKey = config.Key
	}
}
