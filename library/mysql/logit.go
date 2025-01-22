package mysql

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"strings"
	"sync"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/liziwei01/gin-lib/library/logit"
	"github.com/liziwei01/gin-lib/library/net/connector"
	"github.com/liziwei01/gin-lib/library/net/ral"
)

// 这里是日志相关的功能

// NewLogObserverFunc 创建日志打印的方法
func NewLogObserverFunc() ObserverFunc {
	var field1, field2, field3, field4, field5 logit.Field
	var once sync.Once
	var config *Config

	return func(ctx context.Context, client Client, ev Event) {
		logger := client.Logger()
		if logger == nil {
			return
		}

		once.Do(func() {
			field1 = logit.String(ral.LogFieldService, client.Servicer().Name())
			field2 = logit.String(ral.LogFieldCaller, "MySQL")
			field3 = logit.String(ral.LogFieldProtocol, "mysql")

			cc := client.Servicer().Connector()
			if cc != nil {
				if s, ok := cc.(connector.HasStrategy); ok {
					field4 = logit.AutoField(ral.LogFieldBalance, s.Strategy())
				}
			}

			if field4 == nil {
				field4 = logit.String(ral.LogFieldBalance, "unknown")
			}

			// mysql 非常特殊，不能retry
			field5 = logit.String(ral.LogFieldRetry, "0/0")

			var err error
			config, err = getConfig(client.Servicer().Option())
			if err != nil {
				panic("getConfig failed: " + err.Error())
			}
		})

		fields := []logit.Field{
			field1,
			field2,
			field3,
			field4,
			eventTypeField(ev.Type),
			logit.Duration(ral.LogFieldCost, ev.End.Sub(ev.Start)),
			logit.Error(ral.LogFieldErrmsg, ev.Err),
			field5,

			// 请求开始时间，包括连接时间
			logit.Time(ral.LogFieldReqStartTime, ev.Start),

			// 连接耗时
			logit.Duration(ral.LogFieldConnect, ev.ConnectEnd.Sub(ev.ConnectStart)),

			// 连接后，发送请求开始时间
			logit.Time(ral.LogFieldTalkStartTime, ev.ConnectEnd),

			// 获取不到实际发送的长度，所以使用sql的长度
			logit.Int(ral.LogFieldReqLen, len(ev.SQL)),

			logit.String(ral.LogFieldRemoteIP, logitAddrIP(ev.RemoteAddr)),
			// logit.String(ral.LogFieldRemoteIDC, gaddr.RemoteIDC(ev.RemoteAddr)),
		}
		if extFields := logSQLFields(ev, config); len(extFields) > 0 {
			fields = append(fields, extFields...)
		}

		// tx 开启后总使用时间
		if ev.TxID > 0 {
			fields = append(fields,
				logit.Int("tx_id", ev.TxID),
				logit.Duration("tx_duration", ev.End.Sub(ev.TxStart)),
			)
		}

		// callDepth = 4 可以定位到执行到client的方法名
		logger.Output(ctx, logit.NoticeLevel, 4, "", fields...)
		if ev.Err != nil {
			logger.Output(ctx, logit.ErrorLevel, 4, "", fields...)
		}
	}
}

func logitAddrIP(addr net.Addr) string {
	if addr == nil {
		return ""
	}
	return addr.String()
}

const errNoOther uint16 = 1000

func logSQLFields(ev Event, config *Config) []logit.Field {
	var n int
	if ev.Err != nil {
		n++
	}
	if config.SQLLogLen != 0 {
		n++
	}
	if config.SQLArgsLogLen != 0 {
		n++
	}
	if n == 0 {
		return nil
	}

	fields := make([]logit.Field, 0, n)

	if ev.Err != nil {
		var mye *mysqlDriver.MySQLError
		if errors.As(ev.Err, &mye) {
			fields = append(fields, logit.Uint16(ral.LogFieldErrno, mye.Number))
		} else {
			fields = append(fields, logit.Uint16(ral.LogFieldErrno, errNoOther))
		}
	}

	// sql 详情 和 参数按需打印
	// 特别的，当发生错误的时候，总是打印sql日志信息
	if ev.Err != nil || config.SQLLogLen == -1 || (config.SQLLogLen > 0 && config.SQLLogLen >= len(ev.SQL)) {
		fields = append(fields, logit.String("sql", ev.SQL))
	} else if config.SQLLogLen > 0 {
		fields = append(fields, logit.String("sql", ev.SQL[0:config.SQLLogLen]))
	}

	if ev.Err != nil || config.SQLArgsLogLen == -1 || (config.SQLArgsLogLen > 0 && config.SQLArgsLogLen >= len(ev.Args)) {
		fields = append(fields, logit.AutoField("sql_args", ev.Args))
	} else if config.SQLArgsLogLen > 0 {
		fields = append(fields, logit.AutoField("sql_args", ev.Args[0:config.SQLArgsLogLen]))
	}
	return fields
}

// 给sql语句添加注释信息，如logid，可以传递到mysql server
// http://wiki.baidu.com/pages/viewpage.action?pageId=629313927
func sqlWithComment(ctx context.Context, sql string) string {
	cmt := make(map[string]string)
	var sqlNoCmt string
	// 会先去解析sql里的注释，然后将logid字段补充上
	if strings.HasPrefix(sql, "/*") {
		idx := strings.Index(sql, "*/")
		if idx < 0 {
			return sql
		}
		if err := json.Unmarshal([]byte(sql[2:idx]), &cmt); err != nil {
			return sql
		}
		sqlNoCmt = sql[idx+2:]
	}
	field := logit.FindRequestIDField(ctx)
	if field == nil {
		return sql
	}

	logID, ok := field.Value().(string)
	if !ok {
		return sql
	}
	if len(logID) > 64 {
		cmt["log_id"] = logID[0:64]
	} else {
		cmt["log_id"] = logID
	}

	// 字符串，注释标识，必选
	// 以让menu 和 bdproxy识别该注释
	if _, has := cmt["xdb_comment"]; !has {
		cmt["xdb_comment"] = "1"
	}
	bf, err := json.Marshal(cmt)
	if err != nil {
		return sql
	}
	if len(sqlNoCmt) == 0 {
		sqlNoCmt = sql
	}

	var builder strings.Builder
	builder.Grow(3 + 4 + len(bf) + len(sqlNoCmt))
	builder.Write([]byte("/* "))
	builder.Write(bf)
	builder.Write([]byte(" */ "))
	builder.WriteString(sqlNoCmt)
	return builder.String()
}
