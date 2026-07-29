package main

import (
	"context"
	"fmt"
	"io"
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
				 if _,err := w.Write([]byte("iamerick"));err != nil{
						return
				//track and act on consecutive timeouts here
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
    for i, v := range []int64{0,200,300,0,-1,-1,-1}{
         fmt.Printf("Run %d:\n", i+1)
         receivePing(time.Duration(v)*time.Millisecond,r)
    }

    cancel()
    <-done //ensures the pinger exits after cancelling the context
}

func main(){
	pinger()
	
}