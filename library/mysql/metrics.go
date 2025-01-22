
package mysql

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
)

// NewMetricsObserverFunc 注册prometheus的统计指标获取的hook
//
// 会采集 "method", "errno","cost" 这3个指标，类型为Histogram
// 参数 name 是服务的服务名称
// 该功能默认不集成，需要自己通过AddHook添加
func NewMetricsObserverFunc(rg prometheus.Registerer) ObserverFunc {
	var cs *prometheus.SummaryVec
	var once sync.Once

	return func(ctx context.Context, client Client, ev Event) {
		once.Do(func() {
			cs = prometheus.NewSummaryVec(prometheus.SummaryOpts{
				Namespace:   "mysql",
				Subsystem:   client.Name(),
				Name:        "request",
				Help:        "mysql request detail",
				ConstLabels: nil,
				Objectives: map[float64]float64{
					0.8:  0.05,
					0.9:  0.02,
					0.95: 0.01,
					0.99: 0.001,
				},
				MaxAge:     0,
				AgeBuckets: 0,
				BufCap:     0,
			}, []string{"method", "errno"})

			if rg != nil {
				rg.MustRegister(cs)
			} else {
				prometheus.MustRegister(cs)
			}
		})

		var code int
		if ev.Err != nil {
			var mye *mysqlDriver.MySQLError
			if errors.As(ev.Err, &mye) {
				code = int(mye.Number)
			}
		}
		dur := float64(ev.End.Sub(ev.Start).Nanoseconds()) / float64(time.Millisecond)
		cs.WithLabelValues(ev.Type.String(), strconv.Itoa(code)).Observe(dur)
	}
}
