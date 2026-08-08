package main

import (
	"fmt"
	"time"
)

//worker function that processes integers from a channel
//
func worker(id int, jobs <-chan int, results chan<- int){
	for job := range jobs{
		fmt.Println("worker", id , "Processing job", job)
		time.Sleep(time.Second) // simulate a time-consuming task 
		results <- job * 2  //send the processed results 
	}
}

func main(){
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	//create a pool of 3 workers
	for w := 1; w<=3; w++{
		go worker(w, jobs, results)
	}
	//send 9 jobs then close the channel to indicate that's ll the work
	for j := 0; j <= 9; j++{
	jobs <- j
	}
	close(jobs)

	//collect the results 
	for a := 1; a <= 9;a++{
		<-results
	}
}