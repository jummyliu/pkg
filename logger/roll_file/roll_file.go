// file 使用文件方式记录日志
package roll_file

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jummyliu/pkg/datetime"
	"github.com/jummyliu/pkg/file"
	"github.com/jummyliu/pkg/logger"
)

type FileLogger struct {
	Logger       *log.Logger
	MaxLevel     logger.Level
	defaultLFlag int
	fileName     string
	lastTime     string
	out          *os.File

	ctx   context.Context
	chLog chan logmeta
	chOut chan io.Writer
}

type logmeta struct {
	level  logger.Level
	format string
	params []any
}

const defaultLFlag = log.Ldate | log.Ltime | log.Lmicroseconds

const allLFlag = log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.Lshortfile

func New(ctx context.Context, fileName string, maxLevel logger.Level) logger.Logger {
	out, _ := file.GetFile(fileName)
	l := &FileLogger{
		Logger:       log.New(out, "", defaultLFlag),
		MaxLevel:     maxLevel,
		defaultLFlag: defaultLFlag,
		fileName:     fileName,
		lastTime:     datetime.FormatDateWithLayout(time.Now(), "2006-01-02"),
		out:          out,

		ctx:   ctx,
		chLog: make(chan logmeta),
		chOut: make(chan io.Writer),
	}
	go l.run()
	return l
}

func NewWithoutFlag(ctx context.Context, fileName string, maxLevel logger.Level) logger.Logger {
	out, _ := file.GetFile(fileName)
	l := &FileLogger{
		Logger:       log.New(out, "", 0),
		MaxLevel:     maxLevel,
		defaultLFlag: defaultLFlag,
		fileName:     fileName,
		lastTime:     datetime.FormatDateWithLayout(time.Now(), "2006-01-02"),
		out:          out,

		ctx:   ctx,
		chLog: make(chan logmeta),
		chOut: make(chan io.Writer),
	}
	go l.run()
	return l
}

func (l *FileLogger) SetOutput(out io.Writer) {
	select {
	case l.chOut <- out:
	case <-l.ctx.Done():
		return
	}
}

func (l *FileLogger) Log(level logger.Level, format string, params ...any) {
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

func (l *FileLogger) LogEmerg(format string, params ...any) {
	l.Log(logger.LevelEmerg, format, params...)
}
func (l *FileLogger) LogAlter(format string, params ...any) {
	l.Log(logger.LevelAlter, format, params...)
}
func (l *FileLogger) LogCrit(format string, params ...any) {
	l.Log(logger.LevelCrit, format, params...)
}
func (l *FileLogger) LogError(format string, params ...any) {
	l.Log(logger.LevelError, format, params...)
}
func (l *FileLogger) LogWarning(format string, params ...any) {
	l.Log(logger.LevelWarning, format, params...)
}
func (l *FileLogger) LogNotice(format string, params ...any) {
	l.Log(logger.LevelNotice, format, params...)
}
func (l *FileLogger) LogInfo(format string, params ...any) {
	l.Log(logger.LevelInfo, format, params...)
}
func (l *FileLogger) LogDebug(format string, params ...any) {
	l.Log(logger.LevelDebug, format, params...)
}
func (l *FileLogger) CanLog(level logger.Level) bool {
	return level <= l.MaxLevel
}

func (l *FileLogger) run() {
	defer func() {
		// 遇到错误，重启
		if err := recover(); err != nil {
			// 记录错误
			go l.run()
			return
		}
		// 没有错误，则正常退出关闭chan
		close(l.chLog)
		close(l.chOut)
		l.out.Close()
	}()
	// 每 10 分钟判断一次是否创建文件
	d := time.Minute * 10
	timer := time.NewTimer(d)
	defer func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}()
	for {
		select {
		case meta := <-l.chLog:
			l.log(meta.level, meta.format, meta.params...)
		case out := <-l.chOut:
			l.Logger.SetOutput(out)
		case <-timer.C:
			timer.Reset(d)
			// 日期没有变，保留
			lastTime := datetime.FormatDateWithLayout(time.Now(), "2006-01-02")
			if lastTime == l.lastTime {
				continue
			}
			// 1. 重命名
			l.out.Close()
			os.Rename(l.fileName, fmt.Sprintf("%s.%s", l.fileName, l.lastTime))
			// 2. 创建新文件
			out, _ := file.GetFile(l.fileName)
			// 3. 设置新输出
			l.Logger.SetOutput(out)
			l.lastTime = lastTime
			l.out = out
		case <-l.ctx.Done():
			return
		}
	}
}

func (l *FileLogger) log(level logger.Level, format string, params ...any) {
	if l.defaultLFlag == 0 {
		l.Logger.Printf(format, params...)
	} else {
		if level == logger.LevelDebug {
			l.Logger.SetFlags(l.Logger.Flags() | log.Lshortfile)
		} else {
			l.Logger.SetFlags(l.Logger.Flags() & (allLFlag ^ log.Lshortfile))
		}
		l.Logger.Printf(fmt.Sprintf("| [%s]\t| ", logger.LogNameMap[level])+format, params...)
	}
}
