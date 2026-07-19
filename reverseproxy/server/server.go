package server

import (
	"io"
	"log"
	"net"
	"sync"
)
var(
	proxyChan = make(chan string)
	serverString = make(chan string)
    wg *sync.WaitGroup
    
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
func serverMain(wg *sync.WaitGroup){
	//sever listens for a "ping" message and responds with a "pong" message
	// all other messages are echoed back to the client
	server, err := net.Listen("tcp", "127.0.0.1:")//listens for the proxy dialing
	if err != nil{
		panic(err)
	}
	serverString <- server.Addr().String()//sends the serverAddres to dialup
	wg.Add(1)
	go func(){
		defer wg.Done()
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
					if connErr,ok := err.(net.Error); !ok && connErr != io.EOF{//test for read error
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
}

func Proxy(){
	proxyServer, err := net.Listen("tcp", "127.0.0.1:80800")//listen for incoming connections
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
			go func(from net.Conn){
				defer from.Close()

				to, err := net.Dial("tcp", <- serverString)//establish connection to destination server
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

func Start(){
	
	wg.Add(1)
	go serverMain(wg)
	Proxy()
	wg.Wait()
	
}