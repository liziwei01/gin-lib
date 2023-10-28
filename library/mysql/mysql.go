package mysql

import (
	"context"
	"database/sql"

	"github.com/liziwei01/gin-lib/library/logit"

	"github.com/didi/gendry/scanner"
	_ "github.com/go-sql-driver/mysql"
)

func (dao *client) Query(ctx context.Context, tableName string, where map[string]interface{}, columns []string, data interface{}) error {
	builder := NewSelectBuilder(tableName, where, columns)
	err := QueryWithBuilder(ctx, dao, builder, data)
	if err != nil {
		logit.Logger.Warn("[MySQL] [Query] [requestID]=%d, [err]=%s", ctx.Value("requestID"), err.Error())
		return err
	}
	logit.Logger.Info("[MySQL] [Query] [requestID]=%d success", ctx.Value("requestID"))
	return nil
}

func (dao *client) Insert(ctx context.Context, tableName string, data []map[string]interface{}) (sql.Result, error) {
	builder := NewInsertBuilder(tableName, data, insertCommon)
	res, err := ExecWithBuilder(ctx, dao, builder)
	if err != nil {
		logit.Logger.Warn("[MySQL] [Insert] [requestID]=%d, [err]=%s", ctx.Value("requestID"), err.Error())
		return nil, err
	}
	logit.Logger.Info("[MySQL] [Insert] [requestID]=%d success", ctx.Value("requestID"))
	return res, nil
}

func (dao *client) InsertIgnore(ctx context.Context, tableName string, data []map[string]interface{}) (sql.Result, error) {
	builder := NewInsertBuilder(tableName, data, insertIgnore)
	res, err := ExecWithBuilder(ctx, dao, builder)
	if err != nil {
		logit.Logger.Warn("[MySQL] [InsertIgnore] [requestID]=%d, [err]=%s", ctx.Value("requestID"), err.Error())
		return nil, err
	}
	logit.Logger.Info("[MySQL] [InsertIgnore] [requestID]=%d success", ctx.Value("requestID"))
	return res, nil
}

func (dao *client) InsertReplace(ctx context.Context, tableName string, data []map[string]interface{}) (sql.Result, error) {
	builder := NewInsertBuilder(tableName, data, insertReplace)
	res, err := ExecWithBuilder(ctx, dao, builder)
	if err != nil {
		logit.Logger.Warn("[MySQL] [InsertReplace] [requestID]=%d, [err]=%s", ctx.Value("requestID"), err.Error())
		return nil, err
	}
	logit.Logger.Info("[MySQL] [InsertReplace] [requestID]=%d success", ctx.Value("requestID"))
	return res, nil
}

func (dao *client) InsertOnDuplicate(ctx context.Context, tableName string, data []map[string]interface{}, update map[string]interface{}) (sql.Result, error) {
	builder := NewInsertBuilder(tableName, data, insertOnDuplicate, update)
	res, err := ExecWithBuilder(ctx, dao, builder)
	if err != nil {
		logit.Logger.Warn("[MySQL] [InsertOnDuplicate] [requestID]=%d, [err]=%s", ctx.Value("requestID"), err.Error())
		return nil, err
	}
	logit.Logger.Info("[MySQL] [InsertOnDuplicate] [requestID]=%d success", ctx.Value("requestID"))
	return res, nil
}

func (dao *client) Update(ctx context.Context, tableName string, where map[string]interface{}, update map[string]interface{}) (sql.Result, error) {
	builder := NewUpdateBuilder(tableName, where, update)
	res, err := ExecWithBuilder(ctx, dao, builder)
	if err != nil {
		logit.Logger.Warn("[MySQL] [Update] [requestID]=%d, [err]=%s", ctx.Value("requestID"), err.Error())
		return nil, err
	}
	logit.Logger.Info("[MySQL] [Update] [requestID]=%d success", ctx.Value("requestID"))
	return res, nil
}

func (dao *client) Delete(ctx context.Context, tableName string, where map[string]interface{}) (sql.Result, error) {
	builder := NewDeleteBuilder(tableName, where)
	res, err := ExecWithBuilder(ctx, dao, builder)
	if err != nil {
		logit.Logger.Warn("[MySQL] [Delete] [requestID]=%d, [err]=%s", ctx.Value("requestID"), err.Error())
		return nil, err
	}
	logit.Logger.Info("[MySQL] [Delete] success")
	return res, nil
}

func (dao *client) ExecRaw(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	builder := NewRawBuilder(sql, args)
	res, err := ExecWithBuilder(ctx, dao, builder)
	if err != nil {
		logit.Logger.Warn("[MySQL] [ExecRaw] [requestID]=%d, [err]=%s", ctx.Value("requestID"), err.Error())
		return nil, err
	}
	logit.Logger.Info("[MySQL] [ExecRaw] success")
	return res, nil
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
