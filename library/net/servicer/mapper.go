/*
 * @Author: liziwei01
 * @Date: 2023-11-02 02:57:13
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 21:43:08
 * @Description: 存储servicer服务的mapper组件，所有的服务都会注册保存到这里
 */
package servicer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/liziwei01/gin-lib/library/logit"
)

var (
	// ErrDuplicatedService 表示已经有同名服务存在，返回的时候还会带上服务名称，可以用 errors.Is()方法判断。
	// 如果担心自己的服务被别人撞名，可以用private key type的方式：
	// mykey 是一个 unexported type，别人无法创建出这个类型的key
	// type myKey string
	// serviceMap.AddService(myKey("CommonKey"), myServicer)
	// serviceMap.Service("CommonKey") #-> maybe nil maybe other servicer
	// serviceMap.Service(myKey("CommonKey")) #-> myServicer
	ErrDuplicatedService = errors.New("duplicated service name")

	// ErrServiceNotFound 服务未找到
	ErrServiceNotFound = errors.New("service not found")

	// ErrMapperNotFound serviceMapper 不存在 或者是未初始化
	ErrMapperNotFound = errors.New("serviceMapper not found")

	// DefaultMapper 是包级别的全局Mapper。内含ral logger
	DefaultMapper Mapper
)

// Mapper 保存 Service 和 Name 的一一对应关系
type Mapper interface {
	logit.Binder

	// AddService 添加服务，如果服务已存在，添加会失败并返回error
	AddServicer(key interface{}, srv Servicer) error

	// SetService 添加或修改服务，如果服务已存在，旧的服务会被返回，如果不存在旧服务，会返回nil
	SetServicer(key interface{}, srv Servicer) Servicer

	// PopService 将服务取出并返回给调用方，同时从Mapper中删除。如果不存在，返回nil
	PopServicer(key interface{}) Servicer

	// Service 返回名为key的 Servicer，如果不存在，返回 nil
	Servicer(key interface{}) Servicer

	// Range 遍历 Mapper 中的所有服务并调用回调方法，如果方法返回error不为nil，遍历终止。
	Range(func(key interface{}, srv Servicer) error)
}

// NewMapper 创建新的Mapper
func NewMapper() Mapper {
	sm := &serviceMap{}

	if DefaultLogger != nil {
		sm.SetLogger(DefaultLogger)
	}

	return sm
}

// ServiceMap implement Mapper
type serviceMap struct {
	logit.WithLogger

	smap sync.Map
}

// AddServicer 添加服务
func (sm *serviceMap) AddServicer(key interface{}, srv Servicer) error {
	logFields := []logit.Field{
		logit.AutoField("name", key),
		logit.String("servicer_name", srv.Name()),
	}
	if _, ok := sm.smap.Load(key); ok {
		sm.AutoLogger().Warning(context.Background(), "Try to add servicer with duplicated name",
			logFields...)
		return fmt.Errorf("%w: %v", ErrDuplicatedService, key)
	}
	sm.smap.Store(key, srv)
	sm.AutoLogger().Notice(context.Background(), "ServiceMap AddServicer OK", logFields...)
	return nil
}

// SetServicer 添加或修改服务，如果存在旧的同名服务，会返回旧的服务，否则返回nil。
func (sm *serviceMap) SetServicer(key interface{}, srv Servicer) Servicer {
	old := sm.PopServicer(key)
	sm.smap.Store(key, srv)
	sm.AutoLogger().Notice(context.Background(), "ServiceMap SetServicer OK",
		logit.AutoField("name", key), logit.String("servicer_name", srv.Name()))
	return old
}

// PopServicer 取出服务同时从Mapper中删除
func (sm *serviceMap) PopServicer(key interface{}) Servicer {
	srv := sm.Servicer(key)
	if srv == nil {
		sm.AutoLogger().Notice(context.Background(), "ServiceMap PopServicer but servicer not exists",
			logit.AutoField("name", key))
		return nil
	}
	sm.smap.Delete(key)
	sm.AutoLogger().Notice(context.Background(), "ServiceMap PopServicer OK",
		logit.AutoField("name", key), logit.String("servicer_name", srv.Name()))

	return srv
}

// Servicer 查询服务
func (sm *serviceMap) Servicer(key interface{}) Servicer {
	if s, ok := sm.smap.Load(key); ok {
		if srv, tok := s.(Servicer); tok {
			return srv
		}
	}
	return nil
}

// Range 遍历服务
func (sm *serviceMap) Range(cb func(key interface{}, srv Servicer) error) {
	sm.smap.Range(func(k, v interface{}) bool {
		return cb(k, v.(Servicer)) == nil
	})
}
