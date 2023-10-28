/*
 * @Author: liziwei01
 * @Date: 2023-10-28 14:00:38
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-10-28 14:12:11
 * @Description: 读取xlsx文件
 */
package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const (
	CONNECT_SIGN = ","
)

/**
 * @name:  读取xlsx文件 跳过指定行
 * @return {*}
 */
func (u *UXlsx) ReadXlsx(ctx context.Context, fileReader multipart.File, input map[string]interface{}) ([]string, error) {
	var (
		rst []string
	)
	file, err := excelize.OpenReader(fileReader)
	if err != nil {
		return rst, fmt.Errorf("Open file in xlsx with err " + err.Error())
	}
	sheet := input["sheet"].(string)     // xlsx表格名称
	jumpLine := input["jump_line"].(int) // xlsx 需要跳过的（表头）行数
	rows := file.GetRows(sheet)
	for i, row := range rows {
		if i < jumpLine {
			continue
		}
		var tmpStr []string
		for _, col := range row {
			if col == "" {
				break
			}
			tmpStr = append(tmpStr, col)
		}
		if len(tmpStr) == 0 {
			continue
		}
		rst = append(rst, strings.Join(tmpStr, CONNECT_SIGN))
	}
	return rst, nil
}
