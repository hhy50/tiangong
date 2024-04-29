package common

type Runnable interface {
	Start() error
}

type FuncRunable func() error

func (fun FuncRunable) Start() error {
	return fun()
}
