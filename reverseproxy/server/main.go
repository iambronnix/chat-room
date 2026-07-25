package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)
var(
	proxyChan = make(chan string)
	serverString = make(chan string)
    
)

func proxyWorker(from io.Reader, to io.Writer)error{
	fromWriter, fromIsWriter := from.(io.Writer)//test for io.writer implementation
	toReader, toIsReader := to.(io.Reader)//test for io.reader implementation

	if toIsReader && fromIsWriter{
		//sends replies since "from" and  "to" implement necessary interfaces
		go func(){
			_,_ = io.Copy(fromWriter, toReader)
		}()
	}
	_, err := io.Copy(to, from)
	return err
}
func serverMain()string{
	//sever listens for a "ping" message and responds with a "pong" message
	// all other messages are echoed back to the client
	server, err := net.Listen("tcp", "127.0.0.1:45554")//listens for the proxy dialing
	if err != nil{
		panic(err)
	}
	fmt.Println("marker number 2")
	go func(){
		for {
			conn, err := server.Accept()
			if err != nil{
				log.Printf("%v", err)
				return				
			}

			go func(c net.Conn){
				defer conn.Close()
				for {
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if connErr,ok := err.(net.Error); ok && connErr != io.EOF{//test for read error
						log.Printf("%v", connErr)
						return
					
					}

					switch msg := string(buf[:n]); msg{
						case "ping"://reply with pong for ping
						_,err := conn.Write([]byte("pong"))
						if err != nil && err != io.EOF{
							log.Printf("%V", err)
						}
						default:
						_, err := conn.Write(buf[:n])//echo any other message send 
						if err != nil && err != io.EOF{
							log.Printf("%V", err)
						}
					}
					
				}
			}(conn)
		}
	}()

	return server.Addr().String()
}

func Proxy(wg *sync.WaitGroup){
	serverAddr := serverMain()
	proxyServer, err := net.Listen("tcp", "127.0.0.1:8080")//listen for incoming connections
	if err != nil{
		panic(err)
	}
	wg.Add(1)
	go func(){
		defer wg.Done()
		for {
			conn, err := proxyServer.Accept()
			if err != nil{
				log.Printf("%v", err)
				return///retry 
			}
			fmt.Println("waiting for a dial up")
			go func(from net.Conn){
				defer from.Close()

				to, err := net.Dial("tcp", serverAddr)//establish connection to destination server
				if err != nil{
					log.Printf("%v", err)
					return // retry dial to server
				}
				defer to.Close()

				err = proxyWorker(from, to)//the proxyworker that does that thing
				if err != nil && err != io.EOF{
					panic(err)
				}
			}(conn)
		}
	}()

}

func main(){
	wg := & sync.WaitGroup{}
	wg.Add(2)
	fmt.Println("marker number 1")
	 go Proxy(wg)
	fmt.Println("marker number 3")
	wg.Wait()
	
}
