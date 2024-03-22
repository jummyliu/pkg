package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/elastic/elastic-agent-libs/mapstr"
)

// Map 实现任意Map
type Map[T comparable, Q any] map[T]Q

// Scan implements the Scanner interface.
func (m *Map[T, Q]) Scan(src any) error {
	if src == nil {
		return nil
	}
	if src, ok := src.([]byte); ok {
		if err := json.Unmarshal(src, m); err == nil {
			return nil
		}
	}
	return errors.New("try scan to Map[T, Q] failure")
}

// Value implements the driver Valuer interface.
func (m Map[T, Q]) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// String implements flag.Value::String interface.
func (m Map[T, Q]) String() string {
	val, _ := json.Marshal(m)
	return string(val)
}

// Set implements flag.Value::Set interface.
func (m *Map[T, Q]) Set(val string) error {
	return json.Unmarshal([]byte(val), m)
}

func MustGetValue[T any](m map[string]any, key string) (result T) {
	v, err := mapstr.M(m).GetValue(key)
	if v, ok := v.(T); err == nil && ok {
		return v
	}
	return result
}
