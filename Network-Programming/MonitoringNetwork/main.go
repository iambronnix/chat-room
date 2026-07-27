package main

import (
	"io"
	"log"
	"net"
	"os"
)

// monitor embeds a log.Logger  meant fr logging server's network traffic
type Incomming struct{
	*log.Logger
}
type OutGoing struct{
	*log.Logger
}
// write implements the io.Writer interface
func (I *Incomming) Write(p []byte) (int, error) {
	return len(p), I.Output(1, string(p))
}
func (O *OutGoing) Write(p []byte)(int, error) {
	return len(p), O.Output(1, string(p))
}

func ExampleMOnitor() {
	//log instance
	incoming := &Incomming{Logger: log.New(os.Stdout, "clientRequest:", 0)} //logs clients side request
	outGoing := &OutGoing{Logger: log.New(os.Stdout, "serverResponse:", 0)} //logs the serverside response
	
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		incoming.Fatal(err)
	}
	done := make(chan struct{})

	go func() {
		defer close(done)

		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		b := make([]byte, 1024)
		r := io.TeeReader(conn, incoming) //writes clients side requests from the connection 

		n, err := r.Read(b)
		if err != nil && err != io.EOF {
			incoming.Println(err)
			return
		}
		w := io.MultiWriter(conn, outGoing) //logs server response

		_, err = w.Write(b[:n]) //echo message
		if err != nil && err != io.EOF {
			outGoing.Println(err)
			return
		}

	}()

	conn, err := net.Dial("tcp", listener.Addr().String()) //dialer to the server
	if err != nil {
		outGoing.Fatal(err)
	}
	_, err = conn.Write([]byte("iamerick\n"))
	if err != nil {
		outGoing.Fatal(err)
	}
	_ = conn.Close()
	<-done

}

func main() {
	ExampleMOnitor()

}
