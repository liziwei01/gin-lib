/*
 * @Author: liziwei01
 * @Date: 2023-11-03 13:42:00
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-05 00:19:06
 * @Description: 配置解析与回调函数
 */
package servicer

import (
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/liziwei01/gin-lib/library/conf"
	"github.com/liziwei01/gin-lib/library/env"
	"github.com/liziwei01/gin-lib/library/extension/option"
	"github.com/liziwei01/gin-lib/library/logit"
	"github.com/liziwei01/gin-lib/library/net/addresspicker"
	"github.com/liziwei01/gin-lib/library/net/addresspicker/roundrobin"
	"github.com/liziwei01/gin-lib/library/net/connector"
	"github.com/liziwei01/gin-lib/library/net/discoverer"
)

var (
	// defaultAddressPicker 如果没有配置AddressPicker的时候默认的
	defaultAddressPicker = roundrobin.APStrategyRoundRobin
)

// servicer 配置中的key的名称
const (
	// ConfKeyProxy 代理，是否为代理服务器，bool类型，默认为false
	// 如配置为Proxy=true,则 对应的资源位代理服务器的地址
	ConfKeyProxy = "Proxy"
)

// mapToStruct 把一个配置从map转换到struct中
func mapToStruct(dst interface{}, src interface{}) error {
	srcMap, ok := src.(map[string]interface{})
	if !ok {
		return fmt.Errorf("src must be a map[string]interface{}")
	}

	v := reflect.ValueOf(dst).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		if value, ok := srcMap[fieldName]; ok {
			fieldValue := reflect.ValueOf(value)
			if !fieldValue.Type().AssignableTo(v.Field(i).Type()) {
				return fmt.Errorf("cannot assign value of type %s to field %s of type %s", fieldValue.Type(), fieldName, v.Field(i).Type())
			}
			v.Field(i).Set(fieldValue)
		}
	}
	return nil
}

// StrategyConfig 策略配置
type StrategyConfig map[string]interface{}

// Name 名字
func (sc StrategyConfig) Name() string {
	if v, ok := sc["Name"]; ok {
		if name, ok := v.(string); ok {
			return name
		}
	}
	return ""
}

// newDefaultStrategyConfig
func newDefaultStrategyConfig() StrategyConfig {
	sc := map[string]interface{}{
		"Name": defaultAddressPicker,
	}
	return sc
}

// minimalConfig 创建 Servicer 的最小配置
type minimalConfig struct {
	Name     string
	Strategy StrategyConfig
	Resource map[string]interface{}
}

func parseWithMinimalConfig(config map[string]interface{}) (*minimalConfig, error) {
	minConf := &minimalConfig{}
	err := mapToStruct(minConf, config)
	if err != nil {
		return nil, fmt.Errorf("extract minimal config failed: %w", err)
	}

	if len(minConf.Name) == 0 {
		return nil, errors.New("name is empty")
	}

	// Check if stragety is valid
	if minConf.Strategy == nil {
		minConf.Strategy = newDefaultStrategyConfig()
	}

	// Check if resource is valid, and extract it
	if len(minConf.Resource) != 1 {
		return nil, errors.New("resource configured wrong with none or more than one")
	}

	return minConf, nil
}

var configHooks = make([]struct {
	name  string // 名字，需要唯一
	index int    // 排序，小的排前面
	fn    func(config map[string]interface{}) (map[string]interface{}, error)
}, 0)

// RegisterConfigHookFunc 注册配置解析hook func
//
// index 小的排前面执行
// name 不可重复
// 如 mysql、redis的配置解析都是有该功能
func RegisterConfigHookFunc(name string, index int, fn func(config map[string]interface{}) (map[string]interface{}, error)) {
	if name == "" {
		panic("cannot Register with empty name")
	}
	for _, item := range configHooks {
		if item.name == name {
			panic(fmt.Sprintf("cannot Register with duplicate name %q", name))
		}
	}

	item := struct {
		name  string
		index int
		fn    func(config map[string]interface{}) (map[string]interface{}, error)
	}{
		name:  name,
		index: index,
		fn:    fn,
	}
	configHooks = append(configHooks, item)
	sort.SliceStable(configHooks, func(i, j int) bool {
		return configHooks[i].index < configHooks[j].index
	})
}

func executeConfigHookFunc(config map[string]interface{}) (map[string]interface{}, error) {
	if len(configHooks) == 0 {
		return config, nil
	}
	result := make(map[string]interface{}, len(config))
	for k, v := range config {
		result[k] = v
	}
	for _, item := range configHooks {
		var err error
		result, err = item.fn(result)
		if err != nil {
			return nil, fmt.Errorf("execute hook %q with error:%w", item.name, err)
		}
	}
	return result, nil
}

// NewWithConfig 创建新servicer对象
func NewWithConfig(logger logit.Logger, env env.AppEnv, config map[string]interface{}) (Servicer, error) {
	var err error
	// map转换为struct，以获取 Name Strategy Resource三项，用于本函数内配置connector、discoverer
	minConf, err := parseWithMinimalConfig(config)
	if err != nil {
		return nil, err
	}

	// 执行hook func，通知各个组件进行配置解析
	config, err = executeConfigHookFunc(config)
	if err != nil {
		return nil, err
	}

	// 根据配置信息生成组件，Components是一个组件的集合，作为生成Servicer之前的中转
	scomp := &Components{}

	// address picker
	addressPicker, err := addresspicker.New(minConf.Strategy.Name(), func(i interface{}) error {
		name := minConf.Strategy.Name()
		strategy, ok := minConf.Strategy[name]
		if !ok {
			return fmt.Errorf("no strategy name: %q", name)
		}
		return mapToStruct(i, strategy)
	})
	if err != nil {
		return nil, err
	}
	if b, ok := addressPicker.(logit.Binder); ok {
		b.SetLogger(logger)
	}

	// 设置组件connector
	scomp.Connector = connector.NewDefault(addressPicker)
	if b, ok := scomp.Connector.(logit.Binder); ok {
		b.SetLogger(logger)
	}

	// service updater
	var rsrcType string
	for key := range minConf.Resource {
		rsrcType = key
	}

	// 设置组件discoverer: 用于获取Resource
	scomp.Discoverer, err = discoverer.New(rsrcType, env, func(i interface{}) error {
		return mapToStruct(i, minConf.Resource[rsrcType])
	})
	if err != nil {
		return nil, err
	}
	scomp.Discoverer.SetLogger(logger)

	// 将配置信息保存到option中
	opts := option.NewDynamic(nil)
	for key, val := range config {
		opts.Set(key, val)
	}
	scomp.Option = opts

	if err := scomp.TrySetProxy(); err != nil {
		return nil, fmt.Errorf("setProxy with error: %w", err)
	}

	srv := NewDefault(minConf.Name, scomp)
	srv.SetLogger(logger)
	return srv, nil
}

// NewWithConfigName 通过解析配置文件生成一个servicer
func NewWithConfigName(logger logit.Logger, env env.AppEnv, confName string) (Servicer, error) {
	config := make(map[string]interface{})
	if err := conf.Parse(confName, &config); err != nil {
		return nil, err
	}
	return NewWithConfig(logger, env, config)
}
