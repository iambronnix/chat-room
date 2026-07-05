package main

import (
	"crypto/rand"
	"io"
	"net"
	"testing"
)
func TestReadIntoBuffer(t *testing.T){
	payload := make([]byte,1<<24) //16 MB
	_, err := rand.Read(payload)//generate a random payload
	if err != nil{
		t.Fatal(err)
	}
	listener, err := net.Listen("tcp","127.0.0.1:")
	if err != nil{
		t.Fatal(err)
	}

	go func(){//listen for incoming connections 
		conn, err := listener.Accept()
		if err != nil{
			t.Log(err)
			return
		}
		defer conn.Close()
		_, err = conn.Write(payload)//writes entire payload to the network connection
		if err != nil{
			t.Error(err)
		}
	}()
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil{
		t.Fatal(err)
	}
	buf := make([]byte, 1<<19)//512 KB
	for{//iterates through the connection to read all the data
		n, err := conn.Read(buf)//reads upto 512kb 
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}
		t.Logf("read %d bytes", n)// buf[:n] is the data read from conn
	}
	conn.Close()
}