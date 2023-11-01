<!--
 * @Author: liziwei01
 * @Date: 2023-09-13 16:55:50
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 11:55:04
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
	SetServiceLogger(ctx)
	YourFunc(ctx)
}


func YourFunc(ctx context.Context) {
	logit.SvrLogger.Trace(ctx, "test trace", String("DebugTraceKey", "DebugOutput"))
	logit.SvrLogger.Notice(ctx, "test notice", Int("time", 1))

	err := fmt.Errorf("we got error!")
	logit.SvrLogger.Warning(ctx, "test warning", Error("err", err))
	logit.SvrLogger.Error(ctx, "test error", Error("fatal", err))
}

```
