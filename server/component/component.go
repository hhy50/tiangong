package component

type Component interface {
	Start() error
}

type FuncComponent func() error


func (fn FuncComponent) Start() error {
	return fn()
}