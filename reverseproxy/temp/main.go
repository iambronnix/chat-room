package main

import (
	"errors"
	"io"
	"log"
	"net"
	"time"
)
var (
	done = make(chan struct{})	
	requestBuffer = make([]net.Conn,uint32(1<<30))//1gb buffer
	serverAddr = make(chan string, 1)
	i = 7 //network error retries
)
func main(){
	//main code
	
}

func ProxyServer(data io.ReadWriter){
	//listen client requests
  listRequests, err := net.Listen("tcp", "127.0.0.1:")
   if err != nil{
   log.Panic(err)
   }
   defer listRequests.Close()
    //handle serverside temporary errors
   serverConn, serverErr := net.Dial("tcp", <-serverAddr)    
   for ;i >0 ; i--{
   if serverErr!=nil{
      if nErr, ok := serverErr.(net.Error); ok && nErr.Temporary(){//test if serverErr is temporary
          log.Println(nErr)
          time.Sleep(5 * time.Second)//wait as the error is being resolved
          continue//return to for loop
      }
      log.Panic(serverErr) //if error isn't temporary
   }
   break   
   }
   if i == 0{
   log.Fatalf("%s", errors.New("Temporary server failure threshold exceeded"))//temporary server failure exit
   }
   defer serverConn.Close()

   for {
			connRequests, err := listRequests.Accept()
			if err != nil{
				log.Printf("%s",err)
				continue
			}
			go func(){
		     requestBuffer = append(requestBuffer, connRequests)//1gb buffer to store requests and prevent ddos payloads
				for _, j := range requestBuffer{
				_, err := io.Copy(serverConn,j)
				if err != nil{
					return
				}
			}
				done <- struct{}{}				
			}()

			if _, err = io.Copy(connRequests,serverConn);err != nil{
				log.Printf("%s", err)
			}
			
			<- done
			
		}
      
   
}

func server(){
	listener, err := net.Listen("tcp", "127.0.0.1:64100")
	if err != nil{
		log.Panic(err)
	}
	defer listener.Close()
	serverAddr <- listener.Addr().String()
	for {
	conn, err := listener.Accept()//accept proxy connection
	 if err != nil{
			return
		}
		go func(){
			io.Copy()//
		}()
		go handleConnection(conn)
		
	}
	
}
func handleConnection(conn net.Conn){
   payload := "i'm a server encapsulated within a proxy for security purposes"
	 
	_,err := conn.Write([]byte(payload))
	if err != nil{
		log.Panicf("%s", err)
	}
	
}