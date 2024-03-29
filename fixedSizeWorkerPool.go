package workforce

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
)

var (
	// patther: [poolName-type-workerName-ID]
	logMsgPattern = "%s-fixed"
)


type FixedSizeWorkerPool struct {
	// string name of pool to be used in debugging output
	Name        string
	// buffered channel to submit workers
	workerChan chan func() error
	// channel used by pool to send stop signal
	stopChan   chan bool
	numWorkers int
	workersStarted int
}

func NewFixedSizeWorkerPool(name string, workers, buffer int) (*FixedSizeWorkerPool, error) {
	if workers < 1 {
		return nil, errors.New("workers value must be >= 1")
	}

	fsp := &FixedSizeWorkerPool{
		stopChan: make(chan bool),
		workerChan: make(chan func() error, buffer),
		numWorkers: workers,
	}

	if name == "" {
		fsp.Name = uuid.New().String()
	} else {
		fsp.Name = name
	}

	return fsp, nil
}

func (fsp *FixedSizeWorkerPool) startWorker(name string){
	log.Printf("[%s] starting...\n", name)
	for {
		log.Printf("[%s] waiting for work\n", name)
		select {
		case wf := <- fsp.workerChan:
			//log.Printf("[%s] working", name)
			fsp.debug(fmt.Sprintf("[%s] working", name))
			fsp.handleWork(wf)
			//log.Printf("[%s] done", name)
			fsp.debug(fmt.Sprintf("[%s] done", name))
		case <- fsp.stopChan:
			//log.Printf("[%s] stopped working", name)
			fsp.debug(fmt.Sprintf("[%s] stopped working", name))
			return
		}
	}
}

func (fsp *FixedSizeWorkerPool) handleWork(wf func() error) {
	if err := wf(); err != nil {
		//log.Printf("workFunc err: %s\n", err)
		fsp.debug(fmt.Sprintf("workFunc err: %s\n", err))
	}
}

func (fsp *FixedSizeWorkerPool) Start() error {
	if fsp.numWorkers < 1 {
		return errors.New("FixedSizeWorkerPool must have at least 1 worker to start")
	}

	for i := 1; i <= fsp.numWorkers; i++ {
		fsp.workersStarted++
		go fsp.startWorker(fmt.Sprintf("worker-%d", i))
	}

	//log.Printf("Total Workers started: %d\n", fsp.workersStarted)
	fsp.debug(fmt.Sprintf("Total Workers started: %d\n", fsp.workersStarted))
	return nil
}

func (fsp *FixedSizeWorkerPool) Stop() {
	fsp.stopChan <- true
}

func (fsp *FixedSizeWorkerPool) AsyncSubmit(workerFunc func() error) {
	go func() {
		fsp.workerChan <- workerFunc
	}()
}

func (fsp *FixedSizeWorkerPool) SyncSubmit(workerFunc func() error) {
	fsp.workerChan <- workerFunc
}

func (fsp *FixedSizeWorkerPool) debug(msg string){
	log.Printf("["+logMsgPattern+"] %s\n", fsp.Name, msg)
}