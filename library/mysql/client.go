/*
 * @Author: liziwei01
 * @Date: 2022-03-09 19:26:04
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-04-17 15:39:35
 * @Description: file content
 */
package mysql

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/didi/gendry/manager"
)

var (
	// 初始化互斥锁
	mu sync.Mutex
)

type Client interface {
	// Query 查询, 返回数据在data, 内有didi builder, 结构体使用tag `ddb`
	Query(ctx context.Context, tableName string, where map[string]interface{}, columns []string, data interface{}) error

	// Insert 内有didi builder, 结构体使用tag `ddb`
	Insert(ctx context.Context, tableName string, data []map[string]interface{}) (sql.Result, error)

	InsertIgnore(ctx context.Context, tableName string, data []map[string]interface{}) (sql.Result, error)

	InsertReplace(ctx context.Context, tableName string, data []map[string]interface{}) (sql.Result, error)

	InsertOnDuplicate(ctx context.Context, tableName string, data []map[string]interface{}, update map[string]interface{}) (sql.Result, error)

	Update(ctx context.Context, tableName string, where map[string]interface{}, update map[string]interface{}) (sql.Result, error)

	Delete(ctx context.Context, tableName string, where map[string]interface{}) (sql.Result, error)

	// ExecRaw 拼接的原生sql语句
	ExecRaw(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)

	connect(ctx context.Context) (*sql.DB, error)
	open() (*sql.DB, error)

	name() string
	writeTimeOut() int
	readTimeOut() int
	retry() int
	host() string
	port() int
	username() string
	password() string
	dbname() string
	dbdriver() string
	charset() string
	collation() string
	timeout() int
	sqlloglen() int
}

type client struct {
	conf *Config
	db   *sql.DB
}

func (c *client) connect(ctx context.Context) (*sql.DB, error) {
	var err error
	if c.db != nil {
		if err = c.db.PingContext(ctx); err == nil {
			return c.db, nil
		}
	}
	mu.Lock()
	defer mu.Unlock()
	c.db, err = c.open()
	return c.db, err
}

func (c *client) open() (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)
	// 内含 retry 2
	db, err = manager.New(c.dbname(), c.username(), c.password(), c.host()).Set(
		manager.SetCharset(c.charset()),
		manager.SetAllowCleartextPasswords(true),
		manager.SetAllowNativePasswords(true),
		manager.SetInterpolateParams(true),
		manager.SetAllowAllFiles(true),
		manager.SetParseTime(true),
		manager.SetLoc(time.Local.String()),
		manager.SetTimeout(time.Duration(c.timeout())*time.Millisecond),
		manager.SetReadTimeout(time.Duration(c.readTimeOut())*time.Millisecond),
		manager.SetWriteTimeout(time.Duration(c.writeTimeOut())*time.Millisecond),
		manager.SetCollation(c.collation()),
	).Port(c.port()).Open(true)
	return db, err
}

func New(config *Config) Client {
	c := &client{
		conf: config,
		db:   nil,
	}
	return c
}

func NewDefault() Client {
	c := &client{
		conf: nil,
		db:   nil,
	}
	return c
}

func (c *client) name() string {
	return c.conf.Name
}

func (c *client) writeTimeOut() int {
	return c.conf.WriteTimeOut
}

func (c *client) readTimeOut() int {
	return c.conf.ReadTimeOut
}

func (c *client) retry() int {
	return c.conf.Retry
}

func (c *client) dbname() string {
	return c.conf.MySQL.DBName
}

func (c *client) dbdriver() string {
	return c.conf.MySQL.DBDriver
}

func (c *client) charset() string {
	return c.conf.MySQL.Charset
}
func (c *client) collation() string {
	return c.conf.MySQL.Collation
}
func (c *client) timeout() int {
	return c.conf.MySQL.Timeout
}

func (c *client) host() string {
	return c.conf.Resource.Manual.Host
}

func (c *client) port() int {
	return c.conf.Resource.Manual.Port
}

func (c *client) username() string {
	return c.conf.MySQL.Username
}

func (c *client) password() string {
	return c.conf.MySQL.Password
}

func (c *client) sqlloglen() int {
	return c.conf.MySQL.SQLLogLen
}
