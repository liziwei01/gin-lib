/*
 * @Author: liziwei01
 * @Date: 2023-11-04 13:19:54
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 22:21:46
 * @Description: 启动前业务必须调用MustLoad加载全部的servicer配置文件，其他组件内部只需要注册回调函数即可在读取配置文件时收到通知
 */
package servicer

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/liziwei01/gin-lib/library/env"
	"github.com/liziwei01/gin-lib/library/extension/timer"
	"github.com/liziwei01/gin-lib/library/logit"
)

// LoadOption 加载Servicer 的配置选项
type LoadOption interface {
	apply(*loaderOption)
}

type funcLoadOption struct {
	f func(*loaderOption)
}

func (fdo *funcLoadOption) apply(do *loaderOption) {
	fdo.f(do)
}

func newFuncLoadOption(f func(*loaderOption)) *funcLoadOption {
	return &funcLoadOption{
		f: f,
	}
}

// LoadOptMapper 加载到那个mapper,若不设置，将加载到 DefaultMapper
func LoadOptMapper(sm Mapper) LoadOption {
	return newFuncLoadOption(func(option *loaderOption) {
		option.mapper = sm
	})
}

// LoadOptFiles 加载的配置文件
func LoadOptFiles(files []string) LoadOption {
	return newFuncLoadOption(func(option *loaderOption) {
		option.files = files
	})
}

// LoadOptFilesGlob 通过 filepath.Glob(pattern) 方式获取待加载的文件
func LoadOptFilesGlob(pattern string, ignoreErr bool) LoadOption {
	matches, err := filepath.Glob(pattern)
	if err != nil && !ignoreErr {
		panic("LoadOptFilesGlob(" + pattern + ") with error: " + err.Error())
	}
	return LoadOptFiles(matches)
}

// LoadOptAllowStartFail 是否允许在启动阶段Start失败
//
//	默认是允许，如采用bns，短暂的bns服务访问失败，将启动失败
//	启动失败后，将在后台启动任务继续每5秒尝试启动，直到成功
//	若为false，Start失败，会panic
func LoadOptAllowStartFail(allow bool) LoadOption {
	return newFuncLoadOption(func(option *loaderOption) {
		option.allowStartFail = allow
	})
}

// LoadOptTryStartDuration 设置Server Start失败后，异步重试Start的间隔
//
//	默认值是5s
//	若设置为<1的值，同时会将AllowStartFail 设置为false，即不允许启动失败，会panic
func LoadOptTryStartDuration(dur time.Duration) LoadOption {
	return newFuncLoadOption(func(option *loaderOption) {
		if dur < 1 {
			option.allowStartFail = false
		}
		option.tryStartDuration = dur
	})
}

// LoadOptCheckFileName 配置选项：是否检查配置文件名称和配置的Name值是否一致
//
// MustLoad默认是会检查的
func LoadOptCheckFileName(check bool) LoadOption {
	return newFuncLoadOption(func(option *loaderOption) {
		option.checkFileName = check
	})
}

type loaderOption struct {
	mapper           Mapper        // 加载到这个mapper
	files            []string      // 待加载的文件列表
	allowStartFail   bool          // 是否允许启动失败，失败后将异步尝试加载
	tryStartDuration time.Duration // 异步加载的时间间隔
	checkFileName    bool          // 是否检查文件名和 Name是否一致

	onceTask *timer.OnceSuccess
}

func (opt *loaderOption) Logger() logit.Logger {
	logger := opt.mapper.Logger()
	if logger != nil {
		return logger
	}
	return logit.NopLogger
}

