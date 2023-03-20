package mysqlbuilder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Select 查询数据
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

// SelectOne 查询单条数据，不是单条数据会报错
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

// SelectMany 查询总数并返回指定数据
// 	返回的字段名，不能有预处理的参数
func (db *DBConnect) SelectMany(ctx context.Context, dest any, query string, args ...any) (count int64, err error) {
	countStruct := Count{}
	countArgs := args[:]
	countSql := RegCount.ReplaceAllString(query, "${1} COUNT(1) count ${2}")
	hasLimit := RegLimit.FindAllString(countSql, -1)
	if len(hasLimit) != 0 && len(hasLimit[0]) != 0 {
		l := strings.Count(hasLimit[0], "?")
		countArgs = countArgs[0 : len(countArgs)-l]
		countSql = RegLimit.ReplaceAllString(countSql, "")
	}
	err = db.SelectOne(ctx, &countStruct, countSql, countArgs...)
	fmt.Println(countSql, countArgs)
	if err != nil {
		return 0, fmt.Errorf("get data count failure: %s", err)
	}
	_, err = db.Select(ctx, dest, query, args...)
	if err != nil {
		return 0, fmt.Errorf("get data failure: %s", err)
	}
	return countStruct.Count, nil
}

var (
	RegCount = regexp.MustCompile("(?i)^(SELECT).*?(FROM)")
	RegLimit = regexp.MustCompile(`(?i)LIMIT\s+(\d+|\?)(?:\s*,\s*(\d+|\?))*\s*$`)
)
