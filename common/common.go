package common

import "os"

var (
	DateFormat = "2006-01-02 15:04:05"
	LogFilName = "tiangong.log"
)

func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
