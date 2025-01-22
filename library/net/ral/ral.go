/*
 * @Author: liziwei01
 * @Date: 2023-11-02 01:47:53
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 23:03:14
 * @Description: remote access layer 远程过程调用。存储servicer的信息，包括mapper, ral logger, ral-worker logger
 */
package ral

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/liziwei01/gin-lib/library/extension/option"
	"github.com/liziwei01/gin-lib/library/logit"
	"github.com/liziwei01/gin-lib/library/net/connector"
	"github.com/liziwei01/gin-lib/library/net/servicer"
)

var (
	// DefaultRaller is the package level RAL when call ral.RAL or ral.GoRAL actually use
	//
	// 在使用前，需要先初始化
	DefaultRaller Raller

	// ErrNoLogger when IRAL has no logger, many components can't work without logger
	ErrNoLogger = errors.New("raller has no logger, set logger first")

	// ErrExecTimeout MutliRAL has at least one not finished
	ErrExecTimeout = errors.New("ral was not finished before deadline")

	// R 同步调用，是 DefaultRaller.RAL
	R = r

	// RAL 同步调用，是 DefaultRaller.RAL
	RAL = r

	// GoR 异步调用，是 DefaultRaller.GoRAL
	GoR = gor

	// GoRAL 异步调用，是 DefaultRaller.GoRAL
	GoRAL = gor

	// MultiRAL 并发调用，是 DefaultRaller.MultiRAL
	MultiRAL = multiR
)

type (
	// Request 接口定义
	Request interface {
		// Do 完成一个请求，并将返回写到 Response。Ral 会提供 Connector 给请求用于连接下游。
		Do(context.Context, Response, connector.Connector, option.Option) error

		// String 将RalRequest以有意义的string进行描述，用于日志输出等
		String() string
	}

	// Response 接口定义
	Response interface {
		// String 将 RalResponse 以有意义的string进行描述，用于日志输出等
		String() string
	}
)

var errNotInit = errors.New("ral.DefaultRaller has not inited")

// r
func r(ctx context.Context, name interface{}, req Request, resp Response, opts ...ROption) error {
	if DefaultRaller == nil {
		panic(errNotInit.Error())
	}
	return DefaultRaller.RAL(ctx, name, req, resp, opts...)
}

// gor
func gor(ctx context.Context, name interface{}, req Request, resp Response, opts ...ROption) <-chan error {
	if DefaultRaller == nil {
		panic(errNotInit.Error())
	}
	return DefaultRaller.GoRAL(ctx, name, req, resp, opts...)
}

// multiR
func multiR(ctx context.Context, reqs map[interface{}]*RALParam) error {
	if DefaultRaller == nil {
		panic(errNotInit.Error())
	}

	return DefaultRaller.MultiRAL(ctx, reqs)
}

// Raller 支持Ral能力的接口定义
type Raller interface {
	// Logger 获取Logger
	// 用于打印ral组件自己状态的日志
	// 即一般的ral.log
	Logger() logit.Logger

	// 请求日志
	// 即一般的ral-worker.log
	WorkLogger() logit.Logger

	// Ral 同步请求
	RAL(context.Context, interface{}, Request, Response, ...ROption) error

	// MultiRAL 并发执行请求
	MultiRAL(ctx context.Context, reqs map[interface{}]*RALParam) error

	// GoRal 异步请求，通过error channel返回异步请求结果
	GoRAL(context.Context, interface{}, Request, Response, ...ROption) <-chan error

	// ServiceMapper 获取ServiceMapper
	ServiceMapper() servicer.Mapper
}

// Raller 接口的基础实现
type raller struct {
	logit.WithLogger
	withWorkLogger logit.WithLogger

	sm servicer.Mapper
}

// NewRaller 创建新的 IRAL
func NewRaller(sm servicer.Mapper, log logit.Logger, worklog logit.Logger) Raller {
	r := &raller{sm: sm}
	r.SetLogger(log)
	r.setWorkLogger(worklog)
	return r
}