// MustLoad 将指定目录servicer配置加载到对应的mapper，若有异常会panic
//
// 这个是一个默认实现，若业务有别的述求，下来策略不能满足，可自己实现
//
// 可能panic的情况：
//
//	1.若是配置文件异常，导致创建Servicer失败，会panic
//	2.配置文件的文件名 和 配置的name不一致，若不一致会panic (默认会检查，可以关闭：LoadOptCheckFileName(false))
//		推荐开启该检查，这样可维护性更好，同时也能减少出错的机会（如配置了相同的Name）
//	3.服务启动失败，会panic（默认不会panic，若设置 LoadOptAllowStartFail(false)，将panic）
//	4.出现重复的服务名称，如多个配置文件中，配置了相同的Name。
func MustLoad(ctx context.Context, opts ...LoadOption) {
	opt := &loaderOption{
		mapper:           DefaultMapper,
		files:            nil,
		allowStartFail:   true,
		tryStartDuration: 5 * time.Second,
		checkFileName:    true,
	}

	for _, o := range opts {
		o.apply(opt)
	}

	if opt.allowStartFail {
		opt.onceTask = timer.NewOnceSuccess(opt.tryStartDuration)
	}

	logger := opt.Logger()

	logger.Notice(ctx, "load files", logit.AutoField("files", opt.files))

	if len(opt.files) == 0 {
		return
	}

	for _, confPath := range opt.files {
		opt.mustLoadOneFile(ctx, confPath)
	}

	if opt.allowStartFail {
		go opt.onceTask.Start(ctx)
	}
}

// mustLoadOneFile 从文件加载servicer
// 启动失败：
// 1.若允许首次加载失败，会放入后台异步尝试加载直到成功
// 2.若不允许首次加载失败，将panic
func (opt *loaderOption) mustLoadOneFile(ctx context.Context, confPath string) {
	ctx = logit.ForkContext(ctx)

	logger := opt.Logger()
	logit.AddAllLevel(ctx, logit.String("confName", filepath.Base(confPath)))

	srv, err := opt.loadOneServicer(ctx, confPath)
	if err != nil {
		logger.Fatal(ctx, "load server failed", logit.Error("error", err))
		panic(err.Error())
	}

	if err = opt.mapper.AddServicer(srv.Name(), srv); err != nil {
		logger.Fatal(ctx, "AddServicer failed: %w", logit.Error("error", err))
		panic(err.Error())
	}

	if err = srv.Start(ctx); err != nil {
		if !opt.allowStartFail {
			panic("start " + srv.Name() + " failed: " + err.Error())
		}

		ctxJob := logit.ForkContext(ctx)
		// 启动阶段失败了，放入后台去尝试启动，直到成功
		opt.onceTask.AddJob(func() error {
			return opt.asyncLoadAndAddJob(ctxJob, confPath)
		})
	}
}

// asyncLoadAndAddJob 异步加载的任务
//
// 包含动作：1.从文件加载初始化Servicer 2.启动 3.添加到mapper
// 因为这个时候服务加载失败，对应的服务还是不可用，所以打印的日志是 Fatal 等级的
func (opt *loaderOption) asyncLoadAndAddJob(ctx context.Context, confPath string) error {
	ctx = logit.ForkContext(ctx)
	logit.AddAllLevel(ctx, logit.String("task", "asyncLoadAndAdd"))

	logger := opt.Logger()

	// 加载创建一个全新的servicer
	srv, errLoad := opt.loadOneServicer(ctx, confPath)
	if errLoad != nil {
		logger.Fatal(ctx, "async load servicer failed", logit.Error("error", errLoad))
		return errLoad
	}

	logit.AddAllLevel(ctx, logit.String("servicer", srv.Name()))

	if errStart := srv.Start(ctx); errStart != nil {
		logger.Fatal(ctx, "async Start servicer failed", logit.Error("error", errStart))
		return errStart
	}

	if old := opt.mapper.SetServicer(srv.Name(), srv); old != nil {
		errStop := old.Stop()
		if errStop != nil {
			logger.Fatal(ctx, "stop old servicer failed")
		} else {
			logger.Notice(ctx, "stop old servicer success")
		}
	}
	return nil
}

// loadOne 从文件加载一个服务，服务是未启动的
func (opt *loaderOption) loadOneServicer(ctx context.Context, confPath string) (Servicer, error) {
	if fp, err := filepath.Abs(confPath); err == nil {
		confPath = fp
	}
	name := filepath.Base(confPath)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]

	srv, err := NewWithConfigName(opt.Logger(), env.Default, confPath)

	if err != nil {
		return nil, fmt.Errorf(" NewWithConfigName(%s) failed, error:%w", confPath, err)
	}

	if opt.checkFileName {
		// 检查 ral的service 配置里的文件名 和其配置的服务名是否保持一致
		// 名字保持一致，更方便运维管理
		if got := srv.Name(); got != name {
			return nil, fmt.Errorf("config's file name (%q) not eq service name(%q)", name, got)
		}
	}

	return srv, nil
}
