package main

import (
	"io"
	"log"
	"net"
	"time"
)
func main(){
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil{
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil{
			log.Print(err)//connection aborted
			continue
		}
	   go handleConn(conn)//handle connections concurrently
	}
}
func handleConn(c net.Conn){
	defer c.Close()
	for{
		_, err := io.WriteString(c,time.Now().Format("15:04:05\n"))
		if err != nil{
			return//client disconnected
		}
		time.Sleep(1*time.Second)
	}
}