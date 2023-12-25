package common

import "os"

var (
	DateFormat = "2006-01-02 15:04:05"
	LogFilName = "tiangong.log"
)

func FileNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}
