package pool

import (
	"sync"
)

type Job struct {
	task func()
	wg   *sync.WaitGroup

	errCh chan error
}

func NewSliceJobs(len, chunkCount, chunkSize int, wg *sync.WaitGroup, cb func(i int)) ([]Job, chan error) {
	jobs := make([]Job, 0, chunkCount)
	errCh := make(chan error, chunkCount)

	for chunk := range chunkCount {
		start, end := ChunkIndexes(chunk, len, chunkSize)

		jobs = append(jobs, Job{
			task: func() {
				for i := start; i < end; i++ {
					cb(i)
				}
			},
			wg:    wg,
			errCh: errCh,
		})
	}

	return jobs, errCh
}

func ChunkIndexes(chunkIdx, len, chunkSize int) (int, int) {
	start := chunkIdx * chunkSize
	end := min(start+chunkSize, len)

	return start, end
}
