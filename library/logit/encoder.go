/*
 * @Author: liziwei01
 * @Date: 2023-10-30 11:28:48
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 11:37:51
 * @Description: 日志输出器
 */
package logit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"
)

// FieldEncoder is a strongly-typed, encoding-agnostic interface for adding a
// map- or struct-like object to the logging context. Like maps, ObjectEncoders
// aren't safe for concurrent use (though typical use shouldn't require locks).
type FieldEncoder interface {
	io.WriterTo

	// Built-in types.
	AddBinary(key string, value []byte) // for arbitrary bytes
	AddBool(key string, value bool)
	AddByteString(key string, value []byte) // for UTF-8 encoded bytes
	AddDuration(key string, value time.Duration)
	AddFloat64(key string, value float64)
	AddFloat32(key string, value float32)
	AddInt(key string, value int)
	AddInt64(key string, value int64)
	AddInt32(key string, value int32)
	AddInt16(key string, value int16)
	AddInt8(key string, value int8)
	AddString(key, value string)
	AddTime(key string, value time.Time)
	AddUint(key string, value uint)
	AddUint64(key string, value uint64)
	AddUint32(key string, value uint32)
	AddUint16(key string, value uint16)
	AddUint8(key string, value uint8)
	AddUintptr(key string, value uintptr)
	AddError(key string, value error)

	// AddReflected uses reflection to serialize arbitrary objects, so it can be
	// slow and allocation-heavy.
	AddReflected(key string, value interface{}) error

	Reset()
}

// TexEncoderOption 文本encoder的配置
type TexEncoderOption struct {
	KeyPrefix   []byte
	KeySuffix   []byte
	ValuePrefix []byte
	ValueSuffix []byte
	Delim       []byte
	LineBreak   []byte // 换行符
}

// DefaultTextEncoderOption 默认的TextEncoder 选项
var DefaultTextEncoderOption = TexEncoderOption{
	KeyPrefix:   nil,
	KeySuffix:   nil,
	ValuePrefix: []byte("["),
	ValueSuffix: []byte("]"),
	Delim:       []byte(" "),
	LineBreak:   []byte("\n"),
}

// DefaultTextEncoderPool 默认的text encoder pool
var DefaultTextEncoderPool = NewEncoderPool(func() FieldEncoder {
	return NewTextEncoder(DefaultTextEncoderOption)
})

// DefaultJSONEncoderPool 默认json encoder pool
var DefaultJSONEncoderPool = NewEncoderPool(func() FieldEncoder {
	return NewJSONEncoder()
})

// NewTextEncoder 创建text encoder
func NewTextEncoder(opt TexEncoderOption) *TextEncoder {
	return &TextEncoder{
		opt: opt,
	}
}

// TextEncoder 普通的 key-value 格式的 Encoder
type TextEncoder struct {
	opt TexEncoderOption
	buf bytes.Buffer
}

// 实现io.WriterTo接口，将buf写入到另一个io.Writer中
// 实际使用中，一般是写入到另一个bytes.Buffer中
func (e *TextEncoder) WriteTo(w io.Writer) (int64, error) {
	// 去掉自己的最后一个分隔符，默认是一个空格
	if e.buf.Len() > len(e.opt.Delim) {
		e.buf.Truncate(e.buf.Len() - len(e.opt.Delim))
	}
	// 给自己添加换行符
	if len(e.opt.LineBreak) > 0 {
		e.buf.Write(e.opt.LineBreak)
	}
	// 把自己写入到另一个io.Writer中
	return e.buf.WriteTo(w)
}

// AddBinary 二进制字段
func (e *TextEncoder) AddBinary(key string, value []byte) {
	e.write(key, value)
}

// AddByteString bytes字符串
func (e *TextEncoder) AddByteString(key string, value []byte) {
	e.write(key, value)
}

// AddBool bool类型
func (e *TextEncoder) AddBool(key string, value bool) {
	if value {
		e.write(key, []byte("true"))
	} else {
		e.write(key, []byte("false"))
	}
}

// AddDuration 时间间隔
func (e *TextEncoder) AddDuration(key string, value time.Duration) {
	if value < time.Microsecond {
		e.write(key, []byte("0"))
		return
	}
	dur := strconv.FormatFloat(float64(value.Nanoseconds())/float64(time.Millisecond), 'f', 3, 64)
	e.write(key, []byte(dur))
}

// AddFloat64 float64
func (e *TextEncoder) AddFloat64(key string, value float64) {
	e.writeString(key, strconv.FormatFloat(value, 'f', -1, 64))
}

// AddFloat32 Float32
func (e *TextEncoder) AddFloat32(key string, value float32) {
	e.writeString(key, strconv.FormatFloat(float64(value), 'f', -1, 32))
}

