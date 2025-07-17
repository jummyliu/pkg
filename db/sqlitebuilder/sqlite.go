package sqlitebuilder

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	// mysql driver
	"github.com/google/uuid"
	"github.com/jummyliu/pkg/utils"
	"modernc.org/sqlite"
)

type contextKey string

type DBConnect struct {
	*sql.DB
	Options   *Options
	ContextTx contextKey

	cacheMap sync.Map
}

func registerFunc() {
	sqlite.RegisterDeterministicScalarFunction("uuid", 0, func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
		r, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		return r.String(), nil
	})
	sqlite.RegisterDeterministicScalarFunction("sha1", 1, func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
		var b []byte
		switch v := args[0].(type) {
		case []byte:
			b = v
		case string:
			b = []byte(v)
		default:
			return nil, fmt.Errorf("invalid type: %T", v)
		}

		h := sha1.New() //nolint:gosec
		h.Write(b)
		return hex.EncodeToString(h.Sum(nil)), nil
	})
	sqlite.RegisterDeterministicScalarFunction("from_base64", 1, func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
		switch argTyped := args[0].(type) {
		case string:
			return base64.StdEncoding.DecodeString(argTyped)
		default:
			return nil, fmt.Errorf("unsupported type: %T", args[0])
		}
	})
	sqlite.RegisterDeterministicScalarFunction("json_contains", 2, func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
		parse := func(arg driver.Value) (j []any, err error) {
			var data []byte
			switch argTyped := arg.(type) {
			case string:
				data = []byte(argTyped)
			case []byte:
				data = argTyped
			default:
				return nil, fmt.Errorf("unsupported type %T", arg)
			}
			err = json.Unmarshal(data, &j)
			return
		}
		if args[0] == nil || args[1] == nil {
			return nil, nil
		}
		j1, err := parse(args[0])
		if err != nil {
			return nil, err
		}
		j2, err := parse(args[1])
		if err != nil {
			return nil, err
		}
		elements := make(map[any]struct{}, len(j1))
		for _, e := range j1 {
			elements[e] = struct{}{}
		}
		for _, e := range j2 {
			if _, ok := elements[e]; !ok {
				return false, nil
			}
		}
		return true, nil
	})
	sqlite.RegisterDeterministicScalarFunction("regexp", 2, func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
		s1 := args[0].(string)
		s2 := args[1].(string)
		matched, err := regexp.MatchString(s1, s2)
		if err != nil {
			return nil, fmt.Errorf("bad regular expression: %q", err)
		}
		return matched, nil
	})
	sqlite.RegisterDeterministicScalarFunction("concat", -1, func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
		m := make([]string, 0, len(args))
		for _, arg := range args {
			switch argTyped := arg.(type) {
			case string:
				m = append(m, argTyped)
			case []byte:
				m = append(m, string(argTyped))
			default:
				return nil, fmt.Errorf("unsupported type: %T", arg)
			}
		}
		return strings.Join(m, ""), nil
	})
}

// New return a new mysql client, and try ping.
func New(opts ...Option) (*DBConnect, error) {
	options := initOptions(opts...)

	// 注册正则
	registerFunc()
	// driver := BuildDBDriver(options)
	db, err := sql.Open("sqlite", options.DBFilePath)
	if err != nil {
		return nil, err
	}
	// try ping
	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(options.PoolSize / 2)
	db.SetMaxOpenConns(options.PoolSize)

	return &DBConnect{
		DB:        db,
		Options:   options,
		ContextTx: genRandomTx(),
		cacheMap:  sync.Map{},
	}, nil
}

func genRandomTx() contextKey {
	return contextKey("CONTEXT_TX_" + utils.RandomStr(8))
}

// Exec exec query by prepare sql, eg: INSERT, UPDATE, DELETE
func (db *DBConnect) Exec(ctx context.Context, query string, args ...any) (lastInsertId, rowsAffected int64, err error) {
	var stmt *sql.Stmt
	v := ctx.Value(db.ContextTx)
	// use transaction
	if v, ok := v.(*sql.Tx); ok {
		stmt, err = v.PrepareContext(ctx, query)
	} else {
		stmt, err = db.PrepareContext(ctx, query)
	}
	if err != nil {
		return 0, 0, fmt.Errorf("prepare sql failure: %s, err: %s", query, err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, 0, fmt.Errorf("exec sql failure: %s, err: %s", query, err)
	}

	lastInsertId, err = result.LastInsertId()
	if err != nil {
		return 0, 0, fmt.Errorf("get last insert id failure: %s", err)
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return lastInsertId, 0, fmt.Errorf("get rows affected failure: %s", err)
	}

	return lastInsertId, rowsAffected, nil
}

// Query query by prepare sql, eg: SELECT
func (db *DBConnect) Query(ctx context.Context, query string, args ...any) (results []map[string]any, count int64, err error) {
	var stmt *sql.Stmt
	v := ctx.Value(db.ContextTx)
	// use transaction
	if v, ok := v.(*sql.Tx); ok {
		stmt, err = v.PrepareContext(ctx, query)
	} else {
		stmt, err = db.PrepareContext(ctx, query)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("prepare sql failure: %s, err: %s", query, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query sql failure: %s, err: %s", query, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, fmt.Errorf("get column info failure: %s", err)
	}
	var dest []any
	for range columns {
		var item sql.NullString
		dest = append(dest, &item)
	}
	for rows.Next() {
		result := map[string]any{}
		err := rows.Scan(dest...)
		if err != nil {
			return nil, 0, fmt.Errorf("get row data failure: %s", err)
		}
		for index, column := range columns {
			val := dest[index].(*sql.NullString)
			if val.Valid {
				result[column] = val.String
			} else {
				result[column] = ""
			}
		}
		results = append(results, result)
		count++
	}
	return results, count, nil
}

func (db *DBConnect) getColumnMap(columns []string, dest any, ptr bool) ([]any, error) {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return nil, errors.New("must pass a pointer, not a value, to destination")
	}
	if v.IsNil() {
		return nil, errors.New("nil pointer passed to destination")
	}
	t := reflect.TypeOf(dest)
	if v = reflect.Indirect(v); t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("excepts a struct dest")
	}

	var (
		index  map[string][]int
		values = make([]any, 0, len(columns))
	)

	switch idx, ok := db.cacheMap.Load(t); {
	case ok:
		index = idx.(map[string][]int)
	default:
		index = structIdx(t)
		db.cacheMap.Store(t, index)
	}
	for _, name := range columns {
		idx, ok := index[name]
		if !ok {
			return nil, fmt.Errorf("missing destination name %q in %T", name, dest)
		}
		switch field := v.FieldByIndex(idx); {
		case ptr:
			values = append(values, field.Addr().Interface())
		default:
			values = append(values, field.Interface())
		}
	}
	return values, nil
}

func structIdx(t reflect.Type) map[string][]int {
	fields := make(map[string][]int)
	for i := 0; i < t.NumField(); i++ {
		var (
			f    = t.Field(i)
			name = f.Name
		)
		if tn := f.Tag.Get("db"); len(tn) != 0 {
			name = tn
		}
		switch {
		case name == "-", len(f.PkgPath) != 0 && !f.Anonymous:
			continue
		}
		switch {
		case f.Anonymous:
			if f.Type.Kind() != reflect.Ptr {
				for k, idx := range structIdx(f.Type) {
					fields[k] = append(f.Index, idx...)
				}
			}
		default:
			fields[name] = f.Index
		}
	}
	return fields
}
