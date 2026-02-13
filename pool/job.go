package pool

import (
	"sync"
)

type Job struct {
	task func()
	wg   *sync.WaitGroup

	errCh chan error
}

func NewJob(task func(), wg *sync.WaitGroup, errCh chan error) Job {
	return Job{
		task:  task,
		wg:    wg,
		errCh: errCh,
	}
}

func SpawnJob(start, end, maxChunkSize int, wg *sync.WaitGroup, errCh chan error, cb func(i int)) {
	length := end - start
	if length > maxChunkSize {
		mid := start + length/2

		// TODO: Maybe we can improve performance by calling
		// them inside goroutines, but we must ensure that
		// sync.WaitGroup is incremented safely before
		// parIter.Wait() is called.
		SpawnJob(start, mid, maxChunkSize, wg, errCh, cb)
		SpawnJob(mid, end, maxChunkSize, wg, errCh, cb)
		return
	}

	wg.Add(1)

	job := NewJob(func() {
		for i := start; i < end; i++ {
			cb(i)
		}
	}, wg, errCh)

	SpentaPool().SendJob(job)
}
