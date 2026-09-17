package server

import (
	"log"
	"net"
)

func Server()chan string{
	listener, err := net.Listen("tcp4", "127.0.0.1:63500")
	if err != nil{
		log.Panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err !=nil{
			continue
		}
		
	}

}