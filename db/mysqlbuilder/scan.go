package mysqlbuilder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

func (db *DBConnect) Select(ctx context.Context, dest any, query string, args ...any) (count int64, err error) {
	// 1. 校验 dest 是否是有效的 slice
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr {
		return 0, errors.New("must pass a pointer, not a value, to destination")
	}
	if val.IsNil() {
		return 0, errors.New("nil pointer passed to destination")
	}
	direct := reflect.Indirect(val)
	if direct.Kind() != reflect.Slice {
		return 0, errors.New("must pass a slice to destination")
	}
	if direct.Len() != 0 {
		direct.Set(reflect.MakeSlice(direct.Type(), 0, direct.Cap()))
	}
	var (
		base = direct.Type().Elem()
	)

	// 2. 预处理查询
	var stmt *sql.Stmt
	v := ctx.Value(db.ContextTx)
	// use transaction
	if v, ok := v.(*sql.Tx); ok {
		stmt, err = v.PrepareContext(ctx, query)
	} else {
		stmt, err = db.PrepareContext(ctx, query)
	}
	if err != nil {
		return 0, fmt.Errorf("prepare sql failure: %s", query)
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return 0, fmt.Errorf("query sql failure: %s", query)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return 0, fmt.Errorf("get column info failure: %s", err)
	}
	for rows.Next() {
		elem := reflect.New(base)
		data, err := db.getColumnMap(columns, elem.Interface(), true)
		if err != nil {
			return 0, fmt.Errorf("get column map failure: %s", err)
		}
		err = rows.Scan(data...)
		if err != nil {
			return 0, fmt.Errorf("get row data failure")
		}
		direct.Set(reflect.Append(direct, elem.Elem()))
		count++
	}
	return count, nil
}

func (db *DBConnect) SelectOne(ctx context.Context, dest any, query string, args ...any) (err error) {
	// 2. 预处理查询
	var stmt *sql.Stmt
	v := ctx.Value(db.ContextTx)
	// use transaction
	if v, ok := v.(*sql.Tx); ok {
		stmt, err = v.PrepareContext(ctx, query)
	} else {
		stmt, err = db.PrepareContext(ctx, query)
	}
	if err != nil {
		return fmt.Errorf("prepare sql failure: %s", query)
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return fmt.Errorf("query sql failure: %s", query)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("get column info failure: %s", err)
	}
	if rows.Next() {
		data, err := db.getColumnMap(columns, dest, true)
		if err != nil {
			return fmt.Errorf("get column map failure: %s", err)
		}
		err = rows.Scan(data...)
		if err != nil {
			return fmt.Errorf("get row data failure")
		}
		return nil
	}
	return errors.New("get row data failure")
}
