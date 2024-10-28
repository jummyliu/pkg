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

// GetMonthRange 返回月开始（本月第一天零点）和月结束（下月第一天零点）
//
//	firstOfMonth <= date < lastOfMonth
func GetMonthRange(date time.Time) (firstOfMonth, lastOfMonth time.Time) {
	year, month, _ := date.Date()
	firstOfMonth = time.Date(year, month, 1, 0, 0, 0, 0, date.Location())
	lastOfMonth = firstOfMonth.AddDate(0, 1, 0)
	return
}

// GetWeekRange firstDayOfWeek 用于指定每周的第一天，返回周开始（本周第一天零点）和周结束（下一周第一天零点）
//
//	firstOfWeek <= date < lastOfWeek
func GetWeekRange(date time.Time, firstDayOfWeek time.Weekday) (firstOfWeek, lastOfWeek time.Time) {
	year, month, day := date.Date()
	step := firstDayOfWeek - date.Weekday()
	if step > 0 {
		step -= 7
	}
	firstOfWeek = time.Date(year, month, day+int(step), 0, 0, 0, 0, date.Location())
	lastOfWeek = firstOfWeek.AddDate(0, 0, 7)
	return
}
