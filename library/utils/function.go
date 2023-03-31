/*
 * @Author: liziwei01
 * @Date: 2022-06-28 01:12:49
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-06-28 01:12:50
 * @Description: file content
 */
package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// md5Byte 计算字节切片的 MD5 散列值.
func md5Byte(str []byte, length uint8) []byte {
	var res []byte
	h := md5.New()
	h.Write(str)

	hBytes := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(hBytes)))
	hex.Encode(dst, hBytes)
	if length > 0 && length < 32 {
		res = dst[:length]
	} else {
		res = dst
	}

	return res
}

// shaXByte 计算字节切片的 shaX 散列值,x为1/256/512.
func shaXByte(str []byte, x uint16) []byte {
	var h hash.Hash
	switch x {
	case 1:
		h = sha1.New()
		break
	case 256:
		h = sha256.New()
		break
	case 512:
		h = sha512.New()
		break
	default:
		fmt.Print("[shaXByte] x must be in [1, 256, 512]")
		return nil
	}

	h.Write(str)

	hBytes := h.Sum(nil)
	res := make([]byte, hex.EncodedLen(len(hBytes)))
	hex.Encode(res, hBytes)
	return res
}

// getTrimMask 去除mask字符.
func getTrimMask(characterMask []string) string {
	var mask string
	if len(characterMask) == 0 {
		mask = " \t\n\r\v\f\x00　"
	} else {
		mask = strings.Join(characterMask, "")
	}
	return mask
}

// reflectPtr 获取反射的指向.
func reflectPtr(r reflect.Value) reflect.Value {
	// 如果是指针,则获取其所指向的元素
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	return r
}

// creditChecksum 计算身份证校验码,其中id为身份证号码.
func creditChecksum(id string) byte {
	//∑(ai×Wi)(mod 11)
	// 加权因子
	factor := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	// 校验位对应值
	code := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

	leng := len(id)
	sum := 0
	for i, char := range id[:leng-1] {
		num, _ := strconv.Atoi(string(char))
		sum += num * factor[i]
	}

	return code[sum%11]
}

// compareConditionMap 比对数组是否匹配条件.condition为条件字典,arr为要比对的数据数组.
func compareConditionMap(condition map[string]interface{}, arr interface{}) (res interface{}) {
	val := reflect.ValueOf(arr)
	switch val.Kind() {
	case reflect.Map:
		condLen := len(condition)
		chkNum := 0
		if condLen > 0 {
			for _, k := range val.MapKeys() {
				if condVal, ok := condition[k.String()]; ok && reflect.DeepEqual(val.MapIndex(k).Interface(), condVal) {
					chkNum++
				}
			}
		}

		if chkNum == condLen {
			res = arr
		}
	default:
		return
	}

	return
}

// getMethod 获取对象的方法.
func getMethod(t interface{}, method string) reflect.Value {
	m, b := reflect.TypeOf(t).MethodByName(method)
	if !b {
		return reflect.ValueOf(nil)
	}
	return m.Func
}

// ValidFunc 检查是否函数,并且参数个数、类型是否正确.
// 返回有效的函数、有效的参数.
func ValidFunc(f interface{}, args ...interface{}) (vf reflect.Value, vargs []reflect.Value, err error) {
	vf = reflect.ValueOf(f)
	if vf.Kind() != reflect.Func {
		return reflect.ValueOf(nil), nil, fmt.Errorf("[validFunc] %v is not the function", f)
	}

	tf := vf.Type()
	_len := len(args)
	if tf.NumIn() != _len {
		return reflect.ValueOf(nil), nil, fmt.Errorf("[validFunc] %d number of the argument is incorrect", _len)
	}

	vargs = make([]reflect.Value, _len)
	for i := 0; i < _len; i++ {
		typ := tf.In(i).Kind()
		if (typ != reflect.Interface) && (typ != reflect.TypeOf(args[i]).Kind()) {
			return reflect.ValueOf(nil), nil, fmt.Errorf("[validFunc] %d-td argument`s type is incorrect", i+1)
		}
		vargs[i] = reflect.ValueOf(args[i])
	}
	return vf, vargs, nil
}

// CallFunc 动态调用函数.
func CallFunc(f interface{}, args ...interface{}) (results []interface{}, err error) {
	vf, vargs, _err := ValidFunc(f, args...)
	if _err != nil {
		return nil, _err
	}
	ret := vf.Call(vargs)
	_len := len(ret)
	results = make([]interface{}, _len)
	for i := 0; i < _len; i++ {
		results[i] = ret[i].Interface()
	}
	return
}

