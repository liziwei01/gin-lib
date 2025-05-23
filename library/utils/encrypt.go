/*
 * @Author: liziwei01
 * @Date: 2022-06-28 01:08:34
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-10-28 14:10:31
 * @Description: file content
 */
package utils

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"io"
	"math/big"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Base64Encode 使用 MIME base64 对数据进行编码.
func (ue *UEncrypt) Base64EncodeString(str string) string {
	bytes := string(ue.Base64Encode([]byte(str)))
	return bytes
}

// Base64Decode 对使用 MIME base64 编码的数据进行解码.
func (ue *UEncrypt) Base64DecodeString(str string) (string, error) {
	bytes, err := ue.Base64Decode([]byte(str))
	return string(bytes), err
}

// Base64Encode 使用 MIME base64 对数据进行编码.
func (ue *UEncrypt) Base64Encode(str []byte) []byte {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(str)))
	base64.StdEncoding.Encode(buf, str)
	return buf
}

// Base64Decode 对使用 MIME base64 编码的数据进行解码.
func (ue *UEncrypt) Base64Decode(str []byte) ([]byte, error) {
	dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(str)))
	n, err := base64.StdEncoding.Decode(dbuf, str)
	return dbuf[:n], err
}

// Base64UrlEncode url安全的Base64Encode,没有'/'和'+'及结尾的'=' .
func (ue *UEncrypt) Base64UrlEncode(source []byte) []byte {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(source)))
	base64.StdEncoding.Encode(buf, source)

	// Base64 Url Safe is the same as Base64 but does not contain '/' and '+' (replaced by '_' and '-') and trailing '=' are removed.
	buf = bytes.Replace(buf, []byte("/"), []byte("_"), -1)
	buf = bytes.Replace(buf, []byte("+"), []byte("-"), -1)
	buf = bytes.Replace(buf, []byte("="), []byte(""), -1)

	return buf
}

// Base64UrlDecode url安全的Base64Decode.
func (ue *UEncrypt) Base64UrlDecode(data []byte) ([]byte, error) {
	var missing = (4 - len(data)%4) % 4
	data = append(data, bytes.Repeat([]byte("="), missing)...)

	dbuf := make([]byte, base64.URLEncoding.DecodedLen(len(data)))
	n, err := base64.URLEncoding.Decode(dbuf, data)
	return dbuf[:n], err
}

