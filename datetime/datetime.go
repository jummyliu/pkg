package datetime

import "time"

const (
	DatetimeLayout = "2006-01-02 15:04:05"
	ZeroDateStr    = "0001-01-01 00:00:00" // zero date
)

// FormatDate 使用 "2006-01-02 15:04:05" 格式化日期
//
//	如果日期为 zero date（0001-01-01 00:00:00），返回 ""
func FormatDate(date time.Time) string {
	return FormatDateWithLayout(date, DatetimeLayout)
}

// FormatDateWithLayout 使用指定的 日期模版 格式化日期
//
//	如果日期为 zero date，返回 ""
func FormatDateWithLayout(date time.Time, layout string) string {
	if date.IsZero() {
		return ""
	}
	return date.Local().Format(layout)
}

// ParseDate 使用 "2006-01-02 15:04:05" 解析日期
//
//	如果解析失败，返回 time.Time{}
func ParseDate(val any) (result time.Time) {
	return ParseDateWithLayout(val, DatetimeLayout)
}

// ParseDateRFC3339Nano 使用 RFC3339Nano 解析日期
//
//	如果解析失败，返回 time.Time{}
func ParseDateRFC3339Nano(val any) (result time.Time) {
	return ParseDateWithLayout(val, time.RFC3339Nano)
}

// ParseDateWithLayout 使用指定的 日期模版 解析日期
//
//	如果解析失败，返回 time.Time{}
func ParseDateWithLayout(val any, layout string) (result time.Time) {
	var err error
	switch v := val.(type) {
	case string:
		result, err = time.ParseInLocation(layout, v, time.Local)
		if err != nil {
			result = time.Time{}
		}
	case time.Time:
		result = v
	default:
		result = time.Time{}
	}
	return
}
