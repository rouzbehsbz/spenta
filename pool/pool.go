package pool

import (
	"fmt"
	"runtime"
	"sync"
)

var (
	_pool *Pool
	_once sync.Once
)

// Represents a fixed-size worker pool.
type Pool struct {
	jobs chan job

	workers int
}

// Returns a singleton instance of Pool.
func SpentaPool() *Pool {
	_once.Do(func() {
		_pool = &Pool{
			jobs:    make(chan job, 100),
			workers: runtime.NumCPU(),
		}

		// Start worker goroutines.
		for range _pool.workers {
			go _pool.worker()
		}
	})

	return _pool
}

// Submits a job to the pool for execution.
func (p *Pool) SendJob(job job) {
	p.jobs <- job
}

// worker continuously consumes jobs from the job queue
// and executes them.
func (p *Pool) worker() {
	for j := range p.jobs {
		func(j job) {
			defer p.onJobEnd(j)
			j.task()
		}(j)
	}
}

// Handles job completion bookkeeping.
// It recovers from panics during job execution and forwards
// the recovered error to the job's error channel.
func (p *Pool) onJobEnd(job job) {
	if r := recover(); r != nil {
		//TODO: maybe we can produce better errors.
		job.errCh <- fmt.Errorf("%v", r)

		// TODO: need to spawn a new worker after the last one dies
		// but we should ensure its safe.
	}

	job.jobsWg.Done()
}
