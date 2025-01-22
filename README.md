<!--
 * @Author: liziwei01
 * @Date: 2023-09-13 16:55:50
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 03:03:43
 * @Description: file content
-->
# gin-lib

## Description

Go application framework practice based on gin framework
Provides common components and functionality out of the box

## Installation

```bash
go get github.com/liziwei01/gin-lib
```

## Usage

### logit
logit is a log module designed for your back-end to use

./conf/logit/ral.toml needed for modules inside gin-lib
./conf/logit/ral-worker.toml needed for requests sent by your back-end. For example, after using ral to send an http request, record the execution time and other information. If you use redis or mysql client to send a request, record the execution time, status and other information.
./conf/logit/service.toml needed for requests from front-end

```golang
package yourpackage

import (
	"context"
	"github.com/liziwei01/gin-lib/library/logit"
)

// some sample output
// TRACE: 2023-11-01 07:38:36 /Users/liziwei01/Desktop/OpenSource/github.com/gin-lib/library/logit/logit_test.go:29 keyT[valueT] message[test trace]
// NOTICE: 2023-11-01 07:38:36 /Users/liziwei01/Desktop/OpenSource/github.com/gin-lib/library/logit/logit_test.go:30 keyN[valueN] message[test notice]

// WARNING: 2023-11-01 07:38:36 /Users/liziwei01/Desktop/OpenSource/github.com/gin-lib/library/logit/logit_test.go:31 keyW[valueW] message[test warning]
// ERROR: 2023-11-01 07:38:36 /Users/liziwei01/Desktop/OpenSource/github.com/gin-lib/library/logit/logit_test.go:32 keyE[valueE] message[test error]

func main() {
	// Set the default logger before use
	ctx := context.Background()
	logit.SetServiceLogger(ctx)
	YourFunc(ctx)
}


func YourFunc(ctx context.Context) {
	logit.SrvLogger.Trace(ctx, "test trace", String("DebugTraceKey", "DebugOutput"))
	logit.SrvLogger.Notice(ctx, "test notice", Int("time", 1))

	err := fmt.Errorf("we got error!")
	logit.SrvLogger.Warning(ctx, "test warning", Error("err", err))
	logit.SrvLogger.Error(ctx, "test error", Error("fatal", err))
}

```
