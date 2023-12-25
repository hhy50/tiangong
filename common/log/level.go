package log

type Level uint8

const (
	Level_Debug Level = iota
	Level_Info
	Level_Warn
	Level_Error
)

func (l Level) String() string {
	switch l {
	case Level_Debug:
		return "DEBUG"
	case Level_Info:
		return "INFO"
	case Level_Warn:
		return "WARN"
	case Level_Error:
		return "ERROR"
	}
	return ""
}

func LevelValueOf(level string) Level {
	switch level {
	case "DEBUG":
		return Level_Debug
	case "INFO":
		return Level_Info
	case "WARN":
		return Level_Warn
	case "ERROR":
		return Level_Error
	}
	return Level_Info
}
