/*
 * @Author: liziwei01
 * @Date: 2022-03-04 13:52:11
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-03-21 17:43:10
 * @Description: file content
 */
package email

import (
	"context"
	"fmt"

	"gopkg.in/gomail.v2"
)

type Client interface {
	Send(ctx context.Context, to, subject, body string) error

	connect(ctx context.Context) (*gomail.Dialer, error)
}

type client struct {
	conf   *Config
	dialer *gomail.Dialer
}

func New(config *Config) Client {
	c := &client{
		conf: config,
	}
	return c
}

func (c *client) connect(ctx context.Context) (*gomail.Dialer, error) {
	if c.dialer == nil {
		if c.conf.Resource.Manual.Host == "" || c.conf.Resource.Manual.Port == 0 {
			return nil, fmt.Errorf("email resource not set")
		}
		if c.conf.Email.Address == "" || c.conf.Email.Password == "" {
			return nil, fmt.Errorf("email address or password not set")
		}
		c.dialer = gomail.NewDialer(c.conf.Resource.Manual.Host, c.conf.Resource.Manual.Port, c.conf.Email.Address, c.conf.Email.Password)
	}
	return c.dialer, nil
}
