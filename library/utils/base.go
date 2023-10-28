/*
 * @Author: liziwei01
 * @Date: 2022-03-03 19:52:24
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-10-28 14:12:00
 * @Description:
 */
package utils

import (
	"net"
	"regexp"
	"time"
)

type (
	USlice   byte
	UConfig  byte
	URequest byte
	UMd5     byte
	UFile    byte
	UUUID    byte
	UTime    byte
	UEncrypt byte
	UXlsx    byte
)

var (
	Slice   *USlice
	Config  *UConfig
	Request *URequest
	Md5     *UMd5
	File    *UFile
	UUID    *UUUID
	Time    *UTime
	Encrypt *UEncrypt
	Xlsx    *UXlsx
)

type (

	// UFileCover 枚举类型,文件是否覆盖
	UFileCover int8
	// UFileType 枚举类型,文件类型
	UFileType uint8
	// UFileTree 枚举类型,文件树查找类型
	UFileTree uint8
	// URandString 枚举类型,随机字符串类型
	URandString uint8
	// UCaseSwitch 枚举类型,大小写开关
	UCaseSwitch uint8
	// UPadType 枚举类型,字符串填充类型
	UPadType uint8
	// UPKCSType 枚举类型,PKCS填充类型
	UPKCSType int8

	// FileFilter 文件过滤函数
	FileFilter func(string) bool

	// CallBack 回调执行函数,无参数且无返回值
	CallBack func()
)

