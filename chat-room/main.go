package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	dialServer()
}
func dialServer() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	defer func() {
		conn.Close()
	}()
	if err != nil {
		log.Panic(err)
	}
	go clientServer()

}
func clientServer() {
	}
