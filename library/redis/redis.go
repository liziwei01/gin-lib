/*
 * @Author: liziwei01
 * @Date: 2022-03-21 22:36:04
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-04-24 21:37:52
 * @Description: file content
 */
package redis

import (
	"context"
	"time"

	"github.com/liziwei01/gin-lib/library/logit"

	"github.com/gogf/gf/util/gconv"
)

func (c *client) Get(ctx context.Context, key string) (value string, err error) {
	db, err := c.connect(ctx)
	if err != nil {
		return "", err
	}
	ret, err := db.Get(key).Result()
	if err != nil {
		logit.Logger.Error("redis get error: %s", err)
		return "", err
	}
	logit.Logger.Info("redis get key: %s, value: %s", key, ret)
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
		logit.Logger.Error("redis set error: %s", err)
		return err
	}
	logit.Logger.Info("redis set key: %s, value: %s, expireTime: %d", key, value, exp)
	return err
}

func (c *client) Del(ctx context.Context, keys ...string) error {
	db, err := c.connect(ctx)
	if err != nil {
		return err
	}
	err = db.Del(keys...).Err()
	if err != nil {
		logit.Logger.Error("redis del error: %s", err)
		return err
	}
	logit.Logger.Info("redis del key: %s", keys)
	return err
}

func (c *client) Exists(ctx context.Context, keys ...string) (bool, error) {
	db, err := c.connect(ctx)
	if err != nil {
		return false, err
	}
	ret, err := db.Exists(keys...).Result()
	if err != nil {
		logit.Logger.Error("redis exists error: %s", err)
		return false, err
	}
	logit.Logger.Info("redis exists key: %s, ret: %d", keys, ret)
	return gconv.Bool(ret), nil
}