const (

	// FILE_COVER_ALLOW 文件覆盖,允许
	FILE_COVER_ALLOW UFileCover = 1
	// FILE_COVER_IGNORE 文件覆盖,忽略
	FILE_COVER_IGNORE UFileCover = 0
	// FILE_COVER_DENY 文件覆盖,禁止
	FILE_COVER_DENY UFileCover = -1

	// FILE_TYPE_ANY 文件类型-任意
	FILE_TYPE_ANY UFileType = 0
	// FILE_TYPE_LINK 文件类型-链接文件
	FILE_TYPE_LINK UFileType = 1
	// FILE_TYPE_REGULAR 文件类型-常规文件(不包括链接)
	FILE_TYPE_REGULAR UFileType = 2
	// FILE_TYPE_COMMON 文件类型-普通文件(包括常规和链接)
	FILE_TYPE_COMMON UFileType = 3

	// FILE_TREE_ALL 文件树,查找所有(包括目录和文件)
	FILE_TREE_ALL UFileTree = 3
	// FILE_TREE_DIR 文件树,仅查找目录
	FILE_TREE_DIR UFileTree = 2
	// FILE_TREE_FILE 文件树,仅查找文件
	FILE_TREE_FILE UFileTree = 1

	// RAND_STRING_ALPHA 随机字符串类型,字母
	RAND_STRING_ALPHA URandString = 0
	// RAND_STRING_NUMERIC 随机字符串类型,数值
	RAND_STRING_NUMERIC URandString = 1
	// RAND_STRING_ALPHANUM 随机字符串类型,字母+数值
	RAND_STRING_ALPHANUM URandString = 2
	// RAND_STRING_SPECIAL 随机字符串类型,字母+数值+特殊字符
	RAND_STRING_SPECIAL URandString = 3
	// RAND_STRING_CHINESE 随机字符串类型,仅中文
	RAND_STRING_CHINESE URandString = 4

	// CASE_NONE 忽略大小写
	CASE_NONE UCaseSwitch = 0
	// CASE_LOWER 检查小写
	CASE_LOWER UCaseSwitch = 1
	// CASE_UPPER 检查大写
	CASE_UPPER UCaseSwitch = 2

	// PAD_LEFT 左侧填充
	PAD_LEFT UPadType = 0
	// PAD_RIGHT 右侧填充
	PAD_RIGHT UPadType = 1
	// PAD_BOTH 两侧填充
	PAD_BOTH UPadType = 2

	// PKCS_NONE 不进行填充
	PKCS_NONE UPKCSType = -1
	// PKCS_ZERO PKCS 0值填充
	PKCS_ZERO UPKCSType = 0
	// PKCS_SEVEN 即PKCS7
	PKCS_SEVEN UPKCSType = 7

	//默认浮点数精确小数位数
	FLOAT_DECIMAL = 10

	//AuthCode 动态密钥长度,须<32
	DYNAMIC_KEY_LEN = 8

	//检查连接超时的时间
	CHECK_CONNECT_TIMEOUT = time.Second * 5

	// 正则模式-全中文
	PATTERN_CHINESE_ALL = "^[\u4e00-\u9fa5]+$"

	// 正则模式-中文名称
	PATTERN_CHINESE_NAME = "^[\u4e00-\u9fa5][.•·\u4e00-\u9fa5]{0,30}[\u4e00-\u9fa5]$"

	// 正则模式-多字节字符
	PATTERN_MULTIBYTE = "[^\x00-\x7F]"

	// 正则模式-ASCII字符
	PATTERN_ASCII = "^[\x00-\x7F]+$"

	// 正则模式-全角字符
	PATTERN_FULLWIDTH = "[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"

	// 正则模式-半角字符
	PATTERN_HALFWIDTH = "[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"

	// 正则模式-词语,不以下划线开头的中文、英文、数字、下划线
	PATTERN_WORD = "^[a-zA-Z0-9\u4e00-\u9fa5][a-zA-Z0-9_\u4e00-\u9fa5]+$"

	// 正则模式-浮点数
	PATTERN_FLOAT = `^(-?\d+)(\.\d+)?`

	// 正则模式-邮箱
	PATTERN_EMAIL = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

	// 正则模式-用户名-英文
	PATTERN_USERNAMEEN = `^[a-zA-Z0-9_.]+$`

	// 正则模式-大陆手机号
	PATTERN_MOBILECN = `^1[3-9]\d{9}$`

	// 正则模式-固定电话
	PATTERN_TEL_FIX = `^(010|02\d{1}|0[3-9]\d{2})-\d{7,9}(-\d+)?$`

	// 正则模式-400或800
	PATTERN_TEL_4800 = `^[48]00\d?(-?\d{3,4}){2}$`

	// 正则模式-座机号(固定电话或400或800)
	PATTERN_TELEPHONE = `(` + PATTERN_TEL_FIX + `)|(` + PATTERN_TEL_4800 + `)`

	// 正则模式-电话(手机或固话)
	PATTERN_PHONE = `(` + PATTERN_MOBILECN + `)|(` + PATTERN_TEL_FIX + `)`

	// 正则模式-日期时间
	PATTERN_DATETIME = `^[0-9]{4}(|\-[0-9]{2}(|\-[0-9]{2}(|\s+[0-9]{2}(|:[0-9]{2}(|:[0-9]{2})))))$`

	// 正则模式-身份证号码,18位或15位
	PATTERN_CREDIT_NO = `(^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$)`

	// 正则模式-小写英文
	PATTERN_ALPHA_LOWER = `^[a-z]+$`

	// 正则模式-大写英文
	PATTERN_ALPHA_UPPER = `^[A-Z]+$`

	// 正则模式-字母和数字
	PATTERN_ALPHA_NUMERIC = `^[a-zA-Z0-9]+$`

	// 正则模式-十六进制颜色
	PATTERN_HEXCOLOR = `^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`

	// 正则模式-RGB颜色
	PATTERN_RGBCOLOR = "^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$"

	// 正则模式-全空白字符
	PATTERN_WHITESPACE_ALL = "^[[:space:]]+$"

	// 正则模式-带空白字符
	PATTERN_WHITESPACE_HAS = ".*[[:space:]]"

	// 正则模式-连续空白符
	PATTERN_WHITESPACE_DUPLICATE = `[[:space:]]{2,}|[\s\p{Zs}]{2,}`

	// 正则模式-base64字符串
	PATTERN_BASE64 = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"

	// 正则模式-base64编码图片
	PATTERN_BASE64_IMAGE = `^data:\s*(image|img)\/(\w+);base64`

	// 正则模式-HTML标签
	PATTERN_HTML_TAGS = `<(.|\n)*?>`

	// 正则模式-DNS名称
	PATTERN_DNSNAME = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`

	// 正则模式-MD5
	PATTERN_MD5 = `^(?i)([0-9a-h]{32})$`

	// 正则模式-SHA1
	PATTERN_SHA1 = `^(?i)([0-9a-h]{40})$`

	// 正则模式-SHA256
	PATTERN_SHA256 = `^(?i)([0-9a-h]{64})$`

	// 正则模式-SHA512
	PATTERN_SHA512 = `^(?i)([0-9a-h]{128})$`
)

var (
	RegFormatDir             = regexp.MustCompile(`[\/]{2,}`) //连续的"//"或"\\"或"\/"或"/\"
	RegChineseAll            = regexp.MustCompile(PATTERN_CHINESE_ALL)
	RegChineseName           = regexp.MustCompile(PATTERN_CHINESE_NAME)
	RegWord                  = regexp.MustCompile(PATTERN_WORD)
	RegMultiByte             = regexp.MustCompile(PATTERN_MULTIBYTE)
	RegFullWidth             = regexp.MustCompile(PATTERN_FULLWIDTH)
	RegHalfWidth             = regexp.MustCompile(PATTERN_HALFWIDTH)
	RegFloat                 = regexp.MustCompile(PATTERN_FLOAT)
	RegEmail                 = regexp.MustCompile(PATTERN_EMAIL)
	RegMobilecn              = regexp.MustCompile(PATTERN_MOBILECN)
	RegTelephone             = regexp.MustCompile(PATTERN_TELEPHONE)
	RegPhone                 = regexp.MustCompile(PATTERN_PHONE)
	RegDatetime              = regexp.MustCompile(PATTERN_DATETIME)
	RegCreditno              = regexp.MustCompile(PATTERN_CREDIT_NO)
	RegAlphaLower            = regexp.MustCompile(PATTERN_ALPHA_LOWER)
	RegAlphaUpper            = regexp.MustCompile(PATTERN_ALPHA_UPPER)
	RegAlphaNumeric          = regexp.MustCompile(PATTERN_ALPHA_NUMERIC)
	RegHexcolor              = regexp.MustCompile(PATTERN_HEXCOLOR)
	RegRgbcolor              = regexp.MustCompile(PATTERN_RGBCOLOR)
	RegWhitespace            = regexp.MustCompile(`\s`)
	RegWhitespaceAll         = regexp.MustCompile(PATTERN_WHITESPACE_ALL)
	RegWhitespaceHas         = regexp.MustCompile(PATTERN_WHITESPACE_HAS)
	RegWhitespaceDuplicate   = regexp.MustCompile(PATTERN_WHITESPACE_DUPLICATE)
	RegBase64                = regexp.MustCompile(PATTERN_BASE64)
	RegBase64Image           = regexp.MustCompile(PATTERN_BASE64_IMAGE)
	RegHTMLTag               = regexp.MustCompile(PATTERN_HTML_TAGS)
	RegDNSname               = regexp.MustCompile(PATTERN_DNSNAME)
	RegURLBackslashDuplicate = regexp.MustCompile(`([^:])[\/]{2,}`) //URL中连续的"//"或"\\"或"\/"或"/\"
	RegMd5                   = regexp.MustCompile(PATTERN_MD5)
	RegSha1                  = regexp.MustCompile(PATTERN_SHA1)
	RegSha256                = regexp.MustCompile(PATTERN_SHA256)
	RegSha512                = regexp.MustCompile(PATTERN_SHA512)
	RegUsernameen            = regexp.MustCompile(PATTERN_USERNAMEEN)

	// PrivCidrs                 = regexp.MustCompile(PATTERN_ASCII)
	PrivCidrs []*net.IPNet
	// Uptime 当前服务启动时间
	Uptime = time.Now()
)
