package conf

import (
	"os"
	"path/filepath"
	"tiangong/common"
	"tiangong/common/errors"
	"tiangong/common/io"
	"tiangong/common/log"

	"github.com/google/uuid"
	"github.com/magiconair/properties"
)

type Config struct {
	Host     string
	SrvPort  int
	HttpPort int
	UserName string
	Passwd   string
}

var getExecPathFunc = func() string {
	exec, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(exec)
}

func LoadConfig(input string) (*Config, error) {
	if input == "" {
		cur := getExecPathFunc()
		input = filepath.Join(cur, "tiangong.conf")
	}
	log.Debug("find conf file path: %s", input)

	if !common.FileExist(input) {
		return nil, errors.NewError("useage: -conf {path} to specify the configuration file", nil)
	}

	bytes, err := io.ReadFile(input)
	if err != nil {
		return nil, err
	}

	properties, err := properties.Load(bytes, properties.UTF8)
	if err != nil {
		return nil, err
	}
	log.Debug("load config: %+v \n", properties.String())
	config := Config{
		SrvPort:  properties.GetInt(SrvPortDef.First, SrvPortDef.Second),
		HttpPort: properties.GetInt(HttpPortDef.First, HttpPortDef.Second),
		UserName: properties.GetString(UserNameDef.First, UserNameDef.Second),
		Passwd:   properties.GetString(PasswdDef.First, PasswdDef.Second),
	}
	if config.Passwd == "" {
		passwd := uuid.New().String()
		log.Warn("admin.httpPasswd is not set, Generate a random password: %s", passwd)
		config.Passwd = passwd
	}
	return &config, nil
}
