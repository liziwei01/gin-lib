/*
 * @Author: liziwei01
 * @Date: 2022-03-21 22:36:04
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-10-29 12:44:24
 * @Description: file content
 */
package redis

import (
	"context"
	"time"
)

func (c *client) Get(ctx context.Context, key string) (value string, err error) {
	db, err := c.connect(ctx)
	if err != nil {
		return "", err
	}
	ret, err := db.Get(key).Result()
	if err != nil {
		return "", err
	}
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
		return err
	}
	return err
}

func (c *client) Del(ctx context.Context, keys ...string) error {
	db, err := c.connect(ctx)
	if err != nil {
		return err
	}
	err = db.Del(keys...).Err()
	if err != nil {
		return err
	}
	return err
}

func (c *client) Exists(ctx context.Context, keys ...string) (int64, error) {
	db, err := c.connect(ctx)
	if err != nil {
		return 0, err
	}
	ret, err := db.Exists(keys...).Result()
	if err != nil {
		return 0, err
	}
	return ret, nil
}
