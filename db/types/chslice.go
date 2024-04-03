package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// CHSlice 实现任意切片类型，用于clickhouse；与 Slice 不一样的是，修改了 Value() 方法
type CHSlice[T any] []T

// Scan implements the Scanner interface.
func (s *CHSlice[T]) Scan(src any) error {
	if src == nil {
		return nil
	}
	if src, ok := src.([]byte); ok {
		if err := json.Unmarshal(src, s); err == nil {
			return nil
		}
	}
	return errors.New("try scan to CHSlice[T] failure")
}

// Value implements the driver Valuer interface.
func (t CHSlice[T]) Value() (driver.Value, error) {
	if t == nil {
		return []any{}, nil
	}
	arr := []any{}
	for _, item := range t {
		if item, ok := any(item).(driver.Valuer); ok {
			val, err := item.Value()
			if err != nil {
				return nil, err
			}
			arr = append(arr, val)
			continue
		}

		data, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		arr = append(arr, string(data))
	}
	return arr, nil
	// return json.Marshal(t)
}

// String implements flag.Value::String interface.
func (t CHSlice[T]) String() string {
	val, _ := json.Marshal(t)
	return string(val)
}

// Set implements flag.Value::Set interface.
func (t *CHSlice[T]) Set(val string) error {
	return json.Unmarshal([]byte(val), t)
}
