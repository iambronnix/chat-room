package main

import (
	"io"
	"log"
	"net"
	"os"
)

//monitor embeds a log.Logger  meant fr logging server's network traffic
type Monitor struct{
	*log.Logger
}
//write implements the io.Writer interface
func (m *Monitor) Write(p []byte)(int, error){
	return  len(p), m.Output(2, string(p))
}

func ExampleMOnitor(){
	 monitor := &Monitor{Logger: log.New(os.Stdout, "monitor:", 0)} //an instance of Monitor to os.Stdout
		listener, err := net.Listen("tcp", "127.0.0.1:")
		if err != nil{
			monitor.Fatal(err)
		}
		done := make(chan struct{})

		go func(){
			defer close(done)

			conn, err := listener.Accept()
			if err != nil{
				return
			}
			defer conn.Close()

			b := make([]byte, 1024)
			r := io.TeeReader(conn, monitor)

			n, err := r.Read(b)
			if err != nil && err != io.EOF{
				monitor.Println(err)
				return
			}
			w := io.MultiWriter(conn, monitor)

			_, err = w.Write(b[:n])//echo message
			if err != nil && err != io.EOF{
				monitor.Println(err)
				return
			}

	}()

		
}		