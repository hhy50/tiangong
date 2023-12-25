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

type ServerConfig struct {
	Host     string
	TcpPort  int
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

func LoadConfigWithPath(cp string) (*ServerConfig, error) {
	if cp == "" {
		cur := getExecPathFunc()
		cp = filepath.Join(cur, "tiangong.conf")
	}
	log.Debug("find conf file path: %s", cp)

	if !common.FileExist(cp) {
		return nil, errors.NewError("useage: -conf {path} to specify the configuration file", nil)
	}

	bytes, err := io.ReadFile(cp)
	if err != nil {
		return nil, err
	}

	properties, err := properties.Load(bytes, properties.UTF8)
	if err != nil {
		return nil, err
	}
	log.Debug("load config: %+v \n", properties.String())
	config := ServerConfig{
		TcpPort:  properties.GetInt(TcpPort.First, TcpPort.Second),
		HttpPort: properties.GetInt(HttpPort.First, HttpPort.Second),
		UserName: properties.GetString(UserName.First, UserName.Second),
		Passwd:   properties.GetString(UserName.First, UserName.Second),
	}
	if config.Passwd == "" {
		log.Warn("httpPasswd is not set, Generate a random password: %s", uuid.New().String())
	}
	return &config, nil
}
