package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

const defaultPingInterval = 30 * time.Second

func ping (ctx context.Context, w io.Writer, reset <-chan time.Duration){
	var interval time.Duration
	select{
		case <- ctx.Done():
		return
		case interval = <- reset: //pulled initial interval off reset channel
		}
		if interval <= 0 {
			interval = defaultPingInterval
		}

		timer := time.NewTimer(interval)
		defer func(){
			if !timer.Stop(){
				<-timer.C
			}
		}()

		for {
			input := bufio.NewScanner(os.Stdin)
			select{
				case <- ctx.Done():
				     return
				case newInterval := <- reset:
				     if !timer.Stop(){
					    <- timer.C
				}
				if newInterval > 0 {
					interval = newInterval
				}
				case <- timer.C:
				for input.Scan(){
				 _,err := w.Write(input.Bytes())
				 if err != nil{
						panic(err)
				//track and act on consecutive timeouts here
				}
				return
				}
			}
			_ = timer.Reset(interval)
		}
}
func pinger(){
	ctx, cancel := context.WithCancel(context.Background())
	r, w := io.Pipe()
	done := make(chan struct{})
    resetTimer := make(chan time.Duration, 1)
    resetTimer <- time.Second//initial ping interval

    go func(){
    ping(ctx,w, resetTimer)
    close(done)
    }()

    receivePing := func(d time.Duration, r io.Reader){
         if d >= 0{
         fmt.Printf("resetting timer (%s)\n",d)
         resetTimer <- d
         }
         now := time.Now()
         buf := make([]byte, 1024)
         n, err := r.Read(buf)

         if err != nil{
         fmt.Println(err)
         }

         fmt.Printf("received %q (%s)\n", buf[:n], time.Since(now).Round(100*time.Millisecond))
    }
    receivePing(time.Second,r)

    cancel()
    <-done //ensures the pinger exits after cancelling the context
}

func main(){
	pinger()
	
}