package server

import (
	"fmt"
	"net"
	"sync"
)

type ServerConfig struct {
	Host string
	Port string
}

type Book struct {
	Name   string
	Pricae float32
}

type Server struct {
	c         *ServerConfig
	OnlineSet map[string]*User

	Lock    sync.RWMutex
	Message chan string

	// ===========
	Book1 *Book
	Book2 *Book
	Book3 *Book
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.c.Host, s.c.Port))
	if err != nil {
		panic(err)
	}
	fmt.Println("监听主机", s.c.Host, ", 监听端口:", s.c.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		user := NewUser(s, conn)
		s.AddUser(user)

		go func() {
			for {
				msg := <-s.Message
				for _, user := range s.OnlineSet {
					user.C <- []byte(msg)
					// user.conn.Write([]byte(msg))
				}
			}
		}()
	}
}

func (s *Server) AddUser(user *User) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.OnlineSet[user.Addr] = user
	go user.ListenReadMessage()
	go user.ListenWriteMessage()
	fmt.Println(user.Addr, " connect success")

	s.BroadMessage("[" + user.Addr + "]: 已上线")
}

func (s *Server) BroadMessage(msg string) {
	fmt.Println("发送广播消息: [" + msg + "]")
	s.Message <- msg
}

func NewServer(config *ServerConfig) *Server {
	if config == nil {
		config = &ServerConfig{
			Host: "localhost",
			Port: "2023",
		}
	}
	server := &Server{
		c:         config,
		OnlineSet: make(map[string]*User),
		Message:   make(chan string, 10),
	}
	return server
}
