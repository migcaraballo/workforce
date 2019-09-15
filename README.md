# Workforce
Workforce is a simple go module to create and use worker pools.

## Installation
`
$ go get github.com/migcaraballo/workforce
` 

## Example
Below is a simple example. For working examples you can run yourself, see tests and examples dir.
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

##### Etc...
 If you would like to help me improve this project, please fork this repo and submit some pull requests.
