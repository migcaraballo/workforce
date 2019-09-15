package workforce

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const(
	// global ENV setting for debug output from pool
	WORKFORCE_DEBUG_ENV = "WORKFORCE_POOL_DEBUG"
)

// WorkPool represents a worker pool that leverages channels and contains a slice of workers
type WorkerPool struct {
	// string name of pool to be used in debugging output
	Name        string
	// slice of workers that will be doing concurrent work
	workerPool  []Worker
	// buffered channel to submit workers
	workerChan chan Worker
	// buffered channel used by workers to signal when complete
	doneChan   chan bool
	// channel used by pool to send stop signal
	stopChan   chan bool
	// size of buffer for workerChan & doneChan
	buffSize   int
}

// Convenience function to create new pools.
func NewWorkerPool(name string, queSize int) (*WorkerPool, error) {
	if name == "" {
		return nil, errors.New("pool must have a name")
	}

	if queSize < 1 {
		queSize = 1
	}

	return &WorkerPool{
		Name:       name,
		workerChan: make(chan Worker, queSize),
		doneChan:   make(chan bool, queSize),
		stopChan:   make(chan bool),
		workerPool: []Worker{},
		buffSize:   queSize,
	}, nil
}

// Use this function to add workers to the pool before starting
func (wp *WorkerPool) AddWorker(w *Worker) {
	w.donChan = wp.doneChan
	wp.workerPool = append(wp.workerPool, *w)
}

// Use this function to set all workers at once
func (wp *WorkerPool) SetWorkers(wrks []Worker){
	wp.workerPool = wrks
}

// pipeline function to create pool of workers
// should only be called by this pool
func (wp *WorkerPool) startWorker(jchan <- chan Worker) {
	for {
		select {
		case w := <- jchan:
			w.work()
		case <- wp.stopChan:
			return
		}
	}
}

// Public function to start the pool and workers
func (wp *WorkerPool) StartPool(){
	Debug(fmt.Sprintf("[%s] Starting pool with %d workers & que = %d", wp.Name, len(wp.workerPool), wp.buffSize))

	// start the worker channels
	for i := 0; i < len(wp.workerPool); i++ {
		go wp.startWorker(wp.workerChan)
	}

	// send work to the channels
	go func() {
		for _, w := range wp.workerPool {
			wp.workerChan <- w
		}
	}()

	for i := 0; i < len(wp.workerPool); i++ {
		<-wp.doneChan
	}

	Debug(fmt.Sprintf("[%s] workers working: %d", wp.Name, len(wp.workerChan)))
	Debug(fmt.Sprintf("[%s] workers requested: %d", wp.Name, len(wp.workerPool)))
}

// Public function to stop worker pool and all workers
func (wp *WorkerPool) StopPool(){
	Debug(fmt.Sprintf("Stoping [%s] workerPool", wp.Name))
	wp.stopChan <- true
	Debug(fmt.Sprintf("Total workers: %d", len(wp.workerChan)))
}

// If WORKFORCE_POOL_DEBUG environment variable is set to "true", worker pool & worker will log output of work
func Debug(msg string){
	if IsDebug() {
		log.Println(msg)
	}
}

// checks to see if WORKFORCE_POOL_DEBUG is set
func IsDebug() bool {
	if  strings.ToLower(os.Getenv(WORKFORCE_DEBUG_ENV)) == "true" {
		return true
	}

	return false
}