package client

import (
	"fmt"
	"net"
)

type Client struct {
	host string
	port int
	conn net.Conn
}

func (s *Client) connect() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		panic(err)
	}
	s.conn = conn

	go func(c net.Conn) {
		buf := make([]byte, 1024)
		for {
			len, err := c.Read(buf)
			if len == 0 {
				fmt.Println("服务器已停止")
				return
			}
			if err != nil {
				fmt.Println("读取异常")
				return
			}
			msg := string(buf[:len-1])
			fmt.Println("收到服务端消息: ["+msg+"], 消息长度: ", len)
		}
	}(s.conn)
}

func (s *Client) Write(msg []byte) {
	n, err := s.conn.Write(msg)
	if err != nil {
		fmt.Println("数据写入错误,", err)
	}
	fmt.Println("写入数据长度,", n)
}

func NewClient(host string, port int) *Client {
	client := &Client{
		host: host,
		port: port,
	}
	client.connect()
	return client
}
