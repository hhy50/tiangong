package log

import (
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	level Level

	logger = &loggerImpl{
		impl:  log.New(io.MultiWriter(os.Stdout), "", log.Ldate|log.Ltime),
		Level: level,
	}
)

type Logger interface {
	debug(string, ...any)
	info(string, ...any)
	warn(string, ...any)
	error(string, ...any)
}

type loggerImpl struct {
	impl  *log.Logger
	Level Level
}

func (l *loggerImpl) debug(message string, args ...any) {
	if Level_Debug >= l.Level {
		impl := l.impl
		impl.Printf(format(message, Level_Debug), args...)
	}
}

func (l *loggerImpl) info(message string, args ...any) {
	if Level_Info >= l.Level {
		impl := l.impl
		impl.Printf(format(message, Level_Info), args...)
	}
}

func (l *loggerImpl) warn(message string, args ...any) {
	if Level_Warn >= l.Level {
		impl := l.impl
		impl.Printf(format(message, Level_Warn), args...)
	}
}

func (l *loggerImpl) error(message string, args ...any) {
	if Level_Error >= l.Level {
		impl := l.impl
		impl.Printf(formatStackTrace(message), args...)
	}
}

func Debug(message string, v ...any) {
	logger.debug(message, v...)
}

func Info(message string, v ...any) {
	logger.info(message, v...)
}

func Warn(message string, v ...any) {
	logger.warn(message, v...)
}

func Error(message string, v ...any) {
	logger.error(message, v...)
}

func format(message string, level Level) string {
	return strings.Join([]string{"[", level.String(), "] ", message}, "")
}

func formatStackTrace(messgae string) string {
	return format(messgae, Level_Error)
}

func getCaller(depth int) (string, int, bool) {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		split := strings.Split(file, "/")
		last := split[len(split)-1]
		return last, line, true
	}
	return "", -1, false
}
