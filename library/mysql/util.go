package mysql

import (
	"reflect"
	"strings"
)

// Struct2Where 将对象转换为SQLBuilder的查询条件
func Struct2Where(raw interface{}) map[string]interface{} {
	return struct2map(raw, false)
}

// Struct2AssignList 将多个对象转换为 可用于InsertBuilder的二维数组map
func Struct2AssignList(raws ...interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(raws))
	for _, raw := range raws {
		rst = append(rst, struct2map(raw, true))
	}
	return rst
}

// Struct2Assign 将一个对象，转换为 UpdateBuilder 的map
func Struct2Assign(raw interface{}) map[string]interface{} {
	return struct2map(raw, true)
}

func struct2map(raw interface{}, ignoreOpt bool) map[string]interface{} {
	rst := map[string]interface{}{}
	if raw == nil {
		return rst
	}

	structType := reflect.TypeOf(raw)
	if kind := structType.Kind(); kind == reflect.Ptr || kind == reflect.Interface {
		structType = structType.Elem()
	}

	structVal := reflect.ValueOf(raw)
	if structVal.IsZero() {
		return rst
	}
	if structVal.Kind() == reflect.Ptr {
		structVal = structVal.Elem()
	}

	if structType.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < structVal.NumField(); i++ {
		valField := structVal.Field(i)

		valFieldKind := valField.Kind()
		if valFieldKind == reflect.Ptr {
			if valField.IsNil() {
				continue
			}

			valField = valField.Elem()
		}

		if valFieldKind == reflect.Slice {
			if valField.IsZero() {
				continue
			}
		}

		typeField := structType.Field(i)
		dbTag := typeField.Tag.Get(tagName)
		if dbTag == "-" {
			continue
		}
		key, opt := tagSplitter(dbTag)
		if key == "" {
			key = typeField.Name
		}
		if ignoreOpt {
			rst[key] = valField.Interface()
		} else {
			if opt == "" || opt == "=" {
				rst[key] = valField.Interface()
			} else {
				rst[key+" "+opt] = valField.Interface()
			}
		}
	}

	return rst
}

func tagSplitter(dbTag string) (key, opt string) {
	if dbTag == "" {
		return "", ""
	}
	i := strings.Index(dbTag, ",")
	if i == -1 {
		return dbTag, ""
	}
	return strings.TrimSpace(dbTag[:i]), strings.TrimSpace(dbTag[i+1:])
}
