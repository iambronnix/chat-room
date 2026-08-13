package server

import (
	"context"
	"fmt"
	"net"
)

//this piece of code echoes any udp packets received to the sender
func EchoServerUDP(ctx context.Context, addr string)(net.Addr, error){//contexr allows cancellation by the caller
	s, err := net.ListenPacket("udp",addr)
	if err != nil{
		return nil, fmt.Errorf("binding to udp %s: %w", addr, err)
	}

	go func(){
		go func(){//blocks on the context's Done channel
			<-ctx.Done()
			_ = s.Close()
		}()

		buf := make([]byte, 1024)//a byte slice 
        for {
         n, clientAddr, err := s.ReadFrom(buf)// client to server
       if err != nil{
       return 
          }
          _, err = s.WriteTo(buf[:n], clientAddr) //server to client
          if err != nil{
          return
          }
        }
	}()
	return s.LocalAddr(), nil
}