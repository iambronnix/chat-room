package main

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"
)
var (
	done = make(chan struct{})
	requestPipeline = make(chan net.Listener,uint32(1<<30))//1gb pipeline
	requestList = make([]net.Listener, uint64(1<<60))//1gb payload
	serverAddr = make(chan string, 1)
	serverConn = make(chan net.Conn)
	wg = &sync.WaitGroup{}
	i = 7 //network error retries
)
func main(){
	go ProxyServer(data)
	
}

func ProxyServer(data io.ReadWriter){
  listRequests, err := net.Listen("tcp", "127.0.0.1:")//listen for incoming connections
   if err != nil{
   log.Panic(err)
   }
   defer listRequests.Close()

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

   //handle payload requests
   for j := 0; j >=0; j++{//start here where thers is a pipelne handling requests    
   go func(){   
          requestPipeline <- listRequests  //requests pipeline  
        done <- struct{}{}          
     }() 
   requestList = append(requestList, <- requestPipeline)
   
   }
   
  
   <- done   

	
}
func server(){
	listener, err := net.Listen("tcp", "127.0.0.1:64100")
	if err != nil{
		log.Panic(err)
	}
	defer listener.Close()
	serverAddr <- listener.Addr().String()
	for {
	conn, err := listener.Accept()//listen for proxy connection
	 if err != nil{
			return
		}
		go handleConnection(conn)
		
	}
	
}
func handleConnection(io.ReadWriter){
	//handles requests
}