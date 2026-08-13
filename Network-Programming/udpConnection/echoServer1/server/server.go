package server

import (
	"context"
	"fmt"
	"net"
)

func EchoServer(ctx context.Context, addr string)(net.Addr,error){
	s, err := net.ListenPacket("udp", addr)
	if err != nil{
		return  nil, fmt.Errorf("binding to udp %s: %w", addr, err)
		
	}
	go func(){
		go func() {//memory cleanup and leakage 
			<- ctx.Done()
			s.Close()
		}()

		buf := make([]byte, 1024)
		for {
		n, clientAddr, readErr := s.ReadFrom(buf)//read clients message 
		if readErr != nil{
			return 
		}

		_, writeErr := s.WriteTo(buf[:n], clientAddr)
		if writeErr != nil{
			return 
		}
		}
	}()

	return s.LocalAddr(), nil//server address for message accounting
	
}