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

type Pool struct {
	jobs chan Job

	workers int
}

func SpentaPool() *Pool {
	_once.Do(func() {
		_pool = &Pool{
			jobs:    make(chan Job, 100),
			workers: runtime.NumCPU(),
		}

		for range _pool.workers {
			go _pool.worker()
		}
	})

	return _pool
}

func (p *Pool) SendJobs(jobs []Job) {
	for _, job := range jobs {
		p.jobs <- job
	}
}

func (p *Pool) worker() {
	for job := range p.jobs {
		func(j Job) {
			defer p.onJobEnd(j)
			j.task()
		}(job)
	}
}

func (p *Pool) onJobEnd(job Job) {
	if r := recover(); r != nil {
		//TODO: maybe we can produce better errors
		job.errCh <- fmt.Errorf("%v", r)

		// TODO: need to spawn a new worker after the last one dies
		// but we should ensure its safe and wont spawn infinitely
	}

	job.wg.Done()
}
