/*
 * @Author: liziwei01
 * @Date: 2022-03-20 18:17:39
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 11:18:24
 * @Description: file content
 */
package email

import (
	"context"

	"gopkg.in/gomail.v2"
)

func (c *client) Send(ctx context.Context, to, subject, body string) error {
	dialer, err := c.connect(ctx)
	if err != nil {
		return err
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", c.conf.Email.Address)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	return dialer.DialAndSend(msg)
}
