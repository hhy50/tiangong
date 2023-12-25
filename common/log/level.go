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
