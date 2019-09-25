package workforce

import (
	"fmt"
	"time"
)

// Worker represents a unit of work that needs to be done repeatedly based on the WorkHandler function
type Worker struct {
	// ID identifies which worker is doing work
	ID          string
	// channel for worker to send responses to
	donChan     chan bool
	// the function supplied to be executed per iteration of this worker
	WorkHandler func()
}

// Internal function. This should only be called by the WorkerPool. This should not be called by anything else
func (w *Worker) work() {
	if IsDebug(){
		st := time.Now()
		defer Debug(fmt.Sprintf("(%s) done - runtime: %f secs", w.ID, time.Since(st).Seconds()))
	}

	w.WorkHandler()
	w.donChan <- true
}

// Convenience function to create new workers.
func NewWorker(id string) *Worker {
	return &Worker{
		ID: id,
	}
}

// Convenience function to create new workers.
func NewWorkerWithFunc(id string, wrkFunc func()) *Worker {
	return &Worker{
		ID: id,
		WorkHandler: wrkFunc,
	}
}
