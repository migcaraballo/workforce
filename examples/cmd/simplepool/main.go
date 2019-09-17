package main

import (
	"fmt"
	"github.com/migcaraballo/workforce"
	"os"
	"strconv"
	"time"
)

const defaultWorkers = 3

func main(){
	workers := defaultWorkers

	// see if any worker count was passed in
	if len(os.Args) == 2 {
		w, err := strconv.Atoi(os.Args[1])
		if err == nil {
			workers = w
		}
	}

	// create a new worker pool
	pool, err := workforce.NewWorkerPool("sample-pool", 10)

	if err != nil {
		panic(err)
	}

	// create the workers and add them to the pool
	for i := 1; i <= workers; i++ {
		wrk := workforce.NewWorker(fmt.Sprintf("worker-%d", i))

		// give the worker something to do
		wrk.WorkHandler = func() {
			// show that work is starting
			fmt.Printf("%s is working\n", wrk.ID)

			// sleep a little to mimic some work/processing
			time.Sleep(500 * time.Millisecond)

			// show that work is done
			fmt.Printf("%s is done\n", wrk.ID)
		}

		// add the worker to the pool
		pool.Add(wrk)
	}

	// start the pool and defer stopping until all work is done
	pool.Start()
	defer pool.Stop()
}
