package conf

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/log"

	"github.com/magiconair/properties"
)

var (
	PropTag         = "prop"
	DefaultValueTag = "default"
)

var getExecPathFunc = func() string {
	exec, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(exec)
}

type DefaultValueFunc = func(string) string

func EmptyDefaultValueFunc(string) string {
	return ""
}

func LoadConfig(input string, config interface{}, defaultProp DefaultValueFunc) error {
	ptr, ok := common.GetPtr(config)
	if !ok {
		return errors.NewError("param 'config' must be a pointer", nil)
	}
	if common.IsEmpty(input) {
		return errors.NewError("useage: -conf {path} to specify the configuration file", nil)
	}

	if !common.FileExist(input) && !filepath.IsAbs(input) {
		cur := getExecPathFunc()
		log.Debug("find conf file dir: %s", cur)

		input = filepath.Join(cur, input)
	}
	if !common.FileExist(input) {
		return errors.NewError("config file not found!", nil)
	}

	prop, err := properties.LoadFile(input, properties.UTF8)
	if err != nil {
		return err
	}
	log.Debug("load config:\n%+v", prop.String())

	defaultValueMap := common.GetTags(DefaultValueTag, config)
	val := ptr.Elem()
	for fName, tVal := range common.GetTags(PropTag, config) {
		value, ok := prop.Get(tVal)
		if !ok || common.IsEmpty(value) {
			value = defaultProp(tVal)
		}
		if common.IsEmpty(value) {
			if v, f := defaultValueMap[fName]; f {
				value = v
			}
		}
		if common.IsNotEmpty(value) {
			field := val.FieldByName(fName)
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Uint,
				reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				i, err := strconv.Atoi(value)
				if err != nil {
					return err
				}
				field.SetUint(uint64(i))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				i, err := strconv.Atoi(value)
				if err != nil {
					return err
				}
				field.SetInt(int64(i))
			case reflect.Float32:
			case reflect.Float64:
				i, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return err
				}
				field.SetFloat(i)
			}
		}
	}
	return nil
}
