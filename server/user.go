package server

import (
	"fmt"
	"net"
)

type User struct {
	Addr   string
	C      chan []byte
	conn   net.Conn
	server *Server
}

func NewUser(server *Server, conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		Addr:   addr,
		C:      make(chan []byte),
		conn:   conn,
		server: server,
	}
	return user
}

func (s *User) ListenWriteMessage() {
	for {
		msg := <-s.C
		msg = append(msg, byte('\n'))
		s.conn.Write(msg)
	}
}

func (s *User) ListenReadMessage() {
	buf := make([]byte, 1024)
	for {
		len, err := s.conn.Read(buf)
		if len == 0 {
			fmt.Println(s.Addr, "已经下线")
			return
		}

		if err != nil {
			fmt.Println(s.Addr, "读取异常")
			return
		}
		if buf[len-1] == '\n' {
			buf[len-1] = 0x00
		}
		msg := string(buf[:len])
		s.server.BroadMessage(msg)
	}
}
