package common

import (
	"os"
	"reflect"
)

const (
	Zero int = iota
	One
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
)

var (
	DateFormat = "2006-01-02 15:04:05"
)

func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func IsEmpty(obj interface{}) bool {
	return !IsNotEmpty(obj)
}

func IsNotEmpty(obj interface{}) bool {
	s := reflect.Indirect(reflect.ValueOf(obj))
	switch s.Kind() {
	case reflect.Slice:
	case reflect.Array:
	case reflect.Map:
		return s.Len() > 0
	case reflect.String:
		return len(obj.(string)) > 0
	}
	return obj != nil
}

func Min(i1 int, i2 int) int {
	if i1 > i2 {
		return i2
	}
	return i1
}
