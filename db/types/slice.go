package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Slice 实现任意切片类型
type Slice[T any] []T

// Scan implements the Scanner interface.
func (s *Slice[T]) Scan(src any) error {
	if src == nil {
		return nil
	}
	if src, ok := src.([]byte); ok {
		if err := json.Unmarshal(src, s); err == nil {
			return nil
		}
	}
	return errors.New("try scan to Slice[T] failure")
}

// Value implements the driver Valuer interface.
func (t Slice[T]) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return json.Marshal(t)
}
