package main

import (
	"log"
	"net"
	

)

func main(){
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil{
		log.Printf("%v", err)
	}

	msgs := []struct{message, Reply string}{
		{"ping", "pong"},
		{"pong", "pong"},
		{"echo", "echo"},
		{"ping", "pong"},
	}
	for i, m := range msgs{
		_, err = conn.Write([]byte(m.message))
		if err != nil{
			panic(err)
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil{
			panic(err)
		}
		actual := string(buf[:n])
		log.Printf("%q ->proxy -> %q", m.message, actual)
		if actual != m.Reply{
			log.Panicf("%d: expected reply: %q; actual: %q", i, m.Reply, actual)
		}
	}
	_ = conn.Close()
}