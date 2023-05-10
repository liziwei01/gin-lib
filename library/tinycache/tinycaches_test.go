/*
 * @Author: liziwei01
 * @Date: 2021-07-23 15:53:34
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-05-10 00:07:32
 * @Description: tinycache unit test
 * @FilePath: /gdp-config-platform/library/tinycaches/tinycaches_test.go
 */
package tinycache

import (
	"context"
	"testing"
	"time"
)

func TestTinycache(t *testing.T) {
	ctx := context.Background()
	client, err := GetClient(ctx, "asc")
	if err == nil {
		t.Fatalf("create not exist client success")
	}
	client, err = GetClient(ctx, "tnc_lib")
	if err != nil {
		t.Fatalf("create client failed")
	}
	idx1 := "myname"
	val1 := "liziwei01"
	client.Set(idx1, val1)
	
	idx2 := "myage"
	val2 := "21"
	client.Set(idx2, val2, time.Second)

	value1 := client.Get(idx1)
	value2 := client.Get("not exist")
	if value1 == "" {
		t.Error("1!ok")
	}
	if value2 != "" {
		t.Error("2!ok")
	}

	time.Sleep(2 * time.Second)
	value3 := client.Get(idx2)
	if value3 != "" {
		t.Errorf("not expired")
	}
}
