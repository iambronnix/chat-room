package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"sync"
)
var (
	wg = &sync.WaitGroup{}
)
func main(){
	Listener()
}
func Listener(){
	listener, err := net.Listen("tcp","127.0.0.1:8000")
	if err != nil{
		log.Panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err!=nil{
			log.Println(err)
			continue//connection aborted
		}
		go serverChat( conn)
	}
	
		
}

func serverChat(conn net.Conn){
	defer func(){
		wg.Done()
	 conn.Close()
	}()
	for i := 50; i > 1; i--{
		wg.Add(i)
		go func(){
			io.Copy()
		}
	}
	
	
}