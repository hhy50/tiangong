package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"tiangong/common/errors"
	"tiangong/common/file"
)

type Config struct {
	Host     string
	TcpPort  int
	HttpPort int
	Passwd   string
}

type Server struct {
	conf *Config
}

func (s *Server) Start() error {
	conf := s.conf
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", conf.Host, conf.TcpPort))
	if err != nil {
		return err
	}
	fmt.Printf("listener host:%s, listener tcp port:%s, listener http port:%s", conf.Host, conf.TcpPort, conf.HttpPort)
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

func NewServer(cp string) (*Server, error) {
	if cp == "" {
		cur, err := os.Executable()
		if err != nil {
			return nil, errors.NewError("use -conf {path} to specify the configuration file", err)
		}
		cp = filepath.Join(cur, "tiangong.config.json")
	}

	bytes, err := file.ReadAll(cp)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(bytes, config); err != nil {
		return nil, err
	}

	server := &Server{
		conf: config,
	}
	return server, nil
}
