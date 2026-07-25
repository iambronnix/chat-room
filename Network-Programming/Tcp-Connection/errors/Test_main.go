package errors_test
import(
	"net"
	"time"
	"log"
	"errors"
)
var(
 err error
 n int
 i = 7 //maximum number of retries
 )
 func Test()error{
var conn net.Conn
 
 for; i > 0 ; i--{
  n, err = conn.Write([]byte("hello world"))
  		if err != nil{
    	if nErr, ok := err.(net.Error); ok && nErr.Temporary(){//check whether the error is temporary
     		log.Println("temporary error:", nErr)
       		time.Sleep(10* time.Second)
         		continue
     }
     return err
 }
 break
 }
 if i == 0{
 return errors.New("temporary write failure threshhold exceeded")
 }
 log.Printf("wrote %d bytes to %s\n", n, conn.RemoteAddr())
 }