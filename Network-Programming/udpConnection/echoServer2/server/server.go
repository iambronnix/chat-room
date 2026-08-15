package server

import (
	"context"
	"fmt"
	"net"
)
func EchoServer(ctx context.Context, addr string)(net.Addr, error){
	serverAddr, servErr := net.ListenPacket("udp", addr)
	if servErr != nil{
		return nil,fmt.Errorf("error binding %q: to %q", addr, servErr)
	}
	go func(){
		go func(){
			<- ctx.Done()//garbage and memory cleaner
			serverAddr.Close()
		}()

		buf := make([]byte, 1024)
		for{
			n, clientAddr, readErr := serverAddr.ReadFrom(buf)
			if readErr != nil{
				return //retry reading from the connection
			}
			_, writeErr := serverAddr.WriteTo(buf[:n], clientAddr)
			if writeErr != nil{
				return //retry writing to the connection
			}
		}
	}()
	return serverAddr.LocalAddr(), nil
}