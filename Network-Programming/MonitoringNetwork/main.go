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
			r := io.TeeReader(conn, monitor)//writes into monitor and returns an io.reader

			n, err := r.Read(b)
			if err != nil && err != io.EOF{
				monitor.Println(err)
				return
			}
			w := io.MultiWriter(conn, monitor)//duplicates its writes 

			_, err = w.Write(b[:n])//echo message
			if err != nil && err != io.EOF{
				monitor.Println(err)
				return
			}

	}()

		conn, err := net.Dial("tcp",listener.Addr().String())//dialer to the server
		if err != nil{
			monitor.Fatal(err)
		}
		_, err = conn.Write([]byte("iamerick\n"))
		if err != nil{
			monitor.Fatal(err)
		}
		_ = conn.Close()
		<- done		

		
}		

func main(){
	ExampleMOnitor()
	
}