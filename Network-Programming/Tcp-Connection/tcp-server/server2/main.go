package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)
func main(){
	listener, err := net.Listen("tcp","localhost:8080")
	if err != nil{
		log.Fatal(err)
	}
	for {
		conn, connErr := listener.Accept()
		if connErr!=nil{
			return //abort connection
		}
		go handleConnection(conn)
		
	}
}
func handleConnection(c net.Conn){
	input := bufio.NewScanner(c)
	for input.Scan(){
		if scanErr := input.Err();scanErr == io.EOF{
			break
		}
     echo(c,input.Text(),1*time.Second)		
	}
}
func echo(c net.Conn, shout string, delay time.Duration){
	fmt.Fprintln(c, "\t",strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c,"\t", strings.ToLower(shout))
}