package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/liziwei01/gin-lib/library/logit"
	"github.com/liziwei01/gin-lib/library/net/ral"
	"github.com/liziwei01/gin-lib/library/net/servicer"
)

// DefaultLogger 全局的、默认的logger
var DefaultLogger logit.Logger

// Client MySQL 客户端
//
//	一个 Client 会包含 N 个 *sql.DB对象
//	N = 当前服务对应的实例数(如一个 BNS 有10个实例，则 N=10)
//	每次调用的时候，会选择出其中一个 *sql.DB 执行对应的方法
type Client interface {
	// Name 名称
	Name() string

	// BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)

	PingContext(ctx context.Context) error

	PrepareContext(ctx context.Context, query string) (Stmt, error)

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error)

	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row

	// 统计信息
	Stats() sql.DBStats

	// Raw 获取其中一个ip对应的 *sql.DB 实例
	// 注意，获取的实例，最好不要作为单例在之后的程序中一直使用
	// 否则，若 ip 下线后，该实例将不可用
	Raw(ctx context.Context) *sql.DB

	Servicer() servicer.Servicer

	// Logger 获取当前Client的Logger
	Logger() logit.Logger
}

// client mysql client 实现
// 注意，在client 内部不要直接操作添加日志
// 即 不要在client的代码里出现 logit.AddFields的逻辑
type client struct {
	logit.WithLogger

	sqlDB         *sql.DB
	mapper        servicer.Mapper
	servicer      servicer.Servicer
	config        *Config
	observerFuncs []ObserverFunc

	// 每个请求执行前，新生成context的方法
	newContextFunc func(ctx context.Context) context.Context
}

// ClientOption client 选项
type ClientOption interface {
	apply(*client)
}

type funcClientOption struct {
	f func(*client)
}

func (fdo *funcClientOption) apply(do *client) {
	fdo.f(do)
}

func newFuncClientOption(f func(*client)) *funcClientOption {
	return &funcClientOption{
		f: f,
	}
}

// OptServicerMapper 配置选项，设置servicer.Mapper
func OptServicerMapper(sm servicer.Mapper) ClientOption {
	return newFuncClientOption(func(c *client) {
		c.mapper = sm
	})
}

// OptLogger logger 配置选项
func OptLogger(logger logit.Logger) ClientOption {
	return newFuncClientOption(func(c *client) {
		c.SetLogger(logger)
	})
}

// OptObserver 执行结果观察器
func OptObserver(hks ...ObserverFunc) ClientOption {
	return newFuncClientOption(func(c *client) {
		for _, fn := range hks {
			c.addObserverFunc(fn)
		}
	})
}

// OptNewContextFunc Client每次新请求生成context的函数逻辑
//
// 若不设置，会使用默认的逻辑，当前默认逻辑是使用ral的字段进行日志初始化
func OptNewContextFunc(fn func(ctx context.Context) context.Context) ClientOption {
	return newFuncClientOption(func(c *client) {
		c.newContextFunc = fn
	})
}

// DefaultObserverFuncs 默认的观察器，已包含打印日志的组件
var DefaultObserverFuncs = []func() ObserverFunc{
	NewLogObserverFunc,
}

func optDefault() ClientOption {
	return newFuncClientOption(func(c *client) {
		// 优先使用mysql 全局默认 logger
		if c.Logger() == nil && DefaultLogger != nil {
			c.SetLogger(DefaultLogger)
		}

		// 默认使用ral 的 日志
		if c.Logger() == nil && ral.DefaultRaller != nil {
			c.SetLogger(ral.DefaultRaller.WorkLogger())
		}

		if c.mapper == nil {
			c.mapper = servicer.DefaultMapper
		}

		for _, fn := range DefaultObserverFuncs {
			c.addObserverFunc(fn())
		}
	})
}

// NewClient 创建一个新的client
//
// 创建好了之后，client应该作为一个单例使用
func NewClient(serviceName interface{}, opt ...ClientOption) (Client, error) {
	c := &client{}
	for _, o := range opt {
		o.apply(c)
	}

	optDefault().apply(c)

	servicer := c.mapper.Servicer(serviceName)

	if servicer == nil {
		return nil, fmt.Errorf("serviceName=%v not found", serviceName)
	}
	c.servicer = servicer

	config, err := getConfig(servicer.Option())
	if err != nil {
		return nil, err
	}

	mapper.Store(servicer.Name(), servicer)

	c.config = config

	return c, nil

}
func (dbClient *client) Name() string {
	return dbClient.servicer.Name()
}

func (dbClient *client) addObserverFunc(h ObserverFunc) {
	dbClient.observerFuncs = append(dbClient.observerFuncs, h)
}

func (dbClient *client) Raw(ctx context.Context) *sql.DB {
	return dbClient.sqlDB
}

// connect 获取一个数据库连接
func (dbClient *client) connect(ctx context.Context, addr net.Addr) (*sql.DB, error) {
	// sql.Open 在参数正常的情况下，总会返回一个db对象，而且err为nil
	// 需要用ping才知道是否正常的连接到db
	db, err := sql.Open(dbClient.config.DBDriver, dbClient.config.DSN(addr))
	if err != nil {
		dbClient.AutoLogger().Error(ctx, "sql.Open failed", logit.Error("error", err))
		if db != nil {
			db.Close()
		}
		return nil, err
	}

	if db != nil {
		dbClient.config.setToDb(db)
	}
	return db, err
}

