/*
 * @Author: liziwei01
 * @Date: 2022-03-21 17:48:04
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-04-12 13:57:12
 * @Description: file content
 */
package email

import (
	"context"
	"testing"
)

func TestEmail(t *testing.T) {
	ctx := context.Background()
	client, err := GetClient(ctx, "email_lib_user")
	if err != nil {
		t.Error(err)
	}
	subject := "lib Support"
	body := "Hello, your security code is: 123456"
	err = client.Send(ctx, "118010160@link.cuhk.edu.cn", subject, body)
	if err != nil {
		t.Error(err)
	}
}