// AuthCode 授权码编码或解码;encode为true时编码,为false解码;expiry为有效期,秒;返回结果为加密/解密的字符串和有效期时间戳.
func (ue *UEncrypt) AuthCode(str, key []byte, encode bool, expiry int64) ([]byte, int64) {
	// DYNAMIC_KEY_LEN 动态密钥长度，相同的明文会生成不同密文就是依靠动态密钥
	// 加入随机密钥，可以令密文无任何规律，即便是原文和密钥完全相同，加密结果也会每次不同，增大破解难度。
	// 取值越大，密文变动规律越大，密文变化 = 16 的 DYNAMIC_KEY_LEN 次方
	// 当此值为 0 时，则不产生随机密钥

	strLen := len(str)
	if str == nil || strLen == 0 {
		return nil, 0
	} else if !encode && strLen < DYNAMIC_KEY_LEN {
		return nil, 0
	}

	// 密钥
	keyByte := md5Byte(key, 32)

	// 密钥a会参与加解密
	keya := keyByte[:16]

	// 密钥b会用来做数据完整性验证
	keyb := keyByte[16:]

	// 密钥c用于变化生成的密文
	var keyc []byte
	if encode == false {
		keyc = str[:DYNAMIC_KEY_LEN]
	} else {
		now, _ := time.Now().MarshalBinary()
		keycLen := 32 - DYNAMIC_KEY_LEN
		timeBytes := md5Byte(now, 32)
		keyc = timeBytes[keycLen:]
	}

	// 参与运算的密钥
	keyd := md5Byte(append(keya, keyc...), 32)
	cryptkey := append(keya, keyd...)
	cryptkeyLen := len(cryptkey)
	// 明文，前10位用来保存时间戳，解密时验证数据有效性，10到26位用来保存keyb(密钥b)，解密时会通过这个密钥验证数据完整性
	// 如果是解码的话，会从第 DYNAMIC_KEY_LEN 位开始，因为密文前 DYNAMIC_KEY_LEN 位保存 动态密钥，以保证解密正确
	if encode == false { //解密
		var err error
		str, err = ue.Base64UrlDecode(str[DYNAMIC_KEY_LEN:])
		if err != nil {
			return nil, 0
		}
	} else {
		if expiry != 0 {
			expiry = expiry + time.Now().Unix()
		}
		expMd5 := md5Byte(append(str, keyb...), 16)
		str = []byte(fmt.Sprintf("%010d%s%s", expiry, expMd5, str))
		//str = append([]byte(fmt.Sprintf("%010d", expiry)), append(expMd5, str...)...)
	}

	strLen = len(str)
	resdata := make([]byte, 0, strLen)
	var rndkey, box [256]int
	// 产生密钥簿
	h := 0
	i := 0
	j := 0

	for i = 0; i < 256; i++ {
		rndkey[i] = int(cryptkey[i%cryptkeyLen])
		box[i] = i
	}
	// 用固定的算法，打乱密钥簿，增加随机性，好像很复杂，实际上并不会增加密文的强度
	for i = 0; i < 256; i++ {
		j = (j + box[i] + rndkey[i]) % 256
		box[i], box[j] = box[j], box[i]
	}
	// 核心加解密部分
	h = 0
	j = 0
	for i = 0; i < strLen; i++ {
		h = ((h + 1) % 256)
		j = ((j + box[h]) % 256)
		box[h], box[j] = box[j], box[h]
		// 从密钥簿得出密钥进行异或，再转成字符
		resdata = append(resdata, byte(int(str[i])^box[(box[h]+box[j])%256]))
	}
	if encode == false { //解密
		// substr($result, 0, 10) == 0 验证数据有效性
		// substr($result, 0, 10) - time() > 0 验证数据有效性
		// substr($result, 10, 16) == substr(md5(substr($result, 26).$keyb), 0, 16) 验证数据完整性
		// 验证数据有效性，请看未加密明文的格式
		if len(resdata) <= 26 {
			return nil, 0
		}

		expTime, _ := strconv.ParseInt(string(resdata[:10]), 10, 0)
		if (expTime == 0 || expTime-time.Now().Unix() > 0) && string(resdata[10:26]) == string(md5Byte(append(resdata[26:], keyb...), 16)) {
			return resdata[26:], expTime
		} else {
			return nil, expTime
		}
	} else { //加密
		// 把动态密钥保存在密文里，这也是为什么同样的明文，生产不同密文后能解密的原因
		resdata = append(keyc, ue.Base64UrlEncode(resdata)...)
		return resdata, expiry
	}
}

// PasswordHash 创建密码的散列值;costs为算法的cost,范围4~31,默认10;注意:值越大越耗时.
func (ue *UEncrypt) PasswordHashString(password string, costs ...int) (string, error) {
	bytes, err := ue.PasswordHash([]byte(password), costs...)
	return string(bytes), err
}

// PasswordHash 创建密码的散列值;costs为算法的cost,范围4~31,默认10;注意:值越大越耗时.
func (ue *UEncrypt) PasswordHash(password []byte, costs ...int) ([]byte, error) {
	var cost int
	if len(costs) == 0 {
		cost = 10
	} else {
		cost = costs[0]
		if cost < 4 {
			cost = 4
		} else if cost > 31 {
			cost = 31
		}
	}

	res, err := bcrypt.GenerateFromPassword(password, cost)
	return res, err
}

// PasswordVerify 验证密码是否和散列值匹配.
func (ue *UEncrypt) PasswordVerifyString(password, hash string) bool {
	return ue.PasswordVerify([]byte(password), []byte(hash))
}

