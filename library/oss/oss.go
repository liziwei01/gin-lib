/*
 * @Author: liziwei01
 * @Date: 2022-03-20 18:17:39
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-10-28 13:45:09
 * @Description: file content
 */
package oss

import (
	"bytes"
	"context"
	"io"

	"github.com/liziwei01/gin-lib/library/logit"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func (c *client) Get(ctx context.Context, bucket string, objectKey string) (*bytes.Reader, error) {
	ossBucket, err := c.connect(ctx, bucket)
	if err != nil {
		logit.Logger.Warn("[OSS] [connect] [requestID]=%d, [bucket]=%s, [objectKey]=%s, [err]=%s", ctx.Value("requestID"), bucket, objectKey, err.Error())
		return nil, err
	}
	file, err := ossBucket.GetObject(objectKey)
	defer file.Close()
	if err != nil {
		logit.Logger.Warn("[OSS] [GetObject] [requestID]=%d, [bucket]=%s, [objectKey]=%s, [err]=%s", ctx.Value("requestID"), bucket, objectKey, err.Error())
		return nil, err
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logit.Logger.Warn("[OSS] [ReadAll] [requestID]=%d, [bucket]=%s, [objectKey]=%s, [err]=%s", ctx.Value("requestID"), bucket, objectKey, err.Error())
		return nil, err
	}
	logit.Logger.Info("[OSS] [Get] [requestID]=%d, [bucket]=%s, [objectKey]=%s success", ctx.Value("requestID"), bucket, objectKey)
	return bytes.NewReader(fileBytes), nil
}

func (c *client) Put(ctx context.Context, bucket string, objectKey string, fileReader *bytes.Reader) error {
	ossBucket, err := c.connect(ctx, bucket)
	if err != nil {
		logit.Logger.Warn("[OSS] [connect] [requestID]=%d, [bucket]=%s, [objectKey]=%s, [err]=%s", ctx.Value("requestID"), bucket, objectKey, err.Error())
		return err
	}
	err = ossBucket.PutObject(objectKey, fileReader)
	if err != nil {
		logit.Logger.Warn("[OSS] [PutObject] [requestID]=%d, [bucket]=%s, [objectKey]=%s, [err]=%s", ctx.Value("requestID"), bucket, objectKey, err.Error())
		return err
	}
	logit.Logger.Info("[OSS] [Put] [requestID]=%d, [bucket]=%s, [objectKey]=%s success", ctx.Value("requestID"), bucket, objectKey)
	return nil
}

func (c *client) Del(ctx context.Context, bucket string, objectKey string) error {
	ossBucket, err := c.connect(ctx, bucket)
	if err != nil {
		logit.Logger.Warn("[OSS] [connect] [requestID]=%d, [bucket]=%s, [objectKey]=%s, [err]=%s", ctx.Value("requestID"), bucket, objectKey, err.Error())
		return err
	}
	err = ossBucket.DeleteObject(objectKey)
	if err != nil {
		logit.Logger.Warn("[OSS] [DeleteObject] [requestID]=%d, [bucket]=%s, [objectKey]=%s, [err]=%s", ctx.Value("requestID"), bucket, objectKey, err.Error())
		return err
	}
	logit.Logger.Info("[OSS] [Del] [requestID]=%d, [bucket]=%s, [objectKey]=%s success", ctx.Value("requestID"), bucket, objectKey)
	return nil
}

func (c *client) GetURL(ctx context.Context, bucket string, objectKey string) (string, error) {
	ossBucket, err := c.connect(ctx, bucket)
	if err != nil {
		logit.Logger.Warn("[OSS] [connect] [requestID]=%d, [bucket]=%s, [objectKey]=%s, [err]=%s", ctx.Value("requestID"), bucket, objectKey, err.Error())
		return "", err
	}
	url, err := ossBucket.SignURL(objectKey, oss.HTTPGet, 60)
	if err != nil {
		logit.Logger.Warn("[OSS] [SignURL] [requestID]=%d, [bucket]=%s, [objectKey]=%s, [err]=%s", ctx.Value("requestID"), bucket, objectKey, err.Error())
		return "", err
	}
	logit.Logger.Info("[OSS] [GetURL] [requestID]=%d, [bucket]=%s, [objectKey]=%s success", ctx.Value("requestID"), bucket, objectKey)
	return url, nil
}
