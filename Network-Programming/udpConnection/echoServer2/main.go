package main

import (
	"bytes"
	"context"
	server "echoServer2/server"
	"log"
	"net"
	"time"
)
func main(){
	ctx , cancel := context.WithCancel(context.Background())
	serverAddr, servErr := server.EchoServer(ctx, "127.0.0.1:")
	if servErr != nil{
		log.Fatalf("%v", servErr)
	}
	defer cancel()
	client, err := net.Dial("udp", serverAddr.String())
	if err != nil{
		log.Fatalf("%v", err)
	}

	defer func(){ _ = client.Close()}()

	interloper, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil{
		log.Fatalf("%v",err)
	}
	interrupt :=  []byte("pardon me")
	n, err := interloper.WriteTo(interrupt,client.LocalAddr())//send a message from interloping connection
	if err != nil{
		log.Printf("%v", err)
	}
	_ = interloper.Close()

	if l := len(interrupt); l != n{
		log.Printf("wrote %d bytes of %d", n, l)
	}

	ping := []byte("ping")
	_, err = client.Write(ping)
	if err != nil {
		log.Printf("%v",err)
	}

	buff := make([]byte, uint32(1<<30))
	n, err = client.Read(buff)
	if err != nil{
		log.Printf("%v",err)
	}

	if !bytes.Equal(ping, buff[:n]){
	  log.Printf("expected reply %q; actual reply %q", ping, buff[:n])	
	}
	err = client.SetDeadline(time.Now().Add(time.Second))
	if err != nil{
		log.Printf("%v",err)
	}

	_, err = client.Read(buff)
	if err !=nil {
		log.Println("unexpected packet")
	}
	
	
}