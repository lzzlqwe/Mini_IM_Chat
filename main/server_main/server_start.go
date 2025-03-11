package main

import (
	"Mini_IM_Chat/server"
)

func main() {
	server := server.NewServer("127.0.0.1", 8888)
	server.Start()
}
