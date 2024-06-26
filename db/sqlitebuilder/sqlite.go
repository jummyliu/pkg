package sqlitebuilder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"sync"

	// mysql driver
	"github.com/jummyliu/pkg/utils"
	_ "modernc.org/sqlite"
)

type contextKey string

type DBConnect struct {
	*sql.DB
	Options   *Options
	ContextTx contextKey

	cacheMap sync.Map
}

// New return a new mysql client, and try ping.
func New(opts ...Option) (*DBConnect, error) {
	options := initOptions(opts...)
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
