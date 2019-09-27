package workforce

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewFixedSizeWorkerPool_OneWorker(t *testing.T) {
	tasks := 4
	pool, err := NewFixedSizeWorkerPool("", 2, 2)

	if err != nil {
		t.Fatal(err)
	}

	pool.Start()

	for i := 1; i <= tasks; i++ {
		name := fmt.Sprintf("func-%d", i)

		pool.AsyncSubmit(func() error {
			//log.Printf("[%s] doing some work.", name)
			pool.debug(fmt.Sprintf("[%s] doing some work.", name))

			// do some work
			time.Sleep(500 * time.Millisecond)

			//log.Printf("[%s] Finished working...", name)
			pool.debug(fmt.Sprintf("[%s] Finished working...", name))
			return nil
		})
	}

	time.Sleep(5500 * time.Millisecond)
	pool.Stop()
	log.Println("done testing")
}