/*
 * @Author: liziwei01
 * @Date: 2022-03-21 20:25:23
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-05-10 15:30:55
 * @Description: file content
 */
package mysql

import (
	"context"
	"testing"
)

var (
	test_table_name = "tb_user_privare_info"
)

type UserPrivateInfo struct {
	UserID   string `json:"user_id" ddb:"user_id"`
	Nickname string `json:"nickname" ddb:"nickname"`
	Email    string `json:"email" ddb:"email"`
}

func TestQuery(t *testing.T) {
	var (
		ctx        = context.Background()
		returnData []UserPrivateInfo
	)
	client, err := GetClient(ctx, "db_lib_user")
	if err != nil {
		t.Error(err)
	}
	where := map[string]interface{}{
		"user_id": 1,
	}
	columns := []string{"user_id", "nickname", "email"}
	err = client.Query(ctx, test_table_name, where, columns, &returnData)
	if err != nil {
		t.Error(err)
	}
	t.Log(returnData)
}

func TestInsert(t *testing.T) {
	var (
		ctx = context.Background()
	)
	client, err := GetClient(ctx, "db_lib_user")
	if err != nil {
		t.Error(err)
	}
	data := make([]map[string]interface{}, 0)
	user1 := map[string]interface{}{
		"user_id":  1,
		"nickname": "test user 1",
		"email":    "testemail1@163.com",
	}
	user2 := map[string]interface{}{
		"user_id":  2,
		"nickname": "test user 2",
		"email":    "testemail2@163.com",
	}
	data = append(data, user1)
	data = append(data, user2)
	_, err = client.Insert(ctx, test_table_name, data)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	var (
		ctx         = context.Background()
		returnData  []UserPrivateInfo
		updatedNick = "test user 1 updated"
	)
	client, err := GetClient(ctx, "db_lib_user")
	if err != nil {
		t.Error(err)
	}
	where := map[string]interface{}{
		"user_id":  1,
		"nickname": "test user 1",
		"email":    "testemail1@163.com",
	}
	update := map[string]interface{}{
		"nickname": updatedNick,
	}
	_, err = client.Update(ctx, test_table_name, where, update)
	if err != nil {
		t.Error(err)
	}
	columns := []string{"nickname"}
	err = client.Query(ctx, test_table_name, update, columns, &returnData)
	if returnData[0].Nickname != updatedNick {
		t.Errorf("returnData[0].Nickname: %s", returnData[0].Nickname)
	}
}

func TestDelete(t *testing.T) {
	var (
		ctx = context.Background()
	)
	client, err := GetClient(ctx, "db_lib_user")
	if err != nil {
		t.Error(err)
	}
	where := map[string]interface{}{
		"user_id":  1,
		"nickname": "test user 1 updated",
		"email":    "testemail1@163.com",
	}
	res, err := client.Delete(ctx, test_table_name, where)
	if err != nil {
		t.Error(err)
	}
	affectedNum, err := res.RowsAffected()
	if err != nil {
		t.Error(err)
	}
	if affectedNum != 1 {
		t.Errorf("affectedNum = %d", affectedNum)
	}
}

func TestExecRaw(t *testing.T) {
	var (
		ctx = context.Background()
	)
	client, err := GetClient(ctx, "db_lib_user")
	if err != nil {
		t.Error(err)
	}
	res, err := client.ExecRaw(ctx, "delete from tb_user_privare_info where user_id = ? and nickname = ? and email = ?", "2", "test user 2", "testemail2@163.com")
	if err != nil {
		t.Error(err)
	}
	affectedNum, err := res.RowsAffected()
	if err != nil {
		t.Error(err)
	}
	if affectedNum != 1 {
		t.Errorf("affectedNum = %d", affectedNum)
	}
}
