package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)
type client chan<- string //an outgoing message channel

var(
	entering = make(chan client)
	leaving = make(chan client)
	messages = make(chan string) //all incoming client messages
)
func main(){
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil{
		log.Fatal(err)
	}
	go broadCaster()
	for {
		conn, err := listener.Accept()
		if err != nil{
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadCaster(){
	clients := make(map[client]bool) //all connected clients
	for{
		select{
			case msg := <-messages:
			//broadcast incoming message to all
			// client's outgoing message channels
		         for cli := range clients {
					   cli <- msg
				    }
			case cli := <- entering:
		 clients[cli] = true
			case cli := <- leaving:
			delete(clients,cli)
			close(cli)
			
		}
	}
}

func handleConn(conn net.Conn){
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You're " + who
	messages <- who + "\thas arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	if err := input.Err();err != nil{
		log.Fatal(err)
	}
	for input.Scan(){
		messages <- who + ": " + input.Text()
	}
	
	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}
func clientWriter(conn net.Conn, ch <- chan string){
	 for msg := range ch {
			 _, err := fmt.Fprintln(conn, msg)
				if err != nil{
					return
				}
		}
}