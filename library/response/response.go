/*
 * @Author: liziwei01
 * @Date: 2022-03-03 19:50:47
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-06-30 04:47:26
 * @Description: 标准错误函数
 */
package response

import (
	"net/http"

	"github.com/liziwei01/gin-lib/library/utils"

	"github.com/gin-gonic/gin"
)

// data可不传.
// 简单方法—请求时成功返回.
func StdSuccess(ctx *gin.Context, data ...interface{}) {
	StdResponse(ctx, Success, CodeMsgMap[Success], data...)
}

// data可不传.
// 简单方法—请求失败时返回.
func StdFailed(ctx *gin.Context, data ...interface{}) {
	StdResponse(ctx, Failed, CodeMsgMap[Failed], data...)
}

// data可不传.
// 简单方法—参数错误时返回.
func StdInvalidParams(ctx *gin.Context, data ...interface{}) {
	StdResponse(ctx, InvalidParams, CodeMsgMap[InvalidParams], data...)
}

// data可不传.
// 接口MD5校验失败.
func StdSignCheckFailed(ctx *gin.Context, data ...interface{}) {
	StdResponse(ctx, SignCheckFailed, CodeMsgMap[SignCheckFailed], data...)
}

// data可不传.
// 接口MD5校验失败.
func StdTokenCheckFailed(ctx *gin.Context, data ...interface{}) {
	StdResponse(ctx, TokenCheckFailed, CodeMsgMap[TokenCheckFailed], data...)
}

// data可不传.
func StdAuthCheckFailed(ctx *gin.Context, data ...interface{}) {
	StdResponse(ctx, ERR_NO_AUTH, CodeMsgMap[ERR_NO_AUTH], data...)
}

// data可不传.
// 根据错误码返回.
func StdWithCode(ctx *gin.Context, code int, data ...interface{}) {
	StdResponse(ctx, code, CodeMsgMap[code], data...)
}

// data可不传.
// 传入完整错误码、错误信息拼凑返回信息返回.
func StdResponse(ctx *gin.Context, code int, msg string, data ...interface{}) {
	d := utils.Slice.GetFirstDefault(data)
	resp := map[string]interface{}{
		"errno":  code,
		"errmsg": msg,
		"data":   d,
	}
	ctx.JSON(http.StatusOK, resp)
}

// func StdFile(ctx *gin.Context, data []byte) {
// 	ctx.DataFromReader(200, ctx..Response.ContentLength, "application/octet-stream", fileContent, nil)
// }
