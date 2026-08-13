package main

import (
	"bytes"
	"context"
	server "echoServer1/server"
	"log"
	"net"
)

func main(){
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr , err := server.EchoServer(ctx,"127.0.0.1:")//server connection
	if err != nil{
		log.Fatalf("%v", err)
	}
	defer cancel()

	client, err := net.ListenPacket("udp", "127.0.0.1:")//client connection
	if err != nil{
		log.Fatalf("%v",err)
	}
	defer func() { _ = client.Close()}()

	interloper, err := net.ListenPacket("udp","127.0.0.1:")//new udp conn that interlopes on the client and echo server and interrupt the client
	if err != nil{
		log.Fatalf("%v", err)
	}
	interrupt := []byte("pardon me")
	n, err := interloper.WriteTo(interrupt,client.LocalAddr())//message queue's up inthe client's receive buffer
	if err != nil{
		log.Fatalf("%v", err)
	}
	_ = interloper.Close()
	if l := len(interrupt); l != n {
		log.Fatalf("wrote %d bytes of %d", n, l)
	}

	ping := []byte("ping")
	_, err = client.WriteTo(ping, serverAddr)//writes to server a "ping" message 
	if err != nil{
		log.Fatalf("%v", err)
	}
	buf := make([]byte, 1024)
	n, addr, err := client.ReadFrom(buf)//reads promptly an incoming message 
	if err != nil{
		log.Fatalf("%v", err)
	}

	if !bytes.Equal(ping, buf[:n]){
		log.Printf("expected reply %q; actual reply %q\n", ping, buf[:n])
	}

	if addr.String() != serverAddr.String(){//verifies the sender of each packet 
		log.Printf("expected message from %q; actual sender is %q\n", serverAddr, addr)
	}
}