// PasswordVerify 验证密码是否和散列值匹配.
func (ue *UEncrypt) PasswordVerify(password, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

// EasyEncrypt 简单加密.
// data为要加密的原字符串,key为密钥.
func (ue *UEncrypt) EasyEncrypt(data, key []byte) []byte {
	dataLen := len(data)
	if dataLen == 0 {
		return nil
	}

	keyByte := md5Byte(key, 32)
	keyLen := len(keyByte)

	var i, x, c int
	var res []byte
	for i = 0; i < dataLen; i++ {
		if x == keyLen {
			x = 0
		}

		c = (int(data[i]) + int(keyByte[x])) % 256
		res = append(res, byte(c))

		x++
	}

	res = append(keyByte[:DYNAMIC_KEY_LEN], ue.Base64UrlEncode(res)...)
	return res
}

// EasyDecrypt 简单解密.
// val为待解密的字符串,key为密钥.
func (ue *UEncrypt) EasyDecrypt(val, key []byte) []byte {
	if len(val) <= DYNAMIC_KEY_LEN {
		return nil
	}

	data, err := ue.Base64UrlDecode(val[DYNAMIC_KEY_LEN:])
	if err != nil {
		return nil
	}

	keyByte := md5Byte(key, 32)
	if string(val[:DYNAMIC_KEY_LEN]) != string(keyByte[:DYNAMIC_KEY_LEN]) {
		return nil
	}

	dataLen := len(data)
	keyLen := len(keyByte)

	var i, x, c int
	var res []byte
	for i = 0; i < dataLen; i++ {
		if x == keyLen {
			x = 0
		}

		if data[i] < keyByte[x] {
			c = int(data[i]) + 256 - int(keyByte[x])
		} else {
			c = int(data[i]) - int(keyByte[x])
		}
		res = append(res, byte(c))

		x++
	}

	return res
}

// HmacShaX HmacSHA-x加密,x为1/256/512 .
func (ue *UEncrypt) HmacShaX(data, secret []byte, x uint16) []byte {
	// Create a new HMAC by defining the hash type and the key (as byte array)
	var h hash.Hash
	switch x {
	case 1:
		h = hmac.New(sha1.New, secret)
		break
	case 256:
		h = hmac.New(sha256.New, secret)
		break
	case 512:
		h = hmac.New(sha512.New, secret)
		break
	default:
		fmt.Print("[HmacShaX] x must be in [1, 256, 512]")
		return nil
	}

	// Write Data to it
	h.Write(data)

	src := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	return dst
}

// aesEncrypt AES加密.
// clearText为明文;key为密钥,长度16/24/32;
// mode为模式,枚举值(CBC,CFB,CTR,OFB);
// paddingType为填充方式,枚举(PKCS_NONE,PKCS_ZERO,PKCS_SEVEN),默认PKCS_SEVEN.
func (ue *UEncrypt) aesEncrypt(clearText, key []byte, mode string, paddingType ...UPKCSType) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	pt := PKCS_SEVEN
	blockSize := block.BlockSize()
	if len(paddingType) > 0 {
		pt = paddingType[0]
	}
	switch pt {
	case PKCS_ZERO:
		clearText = zeroPadding(clearText, blockSize)
	case PKCS_SEVEN:
		clearText = pkcs7Padding(clearText, blockSize, false)
	}

	cipherText := make([]byte, blockSize+len(clearText))
	//初始化向量
	iv := cipherText[:blockSize]
	_, _ = io.ReadFull(rand.Reader, iv)
	//if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	//	return nil, err
	//}

	switch mode {
	case "CBC":
		cipher.NewCBCEncrypter(block, iv).CryptBlocks(cipherText[blockSize:], clearText)
	case "CFB":
		cipher.NewCFBEncrypter(block, iv).XORKeyStream(cipherText[blockSize:], clearText)
	case "CTR":
		cipher.NewCTR(block, iv).XORKeyStream(cipherText[blockSize:], clearText)
	case "OFB":
		cipher.NewOFB(block, iv).XORKeyStream(cipherText[blockSize:], clearText)
	}

	return cipherText, nil
}

