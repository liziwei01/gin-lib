/*
 * @Author: liziwei01
 * @Date: 2022-03-21 21:22:39
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-03-21 22:38:21
 * @Description: file content
 */
package oss

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/gogf/gf/util/gconv"
)

func TestPut(t *testing.T) {
	ctx := context.Background()
	client, err := GetClient(ctx, "oss_lib_image")
	if err != nil {
		t.Error(err)
	}
	fileReader := bytes.NewReader([]byte("Hello, world!"))
	err = client.Put(ctx, "lib-image", "hw.jpg", fileReader)
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	client, err := GetClient(ctx, "oss_lib_image")
	if err != nil {
		t.Error(err)
	}
	fileReader, err := client.Get(ctx, "lib-image", "hw.jpg")
	if err != nil {
		t.Error(err)
	}
	fileContent, err := ioutil.ReadAll(fileReader)
	if err != nil {
		t.Error(err)
	}
	strContent := gconv.String(fileContent)
	if strContent != "Hello, world!" {
		t.Error(fileContent)
	}
}

func TestGetURL(t *testing.T) {
	ctx := context.Background()
	client, err := GetClient(ctx, "oss_lib_image")
	if err != nil {
		t.Error(err)
	}
	fileURL, err := client.GetURL(ctx, "lib-image", "hw.jpg")
	if err != nil {
		t.Error(err)
	}
	t.Log(fileURL)
}

func TestDel(t *testing.T) {
	ctx := context.Background()
	client, err := GetClient(ctx, "oss_lib_image")
	if err != nil {
		t.Error(err)
	}
	err = client.Del(ctx, "lib-image", "hw.jpg")
	if err != nil {
		t.Error(err)
	}
}