// camelCaseToLowerCase 驼峰转为小写.
func camelCaseToLowerCase(str string, connector rune) string {
	if len(str) == 0 {
		return ""
	}

	buf := &bytes.Buffer{}
	var prev, r0, r1 rune
	var size int

	r0 = connector

	for len(str) > 0 {
		prev = r0
		r0, size = utf8.DecodeRuneInString(str)
		str = str[size:]

		switch {
		case r0 == utf8.RuneError:
			continue

		case unicode.IsUpper(r0):
			if prev != connector && !unicode.IsNumber(prev) {
				buf.WriteRune(connector)
			}

			buf.WriteRune(unicode.ToLower(r0))

			if len(str) == 0 {
				break
			}

			r0, size = utf8.DecodeRuneInString(str)
			str = str[size:]

			if !unicode.IsUpper(r0) {
				buf.WriteRune(r0)
				break
			}

			// find next non-upper-case character and insert connector properly.
			// it's designed to convert `HTTPServer` to `http_server`.
			// if there are more than 2 adjacent upper case characters in a word,
			// treat them as an abbreviation plus a normal word.
			for len(str) > 0 {
				r1 = r0
				r0, size = utf8.DecodeRuneInString(str)
				str = str[size:]

				if r0 == utf8.RuneError {
					buf.WriteRune(unicode.ToLower(r1))
					break
				}

				if !unicode.IsUpper(r0) {
					if isCaseConnector(r0) {
						r0 = connector

						buf.WriteRune(unicode.ToLower(r1))
					} else if unicode.IsNumber(r0) {
						// treat a number as an upper case rune
						// so that both `http2xx` and `HTTP2XX` can be converted to `http_2xx`.
						buf.WriteRune(unicode.ToLower(r1))
						buf.WriteRune(connector)
						buf.WriteRune(r0)
					} else {
						buf.WriteRune(connector)
						buf.WriteRune(unicode.ToLower(r1))
						buf.WriteRune(r0)
					}

					break
				}

				buf.WriteRune(unicode.ToLower(r1))
			}

			if len(str) == 0 || r0 == connector {
				buf.WriteRune(unicode.ToLower(r0))
			}

		case unicode.IsNumber(r0):
			if prev != connector && !unicode.IsNumber(prev) {
				buf.WriteRune(connector)
			}

			buf.WriteRune(r0)

		default:
			if isCaseConnector(r0) {
				r0 = connector
			}

			buf.WriteRune(r0)
		}
	}

	return buf.String()
}

// isCaseConnector 是否字符转换连接符.
func isCaseConnector(r rune) bool {
	return r == '-' || r == '_' || unicode.IsSpace(r)
}

// getPidByInode 根据套接字的inode获取PID.须root权限.
func getPidByInode(inode string, procDirs []string) (pid int) {
	if len(procDirs) == 0 {
		procDirs, _ = filepath.Glob("/proc/[0-9]*/fd/[0-9]*")
	}

	re := regexp.MustCompile(inode)
	for _, item := range procDirs {
		path, _ := os.Readlink(item)
		out := re.FindString(path)
		if len(out) != 0 {
			pid, _ = strconv.Atoi(strings.Split(item, "/")[2])
			break
		}
	}

	return pid
}

// getProcessPathByPid 根据PID获取进程的执行路径.
func getProcessPathByPid(pid int) string {
	exe := fmt.Sprintf("/proc/%d/exe", pid)
	path, _ := os.Readlink(exe)
	return path
}

// pkcs7Padding PKCS7填充.
// cipherText为密文;blockSize为分组长度;isZero是否零填充.
func pkcs7Padding(cipherText []byte, blockSize int, isZero bool) []byte {
	clen := len(cipherText)
	if cipherText == nil || clen == 0 || blockSize <= 0 {
		return nil
	}

	var padtext []byte
	padding := blockSize - clen%blockSize
	if isZero {
		padtext = bytes.Repeat([]byte{0}, padding)
	} else {
		padtext = bytes.Repeat([]byte{byte(padding)}, padding)
	}

	return append(cipherText, padtext...)
}

// pkcs7UnPadding PKCS7拆解.
// origData为源数据;blockSize为分组长度.
func pkcs7UnPadding(origData []byte, blockSize int) []byte {
	olen := len(origData)
	if origData == nil || olen == 0 || blockSize <= 0 || olen%blockSize != 0 {
		return nil
	}

	unpadding := int(origData[olen-1])
	if unpadding == 0 || unpadding > olen {
		return nil
	}

	return origData[:(olen - unpadding)]
}

// zeroPadding PKCS7使用0填充.
func zeroPadding(cipherText []byte, blockSize int) []byte {
	return pkcs7Padding(cipherText, blockSize, true)
}

// zeroUnPadding PKCS7-0拆解.
func zeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

// formatPath 格式化路径
func formatPath(fpath string) string {
	//替换特殊字符
	fpath = strings.NewReplacer(`|`, "", `:`, "", `<`, "", `>`, "", `?`, "", `\`, "/").Replace(fpath)
	// 将"\"替换为"/"
	//fpath = strings.ReplaceAll(fpath, "\\", "/")
	//替换连续斜杠
	fpath = RegFormatDir.ReplaceAllString(fpath, "/")
	return fpath
}