// aesDecrypt AES解密.
// cipherText为密文;key为密钥,长度16/24/32;
// mode为模式,枚举值(CBC,CFB,CTR,OFB);
// paddingType为填充方式,枚举(PKCS_NONE,PKCS_ZERO,PKCS_SEVEN),默认PKCS_SEVEN.
func (ue *UEncrypt) aesDecrypt(cipherText, key []byte, mode string, paddingType ...UPKCSType) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	pt := PKCS_SEVEN
	if len(paddingType) > 0 {
		pt = paddingType[0]
	}

	blockSize := block.BlockSize()
	clen := len(cipherText)
	if clen < blockSize {
		return nil, errors.New("cipherText too short")
	}

	iv := cipherText[:blockSize]
	cipherText = cipherText[blockSize:]

	switch mode {
	case "CBC":
		cipher.NewCBCDecrypter(block, iv).CryptBlocks(cipherText, cipherText)
	case "CFB":
		cipher.NewCFBDecrypter(block, iv).XORKeyStream(cipherText, cipherText)
	case "CTR":
		cipher.NewCTR(block, iv).XORKeyStream(cipherText, cipherText)
	case "OFB":
		cipher.NewOFB(block, iv).XORKeyStream(cipherText, cipherText)
	}

	clen = len(cipherText)
	if pt != PKCS_NONE && clen > 0 && int(cipherText[clen-1]) > clen {
		return nil, fmt.Errorf(fmt.Sprintf("aes [%s] decrypt failed", mode))
	}

	var plainText []byte
	switch pt {
	case PKCS_ZERO:
		plainText = zeroUnPadding(cipherText)
	case PKCS_SEVEN:
		plainText = pkcs7UnPadding(cipherText, blockSize)
	case PKCS_NONE:
		plainText = cipherText
	}

	return plainText, nil
}

// AesCBCEncrypt AES-CBC密码分组链接(Cipher-block chaining)模式加密.加密无法并行,不适合对流数据加密.
// clearText为明文;key为密钥,长16/24/32;paddingType为填充方式,枚举(PKCS_ZERO,PKCS_SEVEN),默认PKCS_SEVEN.
func (ue *UEncrypt) AesCBCEncrypt(clearText, key []byte, paddingType ...UPKCSType) ([]byte, error) {
	return ue.aesEncrypt(clearText, key, "CBC", paddingType...)
}

// AesCBCDecrypt AES-CBC密码分组链接(Cipher-block chaining)模式解密.
// cipherText为密文;key为密钥,长16/24/32;paddingType为填充方式,枚举(PKCS_NONE,PKCS_ZERO,PKCS_SEVEN),默认PKCS_SEVEN.
func (ue *UEncrypt) AesCBCDecrypt(cipherText, key []byte, paddingType ...UPKCSType) ([]byte, error) {
	return ue.aesDecrypt(cipherText, key, "CBC", paddingType...)
}

// AesCFBEncrypt AES-CFB密文反馈(Cipher feedback)模式加密.适合对流数据加密.
// clearText为明文;key为密钥,长16/24/32.
func (ue *UEncrypt) AesCFBEncrypt(clearText, key []byte) ([]byte, error) {
	return ue.aesEncrypt(clearText, key, "CFB", PKCS_NONE)
}

// AesCFBDecrypt AES-CFB密文反馈(Cipher feedback)模式解密.
// cipherText为密文;key为密钥,长16/24/32.
func (ue *UEncrypt) AesCFBDecrypt(cipherText, key []byte) ([]byte, error) {
	return ue.aesDecrypt(cipherText, key, "CFB", PKCS_NONE)
}

