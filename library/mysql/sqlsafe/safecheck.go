/*
 * @Author: liziwei01
 * @Date: 2023-11-04 15:45:37
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-04 15:45:51
 * @Description: file content
 */
package sqlsafe

import (
	"fmt"
	"strings"
)

type (
	// SafeChecker sql防注入check检测
	SafeChecker interface {
		SafeCheck()
	}

	// TableChecker table表名检测
	TableChecker interface {
		CheckTable(table string) error
	}

	// WhereChecker where检测
	WhereChecker interface {
		CheckWhere(where map[string]interface{}) error
	}

	// FieldChecker field列名检测
	FieldChecker interface {
		CheckField([]string) error
	}

	// CheckParams 需要check的参数列表
	CheckParams struct {
		Table  string
		Where  map[string]interface{}
		Fields []string
	}
)

// Check 中介
func Check(sc SafeChecker, p *CheckParams) error {
	if ct, ok := sc.(TableChecker); ok {
		if err := ct.CheckTable(p.Table); err != nil {
			return err
		}
	}
	if cw, ok := sc.(WhereChecker); ok {
		if err := cw.CheckWhere(p.Where); err != nil {
			return err
		}
	}
	if cf, ok := sc.(FieldChecker); ok {
		if err := cf.CheckField(p.Fields); err != nil {
			return err
		}
	}
	return nil
}

// nopChecker SafeCheck无检查规则实现
type nopChecker struct{}

// NopChecker 全局的、不包括任何检查规则的checker
var NopChecker = &nopChecker{}

// SafeCheck sql防注入安全检测
func (nop *nopChecker) SafeCheck() {}

// defaultChecker SafeCheck默认实现
type defaultChecker struct{}

// DefaultChecker 全局的、默认的checker
var DefaultChecker = &defaultChecker{}

// SafeCheck sql防注入安全检测
func (df *defaultChecker) SafeCheck() {}

// CheckTable 默认table检测逻辑
func (df *defaultChecker) CheckTable(Table string) error {
	if Table == "" {
		return fmt.Errorf("table name is empty")
	}
	Table = strings.TrimSpace(Table)
	field := strings.Split(Table, " ")
	if len(field) == 3 {
		if "AS" != strings.ToUpper(field[1]) || nil != checkTableName(field[0]) || nil != checkTableName(field[2]) {
			return fmt.Errorf("table name=%s is unsafe", Table)
		}
		return nil
	} else if len(field) == 1 {
		return checkTableName(field[0])
	}
	return fmt.Errorf("table name=%s is unsafe", Table)
}

// CheckField 默认field检测逻辑
func (df *defaultChecker) CheckField(Fields []string) error {
	return nil
}

// CheckWhere 默认where检测逻辑
// where支持的操作符如下
//
//	where := map[string]interface{}{
//	    "age >": 100,
//	    "_or": []map[string]interface{}{
//	        {
//	            "x1":    11,
//	            "x2 >=": 45,
//	        },
//	    },
//	    "_orderby": "fieldName asc",
//	    "_groupby": "fieldName",
//	    "_having": map[string]interface{}{"foo":"bar",},
//	    "_limit": []uint{offset, row_count},
//	    "_lockMode": "share",
//	}
func (df *defaultChecker) CheckWhere(Where map[string]interface{}) error {
	for k, v := range Where {
		// where条件的key
		err := checkFieldAndOp(k)
		if err != nil {
			return err
		}
		// where条件的value
		switch v.(type) {
		case string:
			if val, ok := Where["_orderby"]; ok {
				err := checkOrderBy(val.(string))
				if err != nil {
					return err
				}
			}
			if val, ok := Where["_groupby"]; ok {
				err := checkOrderBy(val.(string))
				if err != nil {
					return err
				}
			}
		case []map[string]interface{}:
			if val, ok := Where["_or"]; ok {
				err := checkOr(val.([]map[string]interface{}))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func checkTableName(checkStr string) error {
	for _, c := range checkStr {
		if isBasicWhiteChars(c) {
			continue
		}
		return fmt.Errorf("table name=%s is unsafe", checkStr)
	}
	return nil
}

func checkOrderBy(checkStr string) error {
	valSlice := strings.Split(checkStr, ",")
	for _, v := range valSlice {
		v = strings.TrimSpace(v)
		field := strings.Split(v, " ")
		if len(field) > 2 {
			return fmt.Errorf("orderby/groupby statement =%s has Extra space", field)
		} else if len(field) == 2 {
			direction := strings.ToUpper(field[1])
			if direction != "DESC" && direction != "ASC" {
				return fmt.Errorf("orderby/groupby statement =%s is unsafe", field)
			}
		}
		err := checkFieldName(field[0])
		if err != nil {
			return err
		}
	}
	return nil
}

func checkOr(checkArr []map[string]interface{}) error {
	for _, checkMap := range checkArr {
		for fieldName := range checkMap {
			err := checkFieldAndOp(fieldName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkFieldAndOp(checkStr string) error {
	checkStr = strings.TrimSpace(checkStr)
	field := strings.SplitN(checkStr, " ", 2)
	if len(field) == 2 {
		err := checkOp(field[1])
		if err != nil {
			return err
		}
	}
	return checkFieldName(field[0])
}

func checkFieldName(checkStr string) error {
	for _, c := range checkStr {
		if isBasicWhiteChars(c) || c == '.' {
			continue
		}
		return fmt.Errorf("field name=%s is unsafe", checkStr)
	}
	return nil
}

func checkOp(checkStr string) error {
	optArr := []string{"=", ">", "<", "=", "<=", ">=", "!=", "<>", "in", "not in", "like", "not like", "between", "not between"}
	for _, opt := range optArr {
		if opt == checkStr {
			return nil
		}
	}
	return fmt.Errorf("opt name is illegal")
}

func isBasicWhiteChars(c rune) bool {
	if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' || c == '$' {
		return true
	}
	return false
}

// WithSafeChecker 默认检查绑定的实现
type WithSafeChecker struct {
	checker SafeChecker
}

// SetSafeChecker 设置检查checker
func (w *WithSafeChecker) SetSafeChecker(sc SafeChecker) {
	w.checker = sc
}

// SafeChecker 获取检查checker
func (w *WithSafeChecker) SafeChecker() SafeChecker {
	if w.checker == nil {
		return DefaultChecker
	}
	return w.checker
}

var _ SafeChecker = (*defaultChecker)(nil)
var _ SafeChecker = (*nopChecker)(nil)
