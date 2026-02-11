package pool

import (
	"sync"

	"github.com/rouzbehsbz/spenta/internal/share"
)

var (
	_pool *Pool
	_once sync.Once
)

type Pool struct {
	jobs chan Job
}

func SpentaPool() *Pool {
	_once.Do(func() {
		workersCount := share.WorkersCount()

		_pool = &Pool{
			jobs: make(chan Job),
		}

		for range workersCount {
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
		job.task()
	}
}
