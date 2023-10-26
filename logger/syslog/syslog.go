// syslog 方式记录日志
package syslog

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jummyliu/pkg/logger"

	"github.com/papertrail/remote_syslog2/syslog"
)

type SysLogger struct {
	*syslog.Logger
	// gsyslog.Syslogger
	network  string
	host     string
	MaxLevel logger.Level

	ctx    context.Context
	chLog  chan logmeta
	chConn chan string
}

type logmeta struct {
	level  logger.Level
	format string
	params []any
}

const defaultSysPriority = logger.LevelDebug

func NewSyslogger(ctx context.Context, maxLevel logger.Level, addr string) logger.Logger {
	l := &SysLogger{
		MaxLevel: maxLevel,

		ctx:    ctx,
		chLog:  make(chan logmeta),
		chConn: make(chan string),
	}
	go l.run()
	l.chConn <- addr
	return l
}

func (l *SysLogger) SetAddr(addr string) {
	select {
	case l.chConn <- addr:
	case <-l.ctx.Done():
		return
	}
}

// Log log by level
func (l *SysLogger) Log(level logger.Level, format string, params ...any) {
	// 超过配置的日志级别，不输出
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

// LogEmerg log by emerg
func (l *SysLogger) LogEmerg(format string, params ...any) {
	l.Log(logger.LevelEmerg, format, params...)
}

// LogAlter log by Alter
func (l *SysLogger) LogAlter(format string, params ...any) {
	l.Log(logger.LevelAlter, format, params...)
}

// LogCrit log by Crit
func (l *SysLogger) LogCrit(format string, params ...any) {
	l.Log(logger.LevelCrit, format, params...)
}

// LogError log by error
func (l *SysLogger) LogError(format string, params ...any) {
	l.Log(logger.LevelError, format, params...)
}

// LogWarning log by warning
func (l *SysLogger) LogWarning(format string, params ...any) {
	l.Log(logger.LevelWarning, format, params...)
}

// LogNotice log by notice
func (l *SysLogger) LogNotice(format string, params ...any) {
	l.Log(logger.LevelNotice, format, params...)
}

// LogInfo log by info
func (l *SysLogger) LogInfo(format string, params ...any) {
	l.Log(logger.LevelInfo, format, params...)
}

// LogDebug log by debug
func (l *SysLogger) LogDebug(format string, params ...any) {
	l.Log(logger.LevelDebug, format, params...)
}

// CanLog check level can be logged
func (l *SysLogger) CanLog(level logger.Level) bool {
	return level <= l.MaxLevel
}

func (l *SysLogger) run() {
	defer func() {
		// 遇到错误，重启
		if err := recover(); err != nil {
			// 记录错误
			go l.run()
			return
		}
		// 没有错误，则正常退出关闭chan
		close(l.chLog)
		close(l.chConn)
	}()
	for {
		select {
		case meta := <-l.chLog:
			l.log(meta.level, meta.format, meta.params...)
		case addr := <-l.chConn:
			if l.Logger != nil {
				l.Logger.Close()
			}
			u, err := url.Parse(addr)
			if err != nil {
				continue
			}
			l.network = u.Scheme
			l.host = u.Host
			logger, err := syslog.Dial("", u.Scheme, u.Host, nil, time.Second*5, time.Second*5, 999990)
			if err != nil {
				continue
			}
			l.Logger = logger
			// if l.Syslogger != nil {
			// 	l.Syslogger.Close()
			// }
			// u, err := url.Parse(addr)
			// if err != nil {
			// 	continue
			// }
			// l.network = u.Scheme
			// l.host = u.Host
			// logger, err := gsyslog.DialLogger(u.Scheme, u.Host, gsyslog.Priority(defaultSysPriority), "SYSLOG", "")
			// if err != nil {
			// 	continue
			// }
			// l.Syslogger = logger
		case <-l.ctx.Done():
			return
		}
	}
}

func (l *SysLogger) log(level logger.Level, format string, params ...any) {
	if l.Logger == nil {
		return
	}
	l.Write(syslog.Packet{
		Severity: syslog.Priority(level),
		Facility: syslog.LogLocal1,
		Time:     time.Now(),
		Hostname: "syslogclient",
		Tag:      "",
		Message:  fmt.Sprintf(format, params...),
	})
	// if l.Syslogger == nil {
	// 	return
	// }
	// l.WriteLevel(gsyslog.Priority(level), []byte(fmt.Sprintf(format, params...)))
}
