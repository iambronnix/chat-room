package main

import (
	"log"
	"net"
)

func main(){
	//tcpConn, err := conn.(*net.TCPConn)
}

func server()(*net.TCPConn){
	//retrieving *net.TCPconn from the listener
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:")
	if err != nil{
		log.Fatalf("%v", err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil{
		log.Fatalf("%v",err)
	}

	tcpConn, err := listener.AcceptTCP()
	if err != nil{
		log.Fatalf("%v",err)
	}
	return tcpConn
}

func client()*net.TCPConn{
	//retrieving *net.TCPconn from the dialer on clientside
	addr, err := net.ResolveTCPAddr("tcp","www.google.com:https")
	if err != nil{
		log.Fatalf("%v",err)
	}
	tcpConn, err := net.DialTCP("tcp", nil, addr)
	return tcpConn
}