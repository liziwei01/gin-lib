/*
 * @Author: liziwei01
 * @Date: 2022-03-21 22:36:04
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-10-28 13:53:18
 * @Description: file content
 */
package redis

import (
	"context"
	"time"

	"github.com/liziwei01/gin-lib/library/logit"
)

func (c *client) Get(ctx context.Context, key string) (value string, err error) {
	db, err := c.connect(ctx)
	if err != nil {
		return "", err
	}
	ret, err := db.Get(key).Result()
	if err != nil {
		logit.Logger.Warn("[Redis] [Get] [requestID]=%d, [key]=%s, [err]=%s", ctx.Value("requestID"), key, err.Error())
		return "", err
	}
	logit.Logger.Info("[Redis] [Get] [requestID]=%d, [key]=%s, [value]=%s", ctx.Value("requestID"), key, ret)
	return ret, nil
}

func (c *client) Set(ctx context.Context, key string, value string, expireTime ...time.Duration) error {
	var exp time.Duration = time.Hour
	db, err := c.connect(ctx)
	if err != nil {
		return err
	}
	if len(expireTime) > 0 {
		exp = expireTime[0]
	}
	err = db.Set(key, value, exp).Err()
	if err != nil {
		logit.Logger.Warn("[Redis] [Set] [requestID]=%d, [key]=%s, [value]=%s, [err]=%s", ctx.Value("requestID"), key, value, err.Error())
		return err
	}
	logit.Logger.Info("[Redis] [Set] [requestID]=%d, [key]=%s, [value]=%s, [expire]=%d", ctx.Value("requestID"), key, value, exp)
	return err
}

func (c *client) Del(ctx context.Context, keys ...string) error {
	db, err := c.connect(ctx)
	if err != nil {
		return err
	}
	err = db.Del(keys...).Err()
	if err != nil {
		logit.Logger.Warn("[Redis] [Del] [requestID]=%d, [key]=%s, [err]=%s", ctx.Value("requestID"), keys, err.Error())
		return err
	}
	logit.Logger.Info("[Redis] [Del] [requestID]=%d, [key]=%s", ctx.Value("requestID"), keys)
	return err
}

func (c *client) Exists(ctx context.Context, keys ...string) (int64, error) {
	db, err := c.connect(ctx)
	if err != nil {
		return 0, err
	}
	ret, err := db.Exists(keys...).Result()
	if err != nil {
		logit.Logger.Warn("[Redis] [Exists] [requestID]=%d, [keys]=%#v, [err]=%s", ctx.Value("requestID"), keys, err.Error())
		return 0, err
	}
	logit.Logger.Info("[Redis] [Exists] [requestID]=%d, [keys]=%#v, [value]=%s", ctx.Value("requestID"), keys, ret)
	return ret, nil
}