func (dbClient *client) hasHook() bool {
	return len(dbClient.observerFuncs) > 0
}

func (dbClient *client) triggerHook(ctx context.Context, event *Event) {
	if event == nil || !dbClient.hasHook() {
		return
	}
	for i := 0; i < len(dbClient.observerFuncs); i++ {
		dbClient.observerFuncs[i](ctx, dbClient, *event)
	}
}

func (dbClient *client) newContext(ctxOld context.Context) context.Context {
	if dbClient.newContextFunc != nil {
		return dbClient.newContextFunc(ctxOld)
	}
	return newContext(ctxOld)
}

func newContext(ctx context.Context) context.Context {
	ctx = logit.NewContext(ctx)
	ral.InitRalStatisItems(ctx)

	logit.AddAllLevel(ctx, logit.Time(ral.LogFieldReqStartTime, time.Now()))
	return ctx
}

func (dbClient *client) withSTDDb(ctx context.Context, fn func(db *sql.DB) (*Event, error)) error {
	start := time.Now()

	ctx = dbClient.newContext(ctx)

	var event *Event
	if dbClient.hasHook() {
		defer func() {
			dbClient.triggerHook(ctx, event)
		}()
	}

	connectStart := time.Now()
	remoteAddr, db, reuseTimes, err := dbClient.dbPool.GetStdDB(ctx)
	connectStop := time.Now()

	// reuseTimes 复用次数，若是0，说明是新创建的连接

	if err != nil {
		if dbClient.hasHook() {
			event = &Event{
				Start:        start,
				RemoteAddr:   remoteAddr,
				Type:         EventConnect,
				ConnectStart: connectStart,
				ConnectEnd:   connectStop,
				SQL:          "",
				Args:         nil,
				Err:          err,
				End:          time.Now(),
			}
		}
		return err
	}

	event, err = fn(db)

	if event != nil {
		event.Start = start
		event.RemoteAddr = remoteAddr
		event.End = time.Now()
		event.ConnectStart = connectStart
		event.ConnectEnd = connectStop
		event.UsedTimes = reuseTimes
		event.Err = err
	}
	return err
}

// 没有hook的时候，不需要创建event
func (dbClient *client) newEvent(fn func() *Event) *Event {
	if dbClient.hasHook() {
		return fn()
	}
	return nil
}

func (dbClient *client) PingContext(ctx context.Context) error {
	return dbClient.withSTDDb(ctx, func(db *sql.DB) (*Event, error) {
		err := db.PingContext(ctx)
		event := dbClient.newEvent(func() *Event {
			return &Event{
				Type: EventPing,
				SQL:  "",
				Args: nil,
			}
		})
		return event, err
	})
}

func (dbClient *client) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	var rawStmt *sql.Stmt
	err := dbClient.withSTDDb(ctx, func(db *sql.DB) (*Event, error) {
		var err error
		if dbClient.config.LogIDTransport {
			query = sqlWithComment(ctx, query)
		}
		rawStmt, err = db.PrepareContext(ctx, query)
		event := dbClient.newEvent(func() *Event {
			return &Event{
				Type: EventPrepare,
				SQL:  query,
			}
		})
		return event, err
	})
	if err != nil {
		return nil, err
	}
	s := &stmt{
		raw:    rawStmt,
		client: dbClient,
		start:  time.Now(),
		uniqID: rand.Int(),
	}
	return s, nil
}

func (dbClient *client) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var rst sql.Result
	err := dbClient.withSTDDb(ctx, func(db *sql.DB) (*Event, error) {
		var err error
		if dbClient.config.LogIDTransport {
			query = sqlWithComment(ctx, query)
		}
		rst, err = db.ExecContext(ctx, query, args...)
		event := dbClient.newEvent(func() *Event {
			return &Event{
				Type: EventExec,
				SQL:  query,
				Args: args,
			}
		})
		return event, err
	})
	if err != nil {
		return nil, err
	}
	return rst, err
}

func (dbClient *client) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	var rawRows *sql.Rows
	err := dbClient.withSTDDb(ctx, func(db *sql.DB) (*Event, error) {
		var err error
		if dbClient.config.LogIDTransport {
			query = sqlWithComment(ctx, query)
		}
		rawRows, err = db.QueryContext(ctx, query, args...)
		event := dbClient.newEvent(func() *Event {
			return &Event{
				Type: EventQuery,
				SQL:  query,
				Args: args,
			}
		})
		return event, err
	})
	if err != nil {
		return nil, err
	}

	return newRows(rawRows), nil
}

func (dbClient *client) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	var rawRow *sql.Row
	err := dbClient.withSTDDb(ctx, func(db *sql.DB) (*Event, error) {
		if dbClient.config.LogIDTransport {
			query = sqlWithComment(ctx, query)
		}
		rawRow = db.QueryRowContext(ctx, query, args...)
		event := dbClient.newEvent(func() *Event {
			return &Event{
				Type: EventQueryRow,
				SQL:  query,
				Args: args,
			}
		})
		return event, nil
	})

	return &row{
		err: err,
		raw: rawRow,
	}
}

func (dbClient *client) Stats() sql.DBStats {
	return dbClient.dbPool.Stats()
}

func (dbClient *client) Servicer() servicer.Servicer {
	return dbClient.servicer
}

var _ Client = (*client)(nil)

// RawDB 原始 DB 的信息
type RawDB struct {
	DB   *sql.DB
	Addr net.Addr
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
