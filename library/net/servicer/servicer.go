/*
 * @Author: liziwei01
 * @Date: 2023-11-03 13:37:34
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 20:23:12
 * @Description: 载入后端需要使用的服务（如mysql,redis等）用的servicer组件
 */
package servicer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/liziwei01/gin-lib/library/extension/messager"
	"github.com/liziwei01/gin-lib/library/extension/option"
	"github.com/liziwei01/gin-lib/library/logit"

	"github.com/liziwei01/gin-lib/library/net/connector"
	"github.com/liziwei01/gin-lib/library/net/discoverer"
)

// ErrNoConnector 没有connector
var ErrNoConnector = errors.New("service has no connector")

// Servicer interface definition
type Servicer interface {
	Worker
	logit.Binder

	// Name 返回 Servicer 的名字
	Name() string

	// String 返回 Servicer 的描述
	String() string

	// Connector 返回 Servicer 当前提供服务的 Connector
	Connector() connector.Connector

	// Discoverer 返回 Servicer 的 Discoverer
	Discoverer() discoverer.Discoverer

	// Option 返回 Servicer 的配置
	Option() option.Option
}

// Worker 可以后台运行的任务
type Worker interface {
	// Start Worker 开始运行，运行后如果context结束，Worker就结束。
	Start(context.Context) error

	// Stop Worker 停止运行，与直接通过context结束不同之处在于可以返回结束的错误信息。
	Stop() error
}

// DefaultService 是最基础通用的 Servicer 实现，能满足大多数服务需要
type DefaultService struct {
	logit.WithLogger
	stop func()

	name          string
	connector     connector.Connector
	discoverer    discoverer.Discoverer
	option        option.Option
	optionUpdater *option.Updater

	mb *messager.Broker

	mtx sync.Mutex
}

// NewDefault 创建一个名为name并使用对应ServicerComponents的DefaultService
func NewDefault(name string, sc *Components) *DefaultService {
	ds := &DefaultService{
		name:       name,
		connector:  sc.Connector,
		discoverer: sc.Discoverer,
		option:     sc.Option,
	}

	if optset, ok := sc.Option.(option.Setter); ok {
		ds.optionUpdater = &option.Updater{Setter: optset}
	}

	if DefaultLogger != nil {
		ds.SetLogger(DefaultLogger)
	}

	return ds
}

// Name of DefaultService
func (s *DefaultService) Name() string {
	return s.name
}

// String 返回服务的描述
func (s *DefaultService) String() string {
	if s.discoverer != nil {
		return s.Name() + " with " + s.discoverer.String()
	}
	return s.Name()
}

// Connector 获取ServiceConnector
func (s *DefaultService) Connector() connector.Connector {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.connector
}

// Discoverer 获取 Discoverer
func (s *DefaultService) Discoverer() discoverer.Discoverer {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.discoverer
}

// Option 获取 Option
func (s *DefaultService) Option() option.Option {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.option
}

// Start Servicer and it's components
func (s *DefaultService) Start(ctx context.Context) (resultError error) {
	ctx, cancel := context.WithCancel(ctx)

	ctx = logit.ForkContext(ctx)

	logit.AddAllLevel(ctx, logit.String("service_name", s.Name()))

	ts := logit.NewTimeCost("service_start_cost")

	defer func() {
		if resultError != nil {
			s.AutoLogger().Fatal(ctx, "Servicer Start failed", logit.Error("error", resultError), ts())
		}
	}()

	// 先启动 MQ，建立各组件之间的通信关系
	s.mb = messager.NewMessageBroker(0)
	if err := s.startMessageBroker(ctx); err != nil {
		cancel()
		return fmt.Errorf("startMessageBroker failed: %w", err)
	}

	// 先启动后台discoverer，后启动前台connector
	for idx, w := range []interface{}{s.discoverer, s.connector} {
		if worker, ok := w.(Worker); ok {
			if err := worker.Start(ctx); err != nil {
				cancel()
				return fmt.Errorf("worker(%d).Start failed:%w", idx, err)
			}
		}
	}

	// 启动时同步更新一次数据
	if s.option != nil {
		s.mb.Publish([]messager.Messager{option.NewMessager(s.option)})
	}
	if disc := s.Discoverer(); disc != nil {
		msgs, err := disc.Discovery(ctx)
		if err != nil {
			cancel()
			return fmt.Errorf("start Discovery failed: %w", err)
		}
		s.mb.Publish(msgs)
	}

	s.AutoLogger().Notice(ctx, "Servicer Start OK", ts())
	s.stop = cancel
	return nil
}

func (s *DefaultService) startMessageBroker(ctx context.Context) error {
	for _, comp := range []interface{}{s.discoverer, s.connector, s.option, s.optionUpdater} {
		if pub, ok := comp.(messager.Producer); ok {
			s.mb.AddProducer(pub)
		}
		if sub, ok := comp.(messager.Consumer); ok {
			s.mb.AddConsumer(sub)
		}
	}
	s.mb.Run(ctx)
	return nil
}

// Stop 停止 Servicer 以及与其关联的组件
func (s *DefaultService) Stop() error {
	// 停止顺序：与启动顺序相反。
	// stop 中发生 error 不阻碍下一个组件的 stop
	var err error
	for _, comp := range []interface{}{s.optionUpdater, s.option, s.connector, s.discoverer} {
		if worker, ok := comp.(Worker); ok {
			if stopErr := worker.Stop(); stopErr != nil {
				err = fmt.Errorf("with error %w", stopErr)
			}
		}
	}

	// 通过 s.stop 通过所有组件stop
	if s.stop != nil {
		s.stop()
		s.stop = nil
	}

	if err != nil {
		s.AutoLogger().Warning(context.Background(), "Stop Servicer Failed", logit.String("servicer_name", s.Name()), logit.Error("error", err))
		return err
	}
	s.AutoLogger().Notice(context.Background(), "Servicer Stop OK", logit.String("servicer_name", s.Name()))
	return nil
}

var _ Servicer = (*DefaultService)(nil)
var _ Worker = (*DefaultService)(nil)
