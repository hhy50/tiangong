package common

import "context"

var (
	EmptyCtx = context.Background()

	CancelFuncKey = "cancelFunc"
	ServerKey     = "Server"

	ClientKey    = "Client"
	TcpClientKey = "TcpClient"

	ProcessKey = "Processor"
)

func SetProcess(ctx context.Context, p interface{}) context.Context {
	return context.WithValue(ctx, ProcessKey, p)
}

func GetProcess(ctx context.Context) interface{} {
	return ctx.Value(ProcessKey)
}
