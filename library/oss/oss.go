/*
 * @Author: liziwei01
 * @Date: 2022-03-20 18:17:39
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-03-24 23:00:27
 * @Description: file content
 */
package oss

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/liziwei01/gin-lib/library/logit"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func (c *client) Get(ctx context.Context, bucket string, objectKey string) (*bytes.Reader, error) {
	ossBucket, err := c.connect(ctx, bucket)
	if err != nil {
		logit.Logger.Error("oss.connect: %+v", err.Error())
		return nil, err
	}
	file, err := ossBucket.GetObject(objectKey)
	defer file.Close()
	if err != nil {
		logit.Logger.Error("oss.Get: %+v", err.Error())
		return nil, err
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logit.Logger.Error("oss.Get.ReadAll: %+v", err.Error())
		return nil, err
	}
	logit.Logger.Info("oss.Get: %s", objectKey)
	return bytes.NewReader(fileBytes), nil
}

func (c *client) Put(ctx context.Context, bucket string, objectKey string, fileReader *bytes.Reader) error {
	ossBucket, err := c.connect(ctx, bucket)
	if err != nil {
		logit.Logger.Error("oss.connect: %+v", err.Error())
		return err
	}
	err = ossBucket.PutObject(objectKey, fileReader)
	if err != nil {
		logit.Logger.Error("oss.Put: %+v", err.Error())
		return err
	}
	logit.Logger.Info("oss.Put: %s", objectKey)
	return nil
}

func (c *client) Del(ctx context.Context, bucket string, objectKey string) error {
	ossBucket, err := c.connect(ctx, bucket)
	if err != nil {
		logit.Logger.Error("oss.connect: %+v", err.Error())
		return err
	}
	err = ossBucket.DeleteObject(objectKey)
	if err != nil {
		logit.Logger.Error("oss.Del: %+v", err.Error())
		return err
	}
	logit.Logger.Info("oss.Del: %s", objectKey)
	return nil
}

func (c *client) GetURL(ctx context.Context, bucket string, objectKey string) (string, error) {
	ossBucket, err := c.connect(ctx, bucket)
	if err != nil {
		logit.Logger.Error("oss.connect: %+v", err.Error())
		return "", err
	}
	url, err := ossBucket.SignURL(objectKey, oss.HTTPGet, 60)
	if err != nil {
		logit.Logger.Error("oss.GetURL: %+v", err.Error())
		return "", err
	}
	logit.Logger.Info("oss.GetURL: %s", objectKey)
	return url, nil
}
