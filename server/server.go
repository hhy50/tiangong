package server

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"tiangong/common"
	"tiangong/common/errors"
	"tiangong/common/io"
	"tiangong/common/log"
	"tiangong/server/conf"

	"github.com/google/uuid"
	"github.com/magiconair/properties"
)

var DefaultConfPath = func() string {
	exec, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(exec)
}

var getRedomPasswd = func() string {
	exec, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(exec)
}

type Config struct {
	Host     string
	TcpPort  int
	HttpPort int
	UserName string
	Passwd   string
}

type Server struct {
	conf *Config
}

func (s *Server) Start() error {
	conf := s.conf
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.TcpPort))
	if err != nil {
		return err
	}
	fmt.Printf("listener host:%s, listener tcp port:%d, listener http port:%d \n", conf.Host, conf.TcpPort, conf.HttpPort)
	go func() {
		for {
			_, err := listener.Accept()
			if err != nil {

			}
		}
	}()
	return nil
}

func (s *Server) Stop() error {
	return nil
}

func loadConfig(cp string) (*Config, error) {
	if cp == "" {
		cur := DefaultConfPath()
		cp = filepath.Join(cur, "tiangong.conf")
	}
	log.Debug("find conf file path: %s", cp)

	if common.FileNotExist(cp) {
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

	config := Config{
		TcpPort:  properties.GetInt(conf.TcpPort.First, conf.TcpPort.Second),
		HttpPort: properties.GetInt(conf.HttpPort.First, conf.HttpPort.Second),
		UserName: strings.Trim(properties.GetString(conf.UserName.First, conf.UserName.Second), "\""),
		Passwd:   strings.Trim(properties.GetString(conf.Passwd.First, conf.Passwd.Second), "\""),
	}
	if config.Passwd == "" {
		log.Warn("httpPasswd is not set, Generate a random password: %s", uuid.New().String())
	}
	return &config, nil
}

func NewServer(cp string) (*Server, error) {
	conf, err := loadConfig(cp)
	if err != nil {
		return nil, err
	}
	server := &Server{
		conf: conf,
	}
	return server, nil
}
