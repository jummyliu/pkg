package clickhousebuilder

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Select 查询数据
func (db *DBConnect) Select(ctx context.Context, dest any, query string, args ...any) (count int64, err error) {
	err = db.Conn.Select(ctx, dest, query, args...)
	if err != nil {
		return 0, err
	}
	// db.Conn.Select 本身有类型判断，这里不再进行判断
	val := reflect.ValueOf(dest)
	direct := reflect.Indirect(val)
	return int64(direct.Len()), nil
}

// SelectOne 查询单条数据，不是单条数据会报错
func (db *DBConnect) SelectOne(ctx context.Context, dest any, query string, args ...any) (err error) {
	row := db.Conn.QueryRow(ctx, query, args...)
	if row.Err() != nil {
		return row.Err()
	}
	return row.ScanStruct(dest)
}

// SelectMany 查询总数并返回指定数据
func (db *DBConnect) SelectMany(ctx context.Context, dest any, query string, args ...any) (count int64, err error) {
	countStruct := Count{}
	countArgs := args[:]
	countSql := RegCount.ReplaceAllString(query, "${1} COUNT(1) count ${2}")
	hasLimit := RegOrderLimit.FindAllString(countSql, -1)
	if len(hasLimit) != 0 && len(hasLimit[0]) != 0 {
		l := strings.Count(hasLimit[0], "?")
		countArgs = countArgs[0 : len(countArgs)-l]
		countSql = RegOrderLimit.ReplaceAllString(countSql, "")
	}
	err = db.SelectOne(ctx, &countStruct, countSql, countArgs...)
	if err != nil {
		return 0, fmt.Errorf("get data count failure: %s", err)
	}
	_, err = db.Select(ctx, dest, query, args...)
	if err != nil {
		return 0, fmt.Errorf("get data failure: %s", err)
	}
	return int64(countStruct.Count), nil
}

var (
	RegCount      = regexp.MustCompile("(?is)^(SELECT).*?(FROM)")
	RegLimit      = regexp.MustCompile(`(?is)LIMIT\s+(\d+|\?)(?:\s*,\s*(\d+|\?))*\s*$`)
	RegOrderLimit = regexp.MustCompile(`(?is)(ORDER BY \S+(\s+(ASC|DESC))?\s+)?LIMIT\s+(\d+|\?)(?:\s*,\s*(\d+|\?))*\s*$`)
)

// SelectAll 返回所有数据，如果最后有 limit 会删除
func (db *DBConnect) SelectAll(ctx context.Context, dest any, query string, args ...any) (count int64, err error) {
	hasLimit := RegLimit.FindAllString(query, -1)
	if len(hasLimit) != 0 && len(hasLimit[0]) != 0 {
		l := strings.Count(hasLimit[0], "?")
		args = args[0 : len(args)-l]
		query = RegLimit.ReplaceAllString(query, "")
	}
	count, err = db.Select(ctx, dest, query, args...)
	if err != nil {
		return 0, fmt.Errorf("get data failure: %s", err)
	}
	return
}

type MultiSelect func(ctx context.Context, dest any, query string, args ...any) (count int64, err error)
