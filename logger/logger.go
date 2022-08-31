package logger

// Logger 接口
type Logger interface {
	Log(Level, string, ...any)
	LogEmerg(string, ...any)
	LogAlter(string, ...any)
	LogCrit(string, ...any)
	LogError(string, ...any)
	LogWarning(string, ...any)
	LogNotice(string, ...any)
	LogInfo(string, ...any)
	LogDebug(string, ...any)
	CanLog(Level) bool
}

// Level 日志级别
type Level int64

// 8种日志级别
const (
	LevelEmerg Level = iota
	LevelAlter
	LevelCrit
	LevelError
	LevelWarning
	LevelNotice
	LevelInfo
	LevelDebug

	LevelCount // count
)

// LogNameMap 日志级别映射
var LogNameMap = map[Level]string{
	LevelEmerg:   "emerg",
	LevelAlter:   "alert",
	LevelCrit:    "crit",
	LevelError:   "error",
	LevelWarning: "warning",
	LevelNotice:  "notice",
	LevelInfo:    "info",
	LevelDebug:   "debug",
}

const defaultLevel = LevelInfo

// GetLevelByName 过日志级别名称获取日志级别
var GetLevelByName func(string) Level

func init() {
	GetLevelByName = LevelGetter()
}

func LevelGetter() func(string) Level {
	m := map[string]Level{}
	for k, v := range LogNameMap {
		m[v] = k
	}
	return func(levelName string) Level {
		if level, ok := m[levelName]; ok {
			return level
		}
		return defaultLevel
	}
}
