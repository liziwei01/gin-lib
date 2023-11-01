/*
 * @Author: liziwei01
 * @Date: 2022-03-09 19:26:42
 * @LastEditors: liziwei01
 * @LastEditTime: 2023-11-01 11:21:57
 * @Description: file content
 */
package mysql

import (
	"context"

	"github.com/didi/gendry/builder"
)

const (
	// 普通insert
	insertCommon = iota

	// ignore insert
	insertIgnore

	// replace insert
	insertReplace

	// on duplicate key update
	insertOnDuplicate
)

type Builder interface {
	CompileContext(ctx context.Context, c Client) (cond string, values []interface{}, err error)
}

type result struct {
	cond   string
	values []interface{}
	query  string
}

// SelectBuilder 默认的select sql builder
type SelectBuilder struct {
	table  string
	where  map[string]interface{}
	fields []string

	res *result
}

// InsertBuilder 默认的select sql builder
type InsertBuilder struct {
	table  string
	data   []map[string]interface{}
	update map[string]interface{}
	typ    int

	res *result
}

type UpdateBuilder struct {
	table  string
	where  map[string]interface{}
	update map[string]interface{}

	res *result
}

type DeleteBuilder struct {
	table string
	where map[string]interface{}

	res *result
}

type RawBuilder struct {
	table string
	sql   string
	args  []interface{}

	res *result
}

func NewSelectBuilder(table string, where map[string]interface{}, fields []string) *SelectBuilder {
	return &SelectBuilder{
		table:  table,
		where:  where,
		fields: fields,
	}
}

func NewInsertBuilder(table string, data []map[string]interface{}, typ int, update ...map[string]interface{}) *InsertBuilder {
	if len(update) > 0 {
		return &InsertBuilder{
			table:  table,
			data:   data,
			typ:    typ,
			update: update[0],
		}
	}
	return &InsertBuilder{
		table: table,
		data:  data,
		typ:   typ,
	}
}

func NewUpdateBuilder(table string, where map[string]interface{}, update map[string]interface{}) *UpdateBuilder {
	return &UpdateBuilder{
		table:  table,
		where:  where,
		update: update,
	}
}

func NewDeleteBuilder(table string, where map[string]interface{}) *DeleteBuilder {
	return &DeleteBuilder{
		table: table,
		where: where,
	}
}

func NewRawBuilder(sql string, args []interface{}) *RawBuilder {
	return &RawBuilder{
		sql:  sql,
		args: args,
	}
}

// func getResult(c Client, cond string, values []interface{}) *result {
// 	res := &result{
// 		cond:   cond,
// 		values: values,
// 		query:  fmt.Sprintf(cond, values...),
// 	}
// 	if c.sqlloglen() != -1 {
// 		res.query = res.query[0:c.sqlloglen()]
// 	}
// 	return res
// }

// func (b *SelectBuilder) Result() *result {
// 	return b.res
// }

func (b *SelectBuilder) CompileContext(ctx context.Context, c Client) (string, []interface{}, error) {
	cond, values, err := builder.BuildSelect(b.table, b.where, b.fields)
	log(ctx, c, cond, values)
	return cond, values, err
}

// func (b *InsertBuilder) Result() *result {
// 	return b.res
// }

func (b *InsertBuilder) CompileContext(ctx context.Context, c Client) (string, []interface{}, error) {
	var (
		cond   string
		values []interface{}
		err    error
	)
	switch b.typ {
	case insertCommon:
		cond, values, err = builder.BuildInsert(b.table, b.data)
	case insertIgnore:
		cond, values, err = builder.BuildInsertIgnore(b.table, b.data)
	case insertReplace:
		cond, values, err = builder.BuildReplaceInsert(b.table, b.data)
	case insertOnDuplicate:
		cond, values, err = builder.BuildInsertOnDuplicate(b.table, b.data, b.update)
	}
	log(ctx, c, cond, values)
	return cond, values, err
}

// func (b *UpdateBuilder) Result() *result {
// 	return b.res
// }

func (b *UpdateBuilder) CompileContext(ctx context.Context, c Client) (string, []interface{}, error) {
	cond, values, err := builder.BuildUpdate(b.table, b.where, b.update)
	log(ctx, c, cond, values)
	return cond, values, err
}

// func (b *DeleteBuilder) Result() *result {
// 	return b.res
// }

func (b *DeleteBuilder) CompileContext(ctx context.Context, c Client) (string, []interface{}, error) {
	cond, values, err := builder.BuildDelete(b.table, b.where)
	log(ctx, c, cond, values)
	return cond, values, err
}

// func (b *RawBuilder) Result() *result {
// 	return b.res
// }

func (b *RawBuilder) CompileContext(ctx context.Context, c Client) (string, []interface{}, error) {
	log(ctx, c, b.sql, b.args)
	return b.sql, b.args, nil
}

func log(ctx context.Context, c Client, cond string, values []interface{}) {
	// query := fmt.Sprintf(cond, values...)
	// if c.sqlloglen() != -1 {
	// 	query = query[0:c.sqlloglen()]
	// }
	// if len(query) != 0 {
	// 	logit.Logger.Info("[MySQL] [requestID]=%d, [query]=%s", ctx.Value("requestID"), query)
	// }
}
