package errors

type Error struct {
	msg    string
	parent error
}

func (e *Error) Error() string {
	return e.msg
}

func NewError(m string, p error) *Error {
	return &Error{
		msg:    m,
		parent: p,
	}
}
