package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// Multiplying this will speed up the run time by a corresponding amount
var threadLimit = runtime.NumCPU() * 2 * 2 * 2

type job struct {
	id        int
	randomNum int
}

type result struct {
	j           job
	sumOfDigits int
}

var jobs = make(chan job, threadLimit)
var results = make(chan result, threadLimit)

func digits(number int) int {
	sum := 0
	no := number

	for no != 0 {
		digit := no % 10
		sum += digit
		no /= 10
	}

	time.Sleep(50 * time.Millisecond)

	return sum
}

func worker(wg *sync.WaitGroup) {
	for job := range jobs {
		output := result{job, digits(job.randomNum)}
		results <- output
	}

	wg.Done()
}

func createWorkerPool(noOfWorkers int) {
	var wg sync.WaitGroup

	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg)
	}

	wg.Wait()
	close(results)
}

func allocate(noOfJobs int) {
	for i := 0; i < noOfJobs; i++ {
		randomNum := rand.Intn(1000000)
		job := job{i, randomNum}
		jobs <- job
	}

	close(jobs)
}

func showResults(done chan bool) {
	for result := range results {
		fmt.Printf("Job ID %3d, input random number %7d, sum of digits: %2d\n",
			result.j.id, result.j.randomNum, result.sumOfDigits)
	}

	done <- true
}

func main() {
	start := time.Now()
	noOfJobs := 5000
	go allocate(noOfJobs)

	done := make(chan bool)
	go showResults(done)

	createWorkerPool(threadLimit)
	<-done

	diff := time.Now().Sub(start)
	fmt.Printf("\n# of jobs: %3d\nWorkers in pool: %2d\nTotal time taken: %.3f seconds\n",
		noOfJobs, threadLimit, diff.Seconds())
}
