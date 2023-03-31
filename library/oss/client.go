/*
 * @Author: liziwei01
 * @Date: 2022-03-04 13:52:11
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-03-20 19:41:46
 * @Description: file content
 */
package oss

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Client interface {
	// objectName 对应的阿里云中的文件地址
	// fileHeader 文件
	Get(ctx context.Context, bucket string, objectKey string) (*bytes.Reader, error)
	// Put
	Put(ctx context.Context, bucket string, objectKey string, fileReader *bytes.Reader) error
	// Del
	Del(ctx context.Context, bucket string, objectKey string) error
	// Users could access the object directly with this URL without getting the AK.
	GetURL(ctx context.Context, bucket string, objectKey string) (string, error)

	connect(ctx context.Context, bucket string) (*oss.Bucket, error)
}

type client struct {
	conf *Config
}

func New(config *Config) Client {
	c := &client{
		conf: config,
	}
	return c
}

func (c *client) connect(ctx context.Context, bucket string) (*oss.Bucket, error) {
	client, err := oss.New(c.conf.OSS.Endpoint, c.conf.OSS.AccessKeyID, c.conf.OSS.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("oss.New: %w", err)
	}
	ossBucket, err := client.Bucket(bucket)
	if err != nil {
		return nil, fmt.Errorf("client.Bucket: %w", err)
	}
	return ossBucket, nil
}
