package conf

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/haiyanghan/tiangong/common/io"
	"github.com/haiyanghan/tiangong/common/log"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/pelletier/go-toml"
)

var (
	tomlTree *toml.Tree

	ConfigPath      string
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

func init() {
	flag.StringVar(&ConfigPath, "conf", "", "-conf {path}")
}

func Load() {
	if !common.FileExist(ConfigPath) && !filepath.IsAbs(ConfigPath) {
		cur := getExecPathFunc()
		ConfigPath = filepath.Join(cur, ConfigPath)
	}
	if !common.FileExist(ConfigPath) {
		panic(fmt.Sprintf("Conf file not fount ConfigPath: %s", ConfigPath))
	}
	bytes, err := io.ReadFile(ConfigPath)
	if err != nil {
		panic(err)
	}
	tomlTree, err = toml.Load(string(bytes))
	if err != nil {
		panic(err)
	}
	log.Debug("Load config: %+v", tomlTree.String())
}

func LoadConfig(kind string, config interface{}) error {
	if kind == "" {
		return parse(config, ToFlatMap(tomlTree))
	} else if kindVal, ok := tomlTree.Get(kind).(*toml.Tree); ok {
		return parse(config, kindVal.ToMap())
	}
	return nil
}

func LoadToMap(kind string) map[string]interface{} {
	if kind == "" {
		return tomlTree.ToMap()
	} else if kindVal, ok := tomlTree.Get(kind).(*toml.Tree); ok {
		return kindVal.ToMap()
	}
	return nil
}

func GetOrDefault(key, d string) interface{} {
	if key == "" {
		return d
	}
	return tomlTree.GetDefault(key, d)
}

func parse(config interface{}, keyVal map[string]interface{}) error {
	ptr, ok := common.GetPtr(config)
	if !ok {
		return errors.NewError("Param 'config' must be a pointer", nil)
	}
	defaultValueMap := common.GetTags(DefaultValueTag, config)
	val := ptr.Elem()
	for fName, tVal := range common.GetTags(PropTag, config) {
		value, ok := keyVal[tVal]
		if !ok {
			value = keyVal[fName]
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
				field.SetString(value.(string))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if str, ok := value.(string); ok {
					i, err := strconv.Atoi(str)
					if err != nil {
						return err
					}
					field.SetInt(int64(i))
				} else {
					field.SetInt(value.(int64))
				}
			case reflect.Float32, reflect.Float64:
				i, err := strconv.ParseFloat(value.(string), 64)
				if err != nil {
					return err
				}
				field.SetFloat(i)
			}
		}
	}
	return nil
}

func ToFlatMap(tree *toml.Tree) map[string]interface{} {
	mp := map[string]interface{}{}
	for key, val := range tomlTree.ToMap() {
		if submap, ok := val.(map[string]interface{}); ok {
			for skey, sval := range submap {
				mp[key+"."+skey] = sval
			}
		} else {
			mp[key] = val
		}
	}
	return mp
}
