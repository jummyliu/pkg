// fn 使用函数方式记录日志，函数由调用者提供
package fn

import (
	"context"

	"github.com/jummyliu/pkg/logger"
)

type FnLogger struct {
	LogFunc    func(level logger.Level, format string, params ...any)
	CanLogFunc func(level logger.Level) bool

	ctx   context.Context
	chLog chan logmeta
}

type logmeta struct {
	level  logger.Level
	format string
	params []any
}

func New(
	ctx context.Context,
	logFunc func(logger.Level, string, ...any),
	canLogFunc func(logger.Level) bool,
) logger.Logger {
	l := &FnLogger{
		LogFunc:    logFunc,
		CanLogFunc: canLogFunc,

		ctx:   ctx,
		chLog: make(chan logmeta),
	}
	go l.run()
	return l
}

func (l *FnLogger) Log(level logger.Level, format string, params ...any) {
	if !l.CanLog(level) {
		return
	}
	select {
	case l.chLog <- logmeta{
		level:  level,
		format: format,
		params: params,
	}:
	case <-l.ctx.Done():
		return
	}
}

func (l *FnLogger) LogEmerg(format string, params ...any) {
	l.Log(logger.LevelEmerg, format, params...)
}
func (l *FnLogger) LogAlter(format string, params ...any) {
	l.Log(logger.LevelAlter, format, params...)
}
func (l *FnLogger) LogCrit(format string, params ...any) {
	l.Log(logger.LevelCrit, format, params...)
}
func (l *FnLogger) LogError(format string, params ...any) {
	l.Log(logger.LevelError, format, params...)
}
func (l *FnLogger) LogWarning(format string, params ...any) {
	l.Log(logger.LevelWarning, format, params...)
}
func (l *FnLogger) LogNotice(format string, params ...any) {
	l.Log(logger.LevelNotice, format, params...)
}
func (l *FnLogger) LogInfo(format string, params ...any) {
	l.Log(logger.LevelInfo, format, params...)
}
func (l *FnLogger) LogDebug(format string, params ...any) {
	l.Log(logger.LevelDebug, format, params...)
}
func (l *FnLogger) CanLog(level logger.Level) bool {
	return l.CanLogFunc(level)
}

func (l *FnLogger) run() {
	defer func() {
		// 遇到错误，重启
		if err := recover(); err != nil {
			// 记录错误
			go l.run()
			return
		}
		// 没有错误，则正常退出关闭chan
		close(l.chLog)
	}()
	for {
		select {
		case meta := <-l.chLog:
			l.log(meta.level, meta.format, meta.params...)
		case <-l.ctx.Done():
			return
		}
	}
}

func (l *FnLogger) log(level logger.Level, format string, params ...any) {
	l.LogFunc(level, format, params...)
}
