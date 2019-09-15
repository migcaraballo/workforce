# Workforce
Workforce is a simple go module to create and use worker pools.

## Installation
`
$ go get github.com/migcaraballo/workforce
` 

## Example
```
func main(){
	workers := 3

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
```

##### Background
Originally created for my personal use, i didn't want to keep having to re-write this every time i needed a worker pool. However, the real goal is to be able to keep my concurrent thoughts in one place. 

If you would like to help me improve this project, please fork this repo and submit some pull requests.

