package common

import "unsafe"

func String(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
