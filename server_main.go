package main

import (
	"fmt"
	"io"
	"os"
	"tiangong/server"
)

func init() {
	file, _ := os.Open("./server_banner.txt")
	if file != nil {
		defer file.Close()
		bytes, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}
		banner := string(bytes)
		fmt.Println(banner)
	}
}

var config *server.Config

func main() {
	server := server.NewServer(nil)
	server.Start()
}
