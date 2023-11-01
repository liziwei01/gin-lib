package mysql

import (
	"context"
	"database/sql"

	"github.com/didi/gendry/scanner"
	_ "github.com/go-sql-driver/mysql"
)

func (dao *client) Query(ctx context.Context, tableName string, where map[string]interface{}, columns []string, data interface{}) error {
	builder := NewSelectBuilder(tableName, where, columns)
	return QueryWithBuilder(ctx, dao, builder, data)
}

func (dao *client) Insert(ctx context.Context, tableName string, data []map[string]interface{}) (sql.Result, error) {
	builder := NewInsertBuilder(tableName, data, insertCommon)
	return ExecWithBuilder(ctx, dao, builder)
}

func (dao *client) InsertIgnore(ctx context.Context, tableName string, data []map[string]interface{}) (sql.Result, error) {
	builder := NewInsertBuilder(tableName, data, insertIgnore)
	return ExecWithBuilder(ctx, dao, builder)
}

func (dao *client) InsertReplace(ctx context.Context, tableName string, data []map[string]interface{}) (sql.Result, error) {
	builder := NewInsertBuilder(tableName, data, insertReplace)
	return ExecWithBuilder(ctx, dao, builder)
}

func (dao *client) InsertOnDuplicate(ctx context.Context, tableName string, data []map[string]interface{}, update map[string]interface{}) (sql.Result, error) {
	builder := NewInsertBuilder(tableName, data, insertOnDuplicate, update)
	return ExecWithBuilder(ctx, dao, builder)
}

func (dao *client) Update(ctx context.Context, tableName string, where map[string]interface{}, update map[string]interface{}) (sql.Result, error) {
	builder := NewUpdateBuilder(tableName, where, update)
	return ExecWithBuilder(ctx, dao, builder)
}

func (dao *client) Delete(ctx context.Context, tableName string, where map[string]interface{}) (sql.Result, error) {
	builder := NewDeleteBuilder(tableName, where)
	return ExecWithBuilder(ctx, dao, builder)
}

func (dao *client) ExecRaw(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	builder := NewRawBuilder(sql, args)
	return ExecWithBuilder(ctx, dao, builder)
}

// QueryWithBuilder 传入一个 SQLBuilder 并执行 QueryContext
func QueryWithBuilder(ctx context.Context, client Client, builder Builder, data interface{}) error {
	db, err := client.connect(ctx)
	if err != nil {
		return err
	}
	cond, values, err := builder.CompileContext(ctx, client)
	if err != nil {
		return err
	}
	rows, err := db.QueryContext(ctx, cond, values...)
	if err != nil {
		return err
	}
	return scanner.ScanClose(rows, data)
}

func ExecWithBuilder(ctx context.Context, client Client, builder Builder) (sql.Result, error) {
	db, err := client.connect(ctx)
	if err != nil {
		return nil, err
	}
	cond, values, err := builder.CompileContext(ctx, client)
	if err != nil {
		return nil, err
	}
	return db.ExecContext(ctx, cond, values...)
}

func Execraw(ctx context.Context, client Client, builder Builder) (sql.Result, error) {
	db, err := client.connect(ctx)
	if err != nil {
		return nil, err
	}
	cond, values, err := builder.CompileContext(ctx, client)
	if err != nil {
		return nil, err
	}
	return db.ExecContext(ctx, cond, values...)
}

var _ Client = (*client)(nil)
