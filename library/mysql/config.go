
package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/liziwei01/gin-lib/library/extension/option"
	"github.com/liziwei01/gin-lib/library/net/servicer"
)

// 数据库ral配置
// mysql 特有部分
// [MySQL]
// Username    = "gdp_r"
// Password    = "gdp_passwd"
// DBName      = "gdp_demo"
// DBDriver    = "mysql"
// MaxOpenPerIP = 5
// MaxIdlePerIP = 5
// ConnMaxLifeTime = 5000 # 单位 ms
// SQLLogLen  =-1    # 打印sql内容，为0不打印，-1 为全部
// SQLArgsLogLen  =-1 # 打印sql参数内容，为0不打印，-1 为全部
// LogIDTransport = true # 是否sql注释传递logid

// Config 配置
type Config struct {
	Username string // 账号名
	Password string // 密码
	DBName   string // 数据库名称

	// 下面的都是可选项
	DBDriver string // 驱动名称,默认 mysql

	SQLLogLen     int // 日志里打印sql 的长度，-1 为不限制，默认0，不打印
	SQLArgsLogLen int // 日志里打印sql args 的长度，-1 为不限制，默认0，不打印

	// 每个ip最多空闲数
	MaxIdlePerIP int

	// 每个ip最多连接数
	MaxOpenPerIP int

	// 是否通过sql注释 将logid 传递给mysql server
	// http://wiki.baidu.com/pages/viewpage.action?pageId=629313930
	LogIDTransport bool

	// 连接复用时间，ms
	// 若为0 则一直复用
	ConnMaxLifeTime int

	// DSN 的 ? 后面的部分
	// 完整的是 testUserItem:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
	// DSNParams 是 tls=skip-verify&autocommit=true
	DSNParams string

	// 最终的dsn config
	dsnConfig *mysqlDriver.Config
}

// parser 解析配置内容
func (c *Config) parser() error {
	// 优先使用 DSN 反解 DSN 对象
	if c.DSNParams != "" {
		str := "u:p@tcp(host:5555)/db?" + c.DSNParams
		dsn, err := mysqlDriver.ParseDSN(str)
		if err != nil {
			return err
		}
		c.dsnConfig = dsn
	}

	if c.DBName == "" {
		return errors.New("config's DBName is empty")
	}

	if c.Username == "" {
		return errors.New("config's Username is empty")
	}

	if c.DBDriver == "" {
		c.DBDriver = "mysql"
	}

	if c.dsnConfig == nil {
		c.dsnConfig = mysqlDriver.NewConfig()
	}

	c.dsnConfig.User = c.Username
	c.dsnConfig.Passwd = c.Password
	c.dsnConfig.DBName = c.DBName

	return nil
}

// DSN 获取配置的dsn
func (c *Config) DSN(addr net.Addr) string {
	dsn := *c.dsnConfig // copy it
	dsn.Addr = strings.Join([]string{
		dsn.Addr, // 实际是 servicer name
		addr.String(),
	}, dialNameHostSeparator)
	return dsn.FormatDSN()
}

func (c *Config) setToDb(db *sql.DB) {
	if c.MaxOpenPerIP > 0 {
		db.SetMaxOpenConns(c.MaxOpenPerIP)
	}

	if c.MaxIdlePerIP > 0 {
		db.SetMaxIdleConns(c.MaxIdlePerIP)
	}

	if c.ConnMaxLifeTime > 0 {
		db.SetConnMaxLifetime(time.Duration(c.ConnMaxLifeTime) * time.Millisecond)
	}
}

const (
	// 未解析的原始的配置，内容格式是map[string]interface{}
	configKeyOpt = "MySQL"

	// 已解析好的配置,内容格式是*Config
	configKeyObjectOpt = "MySQL_Config"
)

// ErrNotMySQLConfig 不是mysql 配置文件
var ErrNotMySQLConfig = fmt.Errorf("not mysql Config, miss %q section", configKeyOpt)

// init 注册配置解析钩子hook func
// Servicer 会在初始化解析配置文件的时候调用，将
func init() {
	servicer.RegisterConfigHookFunc("gin_mysql", 0, configParserHookFunc)
}

// configParserHookFunc 配置解析hook func
func configParserHookFunc(data map[string]interface{}) (map[string]interface{}, error) {
	val, has := data[configKeyOpt]
	if !has {
		return data, nil
	}

	mapValue, ok := val.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%q not type map[string]interface{}", configKeyOpt)
	}

	bf, err := json.Marshal(mapValue)
	if err != nil {
		return nil, err
	}
	var c *Config
	if err := json.Unmarshal(bf, &c); err != nil {
		return nil, err
	}

	if err := c.parser(); err != nil {
		return nil, err
	}

	c1 := make(map[interface{}]interface{}, len(data))
	for k, v := range data {
		c1[k] = v
	}

	opt := option.NewFixed(nil, c1)

	name, ok := option.String(opt, "Name", "")
	if !ok {
		return nil, fmt.Errorf("connot found Name")
	}

	dsn := c.dsnConfig
	dsn.Timeout = option.ConnTimeout(opt)
	dsn.WriteTimeout = option.WriteTimeout(opt)
	dsn.ReadTimeout = option.ReadTimeout(opt)

	// 详见 dialer.go
	// 注册自定义的Dialer
	dsn.Net = dialNetType

	// 将地址修改为服务名，方便后续自定义拨号
	dsn.Addr = name

	// 将修改后的配置放回去
	c.dsnConfig = dsn

	data[configKeyObjectOpt] = c
	return data, nil
}

// getConfig 从servicer 里面获取配置
func getConfig(opt option.Option) (*Config, error) {
	tmp := opt.Value(configKeyObjectOpt)
	if tmp == nil {
		return nil, ErrNotMySQLConfig
	}
	if c, ok := tmp.(*Config); ok {
		return c, nil
	}
	return nil, fmt.Errorf("%q not *Config: type=%T, value=%v", configKeyOpt, tmp, tmp)
}
