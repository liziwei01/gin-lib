/*
 * @Author: liziwei01
 * @Date: 2022-04-13 23:40:04
 * @LastEditors: liziwei01
 * @LastEditTime: 2022-04-13 23:45:07
 * @Description: file content
 */
package jwtoken

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// set expire time
const TokenExpireDuration = time.Hour * 1

// jwt secret
var JWTSecret = []byte("lib_jwt_secret")

// generate jwt token
func GenToken(username string) (string, error) {
	// generate jwt
	c := MyClaims{
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "lib_official", // 签发人
		},
	}
	// create token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// encode token
	return token.SignedString(JWTSecret)
}
