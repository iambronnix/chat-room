package main

import (
	"io"
	"log"
	"net"
	"os"
)
func main(){
	conn,dialErr := net.Dial("tcp", "127.0.0.1:8080")
	if dialErr !=nil{
	 log.Fatal(dialErr)
	}
	defer conn.Close()	
	go mustCopy(os.Stdout, conn)
	mustCopy(conn,os.Stdout)
}
func mustCopy(dst io.Writer, src io.Reader){
   if _,err := io.Copy(dst, src);err!=nil{
   log.Fatal()
   }
}