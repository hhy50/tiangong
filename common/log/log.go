package log

type Level uint8

const (
	Debug int = iota
	Info
	Error
)

type Message interface {
	String() string
}

type Logger interface {
	SetPattern(string)
	Debug(Message) error
	Info(Message) error
	Error(Message) error
}

type loggerImpl struct {
	
}

func (l *loggerImpl) Debug(message Message) error {
	return nil
}

func (l *loggerImpl) Info(message Message) error {
	return nil
}

func (l *loggerImpl) Error(message Message) error {
	return nil
}
