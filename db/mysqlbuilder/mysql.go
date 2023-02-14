package mysqlbuilder

import (
	"context"
	"database/sql"
	"fmt"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jummyliu/pkg/utils"
)

type contextKey string

type DBConnect struct {
	*sql.DB
	Options   *Options
	ContextTx contextKey
}

// New return a new mysql client, and try ping.
func New(opts ...Option) (*DBConnect, error) {
	options := initOptions(opts...)
	driver := BuildDBDriver(options)
	db, err := sql.Open("mysql", driver)
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
