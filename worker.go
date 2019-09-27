package workforce

import (
	"fmt"
	"time"
)

// Worker represents a unit of work that needs to be done repeatedly based on the WorkHandler function
type Worker struct {
	ID          string
	doneChan    chan bool
	WorkHandler func()
}

// Internal function. This should only be called by the WorkerPool. This should not be called by anything else
func (w *Worker) work() {
	if IsDebug(){
		st := time.Now()
		defer Debug(fmt.Sprintf("(%s) done - runtime: %f secs", w.ID, time.Since(st).Seconds()))
	}

	w.WorkHandler()
	w.doneChan <- true
}

// Convenience function to create new workers.
func NewWorker(id string) *Worker {
	return &Worker{
		ID: id,
		doneChan: make(chan bool),
	}
}

// Convenience function to create new workers.
func NewWorkerWithFunc(id string, wrkFunc func()) *Worker {
	return &Worker{
		ID: id,
		WorkHandler: wrkFunc,
	}
}