// AddInt Int
func (e *TextEncoder) AddInt(key string, value int) {
	e.writeString(key, strconv.FormatInt(int64(value), 10))
}

// AddInt64 Int64
func (e *TextEncoder) AddInt64(key string, value int64) {
	e.writeString(key, strconv.FormatInt(value, 10))
}

// AddInt32 Int32
func (e *TextEncoder) AddInt32(key string, value int32) {
	e.writeString(key, strconv.FormatInt(int64(value), 10))
}

// AddInt16 Int16
func (e *TextEncoder) AddInt16(key string, value int16) {
	e.writeString(key, strconv.FormatInt(int64(value), 10))
}

// AddInt8 Int8
func (e *TextEncoder) AddInt8(key string, value int8) {
	e.writeString(key, strconv.FormatInt(int64(value), 10))
}

// AddString String
func (e *TextEncoder) AddString(key string, value string) {
	e.writeString(key, value)
}

// AddTime 时间类型
func (e *TextEncoder) AddTime(key string, value time.Time) {
	if value.IsZero() {
		e.writeString(key, "0")
	} else {
		e.writeString(key, strconv.FormatInt(value.UnixNano()/int64(time.Millisecond), 10))
	}
}

// AddUint Uint
func (e *TextEncoder) AddUint(key string, value uint) {
	e.writeString(key, strconv.FormatUint(uint64(value), 10))
}

// AddUint64 Uint64
func (e *TextEncoder) AddUint64(key string, value uint64) {
	e.writeString(key, strconv.FormatUint(value, 10))
}

// AddUint32 Uint32
func (e *TextEncoder) AddUint32(key string, value uint32) {
	e.writeString(key, strconv.FormatUint(uint64(value), 10))
}

// AddUint16 Uint16
func (e *TextEncoder) AddUint16(key string, value uint16) {
	e.writeString(key, strconv.FormatUint(uint64(value), 10))
}

// AddUint8 Uint8
func (e *TextEncoder) AddUint8(key string, value uint8) {
	e.writeString(key, strconv.FormatUint(uint64(value), 10))
}

// AddUintptr Uintptr
func (e *TextEncoder) AddUintptr(key string, value uintptr) {
	e.writeString(key, "0x"+strconv.FormatUint(uint64(value), 16))
}

// AddError  Error
func (e *TextEncoder) AddError(key string, value error) {
	if value == nil {
		e.writeString(key, "nil")
	} else {
		e.writeString(key, value.Error())
	}
}

// AddReflected Reflected
func (e *TextEncoder) AddReflected(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil { // 忽略json marshal失败，将错误信息写到error
		e.AddError(key, err)
	}
	e.write(key, b)
	return nil
}

// 写入field的终极实现，将field写成key[value]的形式
// 其中，key的前后缀，value的前后缀，多个kv之间的分隔符，都可以自定义，默认为key[value]，分隔符为一个空格
// 例如：key1[value1] key2[value2]
func (e *TextEncoder) write(key string, val []byte) {
	if len(e.opt.KeyPrefix) > 0 {
		_, _ = e.buf.Write(e.opt.KeyPrefix)
	}
	_, _ = e.buf.WriteString(key)

	if len(e.opt.KeySuffix) > 0 {
		_, _ = e.buf.Write(e.opt.KeySuffix)
	}

	if len(e.opt.ValuePrefix) > 0 {
		_, _ = e.buf.Write(e.opt.ValuePrefix)
	}
	_, _ = e.buf.Write(val)

	if len(e.opt.ValueSuffix) > 0 {
		_, _ = e.buf.Write(e.opt.ValueSuffix)
	}
	_, _ = e.buf.Write(e.opt.Delim)
}

func (e *TextEncoder) writeString(key string, val string) {
	e.write(key, []byte(val))
}

// Reset 重置
func (e *TextEncoder) Reset() {
	e.buf.Reset()
}

var _ FieldEncoder = (*TextEncoder)(nil)

// JSONEncoder 以 {key: value} 格式输出 JSON 格式的Encoder
type JSONEncoder struct {
	kv map[string]interface{}

	LineBreak []byte // 换行符
}

// NewJSONEncoder 打包输出为json格式
func NewJSONEncoder() FieldEncoder {
	return &JSONEncoder{
		// 32:手百日志大概30个字段
		kv: make(map[string]interface{}, 32),

		LineBreak: []byte("\n"),
	}
}

// WriteTo 写入
func (e *JSONEncoder) WriteTo(w io.Writer) (int64, error) {
	b, err := json.Marshal(e.kv)
	if err != nil {
		return 0, err
	}
	if len(e.LineBreak) > 0 {
		b = append(b, e.LineBreak...)
	}
	n, err := w.Write(b)
	return int64(n), err
}

