package pool

import (
	"sync"

	"github.com/rouzbehsbz/spenta/internal/share"
)

type Job struct {
	task  func()
	chunk int
}

func NewSliceJobs(len, chunkCount, chunkSize int, wg *sync.WaitGroup, cb func(i int)) []Job {
	jobs := make([]Job, 0, chunkCount)

	for chunk := range chunkCount {
		start, end := share.ChunkIndexes(chunk, len, chunkSize)

		jobs = append(jobs, Job{
			task: func() {
				defer wg.Done()

				for i := start; i < end; i++ {
					cb(i)
				}
			},
			chunk: chunk,
		})
	}

	return jobs
}
