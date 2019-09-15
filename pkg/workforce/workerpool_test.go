package workforce

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func setup(){
	os.Setenv(WORKFORCE_DEBUG_ENV, "true")
}

// Simple test to fire off 10 workers in a pool.
// Once the workers are done, test exits
func TestNewWorkerPool_Default(t *testing.T) {
	setup()
	workers := 10
	//wp, err := NewWorkerPool("test-workerPool", workers, 1)
	wp, err := NewWorkerPool("test-workerPool", 1)

	if err != nil {
		panic(err)
	}

	for i := 1; i <= workers; i++ {
		w := &Worker{
			ID: fmt.Sprintf("tw-%d", i),
		}

		// mock function that sleeps for 1 second
		w.WorkHandler = func() {
			time.Sleep(1 * time.Second)
			log.Printf("(%s) working...", w.ID)
		}

		wp.Add(w)
	}

	wp.Start()
	defer wp.Stop()
}

// This test spins up a pool of 5 workers which work for 10 seconds. During the 10 seconds, work is being done
// inside the WorkHandler() function. This mimics a longer lasting connection to an external resource.
// calls counter is an approximation since we are not using a mutex to synchronize on the calls variable.
func TestNewWorkerPool_TenSecWorker(t *testing.T) {
	setup()
	runtime := time.Duration(5)
	workers := 5
	pool, _ := NewWorkerPool("topic-test", 10)
	//pool, _ := NewWorkerPool("topic-test", workers, 10)

	for i := 1; i <= workers; i++ {
		l := NewWorker(fmt.Sprintf("wrk-%d", i))

		l.WorkHandler = func() {
			fmt.Printf("(%s) started listening\n", l.ID)
			rand.Seed(time.Now().UnixNano())

			// mimic doing work against some extenal resource.
			calls := 0
			go func() {
				for {
					calls++
					st := time.Now()
					rt := time.Duration(rand.Intn(500))
					time.Sleep(rt * time.Millisecond)
					fmt.Printf("(%s) response: \t ts: %s\n", l.ID, time.Since(st))
				}
			}()

			// test long durations
			time.Sleep(runtime * time.Second)

			fmt.Printf("(%s) calls: %d\n", l.ID, calls)
			fmt.Printf("(%s) stopped listening\n", l.ID)
		}

		pool.Add(l)
	}

	pool.Start()
	defer pool.Stop()
}
