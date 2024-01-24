package log

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	FilName = "tiangong.log"
)

var (
	level string
	path  string

	// system log
	logger     *loggerImpl
	getLogFile = func(path string) (string, error) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return "", nil
		}
		return filepath.Join(path, FilName), nil
	}
)

func init() {
	flag.StringVar(&level, "log.level", "INFO", "log level, options: DEBUG, INFO, WARN, ERROR")
	flag.StringVar(&path, "log.path", "", "log file storage path")
}

func InitLog() {
	write := getLogWriter(path)
	logger = &loggerImpl{
		impl:  log.New(write, "", log.Ldate|log.Ltime),
		Level: LevelValueOf(level),
	}
}

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

func Error(message string, err error, v ...any) {
	logger.error(message+" %+v", append(v, err)...)
}

func format(message string, level Level) string {
	return strings.Join([]string{"[", level.String(), "] ", message}, "")
}

func formatStackTrace(message string) string {
	// TODO
	return format(message, Level_Error)
}

func getLogWriter(path string) io.Writer {
	if path == "" {
		return os.Stdout
	}

	logfile, _ := getLogFile(path)
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	return file
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
