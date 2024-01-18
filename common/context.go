package common

import "context"

var (
	EmptyCtx = context.Background()

	CancelFuncKey = "cancelFunc"
	ServerKey     = "server"

	ClientKey = "Client"
)
