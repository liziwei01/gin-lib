/*
 * @Author: liziwei01
 * @Date: 2022-06-28 00:46:03
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-06-28 01:00:11
 * @Description: file content
 */
package cookie

import (
	"github.com/gorilla/securecookie"
)

var secureCK = securecookie.New(hashKey, nil)

func Encode(name, value string) (string, error) {
	encoded, err := secureCK.Encode(name, value)
	if err != nil {
		return "", err
	}

	return encoded, nil
}

func Decode(name, value string) (string, error) {
	var decoded string
	err := secureCK.Decode(name, value, &decoded)
	if err != nil {
		return "", err
	}

	return decoded, nil
}
