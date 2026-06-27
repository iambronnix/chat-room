package dial1

import (
	"context"
	"net"
	"testing"
	"time"
	"syscall"
)
func TestDialContext(t *testing.T){
	dl := time.Now().Add(5 * time.Second)//context with a deadline of five seconds into the future
	ctx, cancel := context.WithDeadline(context.Background(),dl)//create the context and it cancel function
	defer cancel()//ensures the context is garbage collected as soon as possible

	var d net.Dialer//DialContext is a method on a dialer
	d.Control = func(_,_ string, _ syscall.RawConn) error {//overide the dialer's control()
		//sleep long enough to reach the context's deadline.
		time.Sleep(5*time.Second + time.Millisecond)//delay long enough to make sure you exceed context's deadline
		return nil
	}
	conn, err := d.DialContext(ctx, "tcp","10.0.0.0:80")//pass in context as first argument
	if err == nil{
		conn.Close()
		t.Fatal("connection didn't time out")
	}
	nErr, ok := err.(net.Error)
	if !ok{
		t.Error(err)
	}else{
		if !nErr.Timeout(){
			t.Errorf("error isn't a timeout: %v", err)
		}
	}
	if ctx.Err() != context.DeadlineExceeded{
		t.Errorf("expected deadline exceeded; actual: %v", ctx.Err())
	}
}