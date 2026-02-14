package pool

import (
	"sync"
)

// Represents a unit of work executed by the worker pool.
type job struct {
	task func()

	jobsWg *sync.WaitGroup
	errCh  chan error
}

// Creates a new Job instance.
func newJob(task func(), jobsWg *sync.WaitGroup, errCh chan error) job {
	return job{
		task:   task,
		jobsWg: jobsWg,
		errCh:  errCh,
	}
}

// Recursively splits a workload into smaller chunks and
// submits them to the worker pool based on binary tree conquer and divide.
func SpawnJob(
	start,
	end,
	maxChunkSize,
	minChunkSize int,
	jobsWg *sync.WaitGroup,
	errCh chan error,
	cb func(start, end int),
) {
	length := end - start

	if length > maxChunkSize && length/2 >= minChunkSize {
		mid := start + length/2

		SpawnJob(start, mid, maxChunkSize, minChunkSize, jobsWg, errCh, cb)
		SpawnJob(mid, end, maxChunkSize, minChunkSize, jobsWg, errCh, cb)
		return
	}

	jobsWg.Add(1)

	job := newJob(func() {
		cb(start, end)
	}, jobsWg, errCh)

	SpentaPool().SendJob(job)
}
