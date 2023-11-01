/*
 * @Author: liziwei01
 * @Date: 2023-10-31 21:57:23
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 07:39:07
 * @Description: 用例
 */
package logit

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"
)


// TRACE: 2023-11-01 07:38:36 /Users/liziwei01/Desktop/OpenSource/github.com/gin-lib/library/logit/logit_test.go:29 keyT[valueT] message[test trace]
// NOTICE: 2023-11-01 07:38:36 /Users/liziwei01/Desktop/OpenSource/github.com/gin-lib/library/logit/logit_test.go:30 keyN[valueN] message[test notice]

// WARNING: 2023-11-01 07:38:36 /Users/liziwei01/Desktop/OpenSource/github.com/gin-lib/library/logit/logit_test.go:31 keyW[valueW] message[test warning]
// ERROR: 2023-11-01 07:38:36 /Users/liziwei01/Desktop/OpenSource/github.com/gin-lib/library/logit/logit_test.go:32 keyE[valueE] message[test error]


func TestWarning(t *testing.T) {
	ctx := context.Background()
	logName := "service"
	l, err := NewLogger(ctx, OptConfigFile(filepath.Join(logPath, logName+suffix)))
	if err != nil {
		t.Fatal(err)
	}
	l.Trace(ctx, "test trace", Error("keyT", fmt.Errorf("valueT")))
	l.Notice(ctx, "test notice", Error("keyN", fmt.Errorf("valueN")))
	l.Warning(ctx, "test warning", Error("keyW", fmt.Errorf("valueW")))
	l.Error(ctx, "test error", Error("keyE", fmt.Errorf("valueE")))
	time.Sleep(3 * time.Second)
}
