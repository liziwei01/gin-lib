/*
 * @Author: liziwei01
 * @Date: 2022-03-10 20:45:29
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-03-10 20:57:10
 * @Description: file content
 */
package request

import (
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/liziwei01/gin-lib/library/utils"
)

// 接口数据MD5校验.
// @return (true: valid, false: invalid).
func CheckSignValid(req *http.Request, spcsign string) bool {
	// 获取所有url params.
	urlQueryMap := GetQueryMap(req)
	if urlQueryMap == nil || len(urlQueryMap) == 0 {
		return true
	}
	return checkSign(urlQueryMap, spcsign)
}

// 获取所有的参数.
// get方法获取所有url params.
// post方法获取所有的x-www-form-urlencoded的参数.
func GetQueryMap(req *http.Request) map[string]string {
	method := strings.ToLower(req.Method)
	var params map[string][]string
	if method == "get" {
		// get方法获取所有url params.
		params = req.URL.Query()
	} else if method == "post" {
		// post方法获取所有的x-www-form-urlencoded的参数.
		err := req.ParseForm()
		if err != nil {
			return nil
		}
		params = req.PostForm
	}
	if len(params) == 0 {
		return nil
	}
	res := make(map[string]string)
	for key, value := range params {
		res[key] = value[0]
	}
	return res
}

func checkSign(pars map[string]string, spcsign string) bool {
	// 获取接口参数签名.
	sign := pars["sign"]
	if sign == "" {
		return false
	}
	// 特殊签名不需要检验.
	if sign == spcsign {
		return true
	}
	// 删除sign.
	delete(pars, "sign")
	// 获取时间戳.
	timeStamp := pars["ts"]
	if timeStamp == "" {
		return false
	}
	intTimeStamp, err := strconv.ParseInt(timeStamp, 10, 64)
	if err != nil {
		return false
	}
	// 接口60s内有效.
	curTimeStamp := time.Now().Unix()
	if math.Abs((float64)(curTimeStamp-intTimeStamp)) > 60 {
		return false
	}

	// 添加salt.
	pars["t_secret"] = strconv.FormatInt(intTimeStamp%999983, 10)

	cSign := getSign(pars)

	return cSign == sign
}

func getSign(parsWithTss map[string]string) string {
	// 按照key排序.
	var keys []string
	for k := range parsWithTss {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接字符串.
	var builder strings.Builder
	for _, key := range keys {
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(parsWithTss[key])
	}
	paramStr := builder.String()
	return utils.Md5.Md5String(paramStr)
}
