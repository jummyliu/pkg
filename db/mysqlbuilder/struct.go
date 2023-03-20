package mysqlbuilder

import (
	"encoding/json"
	"errors"
)

type Map map[string]any

// Scan 实现 Scan 接口，用来自动赋值
func (m *Map) Scan(src any) error {
	if src == nil {
		return nil
	}
	if src, ok := src.([]byte); ok {
		if err := json.Unmarshal(src, m); err != nil {
			return nil
		}
	}
	return errors.New("try scan to Map failure")
}

type Count struct {
	Count int64 `db:"count"`
}
