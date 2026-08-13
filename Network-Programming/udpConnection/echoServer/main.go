package main

import (
	"bytes"
	"context"
	server "echoServer/server"
	"fmt"
	"log"
	"net"
)

func main(){
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr, err := server.EchoServerUDP(ctx, "127.0.0.1:")//pass ctx object and local addres
	if err != nil{
		log.Fatalf("%v",err)
	}
	defer cancel()//signals the server to exit and close its goroutines

	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil{
		log.Fatalf("%v", err)
	}
	defer func(){_ = client.Close()}()

	msg := []byte("ping")
	_, err = client.WriteTo(msg, serverAddr)//where to send ur messsage 
	if err != nil {
		log.Fatalf("%v",err)
	}
	fmt.Println("client message: ",string(msg))

	buf := make([]byte, 1024)//byte slice
	n, addr, err := client.ReadFrom(buf)//examine the address returned to confirm echoServer send the message 
	if err != nil{
		log.Fatalf("%v", err)
	}
	if addr.String() != serverAddr.String(){
		log.Fatalf("received reply from %q instead of %q", addr, serverAddr)
	}
	if !bytes.Equal(msg,buf[:n]){
		log.Fatalf("expected reply %q; actual reply %q", msg, buf[:n])
	}
	fmt.Println("server reverb: ",string(buf[:n]))
	
}