// AesCTREncrypt AES-CTR计算器(Counter)模式加密.
// clearText为明文;key为密钥,长16/24/32.
func (ue *UEncrypt) AesCTREncrypt(clearText, key []byte) ([]byte, error) {
	return ue.aesEncrypt(clearText, key, "CTR", PKCS_NONE)
}

// AesCTRDecrypt AES-CTR计算器(Counter)模式解密.
// cipherText为密文;key为密钥,长16/24/32.
func (ue *UEncrypt) AesCTRDecrypt(cipherText, key []byte) ([]byte, error) {
	return ue.aesDecrypt(cipherText, key, "CTR", PKCS_NONE)
}

// AesOFBEncrypt AES-OFB输出反馈(Output feedback)模式加密.适合对流数据加密.
// clearText为明文;key为密钥,长16/24/32.
func (ue *UEncrypt) AesOFBEncrypt(clearText, key []byte) ([]byte, error) {
	return ue.aesEncrypt(clearText, key, "OFB", PKCS_NONE)
}

// AesOFBDecrypt AES-OFB输出反馈(Output feedback)模式解密.
// cipherText为密文;key为密钥,长16/24/32.
func (ue *UEncrypt) AesOFBDecrypt(cipherText, key []byte) ([]byte, error) {
	return ue.aesDecrypt(cipherText, key, "OFB", PKCS_NONE)
}

// GenerateRsaKeys 生成RSA密钥对.bits为密钥位数,通常为1024或2048.
func (ue *UEncrypt) GenerateRsaKeys(bits int) (private []byte, public []byte, err error) {
	// 生成私钥文件
	var privateKey *rsa.PrivateKey
	privateKey, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	privateBuff := new(bytes.Buffer)
	_ = pem.Encode(privateBuff, block)

	// 生成公钥文件
	var derPkix []byte
	publicKey := &privateKey.PublicKey
	derPkix, _ = x509.MarshalPKIXPublicKey(publicKey)
	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}
	publicBuff := new(bytes.Buffer)
	_ = pem.Encode(publicBuff, block)

	private = privateBuff.Bytes()
	public = publicBuff.Bytes()

	return
}

// RsaPublicEncrypt RSA公钥加密.
// clearText为明文,publicKey为公钥.
func (ue *UEncrypt) RsaPublicEncrypt(clearText, publicKey []byte) ([]byte, error) {
	// 解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}

	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 类型断言
	pubKey := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, clearText)
}

// RsaPrivateDecrypt RSA私钥解密.比加密耗时.
// cipherText为密文,privateKey为私钥.
func (ue *UEncrypt) RsaPrivateDecrypt(cipherText, privateKey []byte) ([]byte, error) {
	// 获取私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, fmt.Errorf("private key error! ")
	}

	// 解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
}

// RsaPrivateEncrypt RSA私钥加密.比解密耗时.
// clearText为明文,privateKey为私钥.
func (ue *UEncrypt) RsaPrivateEncrypt(clearText, privateKey []byte) ([]byte, error) {
	// 获取私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, fmt.Errorf("private key error! ")
	}

	// 解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.SignPKCS1v15(nil, priv, crypto.Hash(0), clearText)
}

// RsaPublicDecrypt RSA公钥解密.
// cipherText为密文,publicKey为公钥.
func (ue *UEncrypt) RsaPublicDecrypt(cipherText, publicKey []byte) ([]byte, error) {
	// 解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}

	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 类型断言
	pubKey := pubInterface.(*rsa.PublicKey)

	c := new(big.Int)
	m := new(big.Int)
	m.SetBytes(cipherText)
	e := big.NewInt(int64(pubKey.E))
	c.Exp(m, e, pubKey.N)
	out := c.Bytes()
	olen := len(out)
	skip := 0
	for i := 2; i < olen; i++ {
		if (i+1 < olen) && out[i] == 0xff && out[i+1] == 0 {
			skip = i + 2
			break
		}
	}

	return out[skip:], nil
}
