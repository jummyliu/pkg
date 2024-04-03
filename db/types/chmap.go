package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// CHMap 实现任意Map，用于clickhouse；与 Map 不一样的是，修改了 Value() 方法
type CHMap[T comparable, Q any] map[T]Q

// Scan implements the Scanner interface.
func (m *CHMap[T, Q]) Scan(src any) error {
	if src == nil {
		return nil
	}
	if src, ok := src.([]byte); ok {
		if err := json.Unmarshal(src, m); err == nil {
			return nil
		}
	}
	return errors.New("try scan to CHMap[T, Q] failure")
}

// Value implements the driver Valuer interface.
func (m CHMap[T, Q]) Value() (driver.Value, error) {
	if m == nil {
		return CHMap[T, Q]{}, nil
	}
	val, err := json.Marshal(m)
	return string(val), err
}

// String implements flag.Value::String interface.
func (m CHMap[T, Q]) String() string {
	val, _ := json.Marshal(m)
	return string(val)
}

// Set implements flag.Value::Set interface.
func (m *CHMap[T, Q]) Set(val string) error {
	return json.Unmarshal([]byte(val), m)
}
