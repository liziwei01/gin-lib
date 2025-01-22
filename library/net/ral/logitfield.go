/*
 * @Author: liziwei01
 * @Date: 2023-11-03 13:33:42
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 21:29:55
 * @Description: file content
 */
package ral

import (
	"context"
	"sync"
	"time"

	"github.com/liziwei01/gin-lib/library/env"
	"github.com/liziwei01/gin-lib/library/logit"
)

// ral 日志格式
// 以下为日志的字段名称
const (
	LogFieldAppName       = "appname"
	LogFieldURI           = "uri"
	LogFieldService       = "service"
	LogFieldReqLen        = "req_len"
	LogFieldResLen        = "res_len"
	LogFieldErrno         = "errno"
	LogFieldRetry         = "retry"
	LogFieldCost          = "cost"
	LogFieldAPI           = "api"
	LogFieldLogID         = "logid"
	LogFieldCaller        = "caller"
	LogFieldMethod        = "method"
	LogFieldProtocol      = "protocol"
	LogFieldBalance       = "balance"
	LogFieldUserIP        = "user_ip"
	LogFieldLocalIP       = "local_ip"
	LogFieldRemoteIP      = "remote_ip"
	LogFieldRemoteHost    = "remote_host"
	LogFieldUniqid        = "uniqid"
	LogFieldTalk          = "talk"
	LogFieldConnect       = "connect"
	LogFieldWrite         = "write"
	LogFieldRead          = "read"
	LogFieldPack          = "pack"
	LogFieldUnpack        = "unpack"
	LogFieldReqStartTime  = "req_start_time"
	LogFieldTalkStartTime = "talk_start_time"
	LogFieldErrmsg        = "errmsg"
	LogFieldDNS           = "dnslookup"
)

// statisticsItems 按照采集顺序
var statisItems = []string{
	LogFieldAppName,
	LogFieldURI,
	LogFieldService,
	LogFieldReqLen,
	LogFieldResLen,
	LogFieldErrno,
	LogFieldRetry,
	LogFieldCost,
	LogFieldAPI,
	LogFieldLogID,
	LogFieldCaller,
	LogFieldMethod,
	LogFieldProtocol,
	LogFieldBalance,
	LogFieldUserIP,
	LogFieldLocalIP,
	LogFieldRemoteIP,
	LogFieldRemoteHost,
	LogFieldUniqid,
	LogFieldTalk,
	LogFieldConnect,
	LogFieldWrite,
	LogFieldRead,
	LogFieldPack,
	LogFieldUnpack,
	LogFieldReqStartTime,
	LogFieldTalkStartTime,
	LogFieldErrmsg,
	LogFieldDNS,
}

var fields []logit.Field

func initFields() {
	fields = make([]logit.Field, 0, len(statisItems))
	for _, item := range statisItems {
		var field logit.Field
		switch item {
		case LogFieldAppName:
			field = logit.String(LogFieldAppName, env.AppName())
		case LogFieldLocalIP:
			field = logit.String(LogFieldLocalIP, env.LocalIP())
		case LogFieldReqLen,
			LogFieldResLen,
			LogFieldErrno:
			field = logit.Int32(item, 0)
		case LogFieldReqStartTime,
			LogFieldTalkStartTime:
			field = logit.Time(item, time.Time{})
		case LogFieldTalk,
			LogFieldConnect,
			LogFieldWrite,
			LogFieldRead,
			LogFieldPack,
			LogFieldUnpack,
			LogFieldDNS:
			field = logit.Duration(item, time.Duration(0))
		default:
			field = logit.String(item, "")
		}
		fields = append(fields, field)
	}

	for _, field := range fields {
		field.SetLevel(logit.AllLevels)
	}
}

var initOnce sync.Once

// InitRalStatisItems 日志字段，预初始化
func InitRalStatisItems(ctx context.Context) {
	initOnce.Do(func() {
		initFields()
	})
	logit.AddFields(ctx, fields...)
}