func (r *raller) RAL(ctx context.Context, name interface{}, req Request, resp Response, opts ...ROption) (errRet error) {
	if r.WorkLogger() == nil {
		return ErrNoLogger
	}

	ctx = logit.NewContext(ctx)
	InitRalStatisItems(ctx)

	logit.ReplaceFields(ctx, logit.AutoField(LogFieldService, name))

	sessionCfg := ralOpts(nil, opts...)

	// 会话级别额外的日志字段
	if len(sessionCfg.logFields) > 0 {
		logit.ReplaceFields(ctx, sessionCfg.logFields...)
	}

	wl := r.WorkLogger()

	var noWFLog bool

	// 检查是否需要打印wf日志
	defer func() {
		if errRet == nil || noWFLog {
			return
		}
		// 请求(Request)的类型
		var requestType string
		if req == nil {
			requestType = "nil"
		} else {
			requestType = req.String()
		}

		wl.Warning(ctx, "",
			logit.String("request", requestType),
			logit.Error(LogFieldErrmsg, errRet),
		)
	}()

	if r.sm == nil {
		return servicer.ErrMapperNotFound
	}

	if req == nil {
		return errors.New("request is nil")
	}

	// Check context
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context has error already: %w", err)
	}

	// find Servicer
	srv := r.sm.Servicer(name)
	if srv == nil {
		return servicer.ErrServiceNotFound
	}

	// 一个只读的option
	// 会优先读取当前会话配置的值，若没有则会读取srv.Option()
	opt := sessionCfg.withServicerOpt(srv.Option())

	cnctr := srv.Connector()

	if v := opt.Value(ROptKeyConnector); v != nil {
		if vc, ok := v.(connector.Connector); ok {
			cnctr = vc
		} else {
			return fmt.Errorf("option.%q's type wrong, want 'Connector' type", ROptKeyConnector)
		}
	}

	if cnctr == nil {
		return servicer.ErrNoConnector
	}

	if cs, ok := cnctr.(connector.HasStrategy); ok {
		logit.ReplaceFields(ctx, logit.String(LogFieldBalance, cs.Strategy()))
	}

	// bind work logger and run
	r.tryBindWorkLogger(req)
	r.tryBindWorkLogger(resp)

	errRet = req.Do(ctx, resp, cnctr, opt)
	if errRet != nil {

		// 调用协议自己的 Do 方法，在异常的时候已经打印过日志了
		// 所以不需要再次打印
		noWFLog = true

		return fmt.Errorf("request error: %w", errRet)
	}

	return nil
}

func (r *raller) tryBindWorkLogger(i interface{}) {
	if binder, ok := i.(logit.Binder); ok {
		binder.SetLogger(r.WorkLogger())
	}
}

func (r *raller) GoRAL(ctx context.Context, name interface{}, req Request, resp Response, opts ...ROption) <-chan error {
	cherr := make(chan error)
	go func(cherr chan<- error) {
		defer func() {
			if re := recover(); re != nil {
				cherr <- fmt.Errorf("goral panic %v", re)
			}
		}()
		cherr <- r.RAL(ctx, name, req, resp, opts...)
	}(cherr)
	return cherr
}

// RALParam param for MultiRAL, equal to RAL input params
type RALParam struct {
	Ctx        context.Context
	ServerName interface{}
	Req        Request
	Resp       Response
	Opts       []ROption

	// execError will be nil if everything is ok
	err  error
	lock sync.Mutex
}

// Error will be nil if everything is ok
func (rp *RALParam) Error() error {
	rp.lock.Lock()
	tmp := rp.err
	rp.lock.Unlock()
	return tmp
}

func (rp *RALParam) setError(err error) {
	rp.lock.Lock()
	rp.err = err
	rp.lock.Unlock()
}

const (
	multiRALChanOpen   = true
	multiRALChanClosed = false
)

// MultiRAL Exec RAL concurrently.
// return nil if all req not timeout
// you sholud get the single RAL result by  RALParam.ExecErr
func (r *raller) MultiRAL(ctx context.Context, reqs map[interface{}]*RALParam) error {
	reminderTaskCount := len(reqs)
	if reminderTaskCount == 0 {
		return nil
	}

	waitingChan := make(chan struct{}, reminderTaskCount)
	chanState := multiRALChanOpen
	var mutex sync.Mutex

	for name, param := range reqs {
		if param == nil { // are U joking ?
			reminderTaskCount--
			continue
		}

		go func(param *RALParam, name interface{}) {
			param.setError(ErrExecTimeout)
			err := r.RAL(param.Ctx, param.ServerName, param.Req, param.Resp, param.Opts...)

			mutex.Lock()
			if chanState == multiRALChanOpen {
				waitingChan <- struct{}{}
				// abandon new err, maybe both got rsp and ErrExecTimeout
				param.setError(err)
			}

			mutex.Unlock()
		}(param, name)
	}

	if reminderTaskCount == 0 {
		close(waitingChan)
		return nil
	}

	clear := func() {
		mutex.Lock()
		chanState = multiRALChanClosed
		close(waitingChan)
		mutex.Unlock()
	}

	for {
		select {
		case <-ctx.Done():
			clear()
			return ctx.Err()
		case <-waitingChan:
			reminderTaskCount--
			if reminderTaskCount == 0 {
				clear()
				return nil
			}
		}
	}
}

func (r *raller) ServiceMapper() servicer.Mapper {
	return r.sm
}

func (r *raller) WorkLogger() logit.Logger {
	return r.withWorkLogger.Logger()
}

func (r *raller) setWorkLogger(l logit.Logger) {
	r.withWorkLogger.SetLogger(l)
}

// InitDefault 默认的初始化，用于初始化DefaultRaller
// 在启动阶段，初始化一次即可
func InitDefault(ctx context.Context) error {
	if DefaultRaller != nil {
		return nil
	}
	if err := servicer.InitDefault(ctx); err != nil {
		return err
	}
	DefaultRaller = NewRaller(servicer.DefaultMapper, servicer.DefaultLogger, servicer.DefaultWorkerLogger)
	return nil
}