// AddBinary  Binary
func (e *JSONEncoder) AddBinary(key string, value []byte) {
	e.kv[key] = value
}

// AddBool  Bool
func (e *JSONEncoder) AddBool(key string, value bool) {
	e.kv[key] = value
}

// AddByteString  ByteString
func (e *JSONEncoder) AddByteString(key string, value []byte) {
	e.kv[key] = value
}

// AddDuration duration
func (e *JSONEncoder) AddDuration(key string, value time.Duration) {
	e.kv[key] = value.String()
}

// AddFloat64 Float64
func (e *JSONEncoder) AddFloat64(key string, value float64) {
	e.kv[key] = value
}

// AddFloat32 Float32
func (e *JSONEncoder) AddFloat32(key string, value float32) {
	e.kv[key] = value
}

// AddInt Int
func (e *JSONEncoder) AddInt(key string, value int) {
	e.kv[key] = value
}

// AddInt64 Int64
func (e *JSONEncoder) AddInt64(key string, value int64) {
	e.kv[key] = value
}

// AddInt32 Int32
func (e *JSONEncoder) AddInt32(key string, value int32) {
	e.kv[key] = value
}

// AddInt16 Int16
func (e *JSONEncoder) AddInt16(key string, value int16) {
	e.kv[key] = value
}

// AddInt8 Int8
func (e *JSONEncoder) AddInt8(key string, value int8) {
	e.kv[key] = value
}

// AddString String
func (e *JSONEncoder) AddString(key string, value string) {
	e.kv[key] = value
}

// AddTime Time
func (e *JSONEncoder) AddTime(key string, value time.Time) {
	e.kv[key] = value.Format(time.RFC3339Nano)
}

// AddUint Uint
func (e *JSONEncoder) AddUint(key string, value uint) {
	e.kv[key] = value
}

// AddUint64 Uint64
func (e *JSONEncoder) AddUint64(key string, value uint64) {
	e.kv[key] = value
}

// AddUint32 Uint32
func (e *JSONEncoder) AddUint32(key string, value uint32) {
	e.kv[key] = value
}

// AddUint16 Uint16
func (e *JSONEncoder) AddUint16(key string, value uint16) {
	e.kv[key] = value
}

// AddUint8 Uint8
func (e *JSONEncoder) AddUint8(key string, value uint8) {
	e.kv[key] = value
}

// AddUintptr Uintptr
func (e *JSONEncoder) AddUintptr(key string, value uintptr) {
	e.kv[key] = value
}

// AddReflected Reflected
func (e *JSONEncoder) AddReflected(key string, value interface{}) error {
	e.kv[key] = value
	return nil
}

// AddError  Error
func (e *JSONEncoder) AddError(key string, value error) {
	if value != nil {
		e.kv[key] = value.Error()
		return
	}
	e.kv[key] = nil
}

// Reset 重置
func (e *JSONEncoder) Reset() {
	e.kv = make(map[string]interface{}, len(e.kv))
}

var _ FieldEncoder = (*JSONEncoder)(nil)

// NewEncoderPool 创建一个encoder 对象池
func NewEncoderPool(newFn func() FieldEncoder) EncoderPool {
	return &encoderPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return newFn()
			},
		},
	}
}

// EncoderPool encoder对象池
type EncoderPool interface {
	// 获取一个新对象
	Get() FieldEncoder

	// 返回对象池,会自动调用encoder 的Reset()方法
	Put(enc FieldEncoder)
}

type encoderPool struct {
	pool *sync.Pool
}

func (f *encoderPool) Get() FieldEncoder {
	return f.pool.Get().(FieldEncoder)
}

func (f *encoderPool) Put(enc FieldEncoder) {
	enc.Reset()
	f.pool.Put(enc)
}

var _ EncoderPool = (*encoderPool)(nil)

const (
	encoderPoolNameDefaultText = "default_text"
	encoderPoolNameDefaultJSON = "default_json"
)

var encoderPools = map[interface{}]EncoderPool{
	encoderPoolNameDefaultText: DefaultTextEncoderPool,
	encoderPoolNameDefaultJSON: DefaultJSONEncoderPool,
}

// RegisterEncoderPool 注册一个新的encoder pool
func RegisterEncoderPool(name interface{}, pool EncoderPool) error {
	if _, has := encoderPools[name]; has {
		return fmt.Errorf("name=%v already exists", name)
	}
	encoderPools[name] = pool
	return nil
}

// GetEncoderPool 获取一个encoder pool
// 若不存在，会返回nil
func GetEncoderPool(name interface{}) EncoderPool {
	switch name {
	case encoderPoolNameDefaultText:
		return DefaultTextEncoderPool
	case encoderPoolNameDefaultJSON:
		return DefaultJSONEncoderPool
	}

	return encoderPools[name]
}
