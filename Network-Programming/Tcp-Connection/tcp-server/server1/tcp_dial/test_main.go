package main

import (
	"io"
	"net"
	"testing"
)
func TestDial(t *testing.T){
	//create a listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1")//let Go randomly pick an available port
	if err != nil{
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func(){
		defer func(){
			done <- struct{}{}//signal
		}()
		for {
			conn, err := listener.Accept()//spin off the listener in a goroutine
			if err != nil{
				t.Log(err)
				return
			}
			go func(c net.Conn){
				defer func(){
					c.Close()
					done <- struct{}{}
				}()
				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf)
					if err != nil{
					  if err != io.EOF{
								t.Error(err)
							}
							return
					}
					t.Logf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil{
		t.Fatal(err)
	}
	conn.Close()
	<-done
	listener.Close()
	<-done
}