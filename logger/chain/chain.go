// chain 使用日志链记录日志
package chain

import "github.com/jummyliu/pkg/logger"

// ChainLogger 组合多个 logger.Logger
type ChainLogger struct {
	Loggers []logger.Logger
}

func New(l ...logger.Logger) logger.Logger {
	return &ChainLogger{
		Loggers: l,
	}
}

func (c *ChainLogger) Log(level logger.Level, format string, params ...any) {
	for _, l := range c.Loggers {
		l.Log(level, format, params...)
	}
}

func (c *ChainLogger) LogEmerg(format string, params ...any) {
	for _, l := range c.Loggers {
		l.LogEmerg(format, params...)
	}
}
func (c *ChainLogger) LogAlter(format string, params ...any) {
	for _, l := range c.Loggers {
		l.LogAlter(format, params...)
	}
}
func (c *ChainLogger) LogCrit(format string, params ...any) {
	for _, l := range c.Loggers {
		l.LogCrit(format, params...)
	}
}
func (c *ChainLogger) LogError(format string, params ...any) {
	for _, l := range c.Loggers {
		l.LogError(format, params...)
	}
}
func (c *ChainLogger) LogWarning(format string, params ...any) {
	for _, l := range c.Loggers {
		l.LogWarning(format, params...)
	}
}
func (c *ChainLogger) LogNotice(format string, params ...any) {
	for _, l := range c.Loggers {
		l.LogNotice(format, params...)
	}
}
func (c *ChainLogger) LogInfo(format string, params ...any) {
	for _, l := range c.Loggers {
		l.LogInfo(format, params...)
	}
}
func (c *ChainLogger) LogDebug(format string, params ...any) {
	for _, l := range c.Loggers {
		l.LogDebug(format, params...)
	}
}
func (c *ChainLogger) CanLog(level logger.Level) bool {
	return true
}
