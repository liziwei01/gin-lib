package mysql

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/liziwei01/gin-lib/library/logit"
	"github.com/liziwei01/gin-lib/library/net/ral"
)

// ObserverFunc 执行完后后的回调函数
type ObserverFunc func(ctx context.Context, client Client, ev Event)

// Event 是一次请求的内容的实体
//
// 如 执行一次 Query 查询就会产生一个event
type Event struct {
	Start      time.Time // 当前event 的开始时间
	End        time.Time // 当前event 的结束时间
	RemoteAddr net.Addr  // 远端地址

	ConnectStart time.Time // 获取连接的开始时间
	ConnectEnd   time.Time // 获取连接的结束时间

	Type EventType
	SQL  string
	Args []interface{}
	Err  error

	StmtID    int       // 预处理语句的id，当前只是普通的随机数，每次新生成预处理语句会分配一个新值
	StmtStart time.Time // 当前预处理语句的开始时间

	TxID    int       // 事务的id，当前只是普通的随机数，每次新生成事务会分配一个新值
	TxStart time.Time // 当前事务的开始时间
}

// String 序列化
func (e Event) String() string {
	bf, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(bf)
}

// EventType 按数据库的操作进行分类
type EventType int8

// String 英文名称，如 Ping
func (et EventType) String() string {
	return eventNames[et]
}

// event 事件类型
const (
	EventPing EventType = iota
	EventTxBegin
	EventTxCommit
	EventTxRollback
	EventPrepare
	EventExec
	EventQueryRow
	EventQuery
	EventConnect
)

var eventNames = map[EventType]string{
	EventPing:       "Ping",
	EventTxBegin:    "TxBegin",
	EventTxCommit:   "TxCommit",
	EventTxRollback: "TxRollback",
	EventPrepare:    "Prepare",
	EventExec:       "Exec",
	EventQueryRow:   "QueryRow",
	EventQuery:      "Query",
	EventConnect:    "Connect",
}

var eventTypeFields = map[EventType]logit.Field{}

func eventTypeField(et EventType) logit.Field {
	return eventTypeFields[et]
}

func init() {
	for et, en := range eventNames {
		eventTypeFields[et] = logit.String(ral.LogFieldMethod, en)
	}
}
