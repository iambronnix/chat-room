package main

import (
	"io"
	"log"
	"net"
)
var (
	dialAddr = make(chan string,1)
)

func main(){
	
}
func proxyConn(server string){
	connServer, err := net.Dial("tcp", server)
	if err != nil{
		log.Panic(err)
	}
	defer connServer.Close()
	
	connDial, err := net.Dial("tcp",<-dialAddr)
	if err != nil{
		log.Panic(err)
	}
	connDial.Close()
	go func(){
		io.Copy(connDial, connServer)
	}()
	io.Copy(connServer, connDial)
}
func dial(){
	listener, err := net.Listen("tcp","127.0.0.1")
	if err != nil{
		log.Panic()
	}
	defer listener.Close()
	dialAddr <- listener.Addr().String()
	
	
}
